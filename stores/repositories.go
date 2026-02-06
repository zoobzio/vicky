package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Repositories provides database access for repository records.
type Repositories struct {
	*sum.Database[models.Repository]
}

// NewRepositories creates a new repositories store.
func NewRepositories(db *sqlx.DB, renderer astql.Renderer) (*Repositories, error) {
	database, err := sum.NewDatabase[models.Repository](db, "repositories", renderer)
	if err != nil {
		return nil, err
	}
	return &Repositories{Database: database}, nil
}

// ListByUserID retrieves all repositories registered by a user.
func (s *Repositories) ListByUserID(ctx context.Context, userID int64) ([]*models.Repository, error) {
	return s.Executor().Soy().Query().
		Where("user_id", "=", "user_id").
		Exec(ctx, map[string]any{"user_id": userID})
}

// GetByUserAndGitHubID retrieves a repository by user and GitHub repo ID.
func (s *Repositories) GetByUserAndGitHubID(ctx context.Context, userID, githubID int64) (*models.Repository, error) {
	return s.Executor().Soy().Select().
		Where("user_id", "=", "user_id").
		Where("github_id", "=", "github_id").
		Exec(ctx, map[string]any{"user_id": userID, "github_id": githubID})
}

// GetByUserOwnerAndName retrieves a repository by user, owner, and name.
func (s *Repositories) GetByUserOwnerAndName(ctx context.Context, userID int64, owner, name string) (*models.Repository, error) {
	return s.Executor().Soy().Select().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("name", "=", "name").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "name": name})
}
