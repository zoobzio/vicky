package config

import (
	"errors"
	"strings"
)

// Admin holds vicky-admin specific configuration.
type Admin struct {
	Port int `env:"VICKY_ADMIN_PORT" default:"8081"`

	// AllowedUsers is a comma-separated list of GitHub logins authorized for admin access.
	AllowedUsers string `env:"VICKY_ADMIN_ALLOWED_USERS"`

	// SessionSignKey is the key used to sign session cookies.
	SessionSignKey string `env:"VICKY_ADMIN_SESSION_SIGN_KEY" secret:"vicky-admin/session-sign-key"`
}

// Validate checks Admin configuration for required values.
func (c Admin) Validate() error {
	if c.AllowedUsers == "" {
		return errors.New("VICKY_ADMIN_ALLOWED_USERS is required")
	}
	return nil
}

// AllowedUsersList returns the parsed list of allowed GitHub logins.
func (c Admin) AllowedUsersList() []string {
	if c.AllowedUsers == "" {
		return nil
	}
	logins := strings.Split(c.AllowedUsers, ",")
	for i := range logins {
		logins[i] = strings.TrimSpace(logins[i])
	}
	return logins
}

// IsUserAllowed checks if a GitHub login is in the allowed list.
func (c Admin) IsUserAllowed(login string) bool {
	for _, allowed := range c.AllowedUsersList() {
		if strings.EqualFold(allowed, login) {
			return true
		}
	}
	return false
}
