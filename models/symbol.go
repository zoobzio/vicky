package models

import "time"

// SymbolKind represents the type of code symbol.
type SymbolKind string

// SymbolKind values.
const (
	SymbolKindFunction  SymbolKind = "function"
	SymbolKindMethod    SymbolKind = "method"
	SymbolKindType      SymbolKind = "type"
	SymbolKindInterface SymbolKind = "interface"
	SymbolKindStruct    SymbolKind = "struct"
	SymbolKindConst     SymbolKind = "const"
	SymbolKindVar       SymbolKind = "var"
	SymbolKindField     SymbolKind = "field"
)

// Symbol represents a code entity extracted from source files.
// Used for "mentioned here" queries - linking documentation to API references.
type Symbol struct {
	ID            int64      `json:"id" db:"id" constraints:"primarykey" description:"Internal symbol ID"`
	VersionID     int64      `json:"version_id" db:"version_id" constraints:"notnull" references:"versions(id)" description:"Parent version"`
	UserID        int64      `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Owner         string     `json:"owner" db:"owner" constraints:"notnull" description:"GitHub org or user" example:"octocat"`
	RepoName      string     `json:"repo_name" db:"repo_name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	Tag           string     `json:"tag" db:"tag" constraints:"notnull" description:"Version tag" example:"v1.0.0"`
	Name          string     `json:"name" db:"name" constraints:"notnull" description:"Symbol name" example:"NewClient"`
	QualifiedName string     `json:"qualified_name" db:"qualified_name" constraints:"notnull" description:"Fully qualified name" example:"vicky.NewClient"`
	Kind          SymbolKind `json:"kind" db:"kind" constraints:"notnull" description:"Type of symbol"`
	Signature     *string    `json:"signature,omitempty" db:"signature" description:"Function/method signature" example:"func NewClient(opts ...Option) *Client"`
	Doc           *string    `json:"doc,omitempty" db:"doc" description:"Godoc comment"`
	FilePath      string     `json:"file_path" db:"file_path" constraints:"notnull" description:"Source file path" example:"client.go"`
	StartLine     int        `json:"start_line" db:"start_line" constraints:"notnull" description:"Starting line number"`
	EndLine       int        `json:"end_line" db:"end_line" constraints:"notnull" description:"Ending line number"`
	Exported      bool       `json:"exported" db:"exported" constraints:"notnull" default:"false" description:"Whether symbol is exported"`
	ParentID      *int64     `json:"parent_id,omitempty" db:"parent_id" references:"symbols(id)" description:"Parent symbol for methods/fields"`
	Vector        []float32  `json:"-" db:"vector" description:"Symbol embedding for similarity"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" default:"now()" description:"Extraction time"`
}

// Clone returns a deep copy of the Symbol.
func (s Symbol) Clone() Symbol {
	c := s
	if s.Signature != nil {
		sig := *s.Signature
		c.Signature = &sig
	}
	if s.Doc != nil {
		d := *s.Doc
		c.Doc = &d
	}
	if s.ParentID != nil {
		p := *s.ParentID
		c.ParentID = &p
	}
	if s.Vector != nil {
		c.Vector = make([]float32, len(s.Vector))
		copy(c.Vector, s.Vector)
	}
	return c
}
