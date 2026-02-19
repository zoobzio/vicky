package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/stores"
)

// Users defines the contract for admin user operations.
type Users interface {
	// Get retrieves a user by primary key.
	Get(ctx context.Context, key string) (*models.User, error)
	// Set creates or updates a user.
	Set(ctx context.Context, key string, user *models.User) error
	// Delete removes a user by primary key.
	Delete(ctx context.Context, key string) error
	// List retrieves users with optional filtering and pagination.
	List(ctx context.Context, filter *stores.UserFilter, limit, offset int) ([]*models.User, error)
	// Count returns the total number of users matching the filter.
	Count(ctx context.Context, filter *stores.UserFilter) (int, error)
}
