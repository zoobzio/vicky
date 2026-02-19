package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// Jobs defines the contract for job storage operations.
type Jobs interface {
	// Get retrieves a job by its string key (ID as string).
	Get(ctx context.Context, key string) (*models.Job, error)

	// Set creates or updates a job.
	Set(ctx context.Context, key string, job *models.Job) error

	// ListByVersionID retrieves all jobs for a version (job history).
	ListByVersionID(ctx context.Context, versionID int64) ([]*models.Job, error)

	// LatestByVersionID retrieves the most recent job for a version.
	LatestByVersionID(ctx context.Context, versionID int64) (*models.Job, error)

	// ListByUser retrieves all jobs for a user.
	ListByUser(ctx context.Context, userID int64) ([]*models.Job, error)

	// ListByStatus retrieves jobs with a specific status for a user.
	ListByStatus(ctx context.Context, userID int64, status models.JobStatus) ([]*models.Job, error)

	// UpdateProgress updates the job's stage, progress percentage, and items processed.
	UpdateProgress(ctx context.Context, id int64, stage models.JobStage, progress int, itemsProcessed int) error

	// Start marks the job as running and sets the started timestamp.
	Start(ctx context.Context, id int64) error

	// MarkFailed marks the job as failed with an error message.
	MarkFailed(ctx context.Context, id int64, errMsg string) error

	// MarkCompleted marks the job as completed and sets the completed timestamp.
	MarkCompleted(ctx context.Context, id int64) error

	// MarkCancelled marks a job as cancelled after worker abort.
	MarkCancelled(ctx context.Context, id int64) error

	// IsCancelling checks if a job has been marked for cancellation.
	IsCancelling(ctx context.Context, id int64) (bool, error)
}
