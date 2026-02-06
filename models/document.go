package models

import "time"

// ContentType specifies the type of content in a document.
type ContentType string

// ContentType values.
const (
	ContentTypeCode ContentType = "code"
	ContentTypeDocs ContentType = "docs"
)

// Document represents a file within an ingested version.
type Document struct {
	ID          int64       `json:"id" db:"id" constraints:"primarykey" description:"Internal document ID"`
	VersionID   int64       `json:"version_id" db:"version_id" constraints:"notnull" references:"versions(id)" description:"Parent version"`
	UserID      int64       `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Owner       string      `json:"owner" db:"owner" constraints:"notnull" description:"GitHub org or user" example:"octocat"`
	RepoName    string      `json:"repo_name" db:"repo_name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	Tag         string      `json:"tag" db:"tag" constraints:"notnull" description:"Version tag" example:"v1.0.0"`
	Path        string      `json:"path" db:"path" constraints:"notnull" description:"File path within repository" example:"docs/guide.md"`
	ContentType ContentType `json:"content_type" db:"content_type" constraints:"notnull" description:"Type of content"`
	ContentHash string      `json:"content_hash" db:"content_hash" constraints:"notnull" description:"SHA256 of file content"`
	Vector      []float32   `json:"-" db:"vector" description:"Document-level embedding for similarity"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at" default:"now()" description:"Ingestion time"`
}

// Clone returns a deep copy of the Document.
func (d Document) Clone() Document {
	c := d
	if d.Vector != nil {
		c.Vector = make([]float32, len(d.Vector))
		copy(c.Vector, d.Vector)
	}
	return c
}
