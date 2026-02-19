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
	return s.Query().
		Where("user_id", "=", "user_id").
		Exec(ctx, map[string]any{"user_id": userID})
}

// GetByUserAndGitHubID retrieves a repository by user and GitHub repo ID.
func (s *Repositories) GetByUserAndGitHubID(ctx context.Context, userID, githubID int64) (*models.Repository, error) {
	return s.Select().
		Where("user_id", "=", "user_id").
		Where("github_id", "=", "github_id").
		Exec(ctx, map[string]any{"user_id": userID, "github_id": githubID})
}

// GetByUserOwnerAndName retrieves a repository by user, owner, and name.
func (s *Repositories) GetByUserOwnerAndName(ctx context.Context, userID int64, owner, name string) (*models.Repository, error) {
	return s.Select().
		Where("user_id", "=", "user_id").
		Where("owner", "=", "owner").
		Where("name", "=", "name").
		Exec(ctx, map[string]any{"user_id": userID, "owner": owner, "name": name})
}

// RepositoryFilter defines optional filters for repository queries.
type RepositoryFilter struct {
	Owner  *string
	Name   *string
	UserID *int64
}

// List retrieves repositories with optional filtering and pagination.
// Pass nil filter to list all repositories.
func (s *Repositories) List(ctx context.Context, filter *RepositoryFilter, limit, offset int) ([]*models.Repository, error) {
	builder := s.Query()

	if filter != nil {
		if filter.Owner != nil {
			builder = builder.Where("owner", "ILIKE", "owner")
		}
		if filter.Name != nil {
			builder = builder.Where("name", "ILIKE", "name")
		}
		if filter.UserID != nil {
			builder = builder.Where("user_id", "=", "user_id")
		}
	}

	builder = builder.
		OrderBy("created_at", "DESC").
		Limit(limit).
		Offset(offset)

	params := make(map[string]any)
	if filter != nil {
		if filter.Owner != nil {
			params["owner"] = "%" + *filter.Owner + "%"
		}
		if filter.Name != nil {
			params["name"] = "%" + *filter.Name + "%"
		}
		if filter.UserID != nil {
			params["user_id"] = *filter.UserID
		}
	}

	return builder.Exec(ctx, params)
}

// Count returns the total number of repositories matching the filter.
// Pass nil filter to count all repositories.
func (s *Repositories) Count(ctx context.Context, filter *RepositoryFilter) (int, error) {
	builder := s.Database.Count()

	params := make(map[string]any)
	if filter != nil {
		if filter.Owner != nil {
			builder = builder.Where("owner", "ILIKE", "owner")
			params["owner"] = "%" + *filter.Owner + "%"
		}
		if filter.Name != nil {
			builder = builder.Where("name", "ILIKE", "name")
			params["name"] = "%" + *filter.Name + "%"
		}
		if filter.UserID != nil {
			builder = builder.Where("user_id", "=", "user_id")
			params["user_id"] = *filter.UserID
		}
	}

	count, err := builder.Exec(ctx, params)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
