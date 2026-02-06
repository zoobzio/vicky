package contracts

import (
	"context"

	"github.com/zoobzio/vicky/models"
)

// IngestionConfigs defines the contract for ingestion config storage operations.
type IngestionConfigs interface {
	// Get retrieves a config by its string key (ID as string).
	Get(ctx context.Context, key string) (*models.IngestionConfig, error)

	// Set creates or updates a config.
	Set(ctx context.Context, key string, config *models.IngestionConfig) error

	// GetByRepositoryID retrieves the config for a repository.
	GetByRepositoryID(ctx context.Context, repositoryID int64) (*models.IngestionConfig, error)
}
