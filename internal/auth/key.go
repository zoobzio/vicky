// Package auth provides authentication middleware for vicky services.
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/vicky/api/contracts"
)

const keyPrefix = "vky_"

// keyIdentity implements rocco.Identity for API key-authenticated requests.
type keyIdentity struct {
	userID string
	scopes []string
	meta   map[string]any
}

func (i *keyIdentity) ID() string          { return i.userID }
func (i *keyIdentity) TenantID() string    { return "" }
func (i *keyIdentity) Email() string       { return "" }
func (i *keyIdentity) Scopes() []string    { return i.scopes }
func (i *keyIdentity) Roles() []string     { return nil }
func (i *keyIdentity) Stats() map[string]int { return nil }

func (i *keyIdentity) HasScope(scope string) bool {
	for _, s := range i.scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func (i *keyIdentity) HasRole(_ string) bool { return false }

// KeyExtractor returns a composed identity extractor for use with rocco's
// WithAuthenticator. It checks for a Bearer API key (vky_ prefix) first,
// then falls back to the provided session extractor for cookie-based auth.
//
// The key format is: vky_ followed by 32 random bytes base64url-encoded.
// The raw key is SHA-256 hashed (base64-encoded) and looked up in the store.
func KeyExtractor(
	keys contracts.Keys,
	sessionExtractor func(context.Context, *http.Request) (rocco.Identity, error),
) func(context.Context, *http.Request) (rocco.Identity, error) {
	return func(ctx context.Context, r *http.Request) (rocco.Identity, error) {
		authHeader := r.Header.Get("Authorization")

		// Only handle Bearer tokens with the vky_ prefix.
		// Everything else falls through to the session extractor.
		if strings.HasPrefix(authHeader, "Bearer "+keyPrefix) {
			rawKey := strings.TrimPrefix(authHeader, "Bearer ")
			return authenticateKey(ctx, keys, rawKey)
		}

		// Fall back to session-based authentication.
		return sessionExtractor(ctx, r)
	}
}

// authenticateKey validates a raw API key and returns an Identity.
// It hashes the key, looks it up, validates expiry, and records last use.
func authenticateKey(ctx context.Context, keys contracts.Keys, rawKey string) (rocco.Identity, error) {
	if !strings.HasPrefix(rawKey, keyPrefix) {
		return nil, fmt.Errorf("auth: invalid key format")
	}

	hash := hashKey(rawKey)

	apiKey, err := keys.GetByKeyHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("auth: key not found")
	}

	// Validate expiry.
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("auth: key expired")
	}

	// Record last use asynchronously to avoid blocking the request.
	go func() {
		_ = keys.UpdateLastUsed(context.Background(), apiKey.ID)
	}()

	return &keyIdentity{
		userID: fmt.Sprintf("%d", apiKey.UserID),
		scopes: apiKey.Scopes,
		meta: map[string]any{
			"key_id":     apiKey.ID,
			"key_prefix": apiKey.KeyPrefix,
		},
	}, nil
}

// hashKey returns the SHA-256 base64-encoded hash of a raw API key.
func hashKey(rawKey string) string {
	sum := sha256.Sum256([]byte(rawKey))
	return base64.StdEncoding.EncodeToString(sum[:])
}
