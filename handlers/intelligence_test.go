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

func TestGetDefinition(t *testing.T) {
	occ := &models.SCIPOccurrence{
		ID:          1,
		DocumentID:  100,
		UserID:      1000,
		Owner:       "testorg",
		RepoName:    "testrepo",
		Tag:         "v1.0.0",
		Symbol:      "scip-go gomod example.com/pkg v1.0.0 Connect().",
		SymbolRoles: models.SCIPSymbolRoleDefinition,
		StartLine:   10,
		StartCol:    5,
		EndLine:     10,
		EndCol:      12,
	}
	doc := &models.Document{
		ID:   100,
		Path: "client.go",
	}

	mo := &vickytest.MockSCIPOccurrences{
		OnListDefinitions: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
			return []*models.SCIPOccurrence{occ}, nil
		},
	}
	md := &vickytest.MockDocuments{
		OnGet: func(ctx context.Context, key string) (*models.Document, error) {
			return doc, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPOccurrences(mo),
		vickytest.WithDocuments(md),
	)
	engine.WithHandlers(GetDefinition)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/definition?symbol=Connect", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.DefinitionResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Symbol != "Connect" {
		t.Errorf("Symbol = %q, want %q", resp.Symbol, "Connect")
	}
	if len(resp.Locations) != 1 {
		t.Errorf("len(Locations) = %d, want 1", len(resp.Locations))
	}
	if resp.Locations[0].Path != "client.go" {
		t.Errorf("Path = %q, want %q", resp.Locations[0].Path, "client.go")
	}
}

func TestGetDefinition_MissingSymbol(t *testing.T) {
	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPOccurrences(&vickytest.MockSCIPOccurrences{}),
		vickytest.WithDocuments(&vickytest.MockDocuments{}),
	)
	engine.WithHandlers(GetDefinition)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/definition", nil)
	rtesting.AssertStatus(t, capture, 400)
}

func TestFindReferences(t *testing.T) {
	refs := []*models.SCIPOccurrence{
		{
			ID:          1,
			DocumentID:  100,
			Symbol:      "Connect",
			SymbolRoles: models.SCIPSymbolRoleReadAccess,
			StartLine:   20,
			StartCol:    10,
			EndLine:     20,
			EndCol:      17,
		},
		{
			ID:          2,
			DocumentID:  101,
			Symbol:      "Connect",
			SymbolRoles: models.SCIPSymbolRoleReadAccess,
			StartLine:   30,
			StartCol:    5,
			EndLine:     30,
			EndCol:      12,
		},
	}
	docCache := map[string]*models.Document{
		"100": {ID: 100, Path: "main.go"},
		"101": {ID: 101, Path: "handler.go"},
	}

	mo := &vickytest.MockSCIPOccurrences{
		OnListReferences: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
			return refs, nil
		},
	}
	md := &vickytest.MockDocuments{
		OnGet: func(ctx context.Context, key string) (*models.Document, error) {
			return docCache[key], nil
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPOccurrences(mo),
		vickytest.WithDocuments(md),
	)
	engine.WithHandlers(FindReferences)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/references?symbol=Connect", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.ReferencesResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.References) != 2 {
		t.Errorf("len(References) = %d, want 2", len(resp.References))
	}
	if resp.Total != 2 {
		t.Errorf("Total = %d, want 2", resp.Total)
	}
}

func TestFindReferences_IncludeDefinition(t *testing.T) {
	var calledListBySymbol bool
	mo := &vickytest.MockSCIPOccurrences{
		OnListBySymbol: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
			calledListBySymbol = true
			return nil, nil
		},
	}
	md := &vickytest.MockDocuments{}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPOccurrences(mo),
		vickytest.WithDocuments(md),
	)
	engine.WithHandlers(FindReferences)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/references?symbol=Connect&include_definition=true", nil)
	rtesting.AssertStatus(t, capture, 200)

	if !calledListBySymbol {
		t.Error("expected ListBySymbol to be called with include_definition=true")
	}
}

func TestFindImplementations(t *testing.T) {
	sym := &models.SCIPSymbol{
		ID:       1,
		UserID:   1000,
		Owner:    "testorg",
		RepoName: "testrepo",
		Tag:      "v1.0.0",
		Symbol:   "scip-go gomod example.com/pkg v1.0.0 Reader#",
		Kind:     models.SCIPSymbolKindInterface,
	}
	rels := []*models.SCIPRelationship{
		{
			ID:               1,
			SCIPSymbolID:     1,
			TargetSymbol:     "scip-go gomod example.com/pkg v1.0.0 FileReader#",
			IsImplementation: true,
		},
	}
	implSym := &models.SCIPSymbol{
		ID:       2,
		Symbol:   "scip-go gomod example.com/pkg v1.0.0 FileReader#",
		Kind:     models.SCIPSymbolKindStruct,
	}
	implDisplayName := "FileReader"
	implSym.DisplayName = &implDisplayName

	implDef := &models.SCIPOccurrence{
		ID:         10,
		DocumentID: 100,
		Symbol:     implSym.Symbol,
		StartLine:  50,
		StartCol:   5,
		EndLine:    50,
		EndCol:     15,
	}
	doc := &models.Document{ID: 100, Path: "reader.go"}

	mss := &vickytest.MockSCIPSymbols{
		OnGetBySymbol: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) (*models.SCIPSymbol, error) {
			if symbol == sym.Symbol || symbol == "Reader" {
				return sym, nil
			}
			if symbol == implSym.Symbol {
				return implSym, nil
			}
			return nil, nil
		},
	}
	msr := &vickytest.MockSCIPRelationships{
		OnListImplementations: func(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error) {
			return rels, nil
		},
	}
	mso := &vickytest.MockSCIPOccurrences{
		OnListDefinitions: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
			return []*models.SCIPOccurrence{implDef}, nil
		},
	}
	md := &vickytest.MockDocuments{
		OnGet: func(ctx context.Context, key string) (*models.Document, error) {
			return doc, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPSymbols(mss),
		vickytest.WithSCIPRelationships(msr),
		vickytest.WithSCIPOccurrences(mso),
		vickytest.WithDocuments(md),
	)
	engine.WithHandlers(FindImplementations)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/implementations?symbol=Reader", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.ImplementationsResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Symbol != "Reader" {
		t.Errorf("Symbol = %q, want %q", resp.Symbol, "Reader")
	}
	if len(resp.Implementations) != 1 {
		t.Errorf("len(Implementations) = %d, want 1", len(resp.Implementations))
	}
}

func TestFindImplementations_SymbolNotFound(t *testing.T) {
	mss := &vickytest.MockSCIPSymbols{
		OnGetBySymbol: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) (*models.SCIPSymbol, error) {
			return nil, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPSymbols(mss),
		vickytest.WithSCIPRelationships(&vickytest.MockSCIPRelationships{}),
		vickytest.WithSCIPOccurrences(&vickytest.MockSCIPOccurrences{}),
		vickytest.WithDocuments(&vickytest.MockDocuments{}),
	)
	engine.WithHandlers(FindImplementations)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/implementations?symbol=NonExistent", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.ImplementationsResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Implementations) != 0 {
		t.Errorf("expected empty implementations for unknown symbol")
	}
}

func TestListSymbols(t *testing.T) {
	displayName := "Connect"
	symbols := []*models.SCIPSymbol{
		{
			ID:          1,
			DocumentID:  100,
			Symbol:      "scip-go gomod example.com/pkg v1.0.0 Connect().",
			Kind:        models.SCIPSymbolKindFunction,
			DisplayName: &displayName,
		},
	}
	occ := &models.SCIPOccurrence{
		ID:        10,
		DocumentID: 100,
		StartLine: 15,
		StartCol:  5,
		EndLine:   15,
		EndCol:    12,
	}
	doc := &models.Document{ID: 100, Path: "client.go"}

	mss := &vickytest.MockSCIPSymbols{
		OnListByUserRepoAndTag: func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPSymbol, error) {
			return symbols, nil
		},
	}
	mso := &vickytest.MockSCIPOccurrences{
		OnListDefinitions: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
			return []*models.SCIPOccurrence{occ}, nil
		},
	}
	md := &vickytest.MockDocuments{
		OnGet: func(ctx context.Context, key string) (*models.Document, error) {
			return doc, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPSymbols(mss),
		vickytest.WithSCIPOccurrences(mso),
		vickytest.WithDocuments(md),
	)
	engine.WithHandlers(ListSymbols)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/symbols", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.SymbolListResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Symbols) != 1 {
		t.Errorf("len(Symbols) = %d, want 1", len(resp.Symbols))
	}
	if resp.Total != 1 {
		t.Errorf("Total = %d, want 1", resp.Total)
	}
}

func TestListSymbols_ByPath(t *testing.T) {
	doc := &models.Document{ID: 100, Path: "client.go"}
	displayName := "Connect"
	symbols := []*models.SCIPSymbol{
		{
			ID:          1,
			DocumentID:  100,
			Symbol:      "scip-go gomod example.com/pkg v1.0.0 Connect().",
			Kind:        models.SCIPSymbolKindFunction,
			DisplayName: &displayName,
		},
	}
	occ := &models.SCIPOccurrence{
		ID:        10,
		DocumentID: 100,
		StartLine: 15,
		StartCol:  5,
		EndLine:   15,
		EndCol:    12,
	}

	var calledListByDocument bool
	mss := &vickytest.MockSCIPSymbols{
		OnListByDocument: func(ctx context.Context, documentID int64) ([]*models.SCIPSymbol, error) {
			calledListByDocument = true
			return symbols, nil
		},
	}
	md := &vickytest.MockDocuments{
		OnGetByUserRepoTagAndPath: func(ctx context.Context, userID int64, owner, repoName, tag, path string) (*models.Document, error) {
			return doc, nil
		},
		OnGet: func(ctx context.Context, key string) (*models.Document, error) {
			return doc, nil
		},
	}
	mso := &vickytest.MockSCIPOccurrences{
		OnListDefinitions: func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
			return []*models.SCIPOccurrence{occ}, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t,
		vickytest.WithSCIPSymbols(mss),
		vickytest.WithDocuments(md),
		vickytest.WithSCIPOccurrences(mso),
	)
	engine.WithHandlers(ListSymbols)

	capture := rtesting.ServeRequest(engine, "GET", "/intel/testorg/testrepo/v1.0.0/symbols?path=client.go", nil)
	rtesting.AssertStatus(t, capture, 200)

	if !calledListByDocument {
		t.Error("expected ListByDocument to be called when path is provided")
	}
}
