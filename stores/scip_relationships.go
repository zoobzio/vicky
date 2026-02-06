package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// SCIPRelationships provides database access for SCIP relationship records.
type SCIPRelationships struct {
	*sum.Database[models.SCIPRelationship]
}

// NewSCIPRelationships creates a new SCIP relationships store.
func NewSCIPRelationships(db *sqlx.DB, renderer astql.Renderer) (*SCIPRelationships, error) {
	database, err := sum.NewDatabase[models.SCIPRelationship](db, "scip_relationships", renderer)
	if err != nil {
		return nil, err
	}
	return &SCIPRelationships{Database: database}, nil
}

// ListBySymbol retrieves all relationships for a SCIP symbol.
func (s *SCIPRelationships) ListBySymbol(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error) {
	return s.Executor().Soy().Query().
		Where("scip_symbol_id", "=", "scip_symbol_id").
		Exec(ctx, map[string]any{"scip_symbol_id": symbolID})
}

// ListImplementations retrieves relationships where is_implementation is true.
func (s *SCIPRelationships) ListImplementations(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error) {
	return s.Executor().Soy().Query().
		Where("scip_symbol_id", "=", "scip_symbol_id").
		Where("is_implementation", "=", "is_implementation").
		Exec(ctx, map[string]any{"scip_symbol_id": symbolID, "is_implementation": true})
}

// ListByTargetSymbol retrieves relationships pointing to a target symbol.
func (s *SCIPRelationships) ListByTargetSymbol(ctx context.Context, targetSymbol string) ([]*models.SCIPRelationship, error) {
	return s.Executor().Soy().Query().
		Where("target_symbol", "=", "target_symbol").
		Exec(ctx, map[string]any{"target_symbol": targetSymbol})
}
