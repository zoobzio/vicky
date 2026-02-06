package wire

import "github.com/zoobzio/vicky/models"

// ChunkResult is a search result item.
type ChunkResult struct {
	ID        int64            `json:"id" description:"Chunk ID"`
	Path      string           `json:"path" description:"File path" example:"pkg/client/client.go"`
	Kind      models.ChunkKind `json:"kind" description:"Chunk type" example:"function"`
	Content   string           `json:"content" description:"Chunk content"`
	StartLine int              `json:"start_line" description:"Starting line number" example:"42"`
	EndLine   int              `json:"end_line" description:"Ending line number" example:"67"`
}

// SymbolResult is a symbol search result item.
type SymbolResult struct {
	ID        int64             `json:"id" description:"Symbol ID"`
	Name      string            `json:"name" description:"Symbol name" example:"NewClient"`
	Kind      models.SymbolKind `json:"kind" description:"Symbol type" example:"function"`
	FilePath  string            `json:"file_path" description:"Source file path" example:"client.go"`
	StartLine int               `json:"start_line" description:"Line number" example:"42"`
	Exported  bool              `json:"exported" description:"Whether symbol is exported" example:"true"`
}

// DocumentResult is a document similarity result item.
type DocumentResult struct {
	ID          int64              `json:"id" description:"Document ID"`
	Path        string             `json:"path" description:"File path" example:"docs/guide.md"`
	ContentType models.ContentType `json:"content_type" description:"Content type" example:"markdown"`
}

// SearchResponse is the API response for chunk search.
type SearchResponse struct {
	Query   string        `json:"query" description:"Original search query"`
	Results []ChunkResult `json:"results" description:"Matching chunks ordered by relevance"`
}

// SymbolSearchResponse is the API response for symbol search.
type SymbolSearchResponse struct {
	Query   string         `json:"query" description:"Original search query"`
	Results []SymbolResult `json:"results" description:"Matching symbols ordered by relevance"`
}

// SimilarDocumentsResponse is the API response for similar documents.
type SimilarDocumentsResponse struct {
	Results []DocumentResult `json:"results" description:"Similar documents ordered by similarity"`
}

// Clone returns a deep copy of the ChunkResult.
func (c ChunkResult) Clone() ChunkResult { return c }

// Clone returns a deep copy of the SymbolResult.
func (s SymbolResult) Clone() SymbolResult { return s }

// Clone returns a deep copy of the DocumentResult.
func (d DocumentResult) Clone() DocumentResult { return d }

// Clone returns a deep copy of the SearchResponse.
func (s SearchResponse) Clone() SearchResponse {
	c := s
	if s.Results != nil {
		c.Results = make([]ChunkResult, len(s.Results))
		copy(c.Results, s.Results)
	}
	return c
}

// Clone returns a deep copy of the SymbolSearchResponse.
func (s SymbolSearchResponse) Clone() SymbolSearchResponse {
	c := s
	if s.Results != nil {
		c.Results = make([]SymbolResult, len(s.Results))
		copy(c.Results, s.Results)
	}
	return c
}

// Clone returns a deep copy of the SimilarDocumentsResponse.
func (s SimilarDocumentsResponse) Clone() SimilarDocumentsResponse {
	c := s
	if s.Results != nil {
		c.Results = make([]DocumentResult, len(s.Results))
		copy(c.Results, s.Results)
	}
	return c
}
