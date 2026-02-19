package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/api/transformers"
	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/models"
)

const apiKeyPrefix = "vky_"

// CreateKey creates a new API key for the authenticated user.
// The raw key is returned once in the response and cannot be retrieved again.
var CreateKey = rocco.POST("/api-keys", func(req *rocco.Request[wire.CreateKeyRequest]) (wire.KeyCreatedResponse, error) {
	keys := sum.MustUse[contracts.Keys](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.KeyCreatedResponse{}, err
	}

	// Generate 32 random bytes and encode as base64url for the key body.
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return wire.KeyCreatedResponse{}, err
	}
	rawKey := apiKeyPrefix + base64.RawURLEncoding.EncodeToString(raw)

	// SHA-256 hash of the raw key for storage.
	sum256 := sha256.Sum256([]byte(rawKey))
	keyHash := base64.StdEncoding.EncodeToString(sum256[:])

	// First 8 characters of the raw key (including prefix) for display.
	keyPrefix := rawKey
	if len(keyPrefix) > 8 {
		keyPrefix = keyPrefix[:8]
	}

	key := &models.Key{
		UserID:    userID,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
	}
	transformers.ApplyCreateKeyRequest(req.Body, key)

	if err := keys.Set(req.Context, "", key); err != nil {
		return wire.KeyCreatedResponse{}, err
	}

	return transformers.KeyToCreatedResponse(key, rawKey), nil
}).WithSummary("Create API key").
	WithDescription("Creates a new API key. The raw key is returned once and cannot be retrieved again.").
	WithTags("API Keys").
	WithAuthentication().
	WithSuccessStatus(201)

// ListKeys returns all API keys for the authenticated user.
var ListKeys = rocco.GET("/api-keys", func(req *rocco.Request[rocco.NoBody]) (wire.KeyListResponse, error) {
	keys := sum.MustUse[contracts.Keys](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.KeyListResponse{}, err
	}

	list, err := keys.ListByUserID(req.Context, userID)
	if err != nil {
		return wire.KeyListResponse{}, err
	}

	return transformers.KeysToList(list), nil
}).WithSummary("List API keys").
	WithDescription("Returns all API keys for the authenticated user.").
	WithTags("API Keys").
	WithAuthentication()

// GetKey returns a specific API key by ID.
var GetKey = rocco.GET("/api-keys/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.KeyResponse, error) {
	keyStore := sum.MustUse[contracts.Keys](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.KeyResponse{}, err
	}

	keyID := req.Params.Path["id"]

	key, err := keyStore.Get(req.Context, keyID)
	if err != nil {
		return wire.KeyResponse{}, ErrKeyNotFound
	}

	if key.UserID != userID {
		return wire.KeyResponse{}, ErrKeyForbidden
	}

	return transformers.KeyToResponse(key), nil
}).WithPathParams("id").
	WithSummary("Get API key").
	WithDescription("Returns a specific API key by ID.").
	WithTags("API Keys").
	WithErrors(ErrKeyNotFound, ErrKeyForbidden).
	WithAuthentication()

// DeleteKey revokes and deletes an API key.
var DeleteKey = rocco.DELETE("/api-keys/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	keyStore := sum.MustUse[contracts.Keys](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return rocco.NoBody{}, err
	}

	keyID := req.Params.Path["id"]

	key, err := keyStore.Get(req.Context, keyID)
	if err != nil {
		return rocco.NoBody{}, ErrKeyNotFound
	}

	if key.UserID != userID {
		return rocco.NoBody{}, ErrKeyForbidden
	}

	if err := keyStore.Delete(req.Context, keyID); err != nil {
		return rocco.NoBody{}, err
	}

	return rocco.NoBody{}, nil
}).WithPathParams("id").
	WithSummary("Delete API key").
	WithDescription("Revokes and permanently deletes an API key.").
	WithTags("API Keys").
	WithErrors(ErrKeyNotFound, ErrKeyForbidden).
	WithAuthentication()
