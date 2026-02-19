package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/models"
)

// IngestionConfigs provides database access for ingestion config records.
type IngestionConfigs struct {
	*sum.Database[models.IngestionConfig]
}

// NewIngestionConfigs creates a new ingestion configs store.
func NewIngestionConfigs(db *sqlx.DB, renderer astql.Renderer) (*IngestionConfigs, error) {
	database, err := sum.NewDatabase[models.IngestionConfig](db, "ingestion_configs", renderer)
	if err != nil {
		return nil, err
	}
	return &IngestionConfigs{Database: database}, nil
}

// GetByRepositoryID retrieves the config for a repository.
func (s *IngestionConfigs) GetByRepositoryID(ctx context.Context, repositoryID int64) (*models.IngestionConfig, error) {
	return s.Select().
		Where("repository_id", "=", "repository_id").
		Exec(ctx, map[string]any{"repository_id": repositoryID})
}
