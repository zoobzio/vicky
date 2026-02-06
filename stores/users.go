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
	return s.Executor().Soy().Select().
		Where("login", "=", "login").
		Exec(ctx, map[string]any{"login": login})
}
