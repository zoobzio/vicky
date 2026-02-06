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

func TestSearchChunks(t *testing.T) {
	chunks := vickytest.NewChunks(t, 2)
	mc := &vickytest.MockChunks{
		OnSearch: func(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Chunk, error) {
			return chunks, nil
		},
	}
	me := &vickytest.MockEmbedder{}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithChunks(mc), vickytest.WithEmbedder(me))
	engine.WithHandlers(SearchChunks)

	capture := rtesting.ServeRequest(engine, "GET", "/search/testorg/testrepo/v1.0.0?q=hello", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.SearchResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Query != "hello" {
		t.Errorf("Query = %q, want %q", resp.Query, "hello")
	}
	if len(resp.Results) != 2 {
		t.Errorf("len(Results) = %d, want 2", len(resp.Results))
	}
}

func TestSearchChunks_MissingQuery(t *testing.T) {
	mc := &vickytest.MockChunks{}
	me := &vickytest.MockEmbedder{}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithChunks(mc), vickytest.WithEmbedder(me))
	engine.WithHandlers(SearchChunks)

	capture := rtesting.ServeRequest(engine, "GET", "/search/testorg/testrepo/v1.0.0", nil)
	rtesting.AssertStatus(t, capture, 400)
}

func TestSearchSymbols(t *testing.T) {
	sym := &models.Symbol{
		ID:        1,
		UserID:    1000,
		Owner:     "testorg",
		RepoName:  "testrepo",
		Tag:       "v1.0.0",
		Name:      "Connect",
		Kind:      models.SymbolKindFunction,
		FilePath:  "client.go",
		StartLine: 10,
		Exported:  true,
	}
	ms := &vickytest.MockSymbols{
		OnFindRelated: func(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error) {
			return []*models.Symbol{sym}, nil
		},
	}
	me := &vickytest.MockEmbedder{}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithSymbols(ms), vickytest.WithEmbedder(me))
	engine.WithHandlers(SearchSymbols)

	capture := rtesting.ServeRequest(engine, "GET", "/search/testorg/testrepo/v1.0.0/symbols?q=connect", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.SymbolSearchResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Query != "connect" {
		t.Errorf("Query = %q, want %q", resp.Query, "connect")
	}
	if len(resp.Results) != 1 {
		t.Errorf("len(Results) = %d, want 1", len(resp.Results))
	}
}

func TestSearchSymbols_Exported(t *testing.T) {
	sym := &models.Symbol{
		ID:        1,
		UserID:    1000,
		Owner:     "testorg",
		RepoName:  "testrepo",
		Tag:       "v1.0.0",
		Name:      "Connect",
		Kind:      models.SymbolKindFunction,
		FilePath:  "client.go",
		StartLine: 10,
		Exported:  true,
	}
	var calledExported bool
	ms := &vickytest.MockSymbols{
		OnFindRelatedExported: func(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error) {
			calledExported = true
			return []*models.Symbol{sym}, nil
		},
	}
	me := &vickytest.MockEmbedder{}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithSymbols(ms), vickytest.WithEmbedder(me))
	engine.WithHandlers(SearchSymbols)

	capture := rtesting.ServeRequest(engine, "GET", "/search/testorg/testrepo/v1.0.0/symbols?q=connect&exported=true", nil)
	rtesting.AssertStatus(t, capture, 200)

	if !calledExported {
		t.Error("expected FindRelatedExported to be called")
	}
}

func TestFindSimilarDocuments(t *testing.T) {
	sourceDoc := &models.Document{
		ID:        1,
		VersionID: 10,
		UserID:    1000,
		Owner:     "testorg",
		RepoName:  "testrepo",
		Tag:       "v1.0.0",
		Path:      "main.go",
		Vector:    []float32{0.1, 0.2, 0.3},
	}
	similarDoc := &models.Document{
		ID:        2,
		VersionID: 10,
		UserID:    1000,
		Owner:     "testorg",
		RepoName:  "testrepo",
		Tag:       "v1.0.0",
		Path:      "client.go",
	}

	md := &vickytest.MockDocuments{
		OnGetByUserRepoTagAndPath: func(ctx context.Context, userID int64, owner, repoName, tag, path string) (*models.Document, error) {
			return sourceDoc, nil
		},
		OnFindSimilarInVersion: func(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Document, error) {
			return []*models.Document{sourceDoc, similarDoc}, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithDocuments(md))
	engine.WithHandlers(FindSimilarDocuments)

	capture := rtesting.ServeRequest(engine, "GET", "/search/testorg/testrepo/v1.0.0/similar?path=main.go", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.SimilarDocumentsResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// sourceDoc should be excluded by DocumentsToSimilarResponse
	if len(resp.Results) != 1 {
		t.Errorf("len(Results) = %d, want 1 (source excluded)", len(resp.Results))
	}
}
