package push

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Repository interface {
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

func (r *repository) RegisterToken(ctx context.Context, id uuid.UUID, token Token) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO push.tokens (
			account_id,
			token,
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

