package stores

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Versions provides database access for version records.
type Versions struct {
	*sum.Database[models.Version]
}

// NewVersions creates a new versions store.
func NewVersions(db *sqlx.DB, renderer astql.Renderer) (*Versions, error) {
	database, err := sum.NewDatabase[models.Version](db, "versions", renderer)
	if err != nil {
		return nil, err
	}
	return &Versions{Database: database}, nil
}

// ListByUserAndRepo retrieves all versions for a user's repository.
func (s *Versions) ListByUserAndRepo(ctx context.Context, userID int64, owner, repoName string) ([]*models.Version, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName})
}

// GetByUserRepoAndTag retrieves a specific version by natural identifiers.
func (s *Versions) GetByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) (*models.Version, error) {
	return s.Executor().Soy().Select().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("repo_name", "=", "repo_name").
		Where("tag", "=", "tag").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "repo_name": repoName, "tag": tag})
}

// UpdateStatus updates the ingestion status of a version.
func (s *Versions) UpdateStatus(ctx context.Context, id int64, status models.VersionStatus, versionErr *string) (*models.Version, error) {
	return s.Executor().Soy().Modify().
		Set("status", "status").
		Set("error", "error").
		Set("updated_at", "updated_at").
		Where("id", "=", "id").
		Exec(ctx, map[string]any{
			"id":         id,
			"status":     status,
			"error":      versionErr,
			"updated_at": time.Now(),
		})
}
