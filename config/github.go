package config

import "github.com/zoobzio/check"

// GitHub holds GitHub authentication and API settings.
type GitHub struct {
	ClientID     string `env:"GITHUB_CLIENT_ID"`
	ClientSecret string `env:"GITHUB_CLIENT_SECRET" secret:"vicky/github-client-secret"`
	RedirectURI  string `env:"GITHUB_REDIRECT_URI"`
}

// Validate checks GitHub configuration for required values.
func (c GitHub) Validate() error {
	return check.All(
		check.Str(c.ClientID, "client_id").Required().V(),
		check.Str(c.ClientSecret, "client_secret").Required().V(),
		check.Str(c.RedirectURI, "redirect_uri").Required().V(),
	).Err()
}
