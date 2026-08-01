package email

import (
	"context"
)

type (
	Service interface {
		Send(recipient, subject, body string) error
	}

	postmarkService struct {
		apiKey string
	}
)

func NewPostmarkService(apiKey string) Service {
	return &postmarkService{apiKey: apiKey}
}

