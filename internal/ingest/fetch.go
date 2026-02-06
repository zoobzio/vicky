package ingest

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
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
var FetchStageID = pipz.NewIdentity("fetch", "Fetches repository content from GitHub")

// Pool identity.
var FetchPoolID = pipz.NewIdentity("fetch-pool", "Bounded parallel blob storage")
var fetchFileID = pipz.NewIdentity("fetch-file", "Store single file blob")

// Long-lived pool.
var fetchPool *pipz.WorkerPool[*fetchWork]

// Default configuration.
const (
	defaultFetchWorkers = 8
	defaultFetchTimeout = 30 * time.Second
)

func init() {
	fetchPool = pipz.NewWorkerPool(FetchPoolID, defaultFetchWorkers,
		pipz.Apply(fetchFileID, processFetchFile),
	).WithTimeout(defaultFetchTimeout)
}

// SetFetchConfig updates the fetch pool configuration.
// Called by capacitor when config changes.
func SetFetchConfig(workers int, timeout time.Duration) {
	if workers > 0 {
		fetchPool.SetWorkerCount(workers)
	}
	if timeout > 0 {
		fetchPool.WithTimeout(timeout)
	}
}

// Language file extensions.
var languageExtensions = map[models.Language][]string{
	models.LanguageGo:         {".go"},
	models.LanguageTypeScript: {".ts", ".tsx", ".js", ".jsx"},
}

// fetchWork carries file data for parallel blob storage.
type fetchWork struct {
	// Context
	UserID   int64
	Owner    string
	Repo     string
	Tag      string
	Language string
	JobID    int64

	// File data
	Path    string
	Content []byte
}

func (w *fetchWork) Clone() *fetchWork {
	c := *w
	return &c
}

// processFetchFile stores a single file blob.
func processFetchFile(ctx context.Context, w *fetchWork) (*fetchWork, error) {
	blobs := sum.MustUse[contracts.Blobs](ctx)

	blob := &models.Blob{
		Path:     w.Path,
		Content:  string(w.Content),
		Language: w.Language,
		Owner:    w.Owner,
		Repo:     w.Repo,
		Tag:      w.Tag,
	}

	if err := blobs.PutBlob(ctx, w.UserID, blob); err != nil {
		capitan.Error(ctx, events.FetchBlobErrorSignal,
			events.JobIDKey.Field(w.JobID),
			events.PathKey.Field(w.Path),
			events.ErrorKey.Field(err),
		)
		return w, fmt.Errorf("store blob %s: %w", w.Path, err)
	}

	return w, nil
}

// fetchStage fetches repository content from GitHub.
func fetchStage(ctx context.Context, job *models.Job) (*models.Job, error) {
	events.Ingest.Fetch.Started.Emit(ctx, events.FetchEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
	})

	// Resolve dependencies from registry
	users := sum.MustUse[contracts.Users](ctx)
	github := sum.MustUse[contracts.GitHub](ctx)
	configs := sum.MustUse[contracts.IngestionConfigs](ctx)
	versions := sum.MustUse[contracts.Versions](ctx)

	// Update job stage
	job.Stage = models.JobStageFetch

	// Get user for access token
	user, err := users.Get(ctx, idToKey(job.UserID))
	if err != nil {
		return job, err
	}

	// Get ingestion config
	config, err := configs.GetByRepositoryID(ctx, job.RepositoryID)
	if err != nil {
		return job, err
	}

	// Get version for commit SHA
	version, err := versions.Get(ctx, idToKey(job.VersionID))
	if err != nil {
		return job, err
	}

	// Fetch repository tree
	tree, err := github.GetTree(ctx, user.AccessToken, job.Owner, job.RepoName, version.CommitSHA)
	if err != nil {
		return job, err
	}

	// Filter files based on config
	var filePaths []string
	excludePatterns := config.AllExcludePatterns()
	allowedExts := languageExtensions[config.Language]

	for _, entry := range tree {
		if entry.Type != "blob" {
			continue
		}
		if entry.Size > config.MaxFileSize {
			continue
		}
		if matchesAnyPattern(entry.Path, excludePatterns) {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Path))
		isCode := containsExt(allowedExts, ext)
		isDocs := config.IncludeDocs && containsExt(docExtensions, ext)

		if isCode || isDocs {
			filePaths = append(filePaths, entry.Path)
		}
	}

	job.ItemsTotal = len(filePaths)

	if len(filePaths) == 0 {
		events.Ingest.Fetch.Completed.Emit(ctx, events.FetchEvent{
			RepositoryID: job.RepositoryID,
			VersionID:    job.VersionID,
			ByteCount:    0,
		})
		return job, nil
	}

	// Fetch file contents in batch
	contents, err := github.GetFileContentBatch(ctx, user.AccessToken, job.Owner, job.RepoName, version.CommitSHA, filePaths)
	if err != nil {
		return job, err
	}

	// Process files concurrently via long-lived pool
	var (
		wg       sync.WaitGroup
		errOnce  sync.Once
		firstErr error
		stored   atomic.Int64
	)

	language := string(config.Language)

	for _, content := range contents {
		wg.Add(1)
		go func(path string, data []byte) {
			defer wg.Done()

			_, err := fetchPool.Process(ctx, &fetchWork{
				UserID:   job.UserID,
				Owner:    job.Owner,
				Repo:     job.RepoName,
				Tag:      job.Tag,
				Language: language,
				JobID:    job.ID,
				Path:     path,
				Content:  data,
			})
			if err != nil {
				errOnce.Do(func() { firstErr = err })
				return
			}

			stored.Add(1)
		}(content.Path, content.Content)
	}

	wg.Wait()

	if firstErr != nil {
		return job, firstErr
	}

	job.ItemsProcessed = int(stored.Load())

	events.Ingest.Fetch.Completed.Emit(ctx, events.FetchEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
		ByteCount:    int64(job.ItemsProcessed),
	})

	return job, nil
}
