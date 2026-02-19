package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// Users provides database access for user records.
type Users struct {
	*sum.Database[models.User]
}

// NewUsers creates a new users store.
func NewUsers(db *sqlx.DB, renderer astql.Renderer) (*Users, error) {
	database, err := sum.NewDatabase[models.User](db, "users", renderer)
	if err != nil {
		return nil, err
	}
	return &Users{Database: database}, nil
}

// GetByLogin retrieves a user by their GitHub login.
func (s *Users) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	return s.Select().
		Where("login", "=", "login").
		Exec(ctx, map[string]any{"login": login})
}

// UserFilter defines optional filters for user queries.
type UserFilter struct {
	Login *string
	Email *string
	Name  *string
}

// List retrieves users with optional filtering and pagination.
// Pass nil filter to list all users.
func (s *Users) List(ctx context.Context, filter *UserFilter, limit, offset int) ([]*models.User, error) {
	builder := s.Query()

	if filter != nil {
		if filter.Login != nil {
			builder = builder.Where("login", "ILIKE", "login")
		}
		if filter.Email != nil {
			builder = builder.Where("email", "ILIKE", "email")
		}
		if filter.Name != nil {
			builder = builder.Where("name", "ILIKE", "name")
		}
	}

	builder = builder.
		OrderBy("created_at", "DESC").
		Limit(limit).
		Offset(offset)

	params := make(map[string]any)
	if filter != nil {
		if filter.Login != nil {
			params["login"] = "%" + *filter.Login + "%"
		}
		if filter.Email != nil {
			params["email"] = "%" + *filter.Email + "%"
		}
		if filter.Name != nil {
			params["name"] = "%" + *filter.Name + "%"
		}
	}

	return builder.Exec(ctx, params)
}

// Count returns the total number of users matching the filter.
// Pass nil filter to count all users.
func (s *Users) Count(ctx context.Context, filter *UserFilter) (int, error) {
	builder := s.Database.Count()

	params := make(map[string]any)
	if filter != nil {
		if filter.Login != nil {
			builder = builder.Where("login", "ILIKE", "login")
			params["login"] = "%" + *filter.Login + "%"
		}
		if filter.Email != nil {
			builder = builder.Where("email", "ILIKE", "email")
			params["email"] = "%" + *filter.Email + "%"
		}
		if filter.Name != nil {
			builder = builder.Where("name", "ILIKE", "name")
			params["name"] = "%" + *filter.Name + "%"
		}
	}

	count, err := builder.Exec(ctx, params)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
