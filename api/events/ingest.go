package events

import (
	"time"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
)

// IngestStartedEvent is emitted when ingestion begins for a version.
type IngestStartedEvent struct {
	RepositoryID int64  `json:"repository_id"`
	VersionID    int64  `json:"version_id"`
	Owner        string `json:"owner"`
	RepoName     string `json:"repo_name"`
	VersionTag   string `json:"version_tag"`
	CommitSHA    string `json:"commit_sha"`
}

// IngestCompletedEvent is emitted when ingestion finishes successfully.
type IngestCompletedEvent struct {
	RepositoryID int64         `json:"repository_id"`
	VersionID    int64         `json:"version_id"`
	VersionTag   string        `json:"version_tag"`
	FileCount    int           `json:"file_count"`
	ChunkCount   int           `json:"chunk_count"`
	SymbolCount  int           `json:"symbol_count"`
	Duration     time.Duration `json:"duration"`
}

// IngestFailedEvent is emitted when ingestion fails at any stage.
type IngestFailedEvent struct {
	RepositoryID int64  `json:"repository_id"`
	VersionID    int64  `json:"version_id"`
	VersionTag   string `json:"version_tag"`
	Stage        string `json:"stage"`
	Error        string `json:"error"`
}

// IngestCancelledEvent is emitted when ingestion is cancelled.
type IngestCancelledEvent struct {
	RepositoryID int64  `json:"repository_id"`
	VersionID    int64  `json:"version_id"`
	VersionTag   string `json:"version_tag"`
	Stage        string `json:"stage"`
	Reason       string `json:"reason"`
}

// FetchEvent is emitted during the fetch stage.
type FetchEvent struct {
	RepositoryID int64  `json:"repository_id"`
	VersionID    int64  `json:"version_id"`
	CommitSHA    string `json:"commit_sha"`
	ByteCount    int64  `json:"byte_count,omitempty"`
	Error        string `json:"error,omitempty"`
}

// ParseEvent is emitted during the parse stage.
type ParseEvent struct {
	RepositoryID   int64  `json:"repository_id"`
	VersionID      int64  `json:"version_id"`
	FileCount      int    `json:"file_count,omitempty"`
	SymbolCount    int    `json:"symbol_count,omitempty"`
	ProcessedFiles int    `json:"processed_files,omitempty"`
	SkippedFiles   int    `json:"skipped_files,omitempty"`
	Error          string `json:"error,omitempty"`
}

// ParseFileEvent is emitted per-file during parsing.
type ParseFileEvent struct {
	RepositoryID int64  `json:"repository_id"`
	VersionID    int64  `json:"version_id"`
	FilePath     string `json:"file_path"`
	Language     string `json:"language,omitempty"`
	SymbolCount  int    `json:"symbol_count,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

// ChunkEvent is emitted during the chunk stage.
type ChunkEvent struct {
	RepositoryID int64  `json:"repository_id"`
	VersionID    int64  `json:"version_id"`
	ChunkCount   int    `json:"chunk_count,omitempty"`
	Error        string `json:"error,omitempty"`
}

// EmbedStageEvent is emitted during the embed stage.
type EmbedStageEvent struct {
	RepositoryID int64         `json:"repository_id"`
	VersionID    int64         `json:"version_id"`
	ChunkCount   int           `json:"chunk_count,omitempty"`
	BatchCount   int           `json:"batch_count,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// EmbedBatchEvent is emitted per batch during embedding.
type EmbedBatchEvent struct {
	RepositoryID int64         `json:"repository_id"`
	VersionID    int64         `json:"version_id"`
	BatchNumber  int           `json:"batch_number"`
	BatchSize    int           `json:"batch_size"`
	TotalBatches int           `json:"total_batches"`
	Duration     time.Duration `json:"duration,omitempty"`
}

// StoreEvent is emitted during the store stage.
type StoreEvent struct {
	RepositoryID int64         `json:"repository_id"`
	VersionID    int64         `json:"version_id"`
	ChunkCount   int           `json:"chunk_count,omitempty"`
	SymbolCount  int           `json:"symbol_count,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// Ingest lifecycle signals.
var (
	IngestStartedSignal   = capitan.NewSignal("vicky.ingest.started", "Version ingestion initiated")
	IngestCompletedSignal = capitan.NewSignal("vicky.ingest.completed", "Version ingestion completed")
	IngestFailedSignal    = capitan.NewSignal("vicky.ingest.failed", "Version ingestion failed")
	IngestCancelledSignal = capitan.NewSignal("vicky.ingest.cancelled", "Version ingestion cancelled")
)

// Fetch stage signals.
var (
	FetchStartedSignal   = capitan.NewSignal("vicky.ingest.fetch.started", "Repository fetch initiated")
	FetchCompletedSignal = capitan.NewSignal("vicky.ingest.fetch.completed", "Repository fetch completed")
	FetchFailedSignal    = capitan.NewSignal("vicky.ingest.fetch.failed", "Repository fetch failed")
)

// Parse stage signals.
var (
	ParseStartedSignal       = capitan.NewSignal("vicky.ingest.parse.started", "AST parsing initiated")
	ParseCompletedSignal     = capitan.NewSignal("vicky.ingest.parse.completed", "AST parsing completed")
	ParseFailedSignal        = capitan.NewSignal("vicky.ingest.parse.failed", "AST parsing failed")
	ParseFileCompletedSignal = capitan.NewSignal("vicky.ingest.parse.file.completed", "File parsing completed")
	ParseFileSkippedSignal   = capitan.NewSignal("vicky.ingest.parse.file.skipped", "File skipped during parsing")
)

// Chunk stage signals.
var (
	ChunkStartedSignal   = capitan.NewSignal("vicky.ingest.chunk.started", "Content chunking initiated")
	ChunkCompletedSignal = capitan.NewSignal("vicky.ingest.chunk.completed", "Content chunking completed")
	ChunkFailedSignal    = capitan.NewSignal("vicky.ingest.chunk.failed", "Content chunking failed")
)

// Embed stage signals.
var (
	EmbedStartedSignal        = capitan.NewSignal("vicky.ingest.embed.started", "Embedding generation initiated")
	EmbedCompletedSignal      = capitan.NewSignal("vicky.ingest.embed.completed", "Embedding generation completed")
	EmbedFailedSignal         = capitan.NewSignal("vicky.ingest.embed.failed", "Embedding generation failed")
	EmbedBatchCompletedSignal = capitan.NewSignal("vicky.ingest.embed.batch.completed", "Embedding batch completed")
)

// Store stage signals.
var (
	StoreStartedSignal   = capitan.NewSignal("vicky.ingest.store.started", "Database storage initiated")
	StoreCompletedSignal = capitan.NewSignal("vicky.ingest.store.completed", "Database storage completed")
	StoreFailedSignal    = capitan.NewSignal("vicky.ingest.store.failed", "Database storage failed")
)

// fetchEvents provides access to fetch stage events.
var fetchEvents = struct {
	Started   sum.Event[FetchEvent]
	Completed sum.Event[FetchEvent]
	Failed    sum.Event[FetchEvent]
}{
	Started:   sum.NewDebugEvent[FetchEvent](FetchStartedSignal),
	Completed: sum.NewDebugEvent[FetchEvent](FetchCompletedSignal),
	Failed:    sum.NewErrorEvent[FetchEvent](FetchFailedSignal),
}

// parseEvents provides access to parse stage events.
var parseEvents = struct {
	Started       sum.Event[ParseEvent]
	Completed     sum.Event[ParseEvent]
	Failed        sum.Event[ParseEvent]
	FileCompleted sum.Event[ParseFileEvent]
	FileSkipped   sum.Event[ParseFileEvent]
}{
	Started:       sum.NewDebugEvent[ParseEvent](ParseStartedSignal),
	Completed:     sum.NewDebugEvent[ParseEvent](ParseCompletedSignal),
	Failed:        sum.NewErrorEvent[ParseEvent](ParseFailedSignal),
	FileCompleted: sum.NewDebugEvent[ParseFileEvent](ParseFileCompletedSignal),
	FileSkipped:   sum.NewDebugEvent[ParseFileEvent](ParseFileSkippedSignal),
}

// chunkEvents provides access to chunk stage events.
var chunkEvents = struct {
	Started   sum.Event[ChunkEvent]
	Completed sum.Event[ChunkEvent]
	Failed    sum.Event[ChunkEvent]
}{
	Started:   sum.NewDebugEvent[ChunkEvent](ChunkStartedSignal),
	Completed: sum.NewDebugEvent[ChunkEvent](ChunkCompletedSignal),
	Failed:    sum.NewErrorEvent[ChunkEvent](ChunkFailedSignal),
}

// embedEvents provides access to embed stage events.
var embedEvents = struct {
	Started        sum.Event[EmbedStageEvent]
	Completed      sum.Event[EmbedStageEvent]
	Failed         sum.Event[EmbedStageEvent]
	BatchCompleted sum.Event[EmbedBatchEvent]
}{
	Started:        sum.NewDebugEvent[EmbedStageEvent](EmbedStartedSignal),
	Completed:      sum.NewDebugEvent[EmbedStageEvent](EmbedCompletedSignal),
	Failed:         sum.NewErrorEvent[EmbedStageEvent](EmbedFailedSignal),
	BatchCompleted: sum.NewDebugEvent[EmbedBatchEvent](EmbedBatchCompletedSignal),
}

// storeEvents provides access to store stage events.
var storeEvents = struct {
	Started   sum.Event[StoreEvent]
	Completed sum.Event[StoreEvent]
	Failed    sum.Event[StoreEvent]
}{
	Started:   sum.NewDebugEvent[StoreEvent](StoreStartedSignal),
	Completed: sum.NewDebugEvent[StoreEvent](StoreCompletedSignal),
	Failed:    sum.NewErrorEvent[StoreEvent](StoreFailedSignal),
}

// Ingest provides access to ingestion pipeline events.
var Ingest = struct {
	Started   sum.Event[IngestStartedEvent]
	Completed sum.Event[IngestCompletedEvent]
	Failed    sum.Event[IngestFailedEvent]
	Cancelled sum.Event[IngestCancelledEvent]

	Fetch struct {
		Started   sum.Event[FetchEvent]
		Completed sum.Event[FetchEvent]
		Failed    sum.Event[FetchEvent]
	}
	Parse struct {
		Started       sum.Event[ParseEvent]
		Completed     sum.Event[ParseEvent]
		Failed        sum.Event[ParseEvent]
		FileCompleted sum.Event[ParseFileEvent]
		FileSkipped   sum.Event[ParseFileEvent]
	}
	Chunk struct {
		Started   sum.Event[ChunkEvent]
		Completed sum.Event[ChunkEvent]
		Failed    sum.Event[ChunkEvent]
	}
	Embed struct {
		Started        sum.Event[EmbedStageEvent]
		Completed      sum.Event[EmbedStageEvent]
		Failed         sum.Event[EmbedStageEvent]
		BatchCompleted sum.Event[EmbedBatchEvent]
	}
	Store struct {
		Started   sum.Event[StoreEvent]
		Completed sum.Event[StoreEvent]
		Failed    sum.Event[StoreEvent]
	}
}{
	Started:   sum.NewInfoEvent[IngestStartedEvent](IngestStartedSignal),
	Completed: sum.NewInfoEvent[IngestCompletedEvent](IngestCompletedSignal),
	Failed:    sum.NewErrorEvent[IngestFailedEvent](IngestFailedSignal),
	Cancelled: sum.NewWarnEvent[IngestCancelledEvent](IngestCancelledSignal),

	Fetch: fetchEvents,
	Parse: parseEvents,
	Chunk: chunkEvents,
	Embed: embedEvents,
	Store: storeEvents,
}

