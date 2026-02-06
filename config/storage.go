package config

// Storage holds S3-compatible blob storage settings.
type Storage struct {
	Endpoint  string `env:"VICKY_STORAGE_ENDPOINT" default:"localhost:9000"`
	AccessKey string `env:"VICKY_STORAGE_ACCESS_KEY" secret:"vicky/storage-access-key"`
	SecretKey string `env:"VICKY_STORAGE_SECRET_KEY" secret:"vicky/storage-secret-key"`
	Bucket    string `env:"VICKY_STORAGE_BUCKET" default:"vicky"`
	UseSSL    bool   `env:"VICKY_STORAGE_USE_SSL" default:"false"`
}

// Validate checks Storage configuration for required values.
func (c Storage) Validate() error {
	return nil
}
