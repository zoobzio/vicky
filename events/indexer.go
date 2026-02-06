package events

import (
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/external/indexer"
)

// IndexRequestedEvent is emitted when a job requires SCIP indexing.
type IndexRequestedEvent struct {
	Request indexer.Request `json:"request"`
}

// IndexCompletedEvent is emitted when SCIP indexing completes successfully.
type IndexCompletedEvent struct {
	Result indexer.Result `json:"result"`
}

// IndexFailedEvent is emitted when SCIP indexing fails.
type IndexFailedEvent struct {
	JobID int64  `json:"job_id"`
	Error string `json:"error"`
}

// Indexer lifecycle signals.
var (
	IndexRequestedSignal = capitan.NewSignal("vicky.indexer.requested", "SCIP indexing requested")
	IndexCompletedSignal = capitan.NewSignal("vicky.indexer.completed", "SCIP indexing completed")
	IndexFailedSignal    = capitan.NewSignal("vicky.indexer.failed", "SCIP indexing failed")
)

// Indexer provides access to indexer lifecycle events.
var Indexer = struct {
	Requested sum.Event[IndexRequestedEvent]
	Completed sum.Event[IndexCompletedEvent]
	Failed    sum.Event[IndexFailedEvent]
}{
	Requested: sum.NewInfoEvent[IndexRequestedEvent](IndexRequestedSignal),
	Completed: sum.NewInfoEvent[IndexCompletedEvent](IndexCompletedSignal),
	Failed:    sum.NewErrorEvent[IndexFailedEvent](IndexFailedSignal),
}
