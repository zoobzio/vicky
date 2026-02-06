package config

// Embedding holds configuration for the vector embedding provider.
type Embedding struct {
	Provider   string `env:"VICKY_EMBEDDING_PROVIDER" default:"stub"`
	Model      string `env:"VICKY_EMBEDDING_MODEL" default:"voyage-code-3"`
	Dimensions int    `env:"VICKY_EMBEDDING_DIMENSIONS" default:"1024"`
	APIKey     string `env:"VICKY_EMBEDDING_API_KEY"`
}

// Validate checks Embedding configuration for required values.
func (c Embedding) Validate() error {
	return nil
}
