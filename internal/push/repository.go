package push

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Repository interface {
		// RegisterToken registers a new push token in the db
		RegisterToken(ctx context.Context, accountID uuid.UUID, token DeviceToken) error
		// GetTokens fetches device tokens from db
		GetTokens(ctx context.Context, accountID uuid.UUID) ([]DeviceToken, error)
		// DeleteTokens deletes all tokens registered on a given account
		DeleteTokens(ctx context.Context, accountID uuid.UUID) error
		// DeleteToken deletes the token with a given ID
		DeleteToken(ctx context.Context, accountID, tokenID uuid.UUID) error
	}

	repository struct {
		db *pgxpool.Pool
	}
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) RegisterToken(ctx context.Context, accountID uuid.UUID, token DeviceToken) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO push.device_tokens (
			account_id,
			token,
			platform
		) VALUES ($1, $2, $3)
		ON CONFLICT (token) DO UPDATE
		SET
			account_id = EXCLUDED.account_id,
			platform = EXCLUDED.platform,
			updated_at = NOW()
	`, accountID, token.Value, token.Platform)

	return err
}

func (r *repository) GetTokens(ctx context.Context, accountID uuid.UUID) ([]DeviceToken, error) {
	const query = `SELECT id, token, platform FROM push.device_tokens WHERE account_id = $1`

	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]DeviceToken, 0)

	for rows.Next() {
		var token DeviceToken

		if err := rows.Scan(&token.ID, &token.Value, &token.Platform); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return tokens, nil
	}

	return tokens, nil
}

func (r *repository) DeleteTokens(ctx context.Context, accountID uuid.UUID) error {
	const query = `DELETE FROM push.device_tokens WHERE account_id = $1`
	_, err := r.db.Exec(ctx, query, accountID)
	return err
}

func (r *repository) DeleteToken(ctx context.Context, accountID, tokenID uuid.UUID) error {
	const query = `DELETE FROM push.device_tokens WHERE id = $1 AND account_id = $2`
	_, err := r.db.Exec(ctx, query, accountID, tokenID)
	return err
}

