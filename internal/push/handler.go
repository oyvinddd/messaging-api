package push

import (
	"net/http"
	"encoding/json"
	xauth "github.com/oyvinddd/xhttp/auth"
	"github.com/oyvinddd/xhttp/response"
)

type (
	Handler struct {
		service Service
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendPush(w http.ResponseWriter, r *http.Request) {
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	if err := h.service.Send(r.Context(), message); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.StatusCode(w, http.StatusOK)
}

func (h *Handler) RegisterToken(w http.ResponseWriter, r *http.Request) {
	// we can safely ignore error handling since handler is locked down by mw
	claims, _ := xauth.GetAccessTokenClaims(r.Context())

	var token DeviceToken 
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterToken(r.Context(), claims.ID, token); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.StatusCode(w, http.StatusCreated)
}

func (h *Handler) DeleteTokens(w http.ResponseWriter, r *http.Request) {
	// FIXME: this should prbably be secured by service key and not regular auth MW
	// we can safely ignore error handling since handler is locked down by mw
	claims, _ := xauth.GetAccessTokenClaims(r.Context())
	if err := h.service.DeleteTokens(r.Context(), claims.ID); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.StatusCode(w, http.StatusNoContent)
}

func (h *Handler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	// FIXME: this should prbably be secured by service key and not regular auth MW
	// we can safely ignore error handling since handler is locked down by mw
	/*
	claims, _ := xauth.GetAccessTokenClaims(r.Context())
	if err := h.service.DeleteTokens(r.Context(), claims.ID); err != nil {
		response.WithError(w, err, http.StatusBadRequest)
		return
	}
	response.WithStatusOnly(w, http.StatusNoContent)
	*/
}

