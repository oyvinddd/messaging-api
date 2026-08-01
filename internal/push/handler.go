package push

import (
	"net/http"
	"encoding/json"
	"github.com/oyvinddd/messaging-api/internal/response"
)

type (
	Handler struct {
		service Service
	}

	tokenRequest struct {
		PushToken string `json:"push_token"`
	}
)

func NewHandler(service Service) *Handler {
	&Handler{service: service}
}

func (h Handler) RegisterToken(w http.ResponseWriter, r *http.Request) {
	var request tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.WithError(w, err, http.StatusBadRequest)
		return
	}
	if err := h.service.Register(r.Context(), "", request.PushToken); err != nil {
		response.WithError(w, err, http.StatusBadRequest)
		return
	}
	response.WithStatusOnly(w, http.StatusCreated)
}

func (h Handler) DeleteTokens(w http.ResponseWriter, r *http.Request) {

}
