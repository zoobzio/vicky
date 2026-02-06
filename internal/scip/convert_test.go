package scip

import (
	"testing"

	scipproto "github.com/sourcegraph/scip/bindings/go/scip"
	"github.com/zoobzio/vicky/models"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		name                                       string
		input                                      []int32
		wantStartLine, wantStartCol, wantEndLine, wantEndCol int
	}{
		{"empty", nil, 0, 0, 0, 0},
		{"same line", []int32{5, 10, 20}, 5, 10, 5, 20},
		{"multi line", []int32{1, 0, 3, 5}, 1, 0, 3, 5},
		{"two elements", []int32{7, 3}, 7, 3, 7, 3},
		{"one element", []int32{4}, 4, 0, 4, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl, sc, el, ec := parseRange(tt.input)
			if sl != tt.wantStartLine || sc != tt.wantStartCol || el != tt.wantEndLine || ec != tt.wantEndCol {
				t.Errorf("parseRange(%v) = (%d,%d,%d,%d), want (%d,%d,%d,%d)",
					tt.input, sl, sc, el, ec,
					tt.wantStartLine, tt.wantStartCol, tt.wantEndLine, tt.wantEndCol)
			}
		})
	}
}

func TestParseSymbolName(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
		want   string
	}{
		{"fully qualified", "scip-go gomod github.com/foo/bar v1.0.0 Client#Connect().", "Connect"},
		{"nested", "scip-go gomod example.com/pkg v1.0.0 Server#Handler#ServeHTTP().", "ServeHTTP"},
		{"simple", "scip-go gomod example.com/pkg v1.0.0 main().", "main"},
		{"empty", "", ""},
		{"package level", "scip-go gomod example.com/pkg v1.0.0 DefaultTimeout.", "DefaultTimeout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSymbolName(tt.symbol)
			if got != tt.want {
				t.Errorf("ParseSymbolName(%q) = %q, want %q", tt.symbol, got, tt.want)
			}
		})
	}
}

func TestConvertSymbol(t *testing.T) {
	meta := FileMeta{
		DocumentMeta: DocumentMeta{
			UserID:   1000,
			Owner:    "testorg",
			RepoName: "testrepo",
			Tag:      "v1.0.0",
		},
		DocumentID: 42,
		Path:       "main.go",
	}

	sym := &scipproto.SymbolInformation{
		Symbol:          "scip-go gomod example.com/test v1.0.0 Foo().",
		Kind:            scipproto.SymbolInformation_Function,
		DisplayName:     "Foo",
		Documentation:   []string{"Foo does things."},
		EnclosingSymbol: "Bar",
	}

	result := ConvertSymbol(sym, meta)

	if result.DocumentID != 42 {
		t.Errorf("DocumentID = %d, want 42", result.DocumentID)
	}
	if result.UserID != 1000 {
		t.Errorf("UserID = %d, want 1000", result.UserID)
	}
	if result.Symbol != sym.Symbol {
		t.Errorf("Symbol = %q, want %q", result.Symbol, sym.Symbol)
	}
	if result.Kind != models.SCIPSymbolKind(scipproto.SymbolInformation_Function) {
		t.Errorf("Kind = %d, want %d", result.Kind, scipproto.SymbolInformation_Function)
	}
	if result.DisplayName == nil || *result.DisplayName != "Foo" {
		t.Errorf("DisplayName = %v, want %q", result.DisplayName, "Foo")
	}
	if result.EnclosingSymbol == nil || *result.EnclosingSymbol != "Bar" {
		t.Errorf("EnclosingSymbol = %v, want %q", result.EnclosingSymbol, "Bar")
	}
}

func TestConvertSymbol_NilOptionals(t *testing.T) {
	meta := FileMeta{}
	sym := &scipproto.SymbolInformation{
		Symbol: "test",
		Kind:   scipproto.SymbolInformation_Variable,
	}

	result := ConvertSymbol(sym, meta)

	if result.DisplayName != nil {
		t.Error("expected nil DisplayName for empty string")
	}
	if result.EnclosingSymbol != nil {
		t.Error("expected nil EnclosingSymbol for empty string")
	}
}

func TestConvertOccurrence(t *testing.T) {
	meta := FileMeta{
		DocumentMeta: DocumentMeta{
			UserID:   1000,
			Owner:    "testorg",
			RepoName: "testrepo",
			Tag:      "v1.0.0",
		},
		DocumentID: 42,
	}

	occ := &scipproto.Occurrence{
		Range:          []int32{5, 10, 20},
		Symbol:         "scip-go gomod example.com/test v1.0.0 Foo().",
		SymbolRoles:    int32(scipproto.SymbolRole_Definition),
		SyntaxKind:     scipproto.SyntaxKind_IdentifierFunctionDefinition,
		EnclosingRange: []int32{0, 0, 10, 0},
	}

	result := ConvertOccurrence(occ, meta)

	if result.StartLine != 5 || result.StartCol != 10 || result.EndLine != 5 || result.EndCol != 20 {
		t.Errorf("range = (%d,%d,%d,%d), want (5,10,5,20)",
			result.StartLine, result.StartCol, result.EndLine, result.EndCol)
	}
	if result.SyntaxKind == nil {
		t.Fatal("expected non-nil SyntaxKind")
	}
	if len(result.EnclosingRange) != 4 {
		t.Errorf("EnclosingRange length = %d, want 4", len(result.EnclosingRange))
	}
}

func TestConvertOccurrence_NilOptionals(t *testing.T) {
	meta := FileMeta{}
	occ := &scipproto.Occurrence{
		Range:  []int32{1, 0, 5},
		Symbol: "test",
	}

	result := ConvertOccurrence(occ, meta)

	if result.SyntaxKind != nil {
		t.Error("expected nil SyntaxKind for UnspecifiedSyntaxKind")
	}
	if result.EnclosingRange != nil {
		t.Error("expected nil EnclosingRange when empty")
	}
}

func TestConvertRelationship(t *testing.T) {
	rel := &scipproto.Relationship{
		Symbol:           "scip-go gomod example.com/test v1.0.0 Bar().",
		IsReference:      true,
		IsImplementation: true,
		IsTypeDefinition: false,
		IsDefinition:     false,
	}

	result := ConvertRelationship(rel)

	if result.TargetSymbol != rel.Symbol {
		t.Errorf("TargetSymbol = %q, want %q", result.TargetSymbol, rel.Symbol)
	}
	if !result.IsReference {
		t.Error("IsReference = false, want true")
	}
	if !result.IsImplementation {
		t.Error("IsImplementation = false, want true")
	}
	if result.IsTypeDefinition {
		t.Error("IsTypeDefinition = true, want false")
	}
	if result.IsDefinition {
		t.Error("IsDefinition = true, want false")
	}
}
