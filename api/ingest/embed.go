package ingest

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/pipz"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/api/events"
	"github.com/zoobzio/vicky/models"
)

// Stage identity.
var EmbedStageID = pipz.NewIdentity("embed", "Generates vector embeddings for chunks")

// Pool identity.
var EmbedPoolID = pipz.NewIdentity("embed-pool", "Bounded parallel embedding batches")
var embedBatchID = pipz.NewIdentity("embed-batch", "Process single embedding batch")

// Long-lived pool and configuration.
var (
	embedPool      *pipz.WorkerPool[*embedWork]
	embedBatchSize atomic.Int32
)

// Default configuration.
const (
	defaultEmbedWorkers   = 4
	defaultEmbedBatchSize = 128
	defaultEmbedTimeout   = 30 * time.Second
)

func init() {
	embedBatchSize.Store(defaultEmbedBatchSize)
	embedPool = pipz.NewWorkerPool(EmbedPoolID, defaultEmbedWorkers,
		pipz.Apply(embedBatchID, processEmbedBatch),
	).WithTimeout(defaultEmbedTimeout)
}

// SetEmbedConfig updates the embed pool configuration.
// Called by capacitor when config changes.
func SetEmbedConfig(workers int, batchSize int, timeout time.Duration) {
	if workers > 0 {
		embedPool.SetWorkerCount(workers)
	}
	if batchSize > 0 {
		embedBatchSize.Store(int32(batchSize))
	}
	if timeout > 0 {
		embedPool.WithTimeout(timeout)
	}
}

// embedWork carries batch data for parallel embedding.
type embedWork struct {
	// Batch data
	Batch    []*models.Chunk
	BatchIdx int

	// Context for storing results
	JobID int64
}

func (w *embedWork) Clone() *embedWork {
	c := *w
	// Shallow copy of batch slice is fine - we don't modify the slice itself
	return &c
}

// processEmbedBatch handles a single batch of chunks.
func processEmbedBatch(ctx context.Context, w *embedWork) (*embedWork, error) {
	embedder := sum.MustUse[contracts.Embedder](ctx)
	chunks := sum.MustUse[contracts.Chunks](ctx)

	// Extract content strings
	texts := make([]string, len(w.Batch))
	for i, chunk := range w.Batch {
		texts[i] = chunk.Content
	}

	// Generate embeddings
	vectors, err := embedder.Embed(ctx, texts)
	if err != nil {
		return w, fmt.Errorf("embed batch %d: %w", w.BatchIdx, err)
	}

	if len(vectors) != len(w.Batch) {
		return w, fmt.Errorf("embed batch %d: expected %d vectors, got %d", w.BatchIdx, len(w.Batch), len(vectors))
	}

	// Update chunk vectors
	for i, chunk := range w.Batch {
		chunk.Vector = vectors[i]
		if err := chunks.Set(ctx, idToKey(chunk.ID), chunk); err != nil {
			capitan.Error(ctx, events.EmbedChunkErrorSignal,
				events.JobIDKey.Field(w.JobID),
				events.ChunkIDKey.Field(chunk.ID),
				events.ErrorKey.Field(err),
			)
			return w, fmt.Errorf("update chunk %d: %w", chunk.ID, err)
		}
	}

	return w, nil
}

// embedStage generates vector embeddings for chunks.
func embedStage(ctx context.Context, job *models.Job) (*models.Job, error) {
	events.Ingest.Embed.Started.Emit(ctx, events.EmbedStageEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
	})

	// Update job stage
	job.Stage = models.JobStageEmbed

	// Resolve chunks store
	chunks := sum.MustUse[contracts.Chunks](ctx)

	// List all chunks for this version
	allChunks, err := chunks.ListByUserRepoAndTag(ctx, job.UserID, job.Owner, job.RepoName, job.Tag)
	if err != nil {
		return job, fmt.Errorf("list chunks: %w", err)
	}

	job.ItemsTotal = len(allChunks)

	if len(allChunks) == 0 {
		events.Ingest.Embed.Completed.Emit(ctx, events.EmbedStageEvent{
			RepositoryID: job.RepositoryID,
			VersionID:    job.VersionID,
			ChunkCount:   0,
		})
		return job, nil
	}

	// Create batches using current config
	batchSize := int(embedBatchSize.Load())
	batches := batchChunks(allChunks, batchSize)

	// Process batches concurrently via long-lived pool
	var (
		wg          sync.WaitGroup
		errOnce     sync.Once
		firstErr    error
		totalEmbedded atomic.Int64
	)

	for i, batch := range batches {
		wg.Add(1)
		go func(idx int, b []*models.Chunk) {
			defer wg.Done()

			_, err := embedPool.Process(ctx, &embedWork{
				Batch:    b,
				BatchIdx: idx,
				JobID:    job.ID,
			})
			if err != nil {
				errOnce.Do(func() { firstErr = err })
				return
			}

			totalEmbedded.Add(int64(len(b)))
		}(i, batch)
	}

	wg.Wait()

	if firstErr != nil {
		return job, firstErr
	}

	job.ItemsProcessed = int(totalEmbedded.Load())

	events.Ingest.Embed.Completed.Emit(ctx, events.EmbedStageEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
		ChunkCount:   job.ItemsProcessed,
		BatchCount:   len(batches),
	})

	return job, nil
}

// batchChunks splits a chunk slice into batches of the given size.
func batchChunks(chunks []*models.Chunk, size int) [][]*models.Chunk {
	var batches [][]*models.Chunk
	for i := 0; i < len(chunks); i += size {
		end := i + size
		if end > len(chunks) {
			end = len(chunks)
		}
		batches = append(batches, chunks[i:end])
	}
	return batches
}
