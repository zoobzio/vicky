package stores

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Keys provides database access for API key records.
type Keys struct {
	*sum.Database[models.Key]
}

// NewKeys creates a new keys store.
func NewKeys(db *sqlx.DB, renderer astql.Renderer) (*Keys, error) {
	database, err := sum.NewDatabase[models.Key](db, "keys", renderer)
	if err != nil {
		return nil, err
	}
	return &Keys{Database: database}, nil
}

// GetByKeyHash retrieves an API key by its SHA-256 hash.
// Used during authentication to look up the key record.
func (s *Keys) GetByKeyHash(ctx context.Context, hash string) (*models.Key, error) {
	return s.Select().
		Where("key_hash", "=", "key_hash").
		Exec(ctx, map[string]any{"key_hash": hash})
}

// ListByUserID retrieves all API keys belonging to a user, ordered by creation time.
func (s *Keys) ListByUserID(ctx context.Context, userID int64) ([]*models.Key, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"user_id": userID})
}

// UpdateLastUsed sets the last_used_at timestamp for a key to now.
// This is a targeted single-field update to avoid touching other columns.
func (s *Keys) UpdateLastUsed(ctx context.Context, id int64) error {
	now := time.Now()
	_, err := s.Modify().
		Set("last_used_at", "last_used_at").
		Where("id", "=", "id").
		Exec(ctx, map[string]any{
			"id":           id,
			"last_used_at": &now,
		})
	return err
}
