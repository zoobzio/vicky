package transformers

import (
	"testing"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/api/wire"
)

func TestVersionToResponse(t *testing.T) {
	errMsg := "oops"
	v := &models.Version{
		ID:           10,
		RepositoryID: 100,
		Owner:        "testorg",
		RepoName:     "testrepo",
		Tag:          "v1.0.0",
		CommitSHA:    "abc123",
		Status:       models.VersionStatusReady,
		Error:        &errMsg,
	}

	resp := VersionToResponse(v)

	if resp.ID != 10 {
		t.Errorf("ID = %d, want 10", resp.ID)
	}
	if resp.RepositoryID != 100 {
		t.Errorf("RepositoryID = %d, want 100", resp.RepositoryID)
	}
	if resp.Tag != "v1.0.0" {
		t.Errorf("Tag = %q, want %q", resp.Tag, "v1.0.0")
	}
	if resp.Status != models.VersionStatusReady {
		t.Errorf("Status = %q, want %q", resp.Status, models.VersionStatusReady)
	}
	if resp.Error == nil || *resp.Error != "oops" {
		t.Errorf("Error = %v, want %q", resp.Error, "oops")
	}
}

func TestVersionsToList(t *testing.T) {
	versions := []*models.Version{
		{ID: 1, Tag: "v1.0.0"},
		{ID: 2, Tag: "v2.0.0"},
	}

	resp := VersionsToList(versions)

	if len(resp.Versions) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Versions))
	}
	if resp.Versions[0].Tag != "v1.0.0" {
		t.Errorf("first tag = %q, want %q", resp.Versions[0].Tag, "v1.0.0")
	}
}

func TestApplyIngestRequest(t *testing.T) {
	req := wire.IngestRequest{CommitSHA: "abc123def456"}
	v := &models.Version{}

	ApplyIngestRequest(req, v)

	if v.CommitSHA != "abc123def456" {
		t.Errorf("CommitSHA = %q, want %q", v.CommitSHA, "abc123def456")
	}
}
