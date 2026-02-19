package contracts

import (
	"context"

	"github.com/zoobzio/grub"
	"github.com/zoobzio/vicky/models"
)

// Blobs defines the contract for blob storage operations.
type Blobs interface {
	// GetByPath retrieves a blob by its domain coordinates.
	GetByPath(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error)
	// PutBlob stores a blob using its domain coordinates as the key.
	PutBlob(ctx context.Context, userID int64, blob *models.Blob) error
	// DeleteByPath removes a blob by its domain coordinates.
	DeleteByPath(ctx context.Context, userID int64, owner, repo, tag, path string) error
	// ListByVersion returns object info for all blobs in a version.
	ListByVersion(ctx context.Context, userID int64, owner, repo, tag string, limit int) ([]grub.ObjectInfo, error)
	// ListByRepo returns object info for all blobs in a repository.
	ListByRepo(ctx context.Context, userID int64, owner, repo string, limit int) ([]grub.ObjectInfo, error)
}
