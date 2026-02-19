package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/stores"
)

// Repositories defines the contract for admin repository operations.
type Repositories interface {
	// Get retrieves a repository by primary key.
	Get(ctx context.Context, key string) (*models.Repository, error)
	// Set creates or updates a repository.
	Set(ctx context.Context, key string, repo *models.Repository) error
	// Delete removes a repository by primary key.
	Delete(ctx context.Context, key string) error
	// List retrieves repositories with optional filtering and pagination.
	List(ctx context.Context, filter *stores.RepositoryFilter, limit, offset int) ([]*models.Repository, error)
	// Count returns the total number of repositories matching the filter.
	Count(ctx context.Context, filter *stores.RepositoryFilter) (int, error)
}
