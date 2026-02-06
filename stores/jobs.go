package stores

import (
	"context"
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
	return s.Executor().Soy().Query().
		Where("version_id", "=", "version_id").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"version_id": versionID})
}

// LatestByVersionID retrieves the most recent job for a version.
func (s *Jobs) LatestByVersionID(ctx context.Context, versionID int64) (*models.Job, error) {
	return s.Executor().Soy().Select().
		Where("version_id", "=", "version_id").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"version_id": versionID})
}

// ListByUser retrieves all jobs for a user.
func (s *Jobs) ListByUser(ctx context.Context, userID int64) ([]*models.Job, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"user_id": userID})
}

// ListByStatus retrieves jobs with a specific status for a user.
func (s *Jobs) ListByStatus(ctx context.Context, userID int64, status models.JobStatus) ([]*models.Job, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("status", "=", "status").
		OrderBy("created_at", "DESC").
		Exec(ctx, map[string]any{"user_id": userID, "status": status})
}

// UpdateProgress updates the job's stage, progress percentage, and items processed.
func (s *Jobs) UpdateProgress(ctx context.Context, id int64, stage models.JobStage, progress int, itemsProcessed int) error {
	_, err := s.Executor().Soy().Modify().
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
	_, err := s.Executor().Soy().Modify().
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
	_, err := s.Executor().Soy().Modify().
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
	_, err := s.Executor().Soy().Modify().
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
