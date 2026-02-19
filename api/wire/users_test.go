package wire

import (
	"strings"
	"testing"
)

func TestUserUpdateRequestValidate_Valid(t *testing.T) {
	name := "Valid Name"
	req := &UserUpdateRequest{Name: &name}
	if err := req.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserUpdateRequestValidate_NilName(t *testing.T) {
	req := &UserUpdateRequest{Name: nil}
	if err := req.Validate(); err != nil {
		t.Errorf("unexpected error for nil name: %v", err)
	}
}

func TestUserUpdateRequestValidate_TooLong(t *testing.T) {
	name := strings.Repeat("x", 256)
	req := &UserUpdateRequest{Name: &name}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for name > 255 chars, got nil")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "name")
	}
}
