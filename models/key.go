package models

import (
	"time"

	"github.com/zoobzio/check"
)

// Key represents an API key for programmatic authentication.
// The raw key is never stored; only a SHA-256 hash is persisted.
type Key struct {
	ID         int64      `json:"id" db:"id" constraints:"primarykey" description:"API key ID" example:"1"`
	UserID     int64      `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user ID" example:"12345678"`
	Name       string     `json:"name" db:"name" constraints:"notnull" description:"Key label" example:"Production"`
	KeyHash    string     `json:"-" db:"key_hash" constraints:"notnull,unique" description:"SHA-256 hash of the raw key"`
	KeyPrefix  string     `json:"key_prefix" db:"key_prefix" constraints:"notnull" description:"First 8 characters of the key for display" example:"vky_abc1"`
	Scopes     []string   `json:"scopes" db:"scopes" constraints:"notnull" type:"text[]" description:"Granted permission scopes" example:"[\"search\",\"intel\"]"`
	RateLimit  *int       `json:"rate_limit,omitempty" db:"rate_limit" description:"Requests per minute limit, null for unlimited" example:"60"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at" description:"Expiry timestamp, null for no expiry"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at" description:"Timestamp of most recent authenticated request"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at" default:"now()" description:"Key creation time"`
}

// Clone returns a deep copy of the Key.
func (k Key) Clone() Key {
	c := k
	if k.RateLimit != nil {
		r := *k.RateLimit
		c.RateLimit = &r
	}
	if k.ExpiresAt != nil {
		t := *k.ExpiresAt
		c.ExpiresAt = &t
	}
	if k.LastUsedAt != nil {
		t := *k.LastUsedAt
		c.LastUsedAt = &t
	}
	if k.Scopes != nil {
		c.Scopes = make([]string, len(k.Scopes))
		copy(c.Scopes, k.Scopes)
	}
	return c
}

// Validate validates the Key model fields.
func (k Key) Validate() error {
	return check.All(
		check.Str(k.Name, "name").Required().MaxLen(255).V(),
		check.Str(k.KeyHash, "key_hash").Required().V(),
		check.Str(k.KeyPrefix, "key_prefix").Required().V(),
	).Err()
}
