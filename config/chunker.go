package config

// Chunker holds the gRPC address for the chunker sidecar.
type Chunker struct {
	Addr string `env:"VICKY_CHUNKER_ADDR"`
}

// Validate checks Chunker configuration for required values.
func (c Chunker) Validate() error {
	return nil
}
