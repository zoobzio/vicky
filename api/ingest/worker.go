package ingest

import (
	"context"
	"errors"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/pipz"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/api/events"
	"github.com/zoobzio/vicky/models"
)

// Worker pool configuration.
const (
	// DefaultWorkers is the number of concurrent ingestion jobs.
	DefaultWorkers = 4
)

// Worker pool identity.
var WorkerPoolID = pipz.NewIdentity("ingest-worker-pool", "Worker pool for ingestion jobs")

// Worker manages the ingestion worker pool and event listener.
type Worker struct {
	pool     *pipz.WorkerPool[*models.Job]
	pipeline *pipz.Sequence[*models.Job]
	listener *capitan.Listener
}

// NewWorker creates a new ingestion worker.
// Dependencies are resolved from the sum registry at runtime.
func NewWorker() *Worker {
	pipeline := NewPipeline()

	return &Worker{
		pool:     pipz.NewWorkerPool(WorkerPoolID, DefaultWorkers, pipeline),
		pipeline: pipeline,
	}
}

// Start begins listening for job creation events and processing jobs.
func (w *Worker) Start(ctx context.Context) {
	w.listener = events.Job.Created.Listen(func(ctx context.Context, e events.JobCreatedEvent) {
		// Use the job from the event directly
		w.processJob(ctx, e.Job)
	})

	capitan.Emit(ctx, events.StartupWorkerReady, events.StartupWorkersKey.Field(DefaultWorkers))
}

// Stop gracefully shuts down the worker pool and listener.
func (w *Worker) Stop() error {
	if w.listener != nil {
		w.listener.Close()
	}
	return w.pool.Close()
}

// processJob handles a single ingestion job.
func (w *Worker) processJob(ctx context.Context, job *models.Job) {
	// Resolve jobs contract from registry
	jobs := sum.MustUse[contracts.Jobs](ctx)

	// Mark job as running
	if err := jobs.Start(ctx, job.ID); err != nil {
		events.Job.Failed.Emit(ctx, events.JobProgressEvent{
			JobID:  job.ID,
			Stage:  job.Stage,
			Status: models.JobStatusFailed,
		})
		return
	}

	// Emit started event
	events.Job.Started.Emit(ctx, events.JobProgressEvent{
		JobID:    job.ID,
		Stage:    job.Stage,
		Status:   models.JobStatusRunning,
		Progress: 0,
	})

	// Process through pipeline
	result, err := w.pool.Process(ctx, job)
	if err != nil {
		// Check if cancellation
		if errors.Is(err, ErrJobCancelled) {
			// Mark job as cancelled (not failed)
			_ = jobs.MarkCancelled(ctx, job.ID)

			// Emit cancelled event
			events.Job.Cancelled.Emit(ctx, events.JobProgressEvent{
				JobID:    job.ID,
				Stage:    job.Stage,
				Status:   models.JobStatusCancelled,
				Progress: job.Progress,
			})
			return
		}

		// Mark job as failed
		_ = jobs.MarkFailed(ctx, job.ID, err.Error())

		// Emit failed event
		events.Job.Failed.Emit(ctx, events.JobProgressEvent{
			JobID:    job.ID,
			Stage:    job.Stage,
			Status:   models.JobStatusFailed,
			Progress: job.Progress,
		})
		return
	}

	// Mark job as completed
	_ = jobs.MarkCompleted(ctx, result.ID)

	// Emit completed event
	events.Job.Completed.Emit(ctx, events.JobProgressEvent{
		JobID:    result.ID,
		Stage:    models.JobStageStore,
		Status:   models.JobStatusCompleted,
		Progress: 100,
	})
}
