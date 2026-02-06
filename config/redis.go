package config

// Redis holds Redis connection settings.
type Redis struct {
	Addr     string `env:"VICKY_REDIS_ADDR" default:"localhost:6379"`
	Password string `env:"VICKY_REDIS_PASSWORD" secret:"vicky/redis-password"`
	DB       int    `env:"VICKY_REDIS_DB" default:"0"`
}

// Validate checks Redis configuration for required values.
func (c Redis) Validate() error {
	return nil
}
