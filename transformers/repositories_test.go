package transformers

import (
	"testing"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/wire"
)

func TestRepositoryToResponse(t *testing.T) {
	desc := "A test repo"
	r := &models.Repository{
		ID:            1,
		GitHubID:      999,
		Owner:         "testorg",
		Name:          "testrepo",
		FullName:      "testorg/testrepo",
		Description:   &desc,
		DefaultBranch: "main",
		Private:       true,
		HTMLURL:       "https://github.com/testorg/testrepo",
	}

	resp := RepositoryToResponse(r)

	if resp.ID != 1 {
		t.Errorf("ID = %d, want 1", resp.ID)
	}
	if resp.GitHubID != 999 {
		t.Errorf("GitHubID = %d, want 999", resp.GitHubID)
	}
	if resp.Owner != "testorg" {
		t.Errorf("Owner = %q, want %q", resp.Owner, "testorg")
	}
	if resp.Description == nil || *resp.Description != "A test repo" {
		t.Errorf("Description = %v, want %q", resp.Description, "A test repo")
	}
	if !resp.Private {
		t.Error("Private = false, want true")
	}
}

func TestRepositoriesToList(t *testing.T) {
	repos := []*models.Repository{
		{ID: 1, Owner: "a"},
		{ID: 2, Owner: "b"},
	}

	resp := RepositoriesToList(repos)

	if len(resp.Repositories) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Repositories))
	}
	if resp.Repositories[0].ID != 1 {
		t.Errorf("first ID = %d, want 1", resp.Repositories[0].ID)
	}
	if resp.Repositories[1].ID != 2 {
		t.Errorf("second ID = %d, want 2", resp.Repositories[1].ID)
	}
}

func TestApplyRepositoryRegistration(t *testing.T) {
	desc := "New repo"
	req := wire.RegisterRepositoryRequest{
		GitHubID:      123,
		Owner:         "org",
		Name:          "repo",
		FullName:      "org/repo",
		Description:   &desc,
		DefaultBranch: "develop",
		Private:       true,
		HTMLURL:       "https://github.com/org/repo",
	}

	r := &models.Repository{}
	ApplyRepositoryRegistration(req, r)

	if r.GitHubID != 123 {
		t.Errorf("GitHubID = %d, want 123", r.GitHubID)
	}
	if r.Owner != "org" {
		t.Errorf("Owner = %q, want %q", r.Owner, "org")
	}
	if r.Description == nil || *r.Description != "New repo" {
		t.Errorf("Description = %v, want %q", r.Description, "New repo")
	}
	if r.DefaultBranch != "develop" {
		t.Errorf("DefaultBranch = %q, want %q", r.DefaultBranch, "develop")
	}
}
