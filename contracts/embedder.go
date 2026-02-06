package contracts

import "context"

// Embedder defines the contract for vector embedding generation.
type Embedder interface {
	// Embed generates embeddings for the given texts (document mode).
	Embed(ctx context.Context, texts []string) ([][]float32, error)

	// EmbedQuery generates embeddings optimized for search queries.
	EmbedQuery(ctx context.Context, texts []string) ([][]float32, error)

	// Dimensions returns the vector dimensionality.
	Dimensions() int
}
