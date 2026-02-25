package repository

import (
	"BudgetTracker/backend/internal/model"
	"BudgetTracker/backend/pkg/database"
	"context"
	"time"
)

type APIKeyRepository struct {
	db *database.Database
}

func NewAPIKeyRepository(db *database.Database) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) GetByExchange(ctx context.Context, exchange string) (*model.APIKey, error) {
	query := `
		SELECT id, exchange, api_key, api_secret, is_active, created_at, updated_at
		FROM api_keys
		WHERE exchange = $1
	`

	var apiKey model.APIKey
	err := r.db.Pool.QueryRow(ctx, query, exchange).Scan(
		&apiKey.ID,
		&apiKey.Exchange,
		&apiKey.APIKey,
		&apiKey.APISecret,
		&apiKey.IsActive,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (r *APIKeyRepository) Upsert(ctx context.Context, apiKey *model.APIKey) error {
	query := `
		INSERT INTO api_keys (exchange, api_key, api_secret, is_active, updated_at)
		VALUES ($1, $2, $3, true, $4)
		ON CONFLICT (exchange) DO UPDATE SET
			api_key = EXCLUDED.api_key,
			api_secret = EXCLUDED.api_secret,
			is_active = true,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Pool.Exec(ctx, query,
		apiKey.Exchange,
		apiKey.APIKey,
		apiKey.APISecret,
		time.Now(),
	)

	return err
}

func (r *APIKeyRepository) GetAll(ctx context.Context) ([]model.APIKey, error) {
	query := `
		SELECT id, exchange, api_key, api_secret, is_active, created_at, updated_at
		FROM api_keys
		ORDER BY exchange
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []model.APIKey
	for rows.Next() {
		var apiKey model.APIKey
		err := rows.Scan(
			&apiKey.ID,
			&apiKey.Exchange,
			&apiKey.APIKey,
			&apiKey.APISecret,
			&apiKey.IsActive,
			&apiKey.CreatedAt,
			&apiKey.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}
