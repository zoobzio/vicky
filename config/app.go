// Package config provides configuration loading for vicky.
package config

import "github.com/zoobzio/check"

// App holds vicky-specific configuration.
type App struct {
	Port               int    `env:"VICKY_PORT" default:"8080"`
	SessionSignKey     string `env:"SESSION_SIGN_KEY" secret:"vicky/session-sign-key"`
	SessionRedirectURL string `env:"SESSION_REDIRECT_URL" default:"http://localhost:3000"`
}

// Validate checks App configuration for required values.
func (c App) Validate() error {
	return check.Int(c.Port, "port").Positive().V().Err()
}
