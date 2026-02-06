package indexers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageConfig holds minio connection settings.
type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Storage provides blob fetching from minio.
type Storage struct {
	client *minio.Client
	bucket string
}

// NewStorage creates a new minio storage client.
func NewStorage(cfg StorageConfig) (*Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &Storage{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// blob mirrors the vicky Blob model for JSON deserialization.
type blob struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// FetchBlobs downloads all blobs for a version and writes them as source files
// to the destination directory. Returns the number of files written.
func (s *Storage) FetchBlobs(ctx context.Context, userID int64, owner, repo, tag, destDir string) (int, error) {
	prefix := fmt.Sprintf("%d/%s/%s/%s/", userID, owner, repo, tag)

	objectCh := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	count := 0
	for obj := range objectCh {
		if obj.Err != nil {
			return count, fmt.Errorf("list objects: %w", obj.Err)
		}

		// Download object
		object, err := s.client.GetObject(ctx, s.bucket, obj.Key, minio.GetObjectOptions{})
		if err != nil {
			return count, fmt.Errorf("get object %s: %w", obj.Key, err)
		}

		data, err := io.ReadAll(object)
		_ = object.Close()
		if err != nil {
			return count, fmt.Errorf("read object %s: %w", obj.Key, err)
		}

		// Blobs are stored as JSON with path and content fields
		var b blob
		if err := json.Unmarshal(data, &b); err != nil {
			return count, fmt.Errorf("unmarshal blob %s: %w", obj.Key, err)
		}

		// Write content to file at the original path
		destPath := filepath.Join(destDir, b.Path)
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return count, fmt.Errorf("create dir for %s: %w", b.Path, err)
		}

		if err := os.WriteFile(destPath, []byte(b.Content), 0o644); err != nil {
			return count, fmt.Errorf("write file %s: %w", b.Path, err)
		}

		count++
	}

	return count, nil
}
