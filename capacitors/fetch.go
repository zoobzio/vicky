package capacitors

import (
	"context"
	"log"
	"time"

	"github.com/zoobzio/check"
	"github.com/zoobzio/flux"
	"github.com/zoobzio/vicky/internal/ingest"
)

// Fetch holds operational settings for the fetch stage.
// Hot-reloadable via flux.
type Fetch struct {
	Workers int           `json:"workers"` // pool concurrency
	Timeout time.Duration `json:"timeout"` // per-file timeout
}

// Validate checks Fetch configuration.
// Zero values are allowed and mean "use default".
func (c Fetch) Validate() error {
	return check.All(
		check.NonNegative(c.Workers, "workers"),
		check.Max(c.Workers, 100, "workers"),
		check.DurationNonNegative(c.Timeout, "timeout"),
		check.DurationMax(c.Timeout, 10*time.Minute, "timeout"),
	).Err()
}

// DefaultFetch returns Fetch configuration with sensible defaults.
func DefaultFetch() Fetch {
	return Fetch{
		Workers: 8,
		Timeout: 30 * time.Second,
	}
}

// applyFetch applies config to the fetch pool.
func applyFetch(cfg Fetch) {
	ingest.SetFetchConfig(cfg.Workers, cfg.Timeout)
}

// InitFetch initializes the fetch capacitor with the given watcher.
func InitFetch(ctx context.Context, watcher flux.Watcher) error {
	// Apply defaults
	applyFetch(DefaultFetch())

	c := flux.New[Fetch](
		watcher,
		func(_ context.Context, _, curr Fetch) error {
			applyFetch(curr)
			return nil
		},
	)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("fetch capacitor error: %v", err)
		}
	}()
	return nil
}
