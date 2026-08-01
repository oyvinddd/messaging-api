package mailer

import (
	"net/http"
	"encoding/json"
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
	return Handler{service: service}
}



