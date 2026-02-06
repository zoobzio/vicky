package events

import "github.com/zoobzio/capitan"

// Indexer sidecar signals.
var (
	IndexerStartedSignal     = capitan.NewSignal("vicky.indexer.sidecar.started", "Indexer sidecar started")
	IndexerListeningSignal   = capitan.NewSignal("vicky.indexer.sidecar.listening", "Indexer sidecar listening")
	IndexerRequestSignal     = capitan.NewSignal("vicky.indexer.request", "Indexer received request")
	IndexerFetchedSignal     = capitan.NewSignal("vicky.indexer.fetched", "Indexer fetched blobs")
	IndexerCompletedSignal   = capitan.NewSignal("vicky.indexer.completed", "Indexer completed successfully")
	IndexerErrorSignal       = capitan.NewSignal("vicky.indexer.error", "Indexer encountered error")
)

// Chunker sidecar signals.
var (
	ChunkerStartedSignal   = capitan.NewSignal("vicky.chunker.sidecar.started", "Chunker sidecar started")
	ChunkerListeningSignal = capitan.NewSignal("vicky.chunker.sidecar.listening", "Chunker sidecar listening")
	ChunkerRequestSignal   = capitan.NewSignal("vicky.chunker.request", "Chunker received request")
	ChunkerCompletedSignal = capitan.NewSignal("vicky.chunker.completed", "Chunker completed successfully")
	ChunkerErrorSignal     = capitan.NewSignal("vicky.chunker.error", "Chunker encountered error")
)

// Sidecar field keys.
var (
	AddrKey        = capitan.NewStringKey("addr")
	FilenameKey    = capitan.NewStringKey("filename")
	BytesKey       = capitan.NewInt64Key("bytes")
	IndexBytesKey  = capitan.NewIntKey("index_bytes")
	ChunksKey      = capitan.NewIntKey("chunks")
	WorkDirKey     = capitan.NewStringKey("work_dir")
)
