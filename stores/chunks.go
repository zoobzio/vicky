// Package stores implements vicky's data store contracts.
package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Chunks provides database and vector access for chunk records.
type Chunks struct {
	*sum.Database[models.Chunk]
}

// NewChunks creates a new chunks store.
func NewChunks(db *sqlx.DB, renderer astql.Renderer) (*Chunks, error) {
	database, err := sum.NewDatabase[models.Chunk](db, "chunks", renderer)
	if err != nil {
		return nil, err
	}
	return &Chunks{Database: database}, nil
}

// ListByUserRepoAndTag retrieves all chunks for a version.
func (s *Chunks) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag})
}

// ListByUserRepoTagAndPath retrieves all chunks for a document.
func (s *Chunks) ListByUserRepoTagAndPath(ctx context.Context, userID int64, owner, repoName, tag, path string) ([]*models.Chunk, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("path", "=", "path").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag, "path": path})
}

// Search performs semantic search across chunks in a version.
func (s *Chunks) Search(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Chunk, error) {
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

// SearchByKind performs semantic search filtered by chunk kind.
func (s *Chunks) SearchByKind(ctx context.Context, userID int64, owner, repoName, tag string, kind models.ChunkKind, vector []float32, limit int) ([]*models.Chunk, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Where("kind", "=", "kind").
		OrderByExpr("vector", "<=>", "query_vec", "ASC").
		Limit(limit).
		Exec(ctx, map[string]any{
			"user_id":   userID,
			"owner":     owner,
			"repo_name": repoName,
			"tag":       tag,
			"kind":      kind,
			"query_vec": vector,
		})
}
