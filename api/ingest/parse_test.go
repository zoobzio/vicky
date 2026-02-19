//go:build testing

package ingest

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	scipproto "github.com/sourcegraph/scip/bindings/go/scip"
	"github.com/zoobzio/vicky/external/indexer"
	"github.com/zoobzio/vicky/models"
	vickytest "github.com/zoobzio/vicky/testing"
	"google.golang.org/protobuf/proto"
)

// buildSCIPIndex constructs a minimal SCIP protobuf index for testing.
func buildSCIPIndex(t *testing.T, docs []*scipproto.Document) []byte {
	t.Helper()
	index := &scipproto.Index{
		Metadata: &scipproto.Metadata{
			Version:              scipproto.ProtocolVersion_UnspecifiedProtocolVersion,
			ToolInfo:             &scipproto.ToolInfo{Name: "test-indexer", Version: "1.0"},
			ProjectRoot:          "file:///test",
			TextDocumentEncoding: scipproto.TextEncoding_UTF8,
		},
		Documents: docs,
	}
	data, err := proto.Marshal(index)
	if err != nil {
		t.Fatalf("marshal scip index: %v", err)
	}
	return data
}

func TestParseStage(t *testing.T) {
	version := vickytest.NewVersion(t)

	scipDocs := []*scipproto.Document{
		{
			RelativePath: "main.go",
			Symbols: []*scipproto.SymbolInformation{
				{
					Symbol:      "scip-go gomod example.com/test v1.0.0 main().",
					Kind:        scipproto.SymbolInformation_Function,
					DisplayName: "main",
				},
			},
			Occurrences: []*scipproto.Occurrence{
				{
					Range:       []int32{0, 5, 9},
					Symbol:      "scip-go gomod example.com/test v1.0.0 main().",
					SymbolRoles: int32(scipproto.SymbolRole_Definition),
				},
			},
		},
	}
	indexData := buildSCIPIndex(t, scipDocs)

	mi := &vickytest.MockIndexer{
		OnIndex: func(ctx context.Context, req indexer.Request) (*indexer.Result, error) {
			return &indexer.Result{
				JobID:     req.JobID,
				VersionID: req.VersionID,
				IndexData: indexData,
			}, nil
		},
	}
	mc := &vickytest.MockIngestionConfigs{}
	mv := &vickytest.MockVersions{
		OnGet: func(ctx context.Context, key string) (*models.Version, error) {
			return version, nil
		},
	}

	var docSetCalls int
	md := &vickytest.MockDocuments{
		OnSet: func(ctx context.Context, key string, doc *models.Document) error {
			docSetCalls++
			doc.ID = int64(docSetCalls)
			return nil
		},
	}

	var mu sync.Mutex
	var symbolCalls int
	var occurrenceCalls int

	ms := &vickytest.MockSCIPSymbols{
		OnSet: func(ctx context.Context, key string, symbol *models.SCIPSymbol) error {
			mu.Lock()
			defer mu.Unlock()
			symbolCalls++
			symbol.ID = int64(symbolCalls)
			return nil
		},
	}
	mo := &vickytest.MockSCIPOccurrences{
		OnSet: func(ctx context.Context, key string, occ *models.SCIPOccurrence) error {
			mu.Lock()
			defer mu.Unlock()
			occurrenceCalls++
			return nil
		},
	}
	mr := &vickytest.MockSCIPRelationships{}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithIndexer(mi),
		vickytest.WithIngestionConfigs(mc),
		vickytest.WithVersions(mv),
		vickytest.WithDocuments(md),
		vickytest.WithSCIPSymbols(ms),
		vickytest.WithSCIPOccurrences(mo),
		vickytest.WithSCIPRelationships(mr),
	)

	job := vickytest.NewJob(t)
	result, err := parseStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Stage != models.JobStageParse {
		t.Errorf("Stage = %q, want %q", result.Stage, models.JobStageParse)
	}
	if docSetCalls != 1 {
		t.Errorf("document Set calls = %d, want 1", docSetCalls)
	}

	mu.Lock()
	defer mu.Unlock()
	if symbolCalls != 1 {
		t.Errorf("symbol Set calls = %d, want 1", symbolCalls)
	}
	if occurrenceCalls != 1 {
		t.Errorf("occurrence Set calls = %d, want 1", occurrenceCalls)
	}
}

func TestParseStage_UnsupportedLanguage(t *testing.T) {
	mi := &vickytest.MockIndexer{
		OnSupports: func(language models.Language) bool {
			return false
		},
	}
	mc := &vickytest.MockIngestionConfigs{}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithIndexer(mi),
		vickytest.WithIngestionConfigs(mc),
		vickytest.WithVersions(&vickytest.MockVersions{}),
		vickytest.WithDocuments(&vickytest.MockDocuments{}),
		vickytest.WithSCIPSymbols(&vickytest.MockSCIPSymbols{}),
		vickytest.WithSCIPOccurrences(&vickytest.MockSCIPOccurrences{}),
		vickytest.WithSCIPRelationships(&vickytest.MockSCIPRelationships{}),
	)

	job := vickytest.NewJob(t)
	result, err := parseStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Stage set but no processing
	if result.Stage != models.JobStageParse {
		t.Errorf("Stage = %q, want %q", result.Stage, models.JobStageParse)
	}
}

func TestParseStage_IndexerError(t *testing.T) {
	version := vickytest.NewVersion(t)

	mi := &vickytest.MockIndexer{
		OnIndex: func(ctx context.Context, req indexer.Request) (*indexer.Result, error) {
			return nil, fmt.Errorf("indexer crashed")
		},
	}
	mv := &vickytest.MockVersions{
		OnGet: func(ctx context.Context, key string) (*models.Version, error) {
			return version, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithIndexer(mi),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
		vickytest.WithVersions(mv),
		vickytest.WithDocuments(&vickytest.MockDocuments{}),
		vickytest.WithSCIPSymbols(&vickytest.MockSCIPSymbols{}),
		vickytest.WithSCIPOccurrences(&vickytest.MockSCIPOccurrences{}),
		vickytest.WithSCIPRelationships(&vickytest.MockSCIPRelationships{}),
	)

	job := vickytest.NewJob(t)
	_, err := parseStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "indexer crashed") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "indexer crashed")
	}
}

func TestParseStage_IndexerErrorResult(t *testing.T) {
	version := vickytest.NewVersion(t)

	mi := &vickytest.MockIndexer{
		OnIndex: func(ctx context.Context, req indexer.Request) (*indexer.Result, error) {
			return &indexer.Result{
				JobID:     req.JobID,
				VersionID: req.VersionID,
				Error:     "compilation failed",
			}, nil
		},
	}
	mv := &vickytest.MockVersions{
		OnGet: func(ctx context.Context, key string) (*models.Version, error) {
			return version, nil
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithIndexer(mi),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
		vickytest.WithVersions(mv),
		vickytest.WithDocuments(&vickytest.MockDocuments{}),
		vickytest.WithSCIPSymbols(&vickytest.MockSCIPSymbols{}),
		vickytest.WithSCIPOccurrences(&vickytest.MockSCIPOccurrences{}),
		vickytest.WithSCIPRelationships(&vickytest.MockSCIPRelationships{}),
	)

	job := vickytest.NewJob(t)
	_, err := parseStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "indexer error") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "indexer error")
	}
}

func TestParseStage_DocumentSetError(t *testing.T) {
	version := vickytest.NewVersion(t)

	scipDocs := []*scipproto.Document{
		{RelativePath: "main.go"},
	}
	indexData := buildSCIPIndex(t, scipDocs)

	mi := &vickytest.MockIndexer{
		OnIndex: func(ctx context.Context, req indexer.Request) (*indexer.Result, error) {
			return &indexer.Result{
				JobID:     req.JobID,
				VersionID: req.VersionID,
				IndexData: indexData,
			}, nil
		},
	}
	mv := &vickytest.MockVersions{
		OnGet: func(ctx context.Context, key string) (*models.Version, error) {
			return version, nil
		},
	}
	md := &vickytest.MockDocuments{
		OnSet: func(ctx context.Context, key string, doc *models.Document) error {
			return fmt.Errorf("db constraint violation")
		},
	}

	ctx := vickytest.SetupRegistry(t,
		vickytest.WithIndexer(mi),
		vickytest.WithIngestionConfigs(&vickytest.MockIngestionConfigs{}),
		vickytest.WithVersions(mv),
		vickytest.WithDocuments(md),
		vickytest.WithSCIPSymbols(&vickytest.MockSCIPSymbols{}),
		vickytest.WithSCIPOccurrences(&vickytest.MockSCIPOccurrences{}),
		vickytest.WithSCIPRelationships(&vickytest.MockSCIPRelationships{}),
	)

	job := vickytest.NewJob(t)
	_, err := parseStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "create document") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "create document")
	}
}
