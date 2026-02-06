package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Documents defines the contract for document storage operations.
type Documents interface {
	// Get retrieves a document by primary key.
	Get(ctx context.Context, key string) (*models.Document, error)
	// Set creates or updates a document.
	Set(ctx context.Context, key string, doc *models.Document) error
	// ListByUserRepoAndTag retrieves all documents for a version.
	ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error)
	// GetByUserRepoTagAndPath retrieves a document by natural identifiers.
	GetByUserRepoTagAndPath(ctx context.Context, userID int64, owner, repoName, tag, path string) (*models.Document, error)
	// FindSimilar finds documents similar to the given vector across a user's packages.
	FindSimilar(ctx context.Context, userID int64, vector []float32, limit int) ([]*models.Document, error)
	// FindSimilarInVersion finds documents similar to the given vector within a specific version.
	FindSimilarInVersion(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Document, error)
}
