package push

import (
	"context"
	"github.com/google/uuid"
)

type (
	Service interface {
		// Send sends a push message
		Send(ctx context.Context, message Message) error
		// RegisterToken registers a new push token in the db
		RegisterToken(ctx context.Context, id uuid.UUID, token DeviceToken) error
		// DeleteTokens deletes all tokens registered on a given account
		DeleteTokens(ctx context.Context, id uuid.UUID) error
		// DeleteToken deletes token with a given id and account id
		DeleteToken(ctx context.Context, accountID, tokenID uuid.UUID) error
	}

	service struct {
		repository 	Repository
		provider 	Provider
	}
)

func NewService(repository Repository, provider Provider) Service {
	return &service{repository: repository, provider: provider}
}

func (s *service) Send(ctx context.Context, message Message) error {
	// step 1. -- fetch the device tokens corresponding to the account id
	tokens, err := s.repository.GetTokens(ctx, message.RecipientID)
	if err != nil {
		return err
	}
	// step 2. -- send message to FCM with the correct device ID
	if err := s.provider.Send(ctx, tokens, message); err != nil {
		// TODO: there might be a case here where FCM says I need to remove the token, so handle it
		return err
	}
	return nil
}

func (s *service) RegisterToken(ctx context.Context, id uuid.UUID, token DeviceToken) error {
	return s.repository.RegisterToken(ctx, id, token)
}

func (s *service) DeleteTokens(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteTokens(ctx, id)
}

func (s *service) DeleteToken(ctx context.Context, accountID, tokenID uuid.UUID) error {
	return s.repository.DeleteToken(ctx, accountID, tokenID)
}

