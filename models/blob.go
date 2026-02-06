package models

// Blob represents a code file stored in object storage during repository ingestion.
type Blob struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Language string `json:"language,omitempty"`
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Tag      string `json:"tag"`
}
