package wire

import (
	"errors"
	"time"

	"github.com/zoobzio/check"
)

// CreateKeyRequest is the request body for creating a new API key.
type CreateKeyRequest struct {
	Name      string     `json:"name" description:"Label for this key" example:"Production"`
	Scopes    []string   `json:"scopes" description:"Permission scopes granted to this key" example:"[\"search\",\"intel\"]"`
	RateLimit *int       `json:"rate_limit,omitempty" description:"Requests per minute limit, omit for unlimited" example:"60"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" description:"Expiry timestamp, omit for no expiry"`
}

// Validate validates the CreateKeyRequest.
func (r *CreateKeyRequest) Validate() error {
	if err := check.All(
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
	).Err(); err != nil {
		return err
	}
	if len(r.Scopes) == 0 {
		return errors.New("scopes: at least one scope is required")
	}
	if r.RateLimit != nil && *r.RateLimit <= 0 {
		return errors.New("rate_limit: must be a positive integer")
	}
	return nil
}

// Clone returns a deep copy of the CreateKeyRequest.
func (r CreateKeyRequest) Clone() CreateKeyRequest {
	c := r
	if r.Scopes != nil {
		c.Scopes = make([]string, len(r.Scopes))
		copy(c.Scopes, r.Scopes)
	}
	if r.RateLimit != nil {
		rl := *r.RateLimit
		c.RateLimit = &rl
	}
	if r.ExpiresAt != nil {
		t := *r.ExpiresAt
		c.ExpiresAt = &t
	}
	return c
}

// KeyResponse is the API response for an API key record.
// The raw key is never included; only metadata is returned.
type KeyResponse struct {
	ID         int64      `json:"id" description:"API key ID" example:"1"`
	Name       string     `json:"name" description:"Key label" example:"Production"`
	Prefix     string     `json:"prefix" description:"First 8 characters of the key for identification" example:"vky_abc1"`
	Scopes     []string   `json:"scopes" description:"Granted permission scopes" example:"[\"search\",\"intel\"]"`
	RateLimit  *int       `json:"rate_limit,omitempty" description:"Requests per minute limit, null for unlimited" example:"60"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" description:"Expiry timestamp, null for no expiry"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" description:"Timestamp of most recent authenticated request"`
	CreatedAt  time.Time  `json:"created_at" description:"Key creation time"`
}

// Clone returns a deep copy of the KeyResponse.
func (r KeyResponse) Clone() KeyResponse {
	c := r
	if r.Scopes != nil {
		c.Scopes = make([]string, len(r.Scopes))
		copy(c.Scopes, r.Scopes)
	}
	if r.RateLimit != nil {
		rl := *r.RateLimit
		c.RateLimit = &rl
	}
	if r.ExpiresAt != nil {
		t := *r.ExpiresAt
		c.ExpiresAt = &t
	}
	if r.LastUsedAt != nil {
		t := *r.LastUsedAt
		c.LastUsedAt = &t
	}
	return c
}

// KeyCreatedResponse is returned once when a key is first created.
// It includes the raw key, which is not stored and cannot be retrieved again.
type KeyCreatedResponse struct {
	KeyResponse
	Key string `json:"key" description:"The raw API key. Store this securely; it is shown only once." example:"vky_abc123..."`
}

// Clone returns a deep copy of the KeyCreatedResponse.
func (r KeyCreatedResponse) Clone() KeyCreatedResponse {
	return KeyCreatedResponse{
		KeyResponse: r.KeyResponse.Clone(),
		Key:         r.Key,
	}
}

// KeyListResponse is the API response for listing API keys.
type KeyListResponse struct {
	Keys []KeyResponse `json:"keys" description:"List of API keys"`
}

// Clone returns a deep copy of the KeyListResponse.
func (r KeyListResponse) Clone() KeyListResponse {
	c := r
	if r.Keys != nil {
		c.Keys = make([]KeyResponse, len(r.Keys))
		for i, k := range r.Keys {
			c.Keys[i] = k.Clone()
		}
	}
	return c
}
