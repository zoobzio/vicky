package capacitors

import (
	"context"
	"log"
	"time"

	"github.com/zoobzio/check"
	"github.com/zoobzio/flux"
	"github.com/zoobzio/vicky/internal/ingest"
)

// Parse holds operational settings for the parse stage.
// Hot-reloadable via flux.
type Parse struct {
	Workers int           `json:"workers"` // pool concurrency
	Timeout time.Duration `json:"timeout"` // per-document timeout
}

// Validate checks Parse configuration.
// Zero values are allowed and mean "use default".
func (c Parse) Validate() error {
	return check.All(
		check.NonNegative(c.Workers, "workers"),
		check.Max(c.Workers, 100, "workers"),
		check.DurationNonNegative(c.Timeout, "timeout"),
		check.DurationMax(c.Timeout, 30*time.Minute, "timeout"),
	).Err()
}

// DefaultParse returns Parse configuration with sensible defaults.
func DefaultParse() Parse {
	return Parse{
		Workers: 8,
		Timeout: 60 * time.Second,
	}
}

// applyParse applies config to the parse pool.
func applyParse(cfg Parse) {
	ingest.SetParseConfig(cfg.Workers, cfg.Timeout)
}

// InitParse initializes the parse capacitor with the given watcher.
func InitParse(ctx context.Context, watcher flux.Watcher) error {
	// Apply defaults
	applyParse(DefaultParse())

	c := flux.New[Parse](
		watcher,
		func(_ context.Context, _, curr Parse) error {
			applyParse(curr)
			return nil
		},
	)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("parse capacitor error: %v", err)
		}
	}()
	return nil
}
