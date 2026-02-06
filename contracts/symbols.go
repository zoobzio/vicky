package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Symbols defines the contract for symbol storage operations.
type Symbols interface {
	// Get retrieves a symbol by primary key.
	Get(ctx context.Context, key string) (*models.Symbol, error)
	// Set creates or updates a symbol.
	Set(ctx context.Context, key string, symbol *models.Symbol) error
	// ListByUserRepoAndTag retrieves all symbols for a version.
	ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error)
	// ListExportedByUserRepoAndTag retrieves all exported symbols for a version.
	ListExportedByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error)
	// FindRelated finds symbols related to the given document vector.
	FindRelated(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error)
	// FindRelatedExported finds exported symbols related to the given document vector.
	FindRelatedExported(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error)
}
