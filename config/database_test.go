package config

import (
	"strings"
	"testing"
)

func TestDatabaseValidate_Valid(t *testing.T) {
	c := Database{Host: "localhost", Name: "vicky"}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDatabaseValidate_EmptyHost(t *testing.T) {
	c := Database{Host: "", Name: "vicky"}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for empty host, got nil")
	}
	if !strings.Contains(err.Error(), "host") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "host")
	}
}

func TestDatabaseValidate_EmptyName(t *testing.T) {
	c := Database{Host: "localhost", Name: ""}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "name")
	}
}

func TestDatabaseDSN(t *testing.T) {
	c := Database{
		Host:     "db.example.com",
		Port:     5432,
		Name:     "mydb",
		User:     "admin",
		Password: "secret",
		SSLMode:  "require",
	}
	dsn := c.DSN()

	for _, want := range []string{
		"host=db.example.com",
		"port=5432",
		"dbname=mydb",
		"user=admin",
		"password=secret",
		"sslmode=require",
	} {
		if !strings.Contains(dsn, want) {
			t.Errorf("DSN = %q, want it to contain %q", dsn, want)
		}
	}
}
