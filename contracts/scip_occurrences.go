package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// SCIPOccurrences defines the contract for SCIP occurrence storage operations.
type SCIPOccurrences interface {
	// Get retrieves a SCIP occurrence by primary key.
	Get(ctx context.Context, key string) (*models.SCIPOccurrence, error)
	// Set creates or updates a SCIP occurrence.
	Set(ctx context.Context, key string, occurrence *models.SCIPOccurrence) error
	// ListByUserRepoAndTag retrieves all SCIP occurrences for a version.
	ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPOccurrence, error)
	// ListByDocument retrieves all SCIP occurrences for a document.
	ListByDocument(ctx context.Context, documentID int64) ([]*models.SCIPOccurrence, error)
	// ListBySymbol retrieves all occurrences of a specific symbol within a version.
	ListBySymbol(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error)
	// ListDefinitions retrieves occurrences marked as definitions for a symbol.
	ListDefinitions(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error)
	// ListReferences retrieves occurrences that are references (not definitions) for a symbol.
	ListReferences(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error)
}
