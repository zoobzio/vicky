package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/stores"
)

// Jobs defines the contract for admin job operations.
type Jobs interface {
	// Get retrieves a job by ID.
	Get(ctx context.Context, key string) (*models.Job, error)
	// List retrieves jobs with optional filtering and pagination.
	List(ctx context.Context, filter *stores.JobFilter, limit, offset int) ([]*models.Job, error)
	// Count returns the total number of jobs matching the filter.
	Count(ctx context.Context, filter *stores.JobFilter) (int, error)
	// RequestCancellation requests cancellation of a running/pending job.
	RequestCancellation(ctx context.Context, id int64) error
	// ListByStatus retrieves jobs with a specific status.
	ListByStatus(ctx context.Context, userID int64, status models.JobStatus) ([]*models.Job, error)
	// LatestByVersionID retrieves the most recent job for a version.
	LatestByVersionID(ctx context.Context, versionID int64) (*models.Job, error)
	// CountByStatus returns counts of jobs by status (for stats).
	CountByStatus(ctx context.Context, userID *int64) (map[models.JobStatus]int, error)
	// Set creates or updates a job.
	Set(ctx context.Context, key string, job *models.Job) error
}
