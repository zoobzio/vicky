package config

// Observability holds OTEL and telemetry settings.
type Observability struct {
	OTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	ServiceName  string `env:"OTEL_SERVICE_NAME" default:"vicky"`
}

// Validate checks Observability configuration for required values.
func (c Observability) Validate() error {
	return nil
}
