package mailer

import (
	"net/http"
	"encoding/json"
	"github.com/oyvinddd/xhttp/response"
)

type (
	Handler struct {
		service Service
	}

	Message struct {
		// Sender the email of the sender
		Sender string `json:"sender"`
		// Recipient the email of the recipient
		Recipient string `json:"recipient"`
		// Subject the subject of the email
		Subject string `json:"subject"`
		// Body the content of the email
		Body string `json:"body"`
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		response.Error(w, malformedEmailError, http.StatusBadRequest)
		return
	}
	if err := h.service.Send(r.Context(), message); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.StatusCode(w, http.StatusOK)
}

