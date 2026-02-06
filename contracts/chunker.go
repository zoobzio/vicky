package contracts

import (
	"context"

	"github.com/zoobzio/vicky/external/chunker"
)

// Chunker defines the contract for language-aware content chunking.
// Language is a plain string to support values beyond models.Language
// (e.g. "markdown" for documentation files).
type Chunker interface {
	// Chunk splits file content into semantic chunks using a language-specific parser.
	Chunk(ctx context.Context, language string, filename string, content []byte) ([]chunker.Result, error)

	// Supports returns true if the chunker has a provider for the given language.
	Supports(language string) bool
}
