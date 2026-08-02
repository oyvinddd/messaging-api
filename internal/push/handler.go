package push

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/oyvinddd/messaging-api/internal/response"
)

type (
	Handler struct {
		service Service
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) RegisterToken(w http.ResponseWriter, r *http.Request) {
	var token Token 
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		response.WithError(w, err, http.StatusBadRequest)
		return
	}

	id := uuid.New() //TODO: fix

	if err := h.service.RegisterToken(r.Context(), id, token); err != nil {
		response.WithError(w, err, http.StatusBadRequest)
		return
	}
	response.WithStatusOnly(w, http.StatusCreated)
}

func (h Handler) DeleteTokens(w http.ResponseWriter, r *http.Request) {

}
