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
	messageEventUsecase usecase.MessageEventUsecase
}

// NewMessageController creates a new message controller
func NewMessageController(messageEventUsecase usecase.MessageEventUsecase) *MessageController {
	return &MessageController{
		messageEventUsecase: messageEventUsecase,
	}
}

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

	if err := c.messageEventUsecase.ProcessMessage(r.Context(), &message); err != nil {
		log.Printf("Error processing message: %v", err)
		http.Error(w, "Failed to process message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}
