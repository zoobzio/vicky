package stores

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Jobs provides database access for job records.
type Jobs struct {
	*sum.Database[models.Job]
}

// JobFilter defines filtering options for job queries.
type JobFilter struct {
	UserID       *int64            // Filter by user
	RepositoryID *int64            // Filter by repository
	VersionID    *int64            // Filter by version
	Status       *models.JobStatus // Filter by status
	Owner        *string           // Filter by repository owner
	RepoName     *string           // Filter by repository name
}

// NewJobs creates a new jobs store.
func NewJobs(db *sqlx.DB, renderer astql.Renderer) (*Jobs, error) {
	database, err := sum.NewDatabase[models.Job](db, "jobs", renderer)
	if err != nil {
		return nil, err
	}
	return &Jobs{Database: database}, nil
}

// ListByVersionID retrieves all jobs for a version (job history).
func (s *Jobs) ListByVersionID(ctx context.Context, versionID int64) ([]*models.Job, error) {
	return s.Query().
		Where("version_id", "=", "version_id").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"version_id": versionID})
}

// LatestByVersionID retrieves the most recent job for a version.
func (s *Jobs) LatestByVersionID(ctx context.Context, versionID int64) (*models.Job, error) {
	return s.Select().
		Where("version_id", "=", "version_id").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"version_id": versionID})
}

// ListByUser retrieves all jobs for a user.
func (s *Jobs) ListByUser(ctx context.Context, userID int64) ([]*models.Job, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"user_id": userID})
}

// ListByStatus retrieves jobs with a specific status for a user.
func (s *Jobs) ListByStatus(ctx context.Context, userID int64, status models.JobStatus) ([]*models.Job, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		Where("status", "=", "status").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"user_id": userID, "status": status})
}

// UpdateProgress updates the job's stage, progress percentage, and items processed.
func (s *Jobs) UpdateProgress(ctx context.Context, id int64, stage models.JobStage, progress int, itemsProcessed int) error {
	_, err := s.Modify().
		Set("stage", "stage").
		Set("progress", "progress").
		Set("items_processed", "items_processed").
		Set("updated_at", "updated_at").
		Where("id", "=", "id").
		Exec(ctx, map[string]any{
			"id":              id,
			"stage":           stage,
			"progress":        progress,
			"items_processed": itemsProcessed,
			"updated_at":      time.Now(),
		})
	return err
}

// Start marks the job as running and sets the started timestamp.
func (s *Jobs) Start(ctx context.Context, id int64) error {
	now := time.Now()
	_, err := s.Modify().
		Set("status", "status").
		Set("started_at", "started_at").
		Set("updated_at", "updated_at").
		Where("id", "=", "id").
		Exec(ctx, map[string]any{
			"id":         id,
			"status":     models.JobStatusRunning,
			"started_at": &now,
			"updated_at": now,
		})
	return err
}

// MarkFailed marks the job as failed with an error message.
func (s *Jobs) MarkFailed(ctx context.Context, id int64, errMsg string) error {
	_, err := s.Modify().
		Set("status", "status").
		Set("error", "error").
		Set("updated_at", "updated_at").
		Where("id", "=", "id").
		Exec(ctx, map[string]any{
			"id":         id,
			"status":     models.JobStatusFailed,
			"error":      &errMsg,
			"updated_at": time.Now(),
		})
	return err
}

// MarkCompleted marks the job as completed and sets the completed timestamp.
func (s *Jobs) MarkCompleted(ctx context.Context, id int64) error {
	now := time.Now()
	_, err := s.Modify().
		Set("status", "status").
		Set("progress", "progress").
		Set("completed_at", "completed_at").
		Set("updated_at", "updated_at").
		Where("id", "=", "id").
		Exec(ctx, map[string]any{
			"id":           id,
			"status":       models.JobStatusCompleted,
			"progress":     100,
			"completed_at": &now,
			"updated_at":   now,
		})
	return err
}

// RequestCancellation marks a job for cancellation.
// Returns error if job is already completed, failed, or cancelled.
func (s *Jobs) RequestCancellation(ctx context.Context, id int64) error {
	// First check current status
	job, err := s.Get(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		return err
	}
	if job.Status != models.JobStatusPending && job.Status != models.JobStatusRunning {
		return fmt.Errorf("job cannot be cancelled (not in pending or running state)")
	}

	// Update to cancelling status
	job.Status = models.JobStatusCancelling
	job.UpdatedAt = time.Now()
	return s.Set(ctx, fmt.Sprintf("%d", id), job)
}

// MarkCancelled marks a job as cancelled after worker abort.
func (s *Jobs) MarkCancelled(ctx context.Context, id int64) error {
	now := time.Now()
	_, err := s.Modify().
		Set("status", "status").
		Set("completed_at", "completed_at").
		Set("updated_at", "updated_at").
		Where("id", "=", "id").
		Exec(ctx, map[string]any{
			"id":           id,
			"status":       models.JobStatusCancelled,
			"completed_at": &now,
			"updated_at":   now,
		})
	return err
}

// IsCancelling checks if a job has been marked for cancellation.
func (s *Jobs) IsCancelling(ctx context.Context, id int64) (bool, error) {
	job, err := s.Select().
		Where("id", "=", "id").
		Exec(ctx, map[string]any{"id": id})
	if err != nil {
		return false, err
	}
	return job.Status == models.JobStatusCancelling, nil
}

// List retrieves jobs with optional filtering and pagination.
func (s *Jobs) List(ctx context.Context, filter *JobFilter, limit, offset int) ([]*models.Job, error) {
	q := s.Query().
		OrderBy("created_at", "DESC").
		Limit(limit).
		Offset(offset)

	params := make(map[string]any)

	if filter != nil {
		if filter.UserID != nil {
			q = q.Where("user_id", "=", "user_id")
			params["user_id"] = *filter.UserID
		}
		if filter.RepositoryID != nil {
			q = q.Where("repository_id", "=", "repository_id")
			params["repository_id"] = *filter.RepositoryID
		}
		if filter.VersionID != nil {
			q = q.Where("version_id", "=", "version_id")
			params["version_id"] = *filter.VersionID
		}
		if filter.Status != nil {
			q = q.Where("status", "=", "status")
			params["status"] = *filter.Status
		}
		if filter.Owner != nil {
			q = q.Where("owner", "=", "owner")
			params["owner"] = *filter.Owner
		}
		if filter.RepoName != nil {
			q = q.Where("repo_name", "=", "repo_name")
			params["repo_name"] = *filter.RepoName
		}
	}

	return q.Exec(ctx, params)
}

// Count returns the total number of jobs matching the filter.
func (s *Jobs) Count(ctx context.Context, filter *JobFilter) (int, error) {
	q := s.Query()

	params := make(map[string]any)

	if filter != nil {
		if filter.UserID != nil {
			q = q.Where("user_id", "=", "user_id")
			params["user_id"] = *filter.UserID
		}
		if filter.RepositoryID != nil {
			q = q.Where("repository_id", "=", "repository_id")
			params["repository_id"] = *filter.RepositoryID
		}
		if filter.VersionID != nil {
			q = q.Where("version_id", "=", "version_id")
			params["version_id"] = *filter.VersionID
		}
		if filter.Status != nil {
			q = q.Where("status", "=", "status")
			params["status"] = *filter.Status
		}
		if filter.Owner != nil {
			q = q.Where("owner", "=", "owner")
			params["owner"] = *filter.Owner
		}
		if filter.RepoName != nil {
			q = q.Where("repo_name", "=", "repo_name")
			params["repo_name"] = *filter.RepoName
		}
	}

	jobs, err := q.Exec(ctx, params)
	if err != nil {
		return 0, err
	}
	return len(jobs), nil
}

// CountByStatus returns counts of jobs by status (for stats).
func (s *Jobs) CountByStatus(ctx context.Context, userID *int64) (map[models.JobStatus]int, error) {
	q := s.Query()

	params := make(map[string]any)

	if userID != nil {
		q = q.Where("user_id", "=", "user_id")
		params["user_id"] = *userID
	}

	jobs, err := q.Exec(ctx, params)
	if err != nil {
		return nil, err
	}

	counts := make(map[models.JobStatus]int)
	for _, job := range jobs {
		counts[job.Status]++
	}

	return counts, nil
}
