package contracts

import (
	"context"

	"github.com/zoobzio/vicky/external/indexer"
	"github.com/zoobzio/vicky/models"
)

// Indexer defines the contract for SCIP indexing operations.
// Implementations may be local (CLI), remote (HTTP/gRPC), or event-driven.
type Indexer interface {
	// Index runs the appropriate SCIP indexer for the given request.
	// The implementation fetches source from blob storage, runs the indexer,
	// and returns the raw SCIP index data.
	Index(ctx context.Context, req indexer.Request) (*indexer.Result, error)

	// Supports returns true if this indexer supports the given language.
	Supports(language models.Language) bool
}
