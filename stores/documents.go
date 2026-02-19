package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Documents provides database and vector access for document records.
type Documents struct {
	*sum.Database[models.Document]
}

// NewDocuments creates a new documents store.
func NewDocuments(db *sqlx.DB, renderer astql.Renderer) (*Documents, error) {
	database, err := sum.NewDatabase[models.Document](db, "documents", renderer)
	if err != nil {
		return nil, err
	}
	return &Documents{Database: database}, nil
}

// ListByUserRepoAndTag retrieves all documents for a version.
func (s *Documents) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag})
}

// GetByUserRepoTagAndPath retrieves a document by natural identifiers.
func (s *Documents) GetByUserRepoTagAndPath(ctx context.Context, userID int64, owner, repoName, tag, path string) (*models.Document, error) {
	return s.Select().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("path", "=", "path").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag, "path": path})
}

// FindSimilar finds documents similar to the given vector across a user's packages.
// Used for "more like this" queries.
func (s *Documents) FindSimilar(ctx context.Context, userID int64, vector []float32, limit int) ([]*models.Document, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		OrderByExpr("vector", "<=>", "query_vec", "ASC").
		Limit(limit).
		Exec(ctx, map[string]any{
			"user_id":   userID,
			"query_vec": vector,
		})
}

// FindSimilarInVersion finds documents similar to the given vector within a specific version.
func (s *Documents) FindSimilarInVersion(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Document, error) {
	return s.Query().
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
			"query_vec": vector,
		})
}
