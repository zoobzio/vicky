package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Repositories defines the contract for repository storage operations.
type Repositories interface {
	// Get retrieves a repository by primary key.
	Get(ctx context.Context, key string) (*models.Repository, error)
	// Set creates or updates a repository.
	Set(ctx context.Context, key string, repo *models.Repository) error
	// ListByUserID retrieves all repositories for a user.
	ListByUserID(ctx context.Context, userID int64) ([]*models.Repository, error)
	// GetByUserAndGitHubID retrieves a repository by user and GitHub repo ID.
	GetByUserAndGitHubID(ctx context.Context, userID, githubID int64) (*models.Repository, error)
	// GetByUserOwnerAndName retrieves a repository by user, owner, and name.
	GetByUserOwnerAndName(ctx context.Context, userID int64, owner, name string) (*models.Repository, error)
}
