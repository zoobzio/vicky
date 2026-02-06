package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Users defines the contract for user storage operations.
type Users interface {
	// Get retrieves a user by primary key.
	Get(ctx context.Context, key string) (*models.User, error)
	// Set creates or updates a user.
	Set(ctx context.Context, key string, user *models.User) error
	// GetByLogin retrieves a user by GitHub login.
	GetByLogin(ctx context.Context, login string) (*models.User, error)
}
