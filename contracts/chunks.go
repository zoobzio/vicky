// Package contracts defines interfaces for vicky's data stores.
package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Chunks defines the contract for chunk storage operations.
type Chunks interface {
	// Get retrieves a chunk by primary key.
	Get(ctx context.Context, key string) (*models.Chunk, error)
	// Set creates or updates a chunk.
	Set(ctx context.Context, key string, chunk *models.Chunk) error
	// ListByUserRepoAndTag retrieves all chunks for a version.
	ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error)
	// ListByUserRepoTagAndPath retrieves all chunks for a document.
	ListByUserRepoTagAndPath(ctx context.Context, userID int64, owner, repoName, tag, path string) ([]*models.Chunk, error)
	// Search performs semantic search across chunks in a version.
	Search(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Chunk, error)
	// SearchByKind performs semantic search filtered by chunk kind.
	SearchByKind(ctx context.Context, userID int64, owner, repoName, tag string, kind models.ChunkKind, vector []float32, limit int) ([]*models.Chunk, error)
}
