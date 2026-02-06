package events

import "github.com/zoobzio/capitan"

// Common field keys for direct capitan emissions.
// Use these with capitan.Emit, capitan.Debug, capitan.Error, etc.
var (
	// Job context
	JobIDKey        = capitan.NewInt64Key("job_id")
	RepositoryIDKey = capitan.NewInt64Key("repository_id")
	VersionIDKey    = capitan.NewInt64Key("version_id")
	UserIDKey       = capitan.NewInt64Key("user_id")

	// Repository context
	OwnerKey    = capitan.NewStringKey("owner")
	RepoNameKey = capitan.NewStringKey("repo_name")
	TagKey      = capitan.NewStringKey("tag")

	// File/path context
	PathKey     = capitan.NewStringKey("path")
	LanguageKey = capitan.NewStringKey("language")

	// Counts
	CountKey       = capitan.NewIntKey("count")
	TotalKey       = capitan.NewIntKey("total")
	ProcessedKey   = capitan.NewIntKey("processed")
	BatchKey       = capitan.NewIntKey("batch")
	BatchSizeKey   = capitan.NewIntKey("batch_size")
	SymbolCountKey = capitan.NewIntKey("symbol_count")
	ChunkCountKey  = capitan.NewIntKey("chunk_count")
	FileCountKey   = capitan.NewIntKey("file_count")

	// Errors
	ErrorKey  = capitan.NewErrorKey("error")
	ReasonKey = capitan.NewStringKey("reason")

	// Identifiers
	SymbolKey  = capitan.NewStringKey("symbol")
	ChunkIDKey = capitan.NewInt64Key("chunk_id")
)

// Operational signals for debug/error logging within pipeline stages.
// These are for operational visibility, not lifecycle events.
var (
	// Fetch stage operations
	FetchBlobStoredSignal = capitan.NewSignal("vicky.ingest.fetch.blob.stored", "Blob stored successfully")
	FetchBlobErrorSignal  = capitan.NewSignal("vicky.ingest.fetch.blob.error", "Failed to store blob")

	// Parse stage operations
	ParseDocumentCreatedSignal = capitan.NewSignal("vicky.ingest.parse.document.created", "Document row created")
	ParseSymbolStoredSignal    = capitan.NewSignal("vicky.ingest.parse.symbol.stored", "Symbol stored")
	ParseSymbolErrorSignal     = capitan.NewSignal("vicky.ingest.parse.symbol.error", "Failed to store symbol")
	ParseOccurrenceErrorSignal = capitan.NewSignal("vicky.ingest.parse.occurrence.error", "Failed to store occurrence")
	ParseRelationshipErrorSignal = capitan.NewSignal("vicky.ingest.parse.relationship.error", "Failed to store relationship")

	// Chunk stage operations
	ChunkBlobErrorSignal    = capitan.NewSignal("vicky.ingest.chunk.blob.error", "Failed to fetch blob for chunking")
	ChunkSkippedSignal      = capitan.NewSignal("vicky.ingest.chunk.skipped", "Document skipped - no chunker available")
	ChunkProcessErrorSignal = capitan.NewSignal("vicky.ingest.chunk.process.error", "Failed to chunk document")
	ChunkStoreErrorSignal   = capitan.NewSignal("vicky.ingest.chunk.store.error", "Failed to store chunk")

	// Embed stage operations
	EmbedChunkErrorSignal = capitan.NewSignal("vicky.ingest.embed.chunk.error", "Failed to update chunk with embedding")
)
