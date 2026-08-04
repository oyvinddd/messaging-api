package push

import (
	"errors"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Repository interface {
		// GetTokens fetches device token from db
		GetToken(ctx context.Context, id uuid.UUID) (string, error)
		// RegisterToken registers a new push token in the db
		RegisterToken(ctx context.Context, id uuid.UUID, token Token) error
		// DeleteTokens deletes all tokens registered on a given account
		DeleteTokens(ctx context.Context, id uuid.UUID) error
	}

	repository struct {
		db *pgxpool.Pool
	}
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetToken(ctx context.Context, id uuid.UUID) (string, error) {
	query := `SELECT device_token FROM push.device_tokens WHERE account_id = $1 LIMIT 1`

	var deviceToken string

	err := r.db.QueryRow(ctx, query, id).Scan(&deviceToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", deviceTokenNotFound 
		}
		return "", err
	}

	return deviceToken, nil
}

func (r *repository) RegisterToken(ctx context.Context, id uuid.UUID, token Token) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO push.device_tokens (
			account_id,
			device_token,
			platform
		) VALUES ($1, $2, $3)
		ON CONFLICT (token) DO UPDATE
		SET
			account_id = EXCLUDED.account_id,
			platform = EXCLUDED.platform,
			updated_at = NOW()
	`, id, token.Value, token.Platform)

	return err
}

func (r *repository) DeleteTokens(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM push.tokens WHERE account_id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

