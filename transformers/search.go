package transformers

import (
	"github.com/zoobzio/vicky/wire"
	"github.com/zoobzio/vicky/models"
)

// ChunkToResult transforms a Chunk model to an API result.
func ChunkToResult(c *models.Chunk) wire.ChunkResult {
	return wire.ChunkResult{
		ID:        c.ID,
		Path:      c.Path,
		Kind:      c.Kind,
		Content:   c.Content,
		StartLine: c.StartLine,
		EndLine:   c.EndLine,
	}
}

// ChunksToSearchResponse transforms a slice of Chunk models to a search response.
func ChunksToSearchResponse(query string, chunks []*models.Chunk) wire.SearchResponse {
	resp := wire.SearchResponse{
		Query:   query,
		Results: make([]wire.ChunkResult, len(chunks)),
	}
	for i, c := range chunks {
		resp.Results[i] = ChunkToResult(c)
	}
	return resp
}

// SymbolToResult transforms a Symbol model to an API result.
func SymbolToResult(s *models.Symbol) wire.SymbolResult {
	return wire.SymbolResult{
		ID:        s.ID,
		Name:      s.Name,
		Kind:      s.Kind,
		FilePath:  s.FilePath,
		StartLine: s.StartLine,
		Exported:  s.Exported,
	}
}

// SymbolsToSearchResponse transforms a slice of Symbol models to a search response.
func SymbolsToSearchResponse(query string, symbols []*models.Symbol) wire.SymbolSearchResponse {
	resp := wire.SymbolSearchResponse{
		Query:   query,
		Results: make([]wire.SymbolResult, len(symbols)),
	}
	for i, s := range symbols {
		resp.Results[i] = SymbolToResult(s)
	}
	return resp
}

// DocumentToResult transforms a Document model to an API result.
func DocumentToResult(d *models.Document) wire.DocumentResult {
	return wire.DocumentResult{
		ID:          d.ID,
		Path:        d.Path,
		ContentType: d.ContentType,
	}
}

// DocumentsToSimilarResponse transforms a slice of Document models to a similar documents response.
// Excludes the source document if its ID matches excludeID.
func DocumentsToSimilarResponse(docs []*models.Document, excludeID int64) wire.SimilarDocumentsResponse {
	resp := wire.SimilarDocumentsResponse{
		Results: make([]wire.DocumentResult, 0, len(docs)),
	}
	for _, d := range docs {
		if d.ID == excludeID {
			continue
		}
		resp.Results = append(resp.Results, DocumentToResult(d))
	}
	return resp
}
