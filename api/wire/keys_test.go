//go:build testing

package wire

import (
	"testing"
	"time"
)

// --- CreateKeyRequest.Validate ---

func TestCreateKeyRequest_Validate_Success(t *testing.T) {
	r := &CreateKeyRequest{
		Name:   "Production",
		Scopes: []string{"search"},
	}
	if err := r.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreateKeyRequest_Validate_MissingName(t *testing.T) {
	r := &CreateKeyRequest{
		Name:   "",
		Scopes: []string{"search"},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing name, got nil")
	}
}

func TestCreateKeyRequest_Validate_EmptyScopes(t *testing.T) {
	r := &CreateKeyRequest{
		Name:   "Production",
		Scopes: []string{},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for empty scopes, got nil")
	}
}

func TestCreateKeyRequest_Validate_NilScopes(t *testing.T) {
	r := &CreateKeyRequest{
		Name:   "Production",
		Scopes: nil,
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for nil scopes, got nil")
	}
}

func TestCreateKeyRequest_Validate_ZeroRateLimit(t *testing.T) {
	zero := 0
	r := &CreateKeyRequest{
		Name:      "Production",
		Scopes:    []string{"search"},
		RateLimit: &zero,
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for zero rate_limit, got nil")
	}
}

func TestCreateKeyRequest_Validate_NegativeRateLimit(t *testing.T) {
	neg := -10
	r := &CreateKeyRequest{
		Name:      "Production",
		Scopes:    []string{"search"},
		RateLimit: &neg,
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for negative rate_limit, got nil")
	}
}

func TestCreateKeyRequest_Validate_NilRateLimit(t *testing.T) {
	r := &CreateKeyRequest{
		Name:      "Production",
		Scopes:    []string{"search"},
		RateLimit: nil,
	}
	if err := r.Validate(); err != nil {
		t.Errorf("unexpected error for nil rate_limit: %v", err)
	}
}

// --- CreateKeyRequest.Clone ---

func TestCreateKeyRequest_Clone_Scopes(t *testing.T) {
	r := CreateKeyRequest{
		Name:   "Production",
		Scopes: []string{"search", "intel"},
	}
	c := r.Clone()
	c.Scopes[0] = "mutated"

	if r.Scopes[0] != "search" {
		t.Errorf("original Scopes[0] mutated: got %q", r.Scopes[0])
	}
}

func TestCreateKeyRequest_Clone_RateLimit(t *testing.T) {
	rl := 60
	r := CreateKeyRequest{
		Name:      "Production",
		Scopes:    []string{"search"},
		RateLimit: &rl,
	}
	c := r.Clone()
	*c.RateLimit = 999

	if *r.RateLimit != 60 {
		t.Errorf("original RateLimit mutated: got %d", *r.RateLimit)
	}
}

func TestCreateKeyRequest_Clone_ExpiresAt(t *testing.T) {
	exp := time.Now().Add(24 * time.Hour)
	r := CreateKeyRequest{
		Name:      "Production",
		Scopes:    []string{"search"},
		ExpiresAt: &exp,
	}
	c := r.Clone()
	*c.ExpiresAt = time.Time{}

	if !r.ExpiresAt.Equal(exp) {
		t.Errorf("original ExpiresAt mutated")
	}
}

// --- KeyResponse.Clone ---

func TestKeyResponse_Clone_AllPointers(t *testing.T) {
	rl := 60
	exp := time.Now().Add(24 * time.Hour)
	used := time.Now()

	r := KeyResponse{
		ID:         1,
		Name:       "Production",
		Prefix:     "vky_abc1",
		Scopes:     []string{"search"},
		RateLimit:  &rl,
		ExpiresAt:  &exp,
		LastUsedAt: &used,
	}
	c := r.Clone()

	*c.RateLimit = 999
	*c.ExpiresAt = time.Time{}
	*c.LastUsedAt = time.Time{}
	c.Scopes[0] = "mutated"

	if *r.RateLimit != 60 {
		t.Errorf("original RateLimit mutated: got %d", *r.RateLimit)
	}
	if !r.ExpiresAt.Equal(exp) {
		t.Errorf("original ExpiresAt mutated")
	}
	if !r.LastUsedAt.Equal(used) {
		t.Errorf("original LastUsedAt mutated")
	}
	if r.Scopes[0] != "search" {
		t.Errorf("original Scopes[0] mutated: got %q", r.Scopes[0])
	}
}

// --- KeyCreatedResponse.Clone ---

func TestKeyCreatedResponse_Clone(t *testing.T) {
	r := KeyCreatedResponse{
		KeyResponse: KeyResponse{
			ID:     1,
			Name:   "Production",
			Prefix: "vky_abc1",
			Scopes: []string{"search"},
		},
		Key: "vky_rawsecretkey",
	}
	c := r.Clone()

	if c.Key != "vky_rawsecretkey" {
		t.Errorf("Key = %q, want %q", c.Key, "vky_rawsecretkey")
	}
	c.Scopes[0] = "mutated"
	if r.Scopes[0] != "search" {
		t.Errorf("original Scopes[0] mutated after KeyCreatedResponse clone")
	}
}

// --- KeyListResponse.Clone ---

func TestKeyListResponse_Clone(t *testing.T) {
	r := KeyListResponse{
		Keys: []KeyResponse{
			{ID: 1, Name: "A", Scopes: []string{"search"}},
			{ID: 2, Name: "B", Scopes: []string{"intel"}},
		},
	}
	c := r.Clone()

	if len(c.Keys) != 2 {
		t.Fatalf("Keys len = %d, want 2", len(c.Keys))
	}
	c.Keys[0].Name = "mutated"
	if r.Keys[0].Name != "A" {
		t.Errorf("original Keys[0].Name mutated")
	}
}

func TestKeyListResponse_Clone_NilKeys(t *testing.T) {
	r := KeyListResponse{Keys: nil}
	c := r.Clone()
	if c.Keys != nil {
		t.Errorf("expected nil Keys in clone, got %v", c.Keys)
	}
}
