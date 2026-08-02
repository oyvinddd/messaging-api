package mailer

import (
	"context"
)

type (
	Service interface {
		Send(ctx context.Context, recipient, subject, body string) error
	}

	postmarkService struct {
		apiKey string
	}
)

func NewPostmarkService(apiKey string) Service {
	return &postmarkService{apiKey: apiKey}
}

func (s *postmarkService) Send(ctx context.Context, recipient, subject, body string) error {
	return nil
}

