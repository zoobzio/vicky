package transformers

import (
	"testing"
	"time"

	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/models"
)

func newTestKey() *models.Key {
	return &models.Key{
		ID:        1,
		UserID:    1000,
		Name:      "Production",
		KeyHash:   "shouldnotappear",
		KeyPrefix: "vky_abc1",
		Scopes:    []string{"search", "intel"},
	}
}

// --- KeyToResponse ---

func TestKeyToResponse(t *testing.T) {
	k := newTestKey()
	resp := KeyToResponse(k)

	if resp.ID != k.ID {
		t.Errorf("ID = %d, want %d", resp.ID, k.ID)
	}
	if resp.Name != k.Name {
		t.Errorf("Name = %q, want %q", resp.Name, k.Name)
	}
	if resp.Prefix != k.KeyPrefix {
		t.Errorf("Prefix = %q, want %q", resp.Prefix, k.KeyPrefix)
	}
	if len(resp.Scopes) != len(k.Scopes) {
		t.Errorf("Scopes len = %d, want %d", len(resp.Scopes), len(k.Scopes))
	}
	if resp.RateLimit != nil {
		t.Error("RateLimit should be nil")
	}
	if resp.ExpiresAt != nil {
		t.Error("ExpiresAt should be nil")
	}
	if resp.LastUsedAt != nil {
		t.Error("LastUsedAt should be nil")
	}
}

func TestKeyToResponse_WithOptionalFields(t *testing.T) {
	rate := 60
	exp := time.Now().Add(24 * time.Hour)
	used := time.Now()

	k := newTestKey()
	k.RateLimit = &rate
	k.ExpiresAt = &exp
	k.LastUsedAt = &used

	resp := KeyToResponse(k)

	if resp.RateLimit == nil || *resp.RateLimit != rate {
		t.Errorf("RateLimit = %v, want %d", resp.RateLimit, rate)
	}
	if resp.ExpiresAt == nil || !resp.ExpiresAt.Equal(exp) {
		t.Errorf("ExpiresAt = %v, want %v", resp.ExpiresAt, exp)
	}
	if resp.LastUsedAt == nil || !resp.LastUsedAt.Equal(used) {
		t.Errorf("LastUsedAt = %v, want %v", resp.LastUsedAt, used)
	}
}

// --- KeyToCreatedResponse ---

func TestKeyToCreatedResponse(t *testing.T) {
	k := newTestKey()
	rawKey := "vky_rawsecretkey"

	resp := KeyToCreatedResponse(k, rawKey)

	if resp.Key != rawKey {
		t.Errorf("Key = %q, want %q", resp.Key, rawKey)
	}
	if resp.ID != k.ID {
		t.Errorf("ID = %d, want %d", resp.ID, k.ID)
	}
	if resp.Name != k.Name {
		t.Errorf("Name = %q, want %q", resp.Name, k.Name)
	}
	if resp.Prefix != k.KeyPrefix {
		t.Errorf("Prefix = %q, want %q", resp.Prefix, k.KeyPrefix)
	}
}

// --- KeysToList ---

func TestKeysToList_Empty(t *testing.T) {
	resp := KeysToList([]*models.Key{})

	if resp.Keys == nil {
		t.Error("Keys should not be nil for empty input")
	}
	if len(resp.Keys) != 0 {
		t.Errorf("Keys len = %d, want 0", len(resp.Keys))
	}
}

func TestKeysToList_Multiple(t *testing.T) {
	keys := []*models.Key{
		{ID: 1, Name: "Key A", KeyPrefix: "vky_aaa1", Scopes: []string{"search"}},
		{ID: 2, Name: "Key B", KeyPrefix: "vky_bbb2", Scopes: []string{"intel"}},
		{ID: 3, Name: "Key C", KeyPrefix: "vky_ccc3", Scopes: []string{"search", "intel"}},
	}

	resp := KeysToList(keys)

	if len(resp.Keys) != 3 {
		t.Fatalf("Keys len = %d, want 3", len(resp.Keys))
	}
	for i, k := range keys {
		if resp.Keys[i].ID != k.ID {
			t.Errorf("Keys[%d].ID = %d, want %d", i, resp.Keys[i].ID, k.ID)
		}
		if resp.Keys[i].Name != k.Name {
			t.Errorf("Keys[%d].Name = %q, want %q", i, resp.Keys[i].Name, k.Name)
		}
	}
}

// --- ApplyCreateKeyRequest ---

func TestApplyCreateKeyRequest(t *testing.T) {
	rate := 120
	exp := time.Now().Add(7 * 24 * time.Hour)
	req := wire.CreateKeyRequest{
		Name:      "CI Key",
		Scopes:    []string{"search"},
		RateLimit: &rate,
		ExpiresAt: &exp,
	}

	k := &models.Key{
		ID:        99,
		UserID:    1000,
		KeyHash:   "existinghash",
		KeyPrefix: "vky_exist",
	}

	ApplyCreateKeyRequest(req, k)

	if k.Name != "CI Key" {
		t.Errorf("Name = %q, want %q", k.Name, "CI Key")
	}
	if len(k.Scopes) != 1 || k.Scopes[0] != "search" {
		t.Errorf("Scopes = %v, want [search]", k.Scopes)
	}
	if k.RateLimit == nil || *k.RateLimit != rate {
		t.Errorf("RateLimit = %v, want %d", k.RateLimit, rate)
	}
	if k.ExpiresAt == nil || !k.ExpiresAt.Equal(exp) {
		t.Errorf("ExpiresAt = %v, want %v", k.ExpiresAt, exp)
	}
	// Fields not managed by Apply must remain untouched.
	if k.ID != 99 {
		t.Errorf("ID should not be changed, got %d", k.ID)
	}
	if k.UserID != 1000 {
		t.Errorf("UserID should not be changed, got %d", k.UserID)
	}
	if k.KeyHash != "existinghash" {
		t.Errorf("KeyHash should not be changed, got %q", k.KeyHash)
	}
}

func TestApplyCreateKeyRequest_NilOptionals(t *testing.T) {
	req := wire.CreateKeyRequest{
		Name:   "Minimal",
		Scopes: []string{"search"},
	}
	k := &models.Key{}

	ApplyCreateKeyRequest(req, k)

	if k.RateLimit != nil {
		t.Errorf("RateLimit should be nil, got %v", k.RateLimit)
	}
	if k.ExpiresAt != nil {
		t.Errorf("ExpiresAt should be nil, got %v", k.ExpiresAt)
	}
}
