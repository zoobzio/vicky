package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// SCIPOccurrences provides database access for SCIP occurrence records.
type SCIPOccurrences struct {
	*sum.Database[models.SCIPOccurrence]
}

// NewSCIPOccurrences creates a new SCIP occurrences store.
func NewSCIPOccurrences(db *sqlx.DB, renderer astql.Renderer) (*SCIPOccurrences, error) {
	database, err := sum.NewDatabase[models.SCIPOccurrence](db, "scip_occurrences", renderer)
	if err != nil {
		return nil, err
	}
	return &SCIPOccurrences{Database: database}, nil
}

// ListByUserRepoAndTag retrieves all SCIP occurrences for a version.
func (s *SCIPOccurrences) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPOccurrence, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag})
}

// ListByDocument retrieves all SCIP occurrences for a document.
func (s *SCIPOccurrences) ListByDocument(ctx context.Context, documentID int64) ([]*models.SCIPOccurrence, error) {
	return s.Executor().Soy().Query().
		Where("document_id", "=", "document_id").
		Exec(ctx, map[string]any{"document_id": documentID})
}

// ListBySymbol retrieves all occurrences of a specific symbol within a version.
func (s *SCIPOccurrences) ListBySymbol(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("symbol", "=", "symbol").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag, "symbol": symbol})
}

// ListDefinitions retrieves occurrences marked as definitions for a symbol.
// Filters by checking the Definition bit in symbol_roles bitmask.
func (s *SCIPOccurrences) ListDefinitions(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
	all, err := s.ListBySymbol(ctx, userID, owner, repoName, tag, symbol)
	if err != nil {
		return nil, err
	}

	var definitions []*models.SCIPOccurrence
	for _, occ := range all {
		if occ.SymbolRoles&models.SCIPSymbolRoleDefinition != 0 {
			definitions = append(definitions, occ)
		}
	}
	return definitions, nil
}

// ListReferences retrieves occurrences that are references (not definitions) for a symbol.
// Filters by checking the Definition bit is NOT set in symbol_roles bitmask.
func (s *SCIPOccurrences) ListReferences(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
	all, err := s.ListBySymbol(ctx, userID, owner, repoName, tag, symbol)
	if err != nil {
		return nil, err
	}

	var references []*models.SCIPOccurrence
	for _, occ := range all {
		if occ.SymbolRoles&models.SCIPSymbolRoleDefinition == 0 {
			references = append(references, occ)
		}
	}
	return references, nil
}
