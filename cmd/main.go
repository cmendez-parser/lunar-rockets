package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lunar-rockets/configs"
	"lunar-rockets/db/sqlite"
	httproute "lunar-rockets/http"
	"lunar-rockets/http/controller"
	"lunar-rockets/repository"
	"lunar-rockets/usecase"
)

// @title Lunar Rockets API
// @version 1.0
// @description API for managing lunar rockets and their messages

// @host localhost:8088
// @BasePath /
// @schemes http

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := sqlite.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	rocketRepo := repository.NewRocketRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	rocketStateUsecase := usecase.NewRocketStateUsecase(rocketRepo, messageRepo)
	messageProcessor := usecase.NewRocketMessageUsecase(rocketRepo, messageRepo, rocketStateUsecase)
	rocketUseCase := usecase.NewRocketUseCase(rocketRepo)

	messageController := controller.NewMessageController(messageProcessor)
	rocketController := controller.NewRocketController(rocketUseCase)

	router := httproute.NewRouter(messageController, rocketController)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
