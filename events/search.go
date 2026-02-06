package events

import (
	"time"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
)

// ChunkSearchEvent is emitted during chunk similarity search.
type ChunkSearchEvent struct {
	UserID       int64         `json:"user_id"`
	RepositoryID int64         `json:"repository_id,omitempty"`
	VersionID    int64         `json:"version_id,omitempty"`
	Query        string        `json:"query"`
	Limit        int           `json:"limit"`
	ResultCount  int           `json:"result_count,omitempty"`
	TopScore     float64       `json:"top_score,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// SymbolSearchEvent is emitted during symbol search.
type SymbolSearchEvent struct {
	UserID       int64         `json:"user_id"`
	RepositoryID int64         `json:"repository_id,omitempty"`
	VersionID    int64         `json:"version_id,omitempty"`
	Query        string        `json:"query"`
	Kind         string        `json:"kind,omitempty"`
	Limit        int           `json:"limit"`
	ResultCount  int           `json:"result_count,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// DocumentSearchEvent is emitted during document similarity search.
type DocumentSearchEvent struct {
	UserID       int64         `json:"user_id"`
	RepositoryID int64         `json:"repository_id,omitempty"`
	VersionID    int64         `json:"version_id,omitempty"`
	DocumentID   int64         `json:"document_id"`
	Limit        int           `json:"limit"`
	ResultCount  int           `json:"result_count,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// Chunk search signals.
var (
	ChunkSearchStartedSignal   = capitan.NewSignal("vicky.search.chunk.started", "Chunk search initiated")
	ChunkSearchCompletedSignal = capitan.NewSignal("vicky.search.chunk.completed", "Chunk search completed")
	ChunkSearchFailedSignal    = capitan.NewSignal("vicky.search.chunk.failed", "Chunk search failed")
)

// Symbol search signals.
var (
	SymbolSearchStartedSignal   = capitan.NewSignal("vicky.search.symbol.started", "Symbol search initiated")
	SymbolSearchCompletedSignal = capitan.NewSignal("vicky.search.symbol.completed", "Symbol search completed")
	SymbolSearchFailedSignal    = capitan.NewSignal("vicky.search.symbol.failed", "Symbol search failed")
)

// Document search signals.
var (
	DocumentSearchStartedSignal   = capitan.NewSignal("vicky.search.document.started", "Document search initiated")
	DocumentSearchCompletedSignal = capitan.NewSignal("vicky.search.document.completed", "Document search completed")
	DocumentSearchFailedSignal    = capitan.NewSignal("vicky.search.document.failed", "Document search failed")
)

// chunkSearch provides access to chunk search events.
var chunkSearch = struct {
	Started   sum.Event[ChunkSearchEvent]
	Completed sum.Event[ChunkSearchEvent]
	Failed    sum.Event[ChunkSearchEvent]
}{
	Started:   sum.NewDebugEvent[ChunkSearchEvent](ChunkSearchStartedSignal),
	Completed: sum.NewInfoEvent[ChunkSearchEvent](ChunkSearchCompletedSignal),
	Failed:    sum.NewErrorEvent[ChunkSearchEvent](ChunkSearchFailedSignal),
}

// symbolSearch provides access to symbol search events.
var symbolSearch = struct {
	Started   sum.Event[SymbolSearchEvent]
	Completed sum.Event[SymbolSearchEvent]
	Failed    sum.Event[SymbolSearchEvent]
}{
	Started:   sum.NewDebugEvent[SymbolSearchEvent](SymbolSearchStartedSignal),
	Completed: sum.NewInfoEvent[SymbolSearchEvent](SymbolSearchCompletedSignal),
	Failed:    sum.NewErrorEvent[SymbolSearchEvent](SymbolSearchFailedSignal),
}

// documentSearch provides access to document search events.
var documentSearch = struct {
	Started   sum.Event[DocumentSearchEvent]
	Completed sum.Event[DocumentSearchEvent]
	Failed    sum.Event[DocumentSearchEvent]
}{
	Started:   sum.NewDebugEvent[DocumentSearchEvent](DocumentSearchStartedSignal),
	Completed: sum.NewInfoEvent[DocumentSearchEvent](DocumentSearchCompletedSignal),
	Failed:    sum.NewErrorEvent[DocumentSearchEvent](DocumentSearchFailedSignal),
}

// Search provides access to search events by type.
var Search = struct {
	Chunk struct {
		Started   sum.Event[ChunkSearchEvent]
		Completed sum.Event[ChunkSearchEvent]
		Failed    sum.Event[ChunkSearchEvent]
	}
	Symbol struct {
		Started   sum.Event[SymbolSearchEvent]
		Completed sum.Event[SymbolSearchEvent]
		Failed    sum.Event[SymbolSearchEvent]
	}
	Document struct {
		Started   sum.Event[DocumentSearchEvent]
		Completed sum.Event[DocumentSearchEvent]
		Failed    sum.Event[DocumentSearchEvent]
	}
}{
	Chunk:    chunkSearch,
	Symbol:   symbolSearch,
	Document: documentSearch,
}
