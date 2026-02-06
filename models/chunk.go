// Package models defines domain types for vicky.
package models

import "time"

// ChunkKind represents the type of content segment.
type ChunkKind string

// ChunkKind values for code chunks (aligned with chisel).
const (
	ChunkKindFunction  ChunkKind = "function"
	ChunkKindMethod    ChunkKind = "method"
	ChunkKindClass     ChunkKind = "class"     // struct in Go
	ChunkKindInterface ChunkKind = "interface"
	ChunkKindType      ChunkKind = "type"
	ChunkKindEnum      ChunkKind = "enum"
	ChunkKindConstant  ChunkKind = "constant"
	ChunkKindVariable  ChunkKind = "variable"
	ChunkKindModule    ChunkKind = "module" // package-level
)

// ChunkKind values for documentation chunks.
const (
	ChunkKindSection   ChunkKind = "section"
	ChunkKindParagraph ChunkKind = "paragraph"
	ChunkKindCode      ChunkKind = "code" // Code block within docs
)

// Chunk represents an embedded segment of a document.
type Chunk struct {
	ID         int64     `json:"id" db:"id" constraints:"primarykey" description:"Internal chunk ID"`
	DocumentID int64     `json:"document_id" db:"document_id" constraints:"notnull" references:"documents(id)" description:"Parent document"`
	UserID     int64     `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Owner      string    `json:"owner" db:"owner" constraints:"notnull" description:"GitHub org or user" example:"octocat"`
	RepoName   string    `json:"repo_name" db:"repo_name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	Tag        string    `json:"tag" db:"tag" constraints:"notnull" description:"Version tag" example:"v1.0.0"`
	Path       string    `json:"path" db:"path" constraints:"notnull" description:"Document path" example:"docs/guide.md"`
	Kind       ChunkKind `json:"kind" db:"kind" constraints:"notnull" description:"Type of chunk"`
	StartLine  int       `json:"start_line" db:"start_line" constraints:"notnull" description:"Starting line number"`
	EndLine    int       `json:"end_line" db:"end_line" constraints:"notnull" description:"Ending line number"`
	Symbol     *string   `json:"symbol,omitempty" db:"symbol" description:"Function/type name if applicable" example:"NewClient"`
	Context    []string  `json:"context,omitempty" db:"context" description:"Parent chain for nested symbols" example:"[\"type UserService\", \"method GetUser\"]"`
	Content    string    `json:"content" db:"content" constraints:"notnull" description:"Raw text content"`
	Vector     []float32 `json:"-" db:"vector" description:"Chunk embedding"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" default:"now()" description:"Ingestion time"`
}

// Clone returns a deep copy of the Chunk.
func (c Chunk) Clone() Chunk {
	cl := c
	if c.Symbol != nil {
		s := *c.Symbol
		cl.Symbol = &s
	}
	if c.Context != nil {
		cl.Context = make([]string, len(c.Context))
		copy(cl.Context, c.Context)
	}
	if c.Vector != nil {
		cl.Vector = make([]float32, len(c.Vector))
		copy(cl.Vector, c.Vector)
	}
	return cl
}
