package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Versions defines the contract for version storage operations.
type Versions interface {
	// Get retrieves a version by primary key.
	Get(ctx context.Context, key string) (*models.Version, error)
	// Set creates or updates a version.
	Set(ctx context.Context, key string, version *models.Version) error
	// ListByUserAndRepo retrieves all versions for a repository.
	ListByUserAndRepo(ctx context.Context, userID int64, owner, repoName string) ([]*models.Version, error)
	// GetByUserRepoAndTag retrieves a version by natural identifiers.
	GetByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) (*models.Version, error)
	// UpdateStatus updates the ingestion status of a version.
	UpdateStatus(ctx context.Context, id int64, status models.VersionStatus, versionErr *string) (*models.Version, error)
}
