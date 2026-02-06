package capacitors

import (
	"context"
	"sync/atomic"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/flux"
)

// Events is an alias for capitan.Config.
// Hot-reloadable via flux.
type Events = capitan.Config

// SignalConfig is an alias for capitan.SignalConfig.
type SignalConfig = capitan.SignalConfig

// DefaultEvents returns Events configuration with sensible defaults.
func DefaultEvents() Events {
	return Events{
		Signals: make(map[string]SignalConfig),
	}
}

// eventsCapacitor holds the flux capacitor for events config.
var eventsCapacitor atomic.Pointer[flux.Capacitor[Events]]

// eventsCurrent holds the current events config for fast access.
var eventsCurrent atomic.Pointer[Events]

// EventsCurrent returns the current events configuration.
// Returns default if no config has been loaded yet.
func EventsCurrent() Events {
	if cfg := eventsCurrent.Load(); cfg != nil {
		return *cfg
	}
	return DefaultEvents()
}

// InitEvents initializes the events capacitor with the given watcher.
// Automatically applies config changes to capitan.
func InitEvents(ctx context.Context, watcher flux.Watcher) error {
	// Set default initially
	def := DefaultEvents()
	eventsCurrent.Store(&def)

	c := flux.New[Events](
		watcher,
		func(_ context.Context, _, curr Events) error {
			// Store new config
			eventsCurrent.Store(&curr)
			// Apply to capitan
			return capitan.ApplyConfig(curr)
		},
	)

	eventsCapacitor.Store(c)

	return c.Start(ctx)
}
