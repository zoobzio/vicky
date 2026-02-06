package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// SCIPRelationships defines the contract for SCIP relationship storage operations.
type SCIPRelationships interface {
	// Get retrieves a SCIP relationship by primary key.
	Get(ctx context.Context, key string) (*models.SCIPRelationship, error)
	// Set creates or updates a SCIP relationship.
	Set(ctx context.Context, key string, rel *models.SCIPRelationship) error
	// ListBySymbol retrieves all relationships for a SCIP symbol.
	ListBySymbol(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error)
	// ListImplementations retrieves relationships where is_implementation is true.
	ListImplementations(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error)
	// ListByTargetSymbol retrieves relationships pointing to a target symbol.
	ListByTargetSymbol(ctx context.Context, targetSymbol string) ([]*models.SCIPRelationship, error)
}
