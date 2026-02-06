package config

import "errors"

// Encryption holds encryption configuration for data at rest.
type Encryption struct {
	Key string `env:"VICKY_ENCRYPTION_KEY" secret:"vicky/encryption-key" required:"true"`
}

// Validate checks Encryption configuration for required values.
func (c Encryption) Validate() error {
	if c.Key == "" {
		return errors.New("encryption key is required")
	}
	if len(c.Key) != 64 {
		return errors.New("encryption key must be 64 hex characters (32 bytes for AES-256)")
	}
	return nil
}
