// Package transformers converts between model and API types.
package transformers

import (
	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/models"
)

// OccurrenceToLocation transforms a SCIP occurrence to an API location.
func OccurrenceToLocation(occ *models.SCIPOccurrence, path string) wire.Location {
	return wire.Location{
		Path:      path,
		StartLine: occ.StartLine,
		StartCol:  occ.StartCol,
		EndLine:   occ.EndLine,
		EndCol:    occ.EndCol,
	}
}

// SCIPSymbolToInfo transforms a SCIP symbol to an API symbol info.
func SCIPSymbolToInfo(sym *models.SCIPSymbol, location *wire.Location) wire.SCIPSymbolInfo {
	displayName := ""
	if sym.DisplayName != nil {
		displayName = *sym.DisplayName
	}

	return wire.SCIPSymbolInfo{
		Symbol:        sym.Symbol,
		DisplayName:   displayName,
		Kind:          sym.Kind,
		Documentation: sym.Documentation,
		Location:      location,
	}
}

// OccurrenceToReferenceInfo transforms a SCIP occurrence to an API reference info.
func OccurrenceToReferenceInfo(occ *models.SCIPOccurrence, path string) wire.ReferenceInfo {
	return wire.ReferenceInfo{
		Location: OccurrenceToLocation(occ, path),
		Role:     occ.SymbolRoles,
	}
}

// DefinitionsToResponse transforms definition occurrences to a definition response.
func DefinitionsToResponse(symbol string, occurrences []*models.SCIPOccurrence, pathResolver func(docID int64) string) wire.DefinitionResponse {
	locations := make([]wire.Location, 0, len(occurrences))
	for _, occ := range occurrences {
		path := pathResolver(occ.DocumentID)
		locations = append(locations, OccurrenceToLocation(occ, path))
	}

	return wire.DefinitionResponse{
		Symbol:    symbol,
		Locations: locations,
	}
}

// ReferencesToResponse transforms reference occurrences to a references response.
func ReferencesToResponse(symbol string, occurrences []*models.SCIPOccurrence, pathResolver func(docID int64) string) wire.ReferencesResponse {
	refs := make([]wire.ReferenceInfo, 0, len(occurrences))
	for _, occ := range occurrences {
		path := pathResolver(occ.DocumentID)
		refs = append(refs, OccurrenceToReferenceInfo(occ, path))
	}

	return wire.ReferencesResponse{
		Symbol:     symbol,
		References: refs,
		Total:      len(refs),
	}
}

// RelationshipsToImplementationsResponse transforms relationships to an implementations response.
func RelationshipsToImplementationsResponse(symbol string, relationships []*models.SCIPRelationship, symbolResolver func(targetSymbol string) (*models.SCIPSymbol, *wire.Location)) wire.ImplementationsResponse {
	impls := make([]wire.ImplementationInfo, 0, len(relationships))
	for _, rel := range relationships {
		if !rel.IsImplementation {
			continue
		}

		impl := wire.ImplementationInfo{
			Symbol: rel.TargetSymbol,
		}

		if sym, loc := symbolResolver(rel.TargetSymbol); sym != nil {
			if sym.DisplayName != nil {
				impl.DisplayName = *sym.DisplayName
			}
			impl.Location = loc
		}

		impls = append(impls, impl)
	}

	return wire.ImplementationsResponse{
		Symbol:          symbol,
		Implementations: impls,
		Total:           len(impls),
	}
}

// SCIPSymbolsToListResponse transforms SCIP symbols to a symbol list response.
func SCIPSymbolsToListResponse(symbols []*models.SCIPSymbol, locationResolver func(sym *models.SCIPSymbol) *wire.Location) wire.SymbolListResponse {
	infos := make([]wire.SCIPSymbolInfo, 0, len(symbols))
	for _, sym := range symbols {
		loc := locationResolver(sym)
		infos = append(infos, SCIPSymbolToInfo(sym, loc))
	}

	return wire.SymbolListResponse{
		Symbols: infos,
		Total:   len(infos),
	}
}
