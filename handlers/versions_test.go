//go:build testing

package handlers

import (
	"context"
	"testing"

	rtesting "github.com/zoobzio/rocco/testing"
	vickytest "github.com/zoobzio/vicky/testing"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/wire"
)

func TestListVersions(t *testing.T) {
	versions := []*models.Version{vickytest.NewVersion(t)}
	mv := &vickytest.MockVersions{
		OnListByUserAndRepo: func(ctx context.Context, userID int64, owner, repoName string) ([]*models.Version, error) {
			return versions, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithVersions(mv))
	engine.WithHandlers(ListVersions)

	capture := rtesting.ServeRequest(engine, "GET", "/repositories/testorg/testrepo/versions", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.VersionListResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Versions) != 1 {
		t.Errorf("len = %d, want 1", len(resp.Versions))
	}
}

func TestGetVersion(t *testing.T) {
	version := vickytest.NewVersion(t)
	mv := &vickytest.MockVersions{
		OnGetByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) (*models.Version, error) {
			return version, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithVersions(mv))
	engine.WithHandlers(GetVersion)

	capture := rtesting.ServeRequest(engine, "GET", "/repositories/testorg/testrepo/versions/v1.0.0", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.VersionResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Tag != "v1.0.0" {
		t.Errorf("Tag = %q, want %q", resp.Tag, "v1.0.0")
	}
}

func TestTriggerIngest(t *testing.T) {
	repo := vickytest.NewRepository(t)
	mr := &vickytest.MockRepositories{
		OnGetByUserOwnerAndName: func(ctx context.Context, userID int64, owner, name string) (*models.Repository, error) {
			if owner == repo.Owner && name == repo.Name {
				return repo, nil
			}
			return nil, ErrRepositoryNotFound
		},
	}
	mv := &vickytest.MockVersions{}
	mj := &vickytest.MockJobs{}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithRepositories(mr),
		vickytest.WithVersions(mv),
		vickytest.WithJobs(mj),
	)
	engine.WithHandlers(TriggerIngest)

	body := wire.IngestRequest{CommitSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	capture := rtesting.ServeRequest(engine, "POST", "/repositories/testorg/testrepo/versions/v2.0.0", body)
	rtesting.AssertStatus(t, capture, 202)

	var resp wire.VersionResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Tag != "v2.0.0" {
		t.Errorf("Tag = %q, want %q", resp.Tag, "v2.0.0")
	}
	if resp.Status != models.VersionStatusPending {
		t.Errorf("Status = %q, want %q", resp.Status, models.VersionStatusPending)
	}
}

func TestTriggerIngest_RepoNotFound(t *testing.T) {
	mr := &vickytest.MockRepositories{
		OnGetByUserOwnerAndName: func(ctx context.Context, userID int64, owner, name string) (*models.Repository, error) {
			return nil, ErrRepositoryNotFound
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithRepositories(mr),
		vickytest.WithVersions(&vickytest.MockVersions{}),
		vickytest.WithJobs(&vickytest.MockJobs{}),
	)
	engine.WithHandlers(TriggerIngest)

	body := wire.IngestRequest{CommitSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	capture := rtesting.ServeRequest(engine, "POST", "/repositories/testorg/nonexistent/versions/v1.0.0", body)
	rtesting.AssertStatus(t, capture, 404)
}

