// Package embedding provides a vex-backed implementation of contracts.Embedder.
package embedding

import (
	"context"
	"fmt"
	"time"

	"github.com/zoobzio/vex"
	"github.com/zoobzio/vex/cohere"
	"github.com/zoobzio/vex/gemini"
	"github.com/zoobzio/vex/openai"
	"github.com/zoobzio/vex/voyage"
)

// Client implements contracts.Embedder using vex embedding providers.
type Client struct {
	svc  *vex.Service
	dims int
}

// NewClient creates a new embedding client for the given provider.
// Supported providers: stub, openai, voyage, gemini, cohere.
func NewClient(provider, model, apiKey string, dimensions int) (*Client, error) {
	if provider == "" || provider == "stub" {
		return &Client{dims: dimensions}, nil
	}

	p, err := newProvider(provider, model, apiKey, dimensions)
	if err != nil {
		return nil, err
	}

	svc := vex.NewService(p,
		vex.WithBackoff(3, 100*time.Millisecond),
		vex.WithTimeout(30*time.Second),
		vex.WithCircuitBreaker(5, 30*time.Second),
	)

	return &Client{svc: svc, dims: dimensions}, nil
}

// Embed generates embeddings for the given texts (document mode).
func (c *Client) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	// Stub mode — return zero vectors
	if c.svc == nil {
		return zeroVectors(len(texts), c.dims), nil
	}

	vectors, err := c.svc.Batch(ctx, texts)
	if err != nil {
		return nil, err
	}

	return toFloat32Slices(vectors), nil
}

// EmbedQuery generates embeddings optimized for search queries.
func (c *Client) EmbedQuery(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	// Stub mode — return zero vectors
	if c.svc == nil {
		return zeroVectors(len(texts), c.dims), nil
	}

	vectors, err := c.svc.BatchQuery(ctx, texts)
	if err != nil {
		return nil, err
	}

	return toFloat32Slices(vectors), nil
}

// Dimensions returns the vector dimensionality.
func (c *Client) Dimensions() int {
	return c.dims
}

// newProvider creates a vex.Provider from configuration.
func newProvider(provider, model, apiKey string, dimensions int) (vex.Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("%s provider requires API key", provider)
	}

	switch provider {
	case "openai":
		return openai.New(openai.Config{
			APIKey:     apiKey,
			Model:      model,
			Dimensions: dimensions,
		}), nil

	case "voyage":
		return voyage.New(voyage.Config{
			APIKey:     apiKey,
			Model:      model,
			Dimensions: dimensions,
		}), nil

	case "gemini", "google":
		return gemini.New(gemini.Config{
			APIKey:     apiKey,
			Model:      model,
			Dimensions: dimensions,
		}), nil

	case "cohere":
		return cohere.New(cohere.Config{
			APIKey:     apiKey,
			Model:      model,
			Dimensions: dimensions,
		}), nil

	default:
		return nil, fmt.Errorf("unknown embedding provider: %s", provider)
	}
}

// toFloat32Slices converts vex.Vector slice to [][]float32.
func toFloat32Slices(vectors []vex.Vector) [][]float32 {
	result := make([][]float32, len(vectors))
	for i, vec := range vectors {
		result[i] = []float32(vec)
	}
	return result
}

// zeroVectors returns zero-filled vectors for stub mode.
func zeroVectors(n, dims int) [][]float32 {
	result := make([][]float32, n)
	for i := range result {
		result[i] = make([]float32, dims)
	}
	return result
}
