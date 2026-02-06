package transformers

import (
	"testing"

	"github.com/zoobzio/vicky/models"
)

func TestChunkToResult(t *testing.T) {
	c := &models.Chunk{
		ID:        1,
		Path:      "main.go",
		Kind:      models.ChunkKindFunction,
		Content:   "func main() {}",
		StartLine: 5,
		EndLine:   10,
	}

	result := ChunkToResult(c)

	if result.ID != 1 {
		t.Errorf("ID = %d, want 1", result.ID)
	}
	if result.Path != "main.go" {
		t.Errorf("Path = %q, want %q", result.Path, "main.go")
	}
	if result.Kind != models.ChunkKindFunction {
		t.Errorf("Kind = %q, want %q", result.Kind, models.ChunkKindFunction)
	}
	if result.StartLine != 5 || result.EndLine != 10 {
		t.Errorf("lines = (%d,%d), want (5,10)", result.StartLine, result.EndLine)
	}
}

func TestChunksToSearchResponse(t *testing.T) {
	chunks := []*models.Chunk{
		{ID: 1, Content: "a"},
		{ID: 2, Content: "b"},
	}

	resp := ChunksToSearchResponse("test query", chunks)

	if resp.Query != "test query" {
		t.Errorf("Query = %q, want %q", resp.Query, "test query")
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Results))
	}
	if resp.Results[0].ID != 1 {
		t.Errorf("first ID = %d, want 1", resp.Results[0].ID)
	}
}

func TestDocumentToResult(t *testing.T) {
	d := &models.Document{
		ID:          42,
		Path:        "docs/guide.md",
		ContentType: models.ContentTypeDocs,
	}

	result := DocumentToResult(d)

	if result.ID != 42 {
		t.Errorf("ID = %d, want 42", result.ID)
	}
	if result.Path != "docs/guide.md" {
		t.Errorf("Path = %q, want %q", result.Path, "docs/guide.md")
	}
	if result.ContentType != models.ContentTypeDocs {
		t.Errorf("ContentType = %q, want %q", result.ContentType, models.ContentTypeDocs)
	}
}

func TestDocumentsToSimilarResponse_ExcludesSource(t *testing.T) {
	docs := []*models.Document{
		{ID: 1, Path: "a.go"},
		{ID: 2, Path: "b.go"},
		{ID: 3, Path: "c.go"},
	}

	resp := DocumentsToSimilarResponse(docs, 2) // exclude ID=2

	if len(resp.Results) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Results))
	}
	for _, r := range resp.Results {
		if r.ID == 2 {
			t.Error("source document (ID=2) should be excluded")
		}
	}
}

func TestSymbolToResult(t *testing.T) {
	s := &models.Symbol{
		ID:        10,
		Name:      "NewClient",
		Kind:      models.SymbolKindFunction,
		FilePath:  "client.go",
		StartLine: 42,
		Exported:  true,
	}

	result := SymbolToResult(s)

	if result.ID != 10 {
		t.Errorf("ID = %d, want 10", result.ID)
	}
	if result.Name != "NewClient" {
		t.Errorf("Name = %q, want %q", result.Name, "NewClient")
	}
	if !result.Exported {
		t.Error("Exported = false, want true")
	}
}

func TestSymbolsToSearchResponse(t *testing.T) {
	symbols := []*models.Symbol{
		{ID: 1, Name: "Foo"},
		{ID: 2, Name: "Bar"},
	}

	resp := SymbolsToSearchResponse("search", symbols)

	if resp.Query != "search" {
		t.Errorf("Query = %q, want %q", resp.Query, "search")
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Results))
	}
}
