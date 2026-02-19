package ingest

import (
	"context"
	"errors"
	"time"

	"github.com/zoobzio/pipz"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/models"
)

// Pipeline identities.
var (
	PipelineID = pipz.NewIdentity("ingest-pipeline", "Main ingestion pipeline sequence")

	// Retry wrapper identities
	FetchRetryID = pipz.NewIdentity("fetch-retry", "Retries fetch stage on transient failures")
	ParseRetryID = pipz.NewIdentity("parse-retry", "Retries parse stage on transient failures")
	ChunkRetryID = pipz.NewIdentity("chunk-retry", "Retries chunk stage on transient failures")
	EmbedRetryID = pipz.NewIdentity("embed-retry", "Retries embed stage on transient failures")
	StoreRetryID = pipz.NewIdentity("store-retry", "Retries store stage on transient failures")

	// Timeout wrapper identities
	FetchTimeoutID = pipz.NewIdentity("fetch-timeout", "Timeout for fetch stage")
	ParseTimeoutID = pipz.NewIdentity("parse-timeout", "Timeout for parse stage")
	ChunkTimeoutID = pipz.NewIdentity("chunk-timeout", "Timeout for chunk stage")
	EmbedTimeoutID = pipz.NewIdentity("embed-timeout", "Timeout for embed stage")
	StoreTimeoutID = pipz.NewIdentity("store-timeout", "Timeout for store stage")

	// Cancellation check identity
	CancelCheckID = pipz.NewIdentity("cancel-check", "Checks if job cancellation was requested")
)

// Pipeline configuration.
const (
	// Retry attempts for each stage
	DefaultRetries = 3

	// Timeouts per stage
	FetchTimeout = 5 * time.Minute
	ParseTimeout = 10 * time.Minute
	ChunkTimeout = 5 * time.Minute
	EmbedTimeout = 30 * time.Minute
	StoreTimeout = 10 * time.Minute
)

// ErrJobCancelled is returned when a job is cancelled by user.
var ErrJobCancelled = errors.New("job cancelled by user")

// checkCancellation returns a pipz chainable that checks if the job should be cancelled.
func checkCancellation() pipz.Chainable[*models.Job] {
	return pipz.Apply(CancelCheckID, func(ctx context.Context, job *models.Job) (*models.Job, error) {
		jobs := sum.MustUse[contracts.Jobs](ctx)

		isCancelling, err := jobs.IsCancelling(ctx, job.ID)
		if err != nil {
			return job, err
		}

		if isCancelling {
			return job, ErrJobCancelled
		}

		return job, nil
	})
}

// NewPipeline creates the ingestion pipeline as a pipz.Sequence.
// Each stage is wrapped with timeout and retry for reliability.
// Cancellation checks are inserted between stages.
func NewPipeline() *pipz.Sequence[*models.Job] {
	// Create base stages
	fetch := pipz.Apply(FetchStageID, fetchStage)
	parse := pipz.Apply(ParseStageID, parseStage)
	chunk := pipz.Apply(ChunkStageID, chunkStage)
	embed := pipz.Apply(EmbedStageID, embedStage)
	store := pipz.Apply(StoreStageID, storeStage)

	// Cancellation check
	cancelCheck := checkCancellation()

	// Wrap with timeout
	fetchWithTimeout := pipz.NewTimeout(FetchTimeoutID, fetch, FetchTimeout)
	parseWithTimeout := pipz.NewTimeout(ParseTimeoutID, parse, ParseTimeout)
	chunkWithTimeout := pipz.NewTimeout(ChunkTimeoutID, chunk, ChunkTimeout)
	embedWithTimeout := pipz.NewTimeout(EmbedTimeoutID, embed, EmbedTimeout)
	storeWithTimeout := pipz.NewTimeout(StoreTimeoutID, store, StoreTimeout)

	// Wrap with retry
	fetchReliable := pipz.NewRetry(FetchRetryID, fetchWithTimeout, DefaultRetries)
	parseReliable := pipz.NewRetry(ParseRetryID, parseWithTimeout, DefaultRetries)
	chunkReliable := pipz.NewRetry(ChunkRetryID, chunkWithTimeout, DefaultRetries)
	embedReliable := pipz.NewRetry(EmbedRetryID, embedWithTimeout, DefaultRetries)
	storeReliable := pipz.NewRetry(StoreRetryID, storeWithTimeout, DefaultRetries)

	// Build sequence with cancellation checks between stages:
	// fetch → [check] → parse → [check] → chunk → [check] → embed → [check] → store
	return pipz.NewSequence(PipelineID,
		fetchReliable,
		cancelCheck,
		parseReliable,
		cancelCheck,
		chunkReliable,
		cancelCheck,
		embedReliable,
		cancelCheck,
		storeReliable,
	)
}
