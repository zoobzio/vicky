package config

import (
	"testing"

	"github.com/zoobzio/vicky/models"
)

func TestAddresses_Both(t *testing.T) {
	c := Indexer{GoAddr: "localhost:9090", TsAddr: "localhost:9091"}
	addrs := c.Addresses()

	if got, ok := addrs[models.LanguageGo]; !ok || got != "localhost:9090" {
		t.Errorf("Go address = %q (present=%v), want %q", got, ok, "localhost:9090")
	}
	if got, ok := addrs[models.LanguageTypeScript]; !ok || got != "localhost:9091" {
		t.Errorf("TypeScript address = %q (present=%v), want %q", got, ok, "localhost:9091")
	}
	if len(addrs) != 2 {
		t.Errorf("len(addrs) = %d, want 2", len(addrs))
	}
}

func TestAddresses_GoOnly(t *testing.T) {
	c := Indexer{GoAddr: "localhost:9090"}
	addrs := c.Addresses()

	if _, ok := addrs[models.LanguageTypeScript]; ok {
		t.Error("TypeScript should not be present")
	}
	if len(addrs) != 1 {
		t.Errorf("len(addrs) = %d, want 1", len(addrs))
	}
}

func TestAddresses_Neither(t *testing.T) {
	c := Indexer{}
	addrs := c.Addresses()

	if len(addrs) != 0 {
		t.Errorf("len(addrs) = %d, want 0", len(addrs))
	}
}
