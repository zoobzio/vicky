package capacitors

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/zoobzio/aperture"
	"github.com/zoobzio/flux"
)

// Observability is an alias for aperture.Schema.
// Hot-reloadable via flux.
type Observability = aperture.Schema

// DefaultObservability returns Observability configuration with sensible defaults.
// Metrics and traces are configured via database and hot-reloaded.
func DefaultObservability() Observability {
	return Observability{}
}

// observabilityCapacitor holds the flux capacitor for observability config.
var observabilityCapacitor atomic.Pointer[flux.Capacitor[Observability]]

// observabilityCurrent holds the current observability config for fast access.
var observabilityCurrent atomic.Pointer[Observability]

// ObservabilityCurrent returns the current observability configuration.
// Returns default if no config has been loaded yet.
func ObservabilityCurrent() Observability {
	if cfg := observabilityCurrent.Load(); cfg != nil {
		return *cfg
	}
	return DefaultObservability()
}

// InitObservability initializes the observability capacitor with the given watcher.
// Automatically applies config changes to aperture.
func InitObservability(ctx context.Context, watcher flux.Watcher, ap *aperture.Aperture) error {
	def := DefaultObservability()
	observabilityCurrent.Store(&def)

	if err := ap.Apply(def); err != nil {
		return err
	}

	c := flux.New[Observability](
		watcher,
		func(_ context.Context, _, curr Observability) error {
			observabilityCurrent.Store(&curr)
			return ap.Apply(curr)
		},
	)

	observabilityCapacitor.Store(c)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("observability capacitor error: %v", err)
		}
	}()
	return nil
}
