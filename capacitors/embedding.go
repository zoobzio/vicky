package capacitors

import (
	"context"
	"log"
	"time"

	"github.com/zoobzio/check"
	"github.com/zoobzio/flux"
	"github.com/zoobzio/vicky/internal/ingest"
)

// Embedding holds operational settings for the embed stage.
// Hot-reloadable via flux.
type Embedding struct {
	Workers   int           `json:"workers"`    // pool concurrency
	BatchSize int           `json:"batch_size"` // texts per API call
	Timeout   time.Duration `json:"timeout"`    // per-batch timeout
}

// Validate checks Embedding configuration.
// Zero values are allowed and mean "use default".
func (c Embedding) Validate() error {
	return check.All(
		check.NonNegative(c.Workers, "workers"),
		check.Max(c.Workers, 100, "workers"),
		check.NonNegative(c.BatchSize, "batch_size"),
		check.Max(c.BatchSize, 1000, "batch_size"),
		check.DurationNonNegative(c.Timeout, "timeout"),
		check.DurationMax(c.Timeout, 10*time.Minute, "timeout"),
	).Err()
}

// DefaultEmbedding returns Embedding configuration with sensible defaults.
func DefaultEmbedding() Embedding {
	return Embedding{
		Workers:   4,
		BatchSize: 128,
		Timeout:   30 * time.Second,
	}
}

// applyEmbedding applies config to the embed pool.
func applyEmbedding(cfg Embedding) {
	ingest.SetEmbedConfig(cfg.Workers, cfg.BatchSize, cfg.Timeout)
}

// InitEmbedding initializes the embedding capacitor with the given watcher.
func InitEmbedding(ctx context.Context, watcher flux.Watcher) error {
	// Apply defaults
	applyEmbedding(DefaultEmbedding())

	c := flux.New[Embedding](
		watcher,
		func(_ context.Context, _, curr Embedding) error {
			applyEmbedding(curr)
			return nil
		},
	)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("embedding capacitor error: %v", err)
		}
	}()
	return nil
}
