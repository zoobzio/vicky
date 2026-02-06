//go:build testing

package testing

import (
	"context"
	"testing"

	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/contracts"
	"github.com/zoobzio/vicky/models"
)

func TestNewJob(t *testing.T) {
	job := NewJob(t)

	if job.ID != 1 {
		t.Errorf("ID = %d, want 1", job.ID)
	}
	if job.VersionID != 10 {
		t.Errorf("VersionID = %d, want 10", job.VersionID)
	}
	if job.RepositoryID != 100 {
		t.Errorf("RepositoryID = %d, want 100", job.RepositoryID)
	}
	if job.UserID != 1000 {
		t.Errorf("UserID = %d, want 1000", job.UserID)
	}
	if job.Owner != "testorg" {
		t.Errorf("Owner = %q, want %q", job.Owner, "testorg")
	}
	if job.RepoName != "testrepo" {
		t.Errorf("RepoName = %q, want %q", job.RepoName, "testrepo")
	}
	if job.Tag != "v1.0.0" {
		t.Errorf("Tag = %q, want %q", job.Tag, "v1.0.0")
	}
	if job.Stage != models.JobStageFetch {
		t.Errorf("Stage = %q, want %q", job.Stage, models.JobStageFetch)
	}
	if job.Status != models.JobStatusRunning {
		t.Errorf("Status = %q, want %q", job.Status, models.JobStatusRunning)
	}
}

func TestNewChunk(t *testing.T) {
	chunk := NewChunk(t, 42, "func main() {}")

	if chunk.ID != 42 {
		t.Errorf("ID = %d, want 42", chunk.ID)
	}
	if chunk.Content != "func main() {}" {
		t.Errorf("Content = %q, want %q", chunk.Content, "func main() {}")
	}
	if chunk.UserID != 1000 {
		t.Errorf("UserID = %d, want 1000", chunk.UserID)
	}
	if chunk.Owner != "testorg" {
		t.Errorf("Owner = %q, want %q", chunk.Owner, "testorg")
	}
	if chunk.Kind != models.ChunkKindFunction {
		t.Errorf("Kind = %q, want %q", chunk.Kind, models.ChunkKindFunction)
	}
}

func TestNewChunks(t *testing.T) {
	chunks := NewChunks(t, 5)

	if len(chunks) != 5 {
		t.Fatalf("len = %d, want 5", len(chunks))
	}
	for i, c := range chunks {
		wantID := int64(i + 1)
		if c.ID != wantID {
			t.Errorf("chunks[%d].ID = %d, want %d", i, c.ID, wantID)
		}
		if c.Content == "" {
			t.Errorf("chunks[%d].Content is empty", i)
		}
	}
}

func TestNewVersion(t *testing.T) {
	v := NewVersion(t)

	if v.ID != 10 {
		t.Errorf("ID = %d, want 10", v.ID)
	}
	if v.Status != models.VersionStatusIngesting {
		t.Errorf("Status = %q, want %q", v.Status, models.VersionStatusIngesting)
	}
	if v.CommitSHA != "abc123" {
		t.Errorf("CommitSHA = %q, want %q", v.CommitSHA, "abc123")
	}
}

func TestSetupRegistry(t *testing.T) {
	mv := &MockVersions{}
	me := &MockEmbedder{}
	mc := &MockChunks{}

	ctx := SetupRegistry(t,
		WithVersions(mv),
		WithEmbedder(me),
		WithChunks(mc),
	)

	if ctx == nil {
		t.Fatal("context is nil")
	}

	// Verify mocks are retrievable
	versions := sum.MustUse[contracts.Versions](ctx)
	if versions == nil {
		t.Error("Versions not registered")
	}

	embedder := sum.MustUse[contracts.Embedder](ctx)
	if embedder == nil {
		t.Error("Embedder not registered")
	}

	chunks := sum.MustUse[contracts.Chunks](ctx)
	if chunks == nil {
		t.Error("Chunks not registered")
	}
}

func TestSetupRegistry_Empty(t *testing.T) {
	ctx := SetupRegistry(t)

	if ctx == nil {
		t.Fatal("context is nil")
	}
}

func TestMockVersions_Default(t *testing.T) {
	mv := &MockVersions{}
	ctx := context.Background()

	v, err := mv.UpdateStatus(ctx, 10, models.VersionStatusReady, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.ID != 10 {
		t.Errorf("ID = %d, want 10", v.ID)
	}
	if v.Status != models.VersionStatusReady {
		t.Errorf("Status = %q, want %q", v.Status, models.VersionStatusReady)
	}
}

func TestMockEmbedder_Default(t *testing.T) {
	me := &MockEmbedder{}
	ctx := context.Background()

	vectors, err := me.Embed(ctx, []string{"hello", "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vectors) != 2 {
		t.Fatalf("len = %d, want 2", len(vectors))
	}
	for i, v := range vectors {
		if len(v) != 3 {
			t.Errorf("vectors[%d] len = %d, want 3", i, len(v))
		}
	}

	if me.Dimensions() != 3 {
		t.Errorf("Dimensions = %d, want 3", me.Dimensions())
	}
}

func TestMockChunks_Default(t *testing.T) {
	mc := &MockChunks{}
	ctx := context.Background()

	chunks, err := mc.ListByUserRepoAndTag(ctx, 1, "org", "repo", "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chunks != nil {
		t.Errorf("expected nil, got %v", chunks)
	}

	if err := mc.Set(ctx, "1", &models.Chunk{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
