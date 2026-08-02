package push

import (
	"context"
	"github.com/google/uuid"
)

type (
	Service interface {
		// RegisterToken registers a new push token in the db
		RegisterToken(ctx context.Context, id uuid.UUID, token Token) error
		// DeleteTokens deletes all tokens registered on a given account
		DeleteTokens(ctx context.Context, id uuid.UUID) error
	}

	service struct {
		repository Repository
	}
)

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) RegisterToken(ctx context.Context, id uuid.UUID, token Token) error {
	return s.repository.RegisterToken(ctx, id, token)
}

func (s *service) DeleteTokens(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteTokens(ctx, id)
}

