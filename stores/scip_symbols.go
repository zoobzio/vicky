package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// SCIPSymbols provides database access for SCIP symbol records.
type SCIPSymbols struct {
	*sum.Database[models.SCIPSymbol]
}

// NewSCIPSymbols creates a new SCIP symbols store.
func NewSCIPSymbols(db *sqlx.DB, renderer astql.Renderer) (*SCIPSymbols, error) {
	database, err := sum.NewDatabase[models.SCIPSymbol](db, "scip_symbols", renderer)
	if err != nil {
		return nil, err
	}
	return &SCIPSymbols{Database: database}, nil
}

// ListByUserRepoAndTag retrieves all SCIP symbols for a version.
func (s *SCIPSymbols) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPSymbol, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag})
}

// ListByDocument retrieves all SCIP symbols for a document.
func (s *SCIPSymbols) ListByDocument(ctx context.Context, documentID int64) ([]*models.SCIPSymbol, error) {
	return s.Executor().Soy().Query().
		Where("document_id", "=", "document_id").
		Exec(ctx, map[string]any{"document_id": documentID})
}

// GetBySymbol retrieves a SCIP symbol by its qualified identifier within a version.
func (s *SCIPSymbols) GetBySymbol(ctx context.Context, userID int64, owner, repoName, tag, symbol string) (*models.SCIPSymbol, error) {
	return s.Executor().Soy().Select().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("symbol", "=", "symbol").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag, "symbol": symbol})
}

// ListByKind retrieves SCIP symbols of a specific kind within a version.
func (s *SCIPSymbols) ListByKind(ctx context.Context, userID int64, owner, repoName, tag string, kind models.SCIPSymbolKind) ([]*models.SCIPSymbol, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("kind", "=", "kind").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag, "kind": kind})
}

// ListByEnclosingSymbol retrieves SCIP symbols with the given enclosing symbol.
func (s *SCIPSymbols) ListByEnclosingSymbol(ctx context.Context, userID int64, owner, repoName, tag, enclosingSymbol string) ([]*models.SCIPSymbol, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("enclosing_symbol", "=", "enclosing_symbol").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag, "enclosing_symbol": enclosingSymbol})
}
