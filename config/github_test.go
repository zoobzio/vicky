package config

import (
	"strings"
	"testing"
)

func TestGitHubValidate_Valid(t *testing.T) {
	c := GitHub{
		ClientID:     "abc123",
		ClientSecret: "secret",
		RedirectURI:  "http://localhost/callback",
	}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGitHubValidate_EmptyClientID(t *testing.T) {
	c := GitHub{
		ClientID:     "",
		ClientSecret: "secret",
		RedirectURI:  "http://localhost/callback",
	}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for empty client_id, got nil")
	}
	if !strings.Contains(err.Error(), "client_id") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "client_id")
	}
}

func TestGitHubValidate_EmptyClientSecret(t *testing.T) {
	c := GitHub{
		ClientID:     "abc123",
		ClientSecret: "",
		RedirectURI:  "http://localhost/callback",
	}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for empty client_secret, got nil")
	}
	if !strings.Contains(err.Error(), "client_secret") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "client_secret")
	}
}

func TestGitHubValidate_EmptyRedirectURI(t *testing.T) {
	c := GitHub{
		ClientID:     "abc123",
		ClientSecret: "secret",
		RedirectURI:  "",
	}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for empty redirect_uri, got nil")
	}
	if !strings.Contains(err.Error(), "redirect_uri") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "redirect_uri")
	}
}
