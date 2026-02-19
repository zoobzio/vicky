package config

import (
	"strconv"

	"github.com/zoobzio/check"
)

// Database holds PostgreSQL connection settings.
type Database struct {
	Host     string `env:"VICKY_DB_HOST" default:"localhost"`
	Port     int    `env:"VICKY_DB_PORT" default:"5432"`
	Name     string `env:"VICKY_DB_NAME" default:"vicky"`
	User     string `env:"VICKY_DB_USER" default:"vicky"`
	Password string `env:"VICKY_DB_PASSWORD" secret:"vicky/db-password"`
	SSLMode  string `env:"VICKY_DB_SSLMODE" default:"disable"`
}

// Validate checks Database configuration for required values.
func (c Database) Validate() error {
	return check.All(
		check.Str(c.Host, "host").Required().V(),
		check.Str(c.Name, "name").Required().V(),
	).Err()
}

// DSN returns the PostgreSQL connection string.
func (c Database) DSN() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" dbname=" + c.Name +
		" user=" + c.User +
		" password=" + c.Password +
		" sslmode=" + c.SSLMode
}
