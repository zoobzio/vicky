package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Symbols provides database and vector access for symbol records.
type Symbols struct {
	*sum.Database[models.Symbol]
}

// NewSymbols creates a new symbols store.
func NewSymbols(db *sqlx.DB, renderer astql.Renderer) (*Symbols, error) {
	database, err := sum.NewDatabase[models.Symbol](db, "symbols", renderer)
	if err != nil {
		return nil, err
	}
	return &Symbols{Database: database}, nil
}

// ListByUserRepoAndTag retrieves all symbols for a version.
func (s *Symbols) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag})
}

// ListExportedByUserRepoAndTag retrieves all exported symbols for a version.
func (s *Symbols) ListExportedByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("exported", "=", "exported").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag, "exported": true})
}

// FindRelated finds symbols related to the given document vector.
// Used for "mentioned here" queries - finding API symbols relevant to a document.
func (s *Symbols) FindRelated(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		OrderByExpr("vector", "<=>", "query_vec", "ASC").
		Limit(limit).
		Exec(ctx, map[string]any{
			"user_id":   userID,
			"owner":     owner,
			"repo_name": repoName,
			"tag":       tag,
			"query_vec": docVector,
		})
}

// FindRelatedExported finds exported symbols related to the given document vector.
// Filters to only public API symbols.
func (s *Symbols) FindRelatedExported(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("exported", "=", "exported").
		OrderByExpr("vector", "<=>", "query_vec", "ASC").
		Limit(limit).
		Exec(ctx, map[string]any{
			"user_id":   userID,
			"owner":     owner,
			"repo_name": repoName,
			"tag":       tag,
			"exported":  true,
			"query_vec": docVector,
		})
}
