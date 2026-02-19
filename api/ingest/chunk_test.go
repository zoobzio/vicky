//go:build testing

package ingest

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/zoobzio/grub"
	"github.com/zoobzio/vicky/external/chunker"
	"github.com/zoobzio/vicky/models"
	vickytest "github.com/zoobzio/vicky/testing"
)

func TestChunkStage(t *testing.T) {
	docs := []*models.Document{
		vickytest.NewDocument(t, 1, "main.go"),
		vickytest.NewDocument(t, 2, "utils.go"),
	}

	var mu sync.Mutex
	var chunkCount int

	md := &vickytest.MockDocuments{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error) {
			return docs, nil
		},
	}
	mb := &vickytest.MockBlobs{
		OnGetByPath: func(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error) {
			return &grub.Object[models.Blob]{
				Key:  path,
				Data: models.Blob{Path: path, Content: "package main\n\nfunc Hello() {}\n", Owner: owner, Repo: repo, Tag: tag},
			}, nil
		},
	}
	mch := &vickytest.MockChunker{
		OnChunk: func(ctx context.Context, language string, filename string, content []byte) ([]chunker.Result, error) {
			return []chunker.Result{
				{Content: "func Hello() {}", Kind: models.ChunkKindFunction, StartLine: 3, EndLine: 3},
			}, nil
		},
	}
	mc := &vickytest.MockChunks{
		OnSet: func(ctx context.Context, key string, chunk *models.Chunk) error {
			mu.Lock()
			defer mu.Unlock()
			chunkCount++
			return nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithDocuments(md),
		vickytest.WithBlobs(mb),
		vickytest.WithChunker(mch),
		vickytest.WithChunks(mc),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
	)

	job := vickytest.NewJob(t)
	result, err := chunkStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Stage != models.JobStageChunk {
		t.Errorf("Stage = %q, want %q", result.Stage, models.JobStageChunk)
	}
	if result.ItemsTotal != 2 {
		t.Errorf("ItemsTotal = %d, want 2", result.ItemsTotal)
	}
	if result.ItemsProcessed != 2 {
		t.Errorf("ItemsProcessed = %d, want 2", result.ItemsProcessed)
	}

	mu.Lock()
	defer mu.Unlock()
	// 2 documents * 1 chunk each = 2 chunks
	if chunkCount != 2 {
		t.Errorf("chunks stored = %d, want 2", chunkCount)
	}
}

func TestChunkStage_NoDocs(t *testing.T) {
	md := &vickytest.MockDocuments{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error) {
			return nil, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithDocuments(md),
		vickytest.WithBlobs(&vickytest.MockBlobs{}),
		vickytest.WithChunker(&vickytest.MockChunker{}),
		vickytest.WithChunks(&vickytest.MockChunks{}),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
	)

	job := vickytest.NewJob(t)
	result, err := chunkStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ItemsTotal != 0 {
		t.Errorf("ItemsTotal = %d, want 0", result.ItemsTotal)
	}
}

func TestChunkStage_ChunkerNotSupported(t *testing.T) {
	docs := []*models.Document{
		vickytest.NewDocument(t, 1, "main.go"),
	}

	md := &vickytest.MockDocuments{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error) {
			return docs, nil
		},
	}
	mb := &vickytest.MockBlobs{
		OnGetByPath: func(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error) {
			return &grub.Object[models.Blob]{
				Key:  path,
				Data: models.Blob{Path: path, Content: "content", Owner: owner, Repo: repo, Tag: tag},
			}, nil
		},
	}
	mch := &vickytest.MockChunker{
		OnSupports: func(language string) bool {
			return false
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithDocuments(md),
		vickytest.WithBlobs(mb),
		vickytest.WithChunker(mch),
		vickytest.WithChunks(&vickytest.MockChunks{}),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
	)

	job := vickytest.NewJob(t)
	result, err := chunkStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Document was processed but chunker skipped it
	if result.ItemsProcessed != 1 {
		t.Errorf("ItemsProcessed = %d, want 1", result.ItemsProcessed)
	}
}

func TestChunkStage_BlobFetchError(t *testing.T) {
	docs := []*models.Document{
		vickytest.NewDocument(t, 1, "main.go"),
	}

	md := &vickytest.MockDocuments{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error) {
			return docs, nil
		},
	}
	mb := &vickytest.MockBlobs{
		OnGetByPath: func(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error) {
			return nil, fmt.Errorf("blob not found")
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithDocuments(md),
		vickytest.WithBlobs(mb),
		vickytest.WithChunker(&vickytest.MockChunker{}),
		vickytest.WithChunks(&vickytest.MockChunks{}),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
	)

	job := vickytest.NewJob(t)
	_, err := chunkStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "blob") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "blob")
	}
}

func TestChunkStage_ChunkStoreError(t *testing.T) {
	docs := []*models.Document{
		vickytest.NewDocument(t, 1, "main.go"),
	}

	md := &vickytest.MockDocuments{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error) {
			return docs, nil
		},
	}
	mb := &vickytest.MockBlobs{
		OnGetByPath: func(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error) {
			return &grub.Object[models.Blob]{
				Key:  path,
				Data: models.Blob{Path: path, Content: "content", Owner: owner, Repo: repo, Tag: tag},
			}, nil
		},
	}
	mch := &vickytest.MockChunker{
		OnChunk: func(ctx context.Context, language string, filename string, content []byte) ([]chunker.Result, error) {
			return []chunker.Result{
				{Content: "func main() {}", Kind: models.ChunkKindFunction, StartLine: 1, EndLine: 1},
			}, nil
		},
	}
	mc := &vickytest.MockChunks{
		OnSet: func(ctx context.Context, key string, chunk *models.Chunk) error {
			return fmt.Errorf("disk full")
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithDocuments(md),
		vickytest.WithBlobs(mb),
		vickytest.WithChunker(mch),
		vickytest.WithChunks(mc),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
	)

	job := vickytest.NewJob(t)
	_, err := chunkStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "store chunk") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "store chunk")
	}
}
