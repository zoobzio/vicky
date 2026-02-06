// Package scip provides parsing of SCIP index files.
package scip

import (
	"context"
	"fmt"

	"github.com/sourcegraph/scip/bindings/go/scip"
	"github.com/zoobzio/vicky/models"
	"google.golang.org/protobuf/proto"
)

// DocumentMeta provides context for converting SCIP data to models.
type DocumentMeta struct {
	UserID   int64
	Owner    string
	RepoName string
	Tag      string
}

// FileMeta extends DocumentMeta with document-specific information.
type FileMeta struct {
	DocumentMeta
	DocumentID int64
	Path       string
}

// Result contains parsed SCIP data converted to models.
type Result struct {
	Symbols       []models.SCIPSymbol
	Occurrences   []models.SCIPOccurrence
	Relationships []models.SCIPRelationship
}

// Parser parses SCIP index files.
type Parser struct{}

// New creates a new SCIP parser.
func New() *Parser {
	return &Parser{}
}

// Parse reads a SCIP index from protobuf bytes and returns the raw index.
func (p *Parser) Parse(_ context.Context, data []byte) (*scip.Index, error) {
	var index scip.Index
	if err := proto.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("unmarshal scip index: %w", err)
	}
	return &index, nil
}

// ParseDocument extracts symbols and occurrences from a single SCIP document.
func (p *Parser) ParseDocument(_ context.Context, doc *scip.Document, meta FileMeta) *Result {
	result := &Result{}

	// Convert symbols
	for _, sym := range doc.Symbols {
		symbol := ConvertSymbol(sym, meta)
		result.Symbols = append(result.Symbols, symbol)

		// Convert relationships for this symbol
		for _, rel := range sym.Relationships {
			relationship := ConvertRelationship(rel)
			result.Relationships = append(result.Relationships, relationship)
		}
	}

	// Convert occurrences
	for _, occ := range doc.Occurrences {
		occurrence := ConvertOccurrence(occ, meta)
		result.Occurrences = append(result.Occurrences, occurrence)
	}

	return result
}

// ParseIndex extracts all data from a SCIP index.
// The documentMapper function maps relative paths to FileMeta.
// Documents not found in the mapper are skipped.
func (p *Parser) ParseIndex(ctx context.Context, index *scip.Index, documentMapper func(relativePath string) (FileMeta, bool)) *Result {
	result := &Result{}

	for _, doc := range index.Documents {
		meta, ok := documentMapper(doc.RelativePath)
		if !ok {
			continue
		}

		docResult := p.ParseDocument(ctx, doc, meta)
		result.Symbols = append(result.Symbols, docResult.Symbols...)
		result.Occurrences = append(result.Occurrences, docResult.Occurrences...)
		result.Relationships = append(result.Relationships, docResult.Relationships...)
	}

	// Also process external symbols (defined outside this index)
	for _, sym := range index.ExternalSymbols {
		// External symbols don't have a document - skip for now
		// These are typically stdlib or dependency symbols
		_ = sym
	}

	return result
}
