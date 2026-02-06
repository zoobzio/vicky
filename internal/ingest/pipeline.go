package ingest

import (
	"time"

	"github.com/zoobzio/pipz"
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

// NewPipeline creates the ingestion pipeline as a pipz.Sequence.
// Each stage is wrapped with timeout and retry for reliability.
func NewPipeline() *pipz.Sequence[*models.Job] {
	// Create base stages
	fetch := pipz.Apply(FetchStageID, fetchStage)
	parse := pipz.Apply(ParseStageID, parseStage)
	chunk := pipz.Apply(ChunkStageID, chunkStage)
	embed := pipz.Apply(EmbedStageID, embedStage)
	store := pipz.Apply(StoreStageID, storeStage)

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

	// Build sequence: fetch → parse → chunk → embed → store
	return pipz.NewSequence(PipelineID,
		fetchReliable,
		parseReliable,
		chunkReliable,
		embedReliable,
		storeReliable,
	)
}
