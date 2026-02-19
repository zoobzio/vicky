package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Keys defines the contract for API key storage operations.
type Keys interface {
	// Get retrieves an API key by primary key.
	Get(ctx context.Context, key string) (*models.Key, error)
	// Set creates or updates an API key.
	Set(ctx context.Context, key string, apiKey *models.Key) error
	// Delete removes an API key by primary key.
	Delete(ctx context.Context, key string) error
	// GetByKeyHash retrieves an API key by its SHA-256 hash for authentication.
	GetByKeyHash(ctx context.Context, hash string) (*models.Key, error)
	// ListByUserID retrieves all API keys belonging to a user.
	ListByUserID(ctx context.Context, userID int64) ([]*models.Key, error)
	// UpdateLastUsed sets the last_used_at timestamp for a key to now.
	UpdateLastUsed(ctx context.Context, id int64) error
}
