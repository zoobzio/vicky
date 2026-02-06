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

func TestListRepositories(t *testing.T) {
	repos := []*models.Repository{vickytest.NewRepository(t)}
	mr := &vickytest.MockRepositories{
		OnListByUserID: func(ctx context.Context, userID int64) ([]*models.Repository, error) {
			return repos, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithRepositories(mr))
	engine.WithHandlers(ListRepositories)

	capture := rtesting.ServeRequest(engine, "GET", "/repositories", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.RepositoryListResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Repositories) != 1 {
		t.Errorf("len = %d, want 1", len(resp.Repositories))
	}
}

func TestRegisterRepository(t *testing.T) {
	mr := &vickytest.MockRepositories{
		OnSet: func(ctx context.Context, key string, repo *models.Repository) error {
			repo.ID = 1
			return nil
		},
	}
	mc := &vickytest.MockIngestionConfigs{
		OnSet: func(ctx context.Context, key string, config *models.IngestionConfig) error {
			return nil
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithRepositories(mr),
		vickytest.WithIngestionConfigs(mc),
	)
	engine.WithHandlers(RegisterRepository)

	body := wire.RegisterRepositoryRequest{
		GitHubID:      123,
		Owner:         "testorg",
		Name:          "testrepo",
		FullName:      "testorg/testrepo",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/testorg/testrepo",
		Config:        wire.IngestionConfigRequest{Language: "go"},
	}

	capture := rtesting.ServeRequest(engine, "POST", "/repositories", body)
	rtesting.AssertStatus(t, capture, 201)

	var resp wire.RepositoryResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Owner != "testorg" {
		t.Errorf("Owner = %q, want %q", resp.Owner, "testorg")
	}
}

func TestGetRepository(t *testing.T) {
	repo := vickytest.NewRepository(t)
	mr := &vickytest.MockRepositories{
		OnGetByUserOwnerAndName: func(ctx context.Context, userID int64, owner, name string) (*models.Repository, error) {
			if owner == repo.Owner && name == repo.Name {
				return repo, nil
			}
			return nil, ErrRepositoryNotFound
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithRepositories(mr))
	engine.WithHandlers(GetRepository)

	capture := rtesting.ServeRequest(engine, "GET", "/repositories/testorg/testrepo", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.RepositoryResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Name != "testrepo" {
		t.Errorf("Name = %q, want %q", resp.Name, "testrepo")
	}
}

func TestGetRepository_NotFound(t *testing.T) {
	mr := &vickytest.MockRepositories{
		OnGetByUserOwnerAndName: func(ctx context.Context, userID int64, owner, name string) (*models.Repository, error) {
			return nil, ErrRepositoryNotFound
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithRepositories(mr))
	engine.WithHandlers(GetRepository)

	capture := rtesting.ServeRequest(engine, "GET", "/repositories/testorg/nonexistent", nil)
	rtesting.AssertStatus(t, capture, 404)
}
