package stores

import (
	"context"
	"fmt"

	"github.com/zoobzio/grub"
	"github.com/zoobzio/vicky/models"
)

// Blobs provides blob storage via a grub BucketProvider.
type Blobs struct {
	bucket *grub.Bucket[models.Blob]
}

// NewBlobs creates a new blob store.
func NewBlobs(provider grub.BucketProvider) *Blobs {
	return &Blobs{bucket: grub.NewBucket[models.Blob](provider)}
}

// GetByPath retrieves a blob by its domain coordinates.
func (s *Blobs) GetByPath(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error) {
	return s.bucket.Get(ctx, blobKey(userID, owner, repo, tag, path))
}

// PutBlob stores a blob using its domain coordinates as the key.
func (s *Blobs) PutBlob(ctx context.Context, userID int64, blob *models.Blob) error {
	return s.bucket.Put(ctx, &grub.Object[models.Blob]{
		Key:         blobKey(userID, blob.Owner, blob.Repo, blob.Tag, blob.Path),
		ContentType: "application/json",
		Data:        *blob,
	})
}

// DeleteByPath removes a blob by its domain coordinates.
func (s *Blobs) DeleteByPath(ctx context.Context, userID int64, owner, repo, tag, path string) error {
	return s.bucket.Delete(ctx, blobKey(userID, owner, repo, tag, path))
}

// ListByVersion returns object info for all blobs in a version.
func (s *Blobs) ListByVersion(ctx context.Context, userID int64, owner, repo, tag string, limit int) ([]grub.ObjectInfo, error) {
	return s.bucket.List(ctx, versionPrefix(userID, owner, repo, tag), limit)
}

// ListByRepo returns object info for all blobs in a repository.
func (s *Blobs) ListByRepo(ctx context.Context, userID int64, owner, repo string, limit int) ([]grub.ObjectInfo, error) {
	return s.bucket.List(ctx, repoPrefix(userID, owner, repo), limit)
}

func blobKey(userID int64, owner, repo, tag, path string) string {
	return fmt.Sprintf("%d/%s/%s/%s/%s", userID, owner, repo, tag, path)
}

func versionPrefix(userID int64, owner, repo, tag string) string {
	return fmt.Sprintf("%d/%s/%s/%s/", userID, owner, repo, tag)
}

func repoPrefix(userID int64, owner, repo string) string {
	return fmt.Sprintf("%d/%s/%s/", userID, owner, repo)
}
