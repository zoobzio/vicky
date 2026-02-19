package models

import (
	"testing"
	"time"
)

// --- Clone ---

func TestKeyClone_NilPointers(t *testing.T) {
	k := Key{
		ID:        1,
		UserID:    1000,
		Name:      "Test Key",
		KeyHash:   "hash",
		KeyPrefix: "vky_test",
		Scopes:    []string{"search"},
	}

	c := k.Clone()

	if c.ID != k.ID {
		t.Errorf("ID = %d, want %d", c.ID, k.ID)
	}
	if c.RateLimit != nil {
		t.Error("RateLimit should be nil")
	}
	if c.ExpiresAt != nil {
		t.Error("ExpiresAt should be nil")
	}
	if c.LastUsedAt != nil {
		t.Error("LastUsedAt should be nil")
	}
}

func TestKeyClone_AllPointers(t *testing.T) {
	rate := 60
	expires := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	lastUsed := time.Now().Truncate(time.Second)

	k := Key{
		ID:         1,
		UserID:     1000,
		Name:       "Test Key",
		KeyHash:    "hash",
		KeyPrefix:  "vky_test",
		Scopes:     []string{"search", "intel"},
		RateLimit:  &rate,
		ExpiresAt:  &expires,
		LastUsedAt: &lastUsed,
	}

	c := k.Clone()

	if c.RateLimit == nil || *c.RateLimit != rate {
		t.Errorf("RateLimit = %v, want %d", c.RateLimit, rate)
	}
	if c.ExpiresAt == nil || !c.ExpiresAt.Equal(expires) {
		t.Errorf("ExpiresAt = %v, want %v", c.ExpiresAt, expires)
	}
	if c.LastUsedAt == nil || !c.LastUsedAt.Equal(lastUsed) {
		t.Errorf("LastUsedAt = %v, want %v", c.LastUsedAt, lastUsed)
	}
	if len(c.Scopes) != len(k.Scopes) {
		t.Errorf("Scopes len = %d, want %d", len(c.Scopes), len(k.Scopes))
	}
}

func TestKeyClone_IndependentMutation(t *testing.T) {
	rate := 60
	expires := time.Now().Add(24 * time.Hour)
	lastUsed := time.Now()

	k := Key{
		ID:         1,
		UserID:     1000,
		Name:       "Test Key",
		KeyHash:    "hash",
		KeyPrefix:  "vky_test",
		Scopes:     []string{"search", "intel"},
		RateLimit:  &rate,
		ExpiresAt:  &expires,
		LastUsedAt: &lastUsed,
	}

	c := k.Clone()

	// Mutate the clone; original should be unaffected.
	*c.RateLimit = 999
	*c.ExpiresAt = time.Time{}
	*c.LastUsedAt = time.Time{}
	c.Scopes[0] = "mutated"

	if *k.RateLimit != 60 {
		t.Errorf("original RateLimit mutated: got %d", *k.RateLimit)
	}
	if !k.ExpiresAt.Equal(expires) {
		t.Errorf("original ExpiresAt mutated")
	}
	if !k.LastUsedAt.Equal(lastUsed) {
		t.Errorf("original LastUsedAt mutated")
	}
	if k.Scopes[0] != "search" {
		t.Errorf("original Scopes[0] mutated: got %q", k.Scopes[0])
	}
}

// --- Validate ---

func TestKeyValidate_Success(t *testing.T) {
	k := Key{
		Name:      "Production",
		KeyHash:   "abc123hash",
		KeyPrefix: "vky_abc1",
	}
	if err := k.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestKeyValidate_MissingName(t *testing.T) {
	k := Key{
		Name:      "",
		KeyHash:   "abc123hash",
		KeyPrefix: "vky_abc1",
	}
	if err := k.Validate(); err == nil {
		t.Error("expected error for missing name, got nil")
	}
}

func TestKeyValidate_MissingKeyHash(t *testing.T) {
	k := Key{
		Name:      "Production",
		KeyHash:   "",
		KeyPrefix: "vky_abc1",
	}
	if err := k.Validate(); err == nil {
		t.Error("expected error for missing key_hash, got nil")
	}
}

func TestKeyValidate_MissingKeyPrefix(t *testing.T) {
	k := Key{
		Name:      "Production",
		KeyHash:   "abc123hash",
		KeyPrefix: "",
	}
	if err := k.Validate(); err == nil {
		t.Error("expected error for missing key_prefix, got nil")
	}
}
