package otel

import (
	"context"
	"os"

	"github.com/zoobzio/aperture"
	"github.com/zoobzio/capitan"
)

// SidecarConfig holds sidecar observability configuration.
type SidecarConfig struct {
	ServiceName string
}

// Sidecar holds the observability components for a sidecar service.
type Sidecar struct {
	Providers *Providers
	Aperture  *aperture.Aperture
}

// NewSidecar creates observability for a sidecar service.
// Uses OTEL_EXPORTER_OTLP_ENDPOINT from environment, defaults to localhost:4318.
func NewSidecar(ctx context.Context, cfg SidecarConfig) (*Sidecar, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	providers, err := New(ctx, Config{
		Endpoint:    endpoint,
		ServiceName: cfg.ServiceName,
	})
	if err != nil {
		return nil, err
	}

	ap, err := aperture.New(
		capitan.Default(),
		providers.Log,
		providers.Metric,
		providers.Trace,
	)
	if err != nil {
		providers.Shutdown(ctx)
		return nil, err
	}

	// Apply empty schema - log all signals by default
	if err := ap.Apply(aperture.Schema{}); err != nil {
		ap.Close()
		providers.Shutdown(ctx)
		return nil, err
	}

	return &Sidecar{
		Providers: providers,
		Aperture:  ap,
	}, nil
}

// Shutdown gracefully shuts down observability.
func (s *Sidecar) Shutdown(ctx context.Context) error {
	s.Aperture.Close()
	return s.Providers.Shutdown(ctx)
}
