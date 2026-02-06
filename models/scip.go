package models

import (
	"encoding/json"
	"time"
)

// SCIPSymbolKind represents the kind of symbol from SCIP indexing.
// Values correspond to the SCIP SymbolInformation.Kind enum.
type SCIPSymbolKind int

// SCIPSymbolKind values from the SCIP SymbolInformation.Kind enum.
const (
	SCIPSymbolKindUnspecified SCIPSymbolKind = 0
	SCIPSymbolKindFunction    SCIPSymbolKind = 12
	SCIPSymbolKindMethod      SCIPSymbolKind = 28
	SCIPSymbolKindClass       SCIPSymbolKind = 5
	SCIPSymbolKindInterface   SCIPSymbolKind = 19
	SCIPSymbolKindStruct      SCIPSymbolKind = 62
	SCIPSymbolKindEnum        SCIPSymbolKind = 10
	SCIPSymbolKindConstant    SCIPSymbolKind = 7
	SCIPSymbolKindVariable    SCIPSymbolKind = 70
	SCIPSymbolKindField       SCIPSymbolKind = 11
	SCIPSymbolKindModule      SCIPSymbolKind = 31
	SCIPSymbolKindPackage     SCIPSymbolKind = 36
	SCIPSymbolKindType        SCIPSymbolKind = 66
	SCIPSymbolKindTypeAlias   SCIPSymbolKind = 67
)

// SCIPSymbolRole is a bitmask representing the role of a symbol occurrence.
type SCIPSymbolRole int

// SCIPSymbolRole bitmask values.
const (
	SCIPSymbolRoleDefinition        SCIPSymbolRole = 0x1
	SCIPSymbolRoleImport            SCIPSymbolRole = 0x2
	SCIPSymbolRoleWriteAccess       SCIPSymbolRole = 0x4
	SCIPSymbolRoleReadAccess        SCIPSymbolRole = 0x8
	SCIPSymbolRoleGenerated         SCIPSymbolRole = 0x10
	SCIPSymbolRoleTest              SCIPSymbolRole = 0x20
	SCIPSymbolRoleForwardDefinition SCIPSymbolRole = 0x40
)

// SCIPSyntaxKind represents syntax highlighting classification.
type SCIPSyntaxKind int

// SCIPSymbol represents a symbol definition from SCIP indexing.
type SCIPSymbol struct {
	ID                     int64           `json:"id" db:"id" constraints:"primarykey" description:"Internal SCIP symbol ID"`
	DocumentID             int64           `json:"document_id" db:"document_id" constraints:"notnull" references:"documents(id)" description:"Parent document"`
	UserID                 int64           `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Owner                  string          `json:"owner" db:"owner" constraints:"notnull" description:"GitHub org or user" example:"octocat"`
	RepoName               string          `json:"repo_name" db:"repo_name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	Tag                    string          `json:"tag" db:"tag" constraints:"notnull" description:"Version tag" example:"v1.0.0"`
	Symbol                 string          `json:"symbol" db:"symbol" constraints:"notnull" description:"Fully qualified SCIP symbol identifier" example:"scip-go gomod github.com/foo/bar v1.0.0 Client#Connect()."`
	Kind                   SCIPSymbolKind  `json:"kind" db:"kind" constraints:"notnull" description:"SCIP symbol kind"`
	DisplayName            *string         `json:"display_name,omitempty" db:"display_name" description:"Human-readable symbol name" example:"Connect"`
	Documentation          []string        `json:"documentation,omitempty" db:"documentation" description:"Markdown-formatted documentation"`
	EnclosingSymbol        *string         `json:"enclosing_symbol,omitempty" db:"enclosing_symbol" description:"Parent symbol for nested definitions"`
	SignatureDocumentation json.RawMessage `json:"signature_documentation,omitempty" db:"signature_documentation" description:"Method/type signature details as SCIP Document"`
	CreatedAt              time.Time       `json:"created_at" db:"created_at" default:"now()" description:"Indexing time"`
}

// SCIPOccurrence represents a symbol occurrence from SCIP indexing.
type SCIPOccurrence struct {
	ID             int64          `json:"id" db:"id" constraints:"primarykey" description:"Internal SCIP occurrence ID"`
	DocumentID     int64          `json:"document_id" db:"document_id" constraints:"notnull" references:"documents(id)" description:"Parent document"`
	UserID         int64          `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Owner          string         `json:"owner" db:"owner" constraints:"notnull" description:"GitHub org or user" example:"octocat"`
	RepoName       string         `json:"repo_name" db:"repo_name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	Tag            string         `json:"tag" db:"tag" constraints:"notnull" description:"Version tag" example:"v1.0.0"`
	Symbol         string         `json:"symbol" db:"symbol" constraints:"notnull" description:"SCIP symbol identifier this occurrence refers to"`
	SymbolRoles    SCIPSymbolRole `json:"symbol_roles" db:"symbol_roles" constraints:"notnull" description:"Bitmask of symbol roles"`
	StartLine      int            `json:"start_line" db:"start_line" constraints:"notnull" description:"Starting line number (0-indexed)"`
	StartCol       int            `json:"start_col" db:"start_col" constraints:"notnull" description:"Starting column (0-indexed)"`
	EndLine        int            `json:"end_line" db:"end_line" constraints:"notnull" description:"Ending line number (0-indexed)"`
	EndCol         int            `json:"end_col" db:"end_col" constraints:"notnull" description:"Ending column (0-indexed)"`
	SyntaxKind     *SCIPSyntaxKind `json:"syntax_kind,omitempty" db:"syntax_kind" description:"Syntax highlighting classification"`
	EnclosingRange []int          `json:"enclosing_range,omitempty" db:"enclosing_range" description:"Parent AST node bounds"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at" default:"now()" description:"Indexing time"`
}

// SCIPRelationship represents a relationship between symbols from SCIP indexing.
type SCIPRelationship struct {
	ID               int64  `json:"id" db:"id" constraints:"primarykey" description:"Internal SCIP relationship ID"`
	SCIPSymbolID     int64  `json:"scip_symbol_id" db:"scip_symbol_id" constraints:"notnull" references:"scip_symbols(id)" description:"Source symbol"`
	TargetSymbol     string `json:"target_symbol" db:"target_symbol" constraints:"notnull" description:"Target SCIP symbol identifier"`
	IsReference      bool   `json:"is_reference" db:"is_reference" default:"false" description:"Include in find references results"`
	IsImplementation bool   `json:"is_implementation" db:"is_implementation" default:"false" description:"Include in find implementations results"`
	IsTypeDefinition bool   `json:"is_type_definition" db:"is_type_definition" default:"false" description:"Include in go to type definition results"`
	IsDefinition     bool   `json:"is_definition" db:"is_definition" default:"false" description:"Override definition behavior"`
}

// Clone returns a deep copy of the SCIPSymbol.
func (s SCIPSymbol) Clone() SCIPSymbol {
	c := s
	if s.DisplayName != nil {
		d := *s.DisplayName
		c.DisplayName = &d
	}
	if s.Documentation != nil {
		c.Documentation = make([]string, len(s.Documentation))
		copy(c.Documentation, s.Documentation)
	}
	if s.EnclosingSymbol != nil {
		e := *s.EnclosingSymbol
		c.EnclosingSymbol = &e
	}
	if s.SignatureDocumentation != nil {
		c.SignatureDocumentation = make(json.RawMessage, len(s.SignatureDocumentation))
		copy(c.SignatureDocumentation, s.SignatureDocumentation)
	}
	return c
}

// Clone returns a deep copy of the SCIPOccurrence.
func (o SCIPOccurrence) Clone() SCIPOccurrence {
	c := o
	if o.SyntaxKind != nil {
		sk := *o.SyntaxKind
		c.SyntaxKind = &sk
	}
	if o.EnclosingRange != nil {
		c.EnclosingRange = make([]int, len(o.EnclosingRange))
		copy(c.EnclosingRange, o.EnclosingRange)
	}
	return c
}

// Clone returns a deep copy of the SCIPRelationship.
func (r SCIPRelationship) Clone() SCIPRelationship { return r }
