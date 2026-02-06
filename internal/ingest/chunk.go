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
	"github.com/zoobzio/vicky/contracts"
	"github.com/zoobzio/vicky/events"
	"github.com/zoobzio/vicky/models"
)

// Stage identity.
var ChunkStageID = pipz.NewIdentity("chunk", "Chunks content into embeddable segments")

// Pool identity.
var ChunkPoolID = pipz.NewIdentity("chunk-pool", "Bounded parallel document chunking")
var chunkDocID = pipz.NewIdentity("chunk-doc", "Process single document chunking")

// Long-lived pool.
var chunkPool *pipz.WorkerPool[*chunkWork]

// Default configuration.
const (
	defaultChunkWorkers = 8
	defaultChunkTimeout = 30 * time.Second
)

func init() {
	chunkPool = pipz.NewWorkerPool(ChunkPoolID, defaultChunkWorkers,
		pipz.Apply(chunkDocID, processChunkDoc),
	).WithTimeout(defaultChunkTimeout)
}

// SetChunkConfig updates the chunk pool configuration.
// Called by capacitor when config changes.
func SetChunkConfig(workers int, timeout time.Duration) {
	if workers > 0 {
		chunkPool.SetWorkerCount(workers)
	}
	if timeout > 0 {
		chunkPool.WithTimeout(timeout)
	}
}

// chunkWork carries document data for parallel chunking.
type chunkWork struct {
	// Context
	UserID   int64
	Owner    string
	RepoName string
	Tag      string
	Language string
	JobID    int64

	// Document data
	DocumentID  int64
	Path        string
	ContentType models.ContentType

	// Result counter (shared across goroutines)
	ChunkCount *atomic.Int64
}

func (w *chunkWork) Clone() *chunkWork {
	c := *w
	return &c
}

// processChunkDoc chunks a single document.
func processChunkDoc(ctx context.Context, w *chunkWork) (*chunkWork, error) {
	chunker := sum.MustUse[contracts.Chunker](ctx)
	blobs := sum.MustUse[contracts.Blobs](ctx)
	chunks := sum.MustUse[contracts.Chunks](ctx)

	// Fetch blob content from minio
	obj, err := blobs.GetByPath(ctx, w.UserID, w.Owner, w.RepoName, w.Tag, w.Path)
	if err != nil {
		capitan.Error(ctx, events.ChunkBlobErrorSignal,
			events.JobIDKey.Field(w.JobID),
			events.PathKey.Field(w.Path),
			events.ErrorKey.Field(err),
		)
		return w, fmt.Errorf("blob %s: %w", w.Path, err)
	}

	// Determine chisel language
	lang := w.Language
	if w.ContentType == models.ContentTypeDocs {
		lang = "markdown"
	}

	if !chunker.Supports(lang) {
		capitan.Debug(ctx, events.ChunkSkippedSignal,
			events.JobIDKey.Field(w.JobID),
			events.PathKey.Field(w.Path),
			events.LanguageKey.Field(lang),
		)
		return w, nil
	}

	// Chunk the content
	results, err := chunker.Chunk(ctx, lang, w.Path, []byte(obj.Data.Content))
	if err != nil {
		capitan.Error(ctx, events.ChunkProcessErrorSignal,
			events.JobIDKey.Field(w.JobID),
			events.PathKey.Field(w.Path),
			events.ErrorKey.Field(err),
		)
		return w, fmt.Errorf("chunk %s: %w", w.Path, err)
	}

	// Persist each chunk
	for _, r := range results {
		var symbol *string
		if r.Symbol != "" {
			symbol = &r.Symbol
		}

		chunk := &models.Chunk{
			DocumentID: w.DocumentID,
			UserID:     w.UserID,
			Owner:      w.Owner,
			RepoName:   w.RepoName,
			Tag:        w.Tag,
			Path:       w.Path,
			Kind:       r.Kind,
			StartLine:  r.StartLine,
			EndLine:    r.EndLine,
			Symbol:     symbol,
			Context:    r.Context,
			Content:    r.Content,
		}

		if err := chunks.Set(ctx, "", chunk); err != nil {
			capitan.Error(ctx, events.ChunkStoreErrorSignal,
				events.JobIDKey.Field(w.JobID),
				events.PathKey.Field(w.Path),
				events.ErrorKey.Field(err),
			)
			return w, fmt.Errorf("store chunk in %s: %w", w.Path, err)
		}

		if w.ChunkCount != nil {
			w.ChunkCount.Add(1)
		}
	}

	return w, nil
}

// chunkStage chunks content into embeddable segments.
func chunkStage(ctx context.Context, job *models.Job) (*models.Job, error) {
	events.Ingest.Chunk.Started.Emit(ctx, events.ChunkEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
	})

	// Update job stage
	job.Stage = models.JobStageChunk

	// Resolve dependencies
	documents := sum.MustUse[contracts.Documents](ctx)
	configs := sum.MustUse[contracts.IngestionConfigs](ctx)

	// Get config for language
	config, err := configs.GetByRepositoryID(ctx, job.RepositoryID)
	if err != nil {
		return job, err
	}

	// List documents created in parse stage
	docs, err := documents.ListByUserRepoAndTag(ctx, job.UserID, job.Owner, job.RepoName, job.Tag)
	if err != nil {
		return job, fmt.Errorf("list documents: %w", err)
	}

	job.ItemsTotal = len(docs)

	if len(docs) == 0 {
		events.Ingest.Chunk.Completed.Emit(ctx, events.ChunkEvent{
			RepositoryID: job.RepositoryID,
			VersionID:    job.VersionID,
			ChunkCount:   0,
		})
		return job, nil
	}

	// Process documents concurrently via long-lived pool
	var (
		wg          sync.WaitGroup
		errOnce     sync.Once
		firstErr    error
		totalChunks atomic.Int64
	)

	language := string(config.Language)

	for _, doc := range docs {
		wg.Add(1)
		go func(d *models.Document) {
			defer wg.Done()

			_, err := chunkPool.Process(ctx, &chunkWork{
				UserID:      job.UserID,
				Owner:       job.Owner,
				RepoName:    job.RepoName,
				Tag:         job.Tag,
				Language:    language,
				JobID:       job.ID,
				DocumentID:  d.ID,
				Path:        d.Path,
				ContentType: d.ContentType,
				ChunkCount:  &totalChunks,
			})
			if err != nil {
				errOnce.Do(func() { firstErr = err })
			}
		}(doc)
	}

	wg.Wait()

	if firstErr != nil {
		return job, firstErr
	}

	job.ItemsProcessed = len(docs)

	events.Ingest.Chunk.Completed.Emit(ctx, events.ChunkEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
		ChunkCount:   int(totalChunks.Load()),
	})

	return job, nil
}
