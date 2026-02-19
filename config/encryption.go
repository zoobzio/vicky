package config

import "github.com/zoobzio/check"

// Encryption holds encryption configuration for data at rest.
type Encryption struct {
	Key string `env:"VICKY_ENCRYPTION_KEY" secret:"vicky/encryption-key" required:"true"`
}

// Validate checks Encryption configuration for required values.
func (c Encryption) Validate() error {
	return check.Str(c.Key, "key").Required().Len(64).V().Err()
}
