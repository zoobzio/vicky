package config

import "testing"

func TestAppValidate_Valid(t *testing.T) {
	c := App{Port: 8080}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAppValidate_ZeroPort(t *testing.T) {
	c := App{Port: 0}
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero port, got nil")
	}
}

func TestAppValidate_NegativePort(t *testing.T) {
	c := App{Port: -1}
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative port, got nil")
	}
}
