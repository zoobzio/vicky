package config

import (
	"strings"
	"testing"
)

func TestEncryptionValidate_Valid(t *testing.T) {
	c := Encryption{Key: strings.Repeat("ab", 32)} // 64 hex chars
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEncryptionValidate_Empty(t *testing.T) {
	c := Encryption{Key: ""}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "required")
	}
}

func TestEncryptionValidate_WrongLength(t *testing.T) {
	c := Encryption{Key: "tooshort"}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for wrong length, got nil")
	}
	if !strings.Contains(err.Error(), "64") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "64")
	}
}
