package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"lunar-rockets/domain"
	"lunar-rockets/usecase"
)

//type RocketController interface {
//	GetRocket(w http.ResponseWriter, r *http.Request)
//	ListRockets(w http.ResponseWriter, r *http.Request)
//}

type RocketController struct {
	rocketUseCase usecase.RocketUseCase
}

func NewRocketController(rocketUseCase usecase.RocketUseCase) *RocketController {
	return &RocketController{
		rocketUseCase: rocketUseCase,
	}
}

// @Summary Get a specific rocket
// @Description Retrieve details of a specific rocket by its channel ID
// @Tags rockets
// @Accept json
// @Produce json
// @Param channel path string true "Rocket Channel ID"
// @Success 200 {object} domain.Rocket
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Rocket not found"
// @Router /rockets/{channel} [get]
func (c *RocketController) GetRocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	channel := pathParts[len(pathParts)-1]
	if channel == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	rocket, err := c.rocketUseCase.GetRocket(r.Context(), channel)
	if err != nil {
		log.Printf("Error getting rocket: %v", err)
		if errors.Is(err, domain.ErrRocketNotFound) {
			http.Error(w, "Rocket not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get rocket state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rocket)
}

// @Summary List all rockets
// @Description Retrieve a list of all available rockets with optional sorting
// @Tags rockets
// @Accept json
// @Produce json
// @Param sort query string false "Sort field ('channel','type','speed','mission','status')"
// @Param order query string false "Sort order ('asc' or 'desc')"
// @Success 200 {array} domain.Rocket
// @Router /rockets [get]
func (c *RocketController) ListRockets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sortBy := strings.ToLower(r.URL.Query().Get("sort"))
	order := strings.ToUpper(r.URL.Query().Get("order"))

	rockets, err := c.rocketUseCase.ListRockets(r.Context(), sortBy, order)
	if err != nil {
		log.Printf("Error listing rockets: %v", err)
		http.Error(w, "Failed to get rockets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rockets)
}
