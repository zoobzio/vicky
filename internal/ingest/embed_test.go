//go:build testing

package ingest

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	vickytest "github.com/zoobzio/vicky/testing"

	"github.com/zoobzio/vicky/models"
)

// --- batchChunks tests (pure function) ---

func TestBatchChunks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		total          int
		batchSize      int
		wantBatches    int
		wantLastBatch  int
	}{
		{"empty", 0, 128, 0, 0},
		{"single item", 1, 128, 1, 1},
		{"exact batch", 128, 128, 1, 128},
		{"one over", 129, 128, 2, 1},
		{"two exact", 256, 128, 2, 128},
		{"remainder", 300, 128, 3, 44},
		{"small batch size", 5, 2, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chunks := make([]*models.Chunk, tt.total)
			for i := range chunks {
				chunks[i] = &models.Chunk{ID: int64(i)}
			}

			batches := batchChunks(chunks, tt.batchSize)

			if len(batches) != tt.wantBatches {
				t.Fatalf("len(batches) = %d, want %d", len(batches), tt.wantBatches)
			}

			if tt.wantBatches > 0 {
				lastBatch := batches[len(batches)-1]
				if len(lastBatch) != tt.wantLastBatch {
					t.Errorf("last batch len = %d, want %d", len(lastBatch), tt.wantLastBatch)
				}
			}

			// Verify all items are present
			var totalItems int
			for _, batch := range batches {
				totalItems += len(batch)
			}
			if totalItems != tt.total {
				t.Errorf("total items across batches = %d, want %d", totalItems, tt.total)
			}
		})
	}
}

// --- embedStage tests ---

func TestEmbedStage_NoChunks(t *testing.T) {
	mc := &vickytest.MockChunks{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
			return nil, nil
		},
	}
	me := &vickytest.MockEmbedder{}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithChunks(mc),
		vickytest.WithEmbedder(me),
	)
	job := vickytest.NewJob(t)

	result, err := embedStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ItemsTotal != 0 {
		t.Errorf("ItemsTotal = %d, want 0", result.ItemsTotal)
	}
	if result.Stage != models.JobStageEmbed {
		t.Errorf("Stage = %q, want %q", result.Stage, models.JobStageEmbed)
	}
}

func TestEmbedStage_SingleBatch(t *testing.T) {
	chunks := vickytest.NewChunks(t, 3)

	var mu sync.Mutex
	var embedCalls int
	var setCalls int
	setChunks := make(map[int64][]float32)

	mc := &vickytest.MockChunks{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
			return chunks, nil
		},
		OnSet: func(ctx context.Context, key string, chunk *models.Chunk) error {
			mu.Lock()
			defer mu.Unlock()
			setCalls++
			setChunks[chunk.ID] = chunk.Vector
			return nil
		},
	}

	me := &vickytest.MockEmbedder{
		OnEmbed: func(ctx context.Context, texts []string) ([][]float32, error) {
			mu.Lock()
			defer mu.Unlock()
			embedCalls++
			vectors := make([][]float32, len(texts))
			for i := range texts {
				vectors[i] = []float32{float32(i) + 0.1, float32(i) + 0.2, float32(i) + 0.3}
			}
			return vectors, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithChunks(mc),
		vickytest.WithEmbedder(me),
	)
	job := vickytest.NewJob(t)

	result, err := embedStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if embedCalls != 1 {
		t.Errorf("Embed called %d times, want 1", embedCalls)
	}
	if setCalls != 3 {
		t.Errorf("Set called %d times, want 3", setCalls)
	}
	if result.ItemsTotal != 3 {
		t.Errorf("ItemsTotal = %d, want 3", result.ItemsTotal)
	}
	if result.ItemsProcessed != 3 {
		t.Errorf("ItemsProcessed = %d, want 3", result.ItemsProcessed)
	}

	// Verify vectors were assigned
	for _, c := range chunks {
		if _, ok := setChunks[c.ID]; !ok {
			t.Errorf("chunk %d was not updated via Set", c.ID)
		}
	}
}

func TestEmbedStage_MultipleBatches(t *testing.T) {
	chunks := vickytest.NewChunks(t, 200)

	var mu sync.Mutex
	var embedCalls int
	var setCalls int

	mc := &vickytest.MockChunks{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
			return chunks, nil
		},
		OnSet: func(ctx context.Context, key string, chunk *models.Chunk) error {
			mu.Lock()
			defer mu.Unlock()
			setCalls++
			return nil
		},
	}

	me := &vickytest.MockEmbedder{
		OnEmbed: func(ctx context.Context, texts []string) ([][]float32, error) {
			mu.Lock()
			defer mu.Unlock()
			embedCalls++
			vectors := make([][]float32, len(texts))
			for i := range vectors {
				vectors[i] = []float32{0.1, 0.2, 0.3}
			}
			return vectors, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithChunks(mc),
		vickytest.WithEmbedder(me),
	)
	job := vickytest.NewJob(t)

	result, err := embedStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 200 chunks / 128 batch size = 2 batches
	if embedCalls != 2 {
		t.Errorf("Embed called %d times, want 2", embedCalls)
	}
	if setCalls != 200 {
		t.Errorf("Set called %d times, want 200", setCalls)
	}
	if result.ItemsTotal != 200 {
		t.Errorf("ItemsTotal = %d, want 200", result.ItemsTotal)
	}
	if result.ItemsProcessed != 200 {
		t.Errorf("ItemsProcessed = %d, want 200", result.ItemsProcessed)
	}
}

func TestEmbedStage_EmbedError(t *testing.T) {
	chunks := vickytest.NewChunks(t, 3)

	mc := &vickytest.MockChunks{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
			return chunks, nil
		},
	}

	me := &vickytest.MockEmbedder{
		OnEmbed: func(ctx context.Context, texts []string) ([][]float32, error) {
			return nil, fmt.Errorf("rate limited")
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithChunks(mc),
		vickytest.WithEmbedder(me),
	)
	job := vickytest.NewJob(t)

	_, err := embedStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "embed batch") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "embed batch")
	}
}

func TestEmbedStage_VectorCountMismatch(t *testing.T) {
	chunks := vickytest.NewChunks(t, 3)

	mc := &vickytest.MockChunks{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
			return chunks, nil
		},
	}

	me := &vickytest.MockEmbedder{
		OnEmbed: func(ctx context.Context, texts []string) ([][]float32, error) {
			// Return wrong number of vectors
			return [][]float32{{0.1, 0.2}}, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithChunks(mc),
		vickytest.WithEmbedder(me),
	)
	job := vickytest.NewJob(t)

	_, err := embedStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "expected 3 vectors, got 1") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "expected 3 vectors, got 1")
	}
}

func TestEmbedStage_ChunkSetError(t *testing.T) {
	chunks := vickytest.NewChunks(t, 3)

	mc := &vickytest.MockChunks{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
			return chunks, nil
		},
		OnSet: func(ctx context.Context, key string, chunk *models.Chunk) error {
			return fmt.Errorf("disk full")
		},
	}

	me := &vickytest.MockEmbedder{
		OnEmbed: func(ctx context.Context, texts []string) ([][]float32, error) {
			vectors := make([][]float32, len(texts))
			for i := range vectors {
				vectors[i] = []float32{0.1, 0.2, 0.3}
			}
			return vectors, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithChunks(mc),
		vickytest.WithEmbedder(me),
	)
	job := vickytest.NewJob(t)

	_, err := embedStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "update chunk") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "update chunk")
	}
}
