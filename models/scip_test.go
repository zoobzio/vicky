package models

import (
	"encoding/json"
	"testing"
)

func TestSCIPSymbolClone(t *testing.T) {
	displayName := "Connect"
	enclosing := "Client"
	sigDoc := json.RawMessage(`{"text":"func Connect()"}`)

	orig := SCIPSymbol{
		ID:                     1,
		DisplayName:            &displayName,
		Documentation:          []string{"Connects to server.", "See also: Disconnect."},
		EnclosingSymbol:        &enclosing,
		SignatureDocumentation: sigDoc,
	}
	clone := orig.Clone()

	*clone.DisplayName = "CHANGED"
	clone.Documentation[0] = "CHANGED"
	*clone.EnclosingSymbol = "CHANGED"
	clone.SignatureDocumentation[0] = 'X'

	if *orig.DisplayName != "Connect" {
		t.Error("Clone did not isolate DisplayName")
	}
	if orig.Documentation[0] != "Connects to server." {
		t.Error("Clone did not isolate Documentation")
	}
	if *orig.EnclosingSymbol != "Client" {
		t.Error("Clone did not isolate EnclosingSymbol")
	}
	if orig.SignatureDocumentation[0] != '{' {
		t.Error("Clone did not isolate SignatureDocumentation")
	}
}

func TestSCIPSymbolClone_NilPointers(t *testing.T) {
	orig := SCIPSymbol{ID: 1}
	clone := orig.Clone()

	if clone.DisplayName != nil || clone.Documentation != nil || clone.EnclosingSymbol != nil || clone.SignatureDocumentation != nil {
		t.Error("Clone should preserve nil fields")
	}
}

func TestSCIPOccurrenceClone(t *testing.T) {
	sk := SCIPSyntaxKind(1)
	orig := SCIPOccurrence{
		ID:             1,
		SyntaxKind:     &sk,
		EnclosingRange: []int{0, 0, 10, 0},
	}
	clone := orig.Clone()

	*clone.SyntaxKind = SCIPSyntaxKind(99)
	clone.EnclosingRange[2] = 999

	if *orig.SyntaxKind != SCIPSyntaxKind(1) {
		t.Error("Clone did not isolate SyntaxKind")
	}
	if orig.EnclosingRange[2] != 10 {
		t.Error("Clone did not isolate EnclosingRange")
	}
}

func TestSCIPOccurrenceClone_NilPointers(t *testing.T) {
	orig := SCIPOccurrence{ID: 1}
	clone := orig.Clone()

	if clone.SyntaxKind != nil || clone.EnclosingRange != nil {
		t.Error("Clone should preserve nil fields")
	}
}

func TestSCIPRelationshipClone(t *testing.T) {
	orig := SCIPRelationship{
		ID:               1,
		IsReference:      true,
		IsImplementation: true,
	}
	clone := orig.Clone()

	clone.IsReference = false

	if !orig.IsReference {
		t.Error("Clone should produce independent value copy")
	}
}
