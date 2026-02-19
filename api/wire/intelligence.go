// Package wire defines request and response types for the vicky HTTP API.
package wire

import "github.com/zoobzio/vicky/models"

// Location represents a source code location.
type Location struct {
	Path      string `json:"path" description:"File path" example:"pkg/client/client.go"`
	StartLine int    `json:"start_line" description:"Starting line (0-indexed)" example:"42"`
	StartCol  int    `json:"start_col" description:"Starting column (0-indexed)" example:"5"`
	EndLine   int    `json:"end_line" description:"Ending line (0-indexed)" example:"42"`
	EndCol    int    `json:"end_col" description:"Ending column (0-indexed)" example:"15"`
}

// SCIPSymbolInfo represents a symbol in code intelligence responses.
type SCIPSymbolInfo struct {
	Symbol        string                `json:"symbol" description:"Qualified SCIP symbol identifier"`
	DisplayName   string                `json:"display_name,omitempty" description:"Human-readable name" example:"Connect"`
	Kind          models.SCIPSymbolKind `json:"kind" description:"Symbol kind"`
	Documentation []string              `json:"documentation,omitempty" description:"Symbol documentation"`
	Location      *Location             `json:"location,omitempty" description:"Definition location"`
}

// ReferenceInfo represents a symbol reference.
type ReferenceInfo struct {
	Location Location               `json:"location" description:"Reference location"`
	Role     models.SCIPSymbolRole  `json:"role" description:"Reference role (definition, read, write, etc.)"`
}

// ImplementationInfo represents an implementation of an interface/type.
type ImplementationInfo struct {
	Symbol      string    `json:"symbol" description:"Implementing symbol identifier"`
	DisplayName string    `json:"display_name,omitempty" description:"Human-readable name"`
	Location    *Location `json:"location,omitempty" description:"Implementation location"`
}

// DefinitionResponse is the API response for go-to-definition.
type DefinitionResponse struct {
	Symbol    string     `json:"symbol" description:"Requested symbol identifier"`
	Locations []Location `json:"locations" description:"Definition locations (usually one, multiple for partial classes)"`
}

// ReferencesResponse is the API response for find-references.
type ReferencesResponse struct {
	Symbol     string          `json:"symbol" description:"Requested symbol identifier"`
	References []ReferenceInfo `json:"references" description:"All references to the symbol"`
	Total      int             `json:"total" description:"Total reference count"`
}

// ImplementationsResponse is the API response for find-implementations.
type ImplementationsResponse struct {
	Symbol          string               `json:"symbol" description:"Requested symbol identifier"`
	Implementations []ImplementationInfo `json:"implementations" description:"Types implementing the interface"`
	Total           int                  `json:"total" description:"Total implementation count"`
}

// SymbolListResponse is the API response for listing symbols in a file or version.
type SymbolListResponse struct {
	Symbols []SCIPSymbolInfo `json:"symbols" description:"Symbols in the requested scope"`
	Total   int              `json:"total" description:"Total symbol count"`
}

// Clone returns a deep copy of the Location.
func (l Location) Clone() Location { return l }

// Clone returns a deep copy of the SCIPSymbolInfo.
func (s SCIPSymbolInfo) Clone() SCIPSymbolInfo {
	c := s
	if s.Documentation != nil {
		c.Documentation = make([]string, len(s.Documentation))
		copy(c.Documentation, s.Documentation)
	}
	if s.Location != nil {
		loc := *s.Location
		c.Location = &loc
	}
	return c
}

// Clone returns a deep copy of the ReferenceInfo.
func (r ReferenceInfo) Clone() ReferenceInfo { return r }

// Clone returns a deep copy of the ImplementationInfo.
func (i ImplementationInfo) Clone() ImplementationInfo {
	c := i
	if i.Location != nil {
		loc := *i.Location
		c.Location = &loc
	}
	return c
}

// Clone returns a deep copy of the DefinitionResponse.
func (d DefinitionResponse) Clone() DefinitionResponse {
	c := d
	if d.Locations != nil {
		c.Locations = make([]Location, len(d.Locations))
		copy(c.Locations, d.Locations)
	}
	return c
}

// Clone returns a deep copy of the ReferencesResponse.
func (r ReferencesResponse) Clone() ReferencesResponse {
	c := r
	if r.References != nil {
		c.References = make([]ReferenceInfo, len(r.References))
		copy(c.References, r.References)
	}
	return c
}

// Clone returns a deep copy of the ImplementationsResponse.
func (i ImplementationsResponse) Clone() ImplementationsResponse {
	c := i
	if i.Implementations != nil {
		c.Implementations = make([]ImplementationInfo, len(i.Implementations))
		for idx, impl := range i.Implementations {
			c.Implementations[idx] = impl.Clone()
		}
	}
	return c
}

// Clone returns a deep copy of the SymbolListResponse.
func (s SymbolListResponse) Clone() SymbolListResponse {
	c := s
	if s.Symbols != nil {
		c.Symbols = make([]SCIPSymbolInfo, len(s.Symbols))
		for idx, sym := range s.Symbols {
			c.Symbols[idx] = sym.Clone()
		}
	}
	return c
}
