package http

import (
	"net/http"
	"strings"

	_ "lunar-rockets/docs"
	"lunar-rockets/http/controller"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	messageController *controller.MessageController
	rocketController  *controller.RocketController
}

func NewRouter(messageController *controller.MessageController, rocketController *controller.RocketController) http.Handler {
	router := &Router{
		messageController: messageController,
		rocketController:  rocketController,
	}

	return router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if strings.HasPrefix(path, "/swagger/") {
		httpSwagger.WrapHandler.ServeHTTP(w, req)
		return
	}

	if req.Method == http.MethodPost && path == "/messages" {
		r.messageController.ReceiveMessage(w, req)
		return
	}

	if req.Method == http.MethodGet && path == "/rockets" {
		r.rocketController.ListRockets(w, req)
		return
	}

	if req.Method == http.MethodGet && strings.HasPrefix(path, "/rockets/") {
		r.rocketController.GetRocket(w, req)
		return
	}

	http.NotFound(w, req)
}
