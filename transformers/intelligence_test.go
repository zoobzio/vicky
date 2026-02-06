package transformers

import (
	"testing"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/wire"
)

func TestOccurrenceToLocation(t *testing.T) {
	occ := &models.SCIPOccurrence{
		StartLine: 5,
		StartCol:  10,
		EndLine:   5,
		EndCol:    20,
	}

	loc := OccurrenceToLocation(occ, "main.go")

	if loc.Path != "main.go" {
		t.Errorf("Path = %q, want %q", loc.Path, "main.go")
	}
	if loc.StartLine != 5 || loc.StartCol != 10 || loc.EndLine != 5 || loc.EndCol != 20 {
		t.Errorf("range = (%d,%d,%d,%d), want (5,10,5,20)",
			loc.StartLine, loc.StartCol, loc.EndLine, loc.EndCol)
	}
}

func TestSCIPSymbolToInfo(t *testing.T) {
	displayName := "Connect"
	sym := &models.SCIPSymbol{
		Symbol:        "scip-go gomod example.com v1.0.0 Connect().",
		Kind:          models.SCIPSymbolKindFunction,
		DisplayName:   &displayName,
		Documentation: []string{"Connects to server."},
	}
	loc := &wire.Location{Path: "client.go", StartLine: 10}

	info := SCIPSymbolToInfo(sym, loc)

	if info.Symbol != sym.Symbol {
		t.Errorf("Symbol = %q, want %q", info.Symbol, sym.Symbol)
	}
	if info.DisplayName != "Connect" {
		t.Errorf("DisplayName = %q, want %q", info.DisplayName, "Connect")
	}
	if info.Location == nil || info.Location.Path != "client.go" {
		t.Error("Location not set correctly")
	}
}

func TestSCIPSymbolToInfo_NilDisplayName(t *testing.T) {
	sym := &models.SCIPSymbol{
		Symbol:      "test",
		DisplayName: nil,
	}

	info := SCIPSymbolToInfo(sym, nil)

	if info.DisplayName != "" {
		t.Errorf("DisplayName = %q, want empty string", info.DisplayName)
	}
}

func TestDefinitionsToResponse(t *testing.T) {
	occs := []*models.SCIPOccurrence{
		{DocumentID: 1, StartLine: 10, EndLine: 10, StartCol: 5, EndCol: 15},
		{DocumentID: 2, StartLine: 20, EndLine: 20, StartCol: 0, EndCol: 10},
	}
	pathResolver := func(docID int64) string {
		if docID == 1 {
			return "a.go"
		}
		return "b.go"
	}

	resp := DefinitionsToResponse("test.Symbol", occs, pathResolver)

	if resp.Symbol != "test.Symbol" {
		t.Errorf("Symbol = %q, want %q", resp.Symbol, "test.Symbol")
	}
	if len(resp.Locations) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Locations))
	}
	if resp.Locations[0].Path != "a.go" {
		t.Errorf("first path = %q, want %q", resp.Locations[0].Path, "a.go")
	}
	if resp.Locations[1].Path != "b.go" {
		t.Errorf("second path = %q, want %q", resp.Locations[1].Path, "b.go")
	}
}

func TestReferencesToResponse(t *testing.T) {
	occs := []*models.SCIPOccurrence{
		{DocumentID: 1, StartLine: 5, SymbolRoles: models.SCIPSymbolRoleDefinition},
		{DocumentID: 1, StartLine: 15, SymbolRoles: models.SCIPSymbolRoleReadAccess},
	}
	pathResolver := func(docID int64) string { return "main.go" }

	resp := ReferencesToResponse("sym", occs, pathResolver)

	if resp.Symbol != "sym" {
		t.Errorf("Symbol = %q, want %q", resp.Symbol, "sym")
	}
	if resp.Total != 2 {
		t.Errorf("Total = %d, want 2", resp.Total)
	}
	if len(resp.References) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.References))
	}
	if resp.References[0].Role != models.SCIPSymbolRoleDefinition {
		t.Errorf("first role = %d, want %d", resp.References[0].Role, models.SCIPSymbolRoleDefinition)
	}
}

func TestRelationshipsToImplementationsResponse(t *testing.T) {
	displayName := "ConcreteType"
	rels := []*models.SCIPRelationship{
		{TargetSymbol: "impl.A", IsImplementation: true},
		{TargetSymbol: "ref.B", IsImplementation: false}, // should be filtered
		{TargetSymbol: "impl.C", IsImplementation: true},
	}

	symbolResolver := func(target string) (*models.SCIPSymbol, *wire.Location) {
		if target == "impl.A" {
			return &models.SCIPSymbol{DisplayName: &displayName}, &wire.Location{Path: "a.go"}
		}
		return &models.SCIPSymbol{}, nil
	}

	resp := RelationshipsToImplementationsResponse("iface", rels, symbolResolver)

	if resp.Symbol != "iface" {
		t.Errorf("Symbol = %q, want %q", resp.Symbol, "iface")
	}
	if resp.Total != 2 {
		t.Errorf("Total = %d, want 2 (non-implementation filtered)", resp.Total)
	}
	if resp.Implementations[0].DisplayName != "ConcreteType" {
		t.Errorf("first DisplayName = %q, want %q", resp.Implementations[0].DisplayName, "ConcreteType")
	}
	if resp.Implementations[0].Location == nil || resp.Implementations[0].Location.Path != "a.go" {
		t.Error("first implementation should have location")
	}
}

func TestRelationshipsToImplementationsResponse_NilSymbol(t *testing.T) {
	rels := []*models.SCIPRelationship{
		{TargetSymbol: "missing", IsImplementation: true},
	}

	symbolResolver := func(target string) (*models.SCIPSymbol, *wire.Location) {
		return nil, nil
	}

	resp := RelationshipsToImplementationsResponse("iface", rels, symbolResolver)

	if resp.Total != 1 {
		t.Errorf("Total = %d, want 1", resp.Total)
	}
	if resp.Implementations[0].Symbol != "missing" {
		t.Errorf("Symbol = %q, want %q", resp.Implementations[0].Symbol, "missing")
	}
	if resp.Implementations[0].DisplayName != "" {
		t.Errorf("DisplayName = %q, want empty", resp.Implementations[0].DisplayName)
	}
}
