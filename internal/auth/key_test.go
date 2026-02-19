package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/zoobzio/rocco"
	vickytest "github.com/zoobzio/vicky/testing"
	"github.com/zoobzio/vicky/models"
)

// rawTestKey is a valid vky_ prefixed key used across tests.
const rawTestKey = "vky_dGVzdGtleWZvcnRlc3Rpbmcx"

// testKeyHash returns the SHA-256 base64 hash of rawTestKey.
func testKeyHash() string {
	return hashKey(rawTestKey)
}

func bearerHeader(key string) *http.Request {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+key)
	return r
}

func noAuthHeader() *http.Request {
	r, _ := http.NewRequest("GET", "/", nil)
	return r
}

// --- hashKey ---

func TestHashKey_StableOutput(t *testing.T) {
	h1 := hashKey(rawTestKey)
	h2 := hashKey(rawTestKey)
	if h1 != h2 {
		t.Errorf("hashKey not stable: %q != %q", h1, h2)
	}
}

func TestHashKey_Base64Encoded(t *testing.T) {
	h := hashKey(rawTestKey)
	_, err := base64.StdEncoding.DecodeString(h)
	if err != nil {
		t.Errorf("hashKey output is not valid base64: %v", err)
	}
}

func TestHashKey_MatchesSHA256(t *testing.T) {
	sum := sha256.Sum256([]byte(rawTestKey))
	want := base64.StdEncoding.EncodeToString(sum[:])
	got := hashKey(rawTestKey)
	if got != want {
		t.Errorf("hashKey = %q, want %q", got, want)
	}
}

// --- KeyExtractor: Bearer vky_ token ---

func TestKeyExtractor_ValidKey(t *testing.T) {
	hash := testKeyHash()
	mk := &vickytest.MockKeys{
		OnGetByKeyHash: func(ctx context.Context, h string) (*models.Key, error) {
			if h != hash {
				return nil, fmt.Errorf("wrong hash")
			}
			return &models.Key{
				ID:     1,
				UserID: 1000,
				Scopes: []string{"search", "intel"},
			}, nil
		},
	}

	extractor := KeyExtractor(mk, func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		t.Error("session extractor should not be called")
		return nil, nil
	})

	identity, err := extractor(context.Background(), bearerHeader(rawTestKey))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity.ID() != "1000" {
		t.Errorf("ID = %q, want %q", identity.ID(), "1000")
	}
}

func TestKeyExtractor_Identity_Scopes(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnGetByKeyHash: func(ctx context.Context, hash string) (*models.Key, error) {
			return &models.Key{
				ID:     1,
				UserID: 42,
				Scopes: []string{"search", "intel"},
			}, nil
		},
	}

	extractor := KeyExtractor(mk, nil)
	identity, err := extractor(context.Background(), bearerHeader(rawTestKey))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	scopes := identity.Scopes()
	if len(scopes) != 2 {
		t.Fatalf("Scopes len = %d, want 2", len(scopes))
	}
	if scopes[0] != "search" || scopes[1] != "intel" {
		t.Errorf("Scopes = %v, want [search intel]", scopes)
	}
}

func TestKeyExtractor_KeyNotFound(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnGetByKeyHash: func(ctx context.Context, hash string) (*models.Key, error) {
			return nil, fmt.Errorf("not found")
		},
	}

	extractor := KeyExtractor(mk, func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		t.Error("session extractor should not be called for vky_ tokens")
		return nil, nil
	})

	_, err := extractor(context.Background(), bearerHeader(rawTestKey))
	if err == nil {
		t.Error("expected error for unknown key, got nil")
	}
}

func TestKeyExtractor_ExpiredKey(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	mk := &vickytest.MockKeys{
		OnGetByKeyHash: func(ctx context.Context, hash string) (*models.Key, error) {
			return &models.Key{
				ID:        1,
				UserID:    1000,
				Scopes:    []string{"search"},
				ExpiresAt: &past,
			}, nil
		},
	}

	extractor := KeyExtractor(mk, func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		t.Error("session extractor should not be called for vky_ tokens")
		return nil, nil
	})

	_, err := extractor(context.Background(), bearerHeader(rawTestKey))
	if err == nil {
		t.Error("expected error for expired key, got nil")
	}
}

func TestKeyExtractor_ValidKey_NotExpired(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	mk := &vickytest.MockKeys{
		OnGetByKeyHash: func(ctx context.Context, hash string) (*models.Key, error) {
			return &models.Key{
				ID:        1,
				UserID:    1000,
				Scopes:    []string{"search"},
				ExpiresAt: &future,
			}, nil
		},
	}

	extractor := KeyExtractor(mk, nil)
	_, err := extractor(context.Background(), bearerHeader(rawTestKey))
	if err != nil {
		t.Errorf("unexpected error for non-expired key: %v", err)
	}
}

// --- KeyExtractor: fallback to session extractor ---

func TestKeyExtractor_NoAuthHeader_FallsBack(t *testing.T) {
	sessionCalled := false
	extractor := KeyExtractor(&vickytest.MockKeys{}, func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		sessionCalled = true
		return nil, fmt.Errorf("no session")
	})

	_, _ = extractor(context.Background(), noAuthHeader())
	if !sessionCalled {
		t.Error("expected session extractor to be called when no Authorization header present")
	}
}

func TestKeyExtractor_NonVkyBearer_FallsBack(t *testing.T) {
	sessionCalled := false
	extractor := KeyExtractor(&vickytest.MockKeys{}, func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		sessionCalled = true
		return nil, fmt.Errorf("no session")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer some_oauth_token")
	_, _ = extractor(context.Background(), r)

	if !sessionCalled {
		t.Error("expected session extractor to be called for non-vky_ Bearer token")
	}
}
