// Package config provides configuration loading for vicky.
package config

import "errors"

// App holds vicky-specific configuration.
type App struct {
	Port           int    `env:"VICKY_PORT" default:"8080"`
	SessionSignKey string `env:"SESSION_SIGN_KEY" secret:"vicky/session-sign-key"`
	GitHub         GitHub
}

// GitHub holds GitHub authentication and API settings.
type GitHub struct {
	ClientID     string `env:"GITHUB_CLIENT_ID"`
	ClientSecret string `env:"GITHUB_CLIENT_SECRET" secret:"vicky/github-client-secret"`
	RedirectURI  string `env:"GITHUB_REDIRECT_URI"`
}

// Validate checks App configuration for required values.
func (c App) Validate() error {
	if c.Port <= 0 {
		return errors.New("port must be positive")
	}
	return nil
}
