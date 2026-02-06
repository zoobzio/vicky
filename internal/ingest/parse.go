package ingest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/sourcegraph/scip/bindings/go/scip"
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/pipz"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/contracts"
	"github.com/zoobzio/vicky/events"
	"github.com/zoobzio/vicky/external/indexer"
	vickyScip "github.com/zoobzio/vicky/internal/scip"
	"github.com/zoobzio/vicky/models"
)

// Stage identity.
var ParseStageID = pipz.NewIdentity("parse", "Parses source files and extracts SCIP data")

// Pool identity.
var ParsePoolID = pipz.NewIdentity("parse-pool", "Bounded parallel SCIP document processing")
var parseDocID = pipz.NewIdentity("parse-doc", "Process single SCIP document")

// Long-lived pool.
var parsePool *pipz.WorkerPool[*parseWork]

// Default configuration.
const (
	defaultParseWorkers = 8
	defaultParseTimeout = 60 * time.Second
)

func init() {
	parsePool = pipz.NewWorkerPool(ParsePoolID, defaultParseWorkers,
		pipz.Apply(parseDocID, processParseDoc),
	).WithTimeout(defaultParseTimeout)
}

// SetParseConfig updates the parse pool configuration.
// Called by capacitor when config changes.
func SetParseConfig(workers int, timeout time.Duration) {
	if workers > 0 {
		parsePool.SetWorkerCount(workers)
	}
	if timeout > 0 {
		parsePool.WithTimeout(timeout)
	}
}

// parseWork carries SCIP document data for parallel processing.
type parseWork struct {
	// Context
	UserID   int64
	Owner    string
	RepoName string
	Tag      string
	JobID    int64

	// Document data
	DocumentID int64
	Document   *scip.Document
}

func (w *parseWork) Clone() *parseWork {
	c := *w
	return &c
}

// processParseDoc processes a single SCIP document's symbols and occurrences.
func processParseDoc(ctx context.Context, w *parseWork) (*parseWork, error) {
	scipSymbols := sum.MustUse[contracts.SCIPSymbols](ctx)
	scipOccurrences := sum.MustUse[contracts.SCIPOccurrences](ctx)
	scipRelationships := sum.MustUse[contracts.SCIPRelationships](ctx)

	meta := vickyScip.FileMeta{
		DocumentMeta: vickyScip.DocumentMeta{
			UserID:   w.UserID,
			Owner:    w.Owner,
			RepoName: w.RepoName,
			Tag:      w.Tag,
		},
		DocumentID: w.DocumentID,
		Path:       w.Document.RelativePath,
	}

	// Process symbols and their relationships
	for _, sym := range w.Document.Symbols {
		symbol := vickyScip.ConvertSymbol(sym, meta)
		if err := scipSymbols.Set(ctx, "", &symbol); err != nil {
			capitan.Error(ctx, events.ParseSymbolErrorSignal,
				events.JobIDKey.Field(w.JobID),
				events.SymbolKey.Field(sym.Symbol),
				events.ErrorKey.Field(err),
			)
			continue
		}

		for _, rel := range sym.Relationships {
			relationship := vickyScip.ConvertRelationship(rel)
			relationship.SCIPSymbolID = symbol.ID
			if err := scipRelationships.Set(ctx, "", &relationship); err != nil {
				capitan.Error(ctx, events.ParseRelationshipErrorSignal,
					events.JobIDKey.Field(w.JobID),
					events.ErrorKey.Field(err),
				)
				continue
			}
		}
	}

	// Process occurrences
	for _, occ := range w.Document.Occurrences {
		occurrence := vickyScip.ConvertOccurrence(occ, meta)
		if err := scipOccurrences.Set(ctx, "", &occurrence); err != nil {
			capitan.Error(ctx, events.ParseOccurrenceErrorSignal,
				events.JobIDKey.Field(w.JobID),
				events.ErrorKey.Field(err),
			)
			continue
		}
	}

	return w, nil
}

// parseStage parses source files and extracts SCIP data.
// Delegates to the Indexer contract which may be local or remote.
func parseStage(ctx context.Context, job *models.Job) (*models.Job, error) {
	events.Ingest.Parse.Started.Emit(ctx, events.ParseEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
	})

	// Update job stage
	job.Stage = models.JobStageParse

	// Resolve dependencies
	idx := sum.MustUse[contracts.Indexer](ctx)
	configs := sum.MustUse[contracts.IngestionConfigs](ctx)
	versions := sum.MustUse[contracts.Versions](ctx)
	documents := sum.MustUse[contracts.Documents](ctx)

	// Get config for language
	config, err := configs.GetByRepositoryID(ctx, job.RepositoryID)
	if err != nil {
		return job, err
	}

	// Check if indexer supports this language
	if !idx.Supports(config.Language) {
		events.Ingest.Parse.FileSkipped.Emit(ctx, events.ParseFileEvent{
			RepositoryID: job.RepositoryID,
			VersionID:    job.VersionID,
			Language:     string(config.Language),
			Reason:       "no indexer configured",
		})
		return job, nil
	}

	// Get version for commit SHA
	version, err := versions.Get(ctx, idToKey(job.VersionID))
	if err != nil {
		return job, err
	}

	// Build index request
	req := indexer.Request{
		JobID:        job.ID,
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
		UserID:       job.UserID,
		Owner:        job.Owner,
		RepoName:     job.RepoName,
		Tag:          job.Tag,
		CommitSHA:    version.CommitSHA,
		Language:     config.Language,
	}

	// Call indexer (may be local CLI or remote service)
	result, err := idx.Index(ctx, req)
	if err != nil {
		return job, err
	}

	if result.Error != "" {
		return job, fmt.Errorf("indexer error: %s", result.Error)
	}

	// Parse the SCIP index data
	parser := vickyScip.New()
	index, err := parser.Parse(ctx, result.IndexData)
	if err != nil {
		return job, fmt.Errorf("parse scip index: %w", err)
	}

	// Create document rows sequentially — we need the docIDs map before concurrent processing
	docIDs := make(map[string]int64, len(index.Documents))
	for _, doc := range index.Documents {
		contentType := contentTypeForPath(doc.RelativePath)
		h := sha256.Sum256([]byte(doc.RelativePath + ":" + job.Tag))

		d := &models.Document{
			VersionID:   job.VersionID,
			UserID:      job.UserID,
			Owner:       job.Owner,
			RepoName:    job.RepoName,
			Tag:         job.Tag,
			Path:        doc.RelativePath,
			ContentType: contentType,
			ContentHash: hex.EncodeToString(h[:]),
		}

		if err := documents.Set(ctx, "", d); err != nil {
			return job, fmt.Errorf("create document %s: %w", doc.RelativePath, err)
		}

		docIDs[doc.RelativePath] = d.ID
	}

	if len(docIDs) == 0 {
		return job, nil
	}

	// Process documents concurrently via long-lived pool
	var (
		wg       sync.WaitGroup
		errOnce  sync.Once
		firstErr error
	)

	for _, doc := range index.Documents {
		docID, ok := docIDs[doc.RelativePath]
		if !ok {
			continue
		}

		wg.Add(1)
		go func(d *scip.Document, id int64) {
			defer wg.Done()

			_, err := parsePool.Process(ctx, &parseWork{
				UserID:     job.UserID,
				Owner:      job.Owner,
				RepoName:   job.RepoName,
				Tag:        job.Tag,
				JobID:      job.ID,
				DocumentID: id,
				Document:   d,
			})
			if err != nil {
				errOnce.Do(func() { firstErr = err })
			}
		}(doc, docID)
	}

	wg.Wait()

	if firstErr != nil {
		return job, firstErr
	}

	events.Ingest.Parse.Completed.Emit(ctx, events.ParseEvent{
		RepositoryID:   job.RepositoryID,
		VersionID:      job.VersionID,
		FileCount:      len(docIDs),
		ProcessedFiles: len(docIDs),
	})

	return job, nil
}
