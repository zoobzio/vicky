package capacitors

import (
	"context"
	"log"
	"time"

	"github.com/zoobzio/check"
	"github.com/zoobzio/flux"
	"github.com/zoobzio/vicky/internal/ingest"
)

// Chunk holds operational settings for the chunk stage.
// Hot-reloadable via flux.
type Chunk struct {
	Workers int           `json:"workers"` // pool concurrency
	Timeout time.Duration `json:"timeout"` // per-document timeout
}

// Validate checks Chunk configuration.
// Zero values are allowed and mean "use default".
func (c Chunk) Validate() error {
	return check.All(
		check.NonNegative(c.Workers, "workers"),
		check.Max(c.Workers, 100, "workers"),
		check.DurationNonNegative(c.Timeout, "timeout"),
		check.DurationMax(c.Timeout, 10*time.Minute, "timeout"),
	).Err()
}

// DefaultChunk returns Chunk configuration with sensible defaults.
func DefaultChunk() Chunk {
	return Chunk{
		Workers: 8,
		Timeout: 30 * time.Second,
	}
}

// applyChunk applies config to the chunk pool.
func applyChunk(cfg Chunk) {
	ingest.SetChunkConfig(cfg.Workers, cfg.Timeout)
}

// InitChunk initializes the chunk capacitor with the given watcher.
func InitChunk(ctx context.Context, watcher flux.Watcher) error {
	// Apply defaults
	applyChunk(DefaultChunk())

	c := flux.New[Chunk](
		watcher,
		func(_ context.Context, _, curr Chunk) error {
			applyChunk(curr)
			return nil
		},
	)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("chunk capacitor error: %v", err)
		}
	}()
	return nil
}
