package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"lunar-rockets/domain"
	"lunar-rockets/usecase"
)

// MessageController handles HTTP requests for rocket messages
type MessageController struct {
	rocketMessageUsecase usecase.RocketMessageUsecase
}

// NewMessageController creates a new message controller
func NewMessageController(rocketMessageUsecase usecase.RocketMessageUsecase) *MessageController {
	return &MessageController{
		rocketMessageUsecase: rocketMessageUsecase,
	}
}

// @Summary Receive a message
// @Description Process and store a new rocket message
// @Tags messages
// @Accept json
// @Produce json
// @Param message body domain.RocketMessage true "Message to be processed"
// @Success 202 {object} map[string]string "Message accepted"
// @Failure 400 {string} string "Invalid request"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /messages [post]
func (c *MessageController) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var message domain.RocketMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		log.Printf("Error decoding message: %v", err)
		http.Error(w, "Invalid message format", http.StatusBadRequest)
		return
	}

	if message.Metadata.Channel == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	if message.Metadata.MessageType == "" {
		http.Error(w, "Missing message type", http.StatusBadRequest)
		return
	}

	if err := c.rocketMessageUsecase.ProcessMessage(r.Context(), &message); err != nil {
		log.Printf("Error processing message: %v", err)
		http.Error(w, "Failed to process message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}
