//go:build testing

package ingest

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/zoobzio/vicky/external/github"
	"github.com/zoobzio/vicky/models"
	vickytest "github.com/zoobzio/vicky/testing"
)

func TestFetchStage(t *testing.T) {
	version := vickytest.NewVersion(t)

	mu := &vickytest.MockUsers{}
	mg := &vickytest.MockGitHub{
		OnGetTree: func(ctx context.Context, token, owner, repo, ref string) ([]github.TreeEntry, error) {
			return []github.TreeEntry{
				{Path: "main.go", Type: "blob", Size: 100},
				{Path: "README.md", Type: "blob", Size: 200},
				{Path: "vendor/lib.go", Type: "blob", Size: 50},       // excluded by pattern
				{Path: "src", Type: "tree"},                            // not a blob
				{Path: "huge.go", Type: "blob", Size: 2 * 1024 * 1024}, // too large
				{Path: "utils.go", Type: "blob", Size: 80},
			}, nil
		},
		OnGetFileContentBatch: func(ctx context.Context, token, owner, repo, ref string, paths []string) ([]*github.FileContent, error) {
			result := make([]*github.FileContent, len(paths))
			for i, p := range paths {
				result[i] = &github.FileContent{Path: p, Content: []byte("package main")}
			}
			return result, nil
		},
	}
	mc := &vickytest.MockIngestionConfigs{}
	mv := &vickytest.MockVersions{
		OnGet: func(ctx context.Context, key string) (*models.Version, error) {
			return version, nil
		},
	}

	var putMu sync.Mutex
	var putPaths []string
	mb := &vickytest.MockBlobs{
		OnPutBlob: func(ctx context.Context, userID int64, blob *models.Blob) error {
			putMu.Lock()
			defer putMu.Unlock()
			putPaths = append(putPaths, blob.Path)
			return nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithUsers(mu),
		vickytest.WithGitHub(mg),
		vickytest.WithIngestionConfigs(mc),
		vickytest.WithVersions(mv),
		vickytest.WithBlobs(mb),
	)

	job := vickytest.NewJob(t)
	result, err := fetchStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Stage != models.JobStageFetch {
		t.Errorf("Stage = %q, want %q", result.Stage, models.JobStageFetch)
	}

	// main.go, README.md (docs included), utils.go should pass; vendor excluded, tree skipped, huge excluded
	if result.ItemsTotal != 3 {
		t.Errorf("ItemsTotal = %d, want 3", result.ItemsTotal)
	}
	if result.ItemsProcessed != 3 {
		t.Errorf("ItemsProcessed = %d, want 3", result.ItemsProcessed)
	}

	// Verify stored paths
	putMu.Lock()
	defer putMu.Unlock()
	wantPaths := map[string]bool{"main.go": true, "README.md": true, "utils.go": true}
	for _, p := range putPaths {
		if !wantPaths[p] {
			t.Errorf("unexpected blob stored: %s", p)
		}
		delete(wantPaths, p)
	}
	for p := range wantPaths {
		t.Errorf("expected blob not stored: %s", p)
	}
}

func TestFetchStage_EmptyTree(t *testing.T) {
	version := vickytest.NewVersion(t)

	mg := &vickytest.MockGitHub{
		OnGetTree: func(ctx context.Context, token, owner, repo, ref string) ([]github.TreeEntry, error) {
			return nil, nil
		},
	}
	mv := &vickytest.MockVersions{
		OnGet: func(ctx context.Context, key string) (*models.Version, error) {
			return version, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithUsers(&vickytest.MockUsers{}),
		vickytest.WithGitHub(mg),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
		vickytest.WithVersions(mv),
		vickytest.WithBlobs(&vickytest.MockBlobs{}),
	)

	job := vickytest.NewJob(t)
	result, err := fetchStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ItemsTotal != 0 {
		t.Errorf("ItemsTotal = %d, want 0", result.ItemsTotal)
	}
}

func TestFetchStage_GetUserError(t *testing.T) {
	mu := &vickytest.MockUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			return nil, fmt.Errorf("user not found")
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithUsers(mu),
		vickytest.WithGitHub(&vickytest.MockGitHub{}),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
		vickytest.WithVersions(&vickytest.MockVersions{}),
		vickytest.WithBlobs(&vickytest.MockBlobs{}),
	)

	job := vickytest.NewJob(t)
	_, err := fetchStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "user not found") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "user not found")
	}
}

func TestFetchStage_BlobStorageError(t *testing.T) {
	version := vickytest.NewVersion(t)

	mg := &vickytest.MockGitHub{
		OnGetTree: func(ctx context.Context, token, owner, repo, ref string) ([]github.TreeEntry, error) {
			return []github.TreeEntry{
				{Path: "main.go", Type: "blob", Size: 100},
			}, nil
		},
		OnGetFileContentBatch: func(ctx context.Context, token, owner, repo, ref string, paths []string) ([]*github.FileContent, error) {
			return []*github.FileContent{{Path: "main.go", Content: []byte("package main")}}, nil
		},
	}
	mv := &vickytest.MockVersions{
		OnGet: func(ctx context.Context, key string) (*models.Version, error) {
			return version, nil
		},
	}
	mb := &vickytest.MockBlobs{
		OnPutBlob: func(ctx context.Context, userID int64, blob *models.Blob) error {
			return fmt.Errorf("storage unavailable")
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithUsers(&vickytest.MockUsers{}),
		vickytest.WithGitHub(mg),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
		vickytest.WithVersions(mv),
		vickytest.WithBlobs(mb),
	)

	job := vickytest.NewJob(t)
	_, err := fetchStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "store blob") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "store blob")
	}
}
