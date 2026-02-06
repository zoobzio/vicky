package scip

import (
	"encoding/json"
	"strings"

	"github.com/sourcegraph/scip/bindings/go/scip"
	"github.com/zoobzio/vicky/models"
)

// ConvertSymbol converts a SCIP SymbolInformation to models.SCIPSymbol.
func ConvertSymbol(sym *scip.SymbolInformation, meta FileMeta) models.SCIPSymbol {
	var displayName *string
	if sym.DisplayName != "" {
		displayName = &sym.DisplayName
	}

	var enclosingSymbol *string
	if sym.EnclosingSymbol != "" {
		enclosingSymbol = &sym.EnclosingSymbol
	}

	var sigDoc json.RawMessage
	if sym.SignatureDocumentation != nil {
		// Marshal the signature documentation to JSON for storage
		if data, err := json.Marshal(sym.SignatureDocumentation); err == nil {
			sigDoc = data
		}
	}

	return models.SCIPSymbol{
		DocumentID:             meta.DocumentID,
		UserID:                 meta.UserID,
		Owner:                  meta.Owner,
		RepoName:               meta.RepoName,
		Tag:                    meta.Tag,
		Symbol:                 sym.Symbol,
		Kind:                   models.SCIPSymbolKind(sym.Kind),
		DisplayName:            displayName,
		Documentation:          sym.Documentation,
		EnclosingSymbol:        enclosingSymbol,
		SignatureDocumentation: sigDoc,
	}
}

// ConvertOccurrence converts a SCIP Occurrence to models.SCIPOccurrence.
func ConvertOccurrence(occ *scip.Occurrence, meta FileMeta) models.SCIPOccurrence {
	startLine, startCol, endLine, endCol := parseRange(occ.Range)

	var syntaxKind *models.SCIPSyntaxKind
	if occ.SyntaxKind != scip.SyntaxKind_UnspecifiedSyntaxKind {
		sk := models.SCIPSyntaxKind(occ.SyntaxKind)
		syntaxKind = &sk
	}

	var enclosingRange []int
	if len(occ.EnclosingRange) > 0 {
		enclosingRange = make([]int, len(occ.EnclosingRange))
		for i, v := range occ.EnclosingRange {
			enclosingRange[i] = int(v)
		}
	}

	return models.SCIPOccurrence{
		DocumentID:     meta.DocumentID,
		UserID:         meta.UserID,
		Owner:          meta.Owner,
		RepoName:       meta.RepoName,
		Tag:            meta.Tag,
		Symbol:         occ.Symbol,
		SymbolRoles:    models.SCIPSymbolRole(occ.SymbolRoles),
		StartLine:      startLine,
		StartCol:       startCol,
		EndLine:        endLine,
		EndCol:         endCol,
		SyntaxKind:     syntaxKind,
		EnclosingRange: enclosingRange,
	}
}

// ConvertRelationship converts a SCIP Relationship to models.SCIPRelationship.
// Note: SCIPSymbolID must be set after the symbol is inserted.
func ConvertRelationship(rel *scip.Relationship) models.SCIPRelationship {
	return models.SCIPRelationship{
		TargetSymbol:     rel.Symbol,
		IsReference:      rel.IsReference,
		IsImplementation: rel.IsImplementation,
		IsTypeDefinition: rel.IsTypeDefinition,
		IsDefinition:     rel.IsDefinition,
	}
}

// parseRange extracts line/column coordinates from a SCIP range.
// SCIP ranges are encoded as [startLine, startChar, endLine, endChar]
// or [startLine, startChar, endChar] when on the same line.
func parseRange(r []int32) (startLine, startCol, endLine, endCol int) {
	if len(r) == 0 {
		return 0, 0, 0, 0
	}

	startLine = int(r[0])
	if len(r) >= 2 {
		startCol = int(r[1])
	}

	switch {
	case len(r) == 3:
		// Same line: [startLine, startChar, endChar]
		endLine = startLine
		endCol = int(r[2])
	case len(r) >= 4:
		// Different lines: [startLine, startChar, endLine, endChar]
		endLine = int(r[2])
		endCol = int(r[3])
	default:
		endLine = startLine
		endCol = startCol
	}

	return
}

// ParseSymbolName extracts a human-readable name from a SCIP symbol identifier.
// SCIP symbols follow a grammar like: "scip-go gomod github.com/foo/bar v1.0.0 Client#Connect()."
func ParseSymbolName(symbol string) string {
	// Find the last segment after the version
	parts := strings.Fields(symbol)
	if len(parts) == 0 {
		return symbol
	}

	// The descriptor is typically the last part
	descriptor := parts[len(parts)-1]

	// Remove trailing punctuation (method suffix, etc.)
	descriptor = strings.TrimRight(descriptor, "().#/")

	// Handle nested symbols (e.g., "Client#Connect")
	if idx := strings.LastIndex(descriptor, "#"); idx >= 0 {
		descriptor = descriptor[idx+1:]
	}
	if idx := strings.LastIndex(descriptor, "/"); idx >= 0 {
		descriptor = descriptor[idx+1:]
	}

	return descriptor
}
