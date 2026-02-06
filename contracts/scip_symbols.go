package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// SCIPSymbols defines the contract for SCIP symbol storage operations.
type SCIPSymbols interface {
	// Get retrieves a SCIP symbol by primary key.
	Get(ctx context.Context, key string) (*models.SCIPSymbol, error)
	// Set creates or updates a SCIP symbol.
	Set(ctx context.Context, key string, symbol *models.SCIPSymbol) error
	// ListByUserRepoAndTag retrieves all SCIP symbols for a version.
	ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPSymbol, error)
	// ListByDocument retrieves all SCIP symbols for a document.
	ListByDocument(ctx context.Context, documentID int64) ([]*models.SCIPSymbol, error)
	// GetBySymbol retrieves a SCIP symbol by its qualified identifier within a version.
	GetBySymbol(ctx context.Context, userID int64, owner, repoName, tag, symbol string) (*models.SCIPSymbol, error)
	// ListByKind retrieves SCIP symbols of a specific kind within a version.
	ListByKind(ctx context.Context, userID int64, owner, repoName, tag string, kind models.SCIPSymbolKind) ([]*models.SCIPSymbol, error)
	// ListByEnclosingSymbol retrieves SCIP symbols with the given enclosing symbol.
	ListByEnclosingSymbol(ctx context.Context, userID int64, owner, repoName, tag, enclosingSymbol string) ([]*models.SCIPSymbol, error)
}
