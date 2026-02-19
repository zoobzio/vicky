package events

import (
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// JobCreatedEvent is emitted when a new ingestion job is created.
type JobCreatedEvent struct {
	Job *models.Job `json:"job"`
}

// JobProgressEvent is emitted when job progress updates.
type JobProgressEvent struct {
	JobID          int64            `json:"job_id"`
	Stage          models.JobStage  `json:"stage"`
	Status         models.JobStatus `json:"status"`
	Progress       int              `json:"progress"`
	ItemsTotal     int              `json:"items_total"`
	ItemsProcessed int              `json:"items_processed"`
}

// Job lifecycle signals.
var (
	JobCreatedSignal   = capitan.NewSignal("vicky.job.created", "Ingestion job created and queued")
	JobStartedSignal   = capitan.NewSignal("vicky.job.started", "Ingestion job processing started")
	JobProgressSignal  = capitan.NewSignal("vicky.job.progress", "Ingestion job progress update")
	JobCompletedSignal = capitan.NewSignal("vicky.job.completed", "Ingestion job completed successfully")
	JobFailedSignal    = capitan.NewSignal("vicky.job.failed", "Ingestion job failed")
	JobCancelledSignal = capitan.NewSignal("vicky.job.cancelled", "Ingestion job cancelled by user")
)

// Job provides access to job lifecycle events.
var Job = struct {
	Created   sum.Event[JobCreatedEvent]
	Started   sum.Event[JobProgressEvent]
	Progress  sum.Event[JobProgressEvent]
	Completed sum.Event[JobProgressEvent]
	Failed    sum.Event[JobProgressEvent]
	Cancelled sum.Event[JobProgressEvent]
}{
	Created:   sum.NewInfoEvent[JobCreatedEvent](JobCreatedSignal),
	Started:   sum.NewInfoEvent[JobProgressEvent](JobStartedSignal),
	Progress:  sum.NewDebugEvent[JobProgressEvent](JobProgressSignal),
	Completed: sum.NewInfoEvent[JobProgressEvent](JobCompletedSignal),
	Failed:    sum.NewErrorEvent[JobProgressEvent](JobFailedSignal),
	Cancelled: sum.NewWarnEvent[JobProgressEvent](JobCancelledSignal),
}
