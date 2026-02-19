package transformers

import (
	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/models"
)

// KeyToResponse transforms a Key model to an API response.
// The raw key and key hash are never included in the response.
func KeyToResponse(k *models.Key) wire.KeyResponse {
	return wire.KeyResponse{
		ID:         k.ID,
		Name:       k.Name,
		Prefix:     k.KeyPrefix,
		Scopes:     k.Scopes,
		RateLimit:  k.RateLimit,
		ExpiresAt:  k.ExpiresAt,
		LastUsedAt: k.LastUsedAt,
		CreatedAt:  k.CreatedAt,
	}
}

// KeyToCreatedResponse transforms a Key model and its raw key to a creation response.
// rawKey is the plaintext key returned once at creation time and never again.
func KeyToCreatedResponse(k *models.Key, rawKey string) wire.KeyCreatedResponse {
	return wire.KeyCreatedResponse{
		KeyResponse: KeyToResponse(k),
		Key:         rawKey,
	}
}

// KeysToList transforms a slice of Key models to an API list response.
func KeysToList(keys []*models.Key) wire.KeyListResponse {
	resp := wire.KeyListResponse{
		Keys: make([]wire.KeyResponse, len(keys)),
	}
	for i, k := range keys {
		resp.Keys[i] = KeyToResponse(k)
	}
	return resp
}

// ApplyCreateKeyRequest applies a CreateKeyRequest to a Key model.
// It sets all user-supplied fields; the caller must populate KeyHash,
// KeyPrefix, UserID, and ID separately.
func ApplyCreateKeyRequest(req wire.CreateKeyRequest, k *models.Key) {
	k.Name = req.Name
	k.Scopes = req.Scopes
	k.RateLimit = req.RateLimit
	k.ExpiresAt = req.ExpiresAt
}
