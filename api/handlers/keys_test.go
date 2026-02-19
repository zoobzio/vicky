//go:build testing

package handlers

import (
	"context"
	"fmt"
	"testing"

	rtesting "github.com/zoobzio/rocco/testing"
	vickytest "github.com/zoobzio/vicky/testing"
	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/models"
)

// --- CreateKey ---

func TestCreateKey_Success(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnSet: func(ctx context.Context, key string, apiKey *models.Key) error {
			apiKey.ID = 1
			return nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(CreateKey)

	body := wire.CreateKeyRequest{
		Name:   "Production",
		Scopes: []string{"search"},
	}
	capture := rtesting.ServeRequest(engine, "POST", "/api-keys", body)
	rtesting.AssertStatus(t, capture, 201)

	var resp wire.KeyCreatedResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Key == "" {
		t.Error("expected raw key in response, got empty string")
	}
	if resp.Name != "Production" {
		t.Errorf("Name = %q, want %q", resp.Name, "Production")
	}
}

func TestCreateKey_StoreError(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnSet: func(ctx context.Context, key string, apiKey *models.Key) error {
			return fmt.Errorf("db write failed")
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(CreateKey)

	body := wire.CreateKeyRequest{
		Name:   "Production",
		Scopes: []string{"search"},
	}
	capture := rtesting.ServeRequest(engine, "POST", "/api-keys", body)
	rtesting.AssertStatus(t, capture, 500)
}

// --- ListKeys ---

func TestListKeys_Success(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnListByUserID: func(ctx context.Context, userID int64) ([]*models.Key, error) {
			return []*models.Key{
				{ID: 1, UserID: userID, Name: "Key A", KeyPrefix: "vky_aaa1", Scopes: []string{"search"}},
				{ID: 2, UserID: userID, Name: "Key B", KeyPrefix: "vky_bbb2", Scopes: []string{"intel"}},
			}, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(ListKeys)

	capture := rtesting.ServeRequest(engine, "GET", "/api-keys", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.KeyListResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Keys) != 2 {
		t.Errorf("Keys len = %d, want 2", len(resp.Keys))
	}
}

func TestListKeys_Empty(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnListByUserID: func(ctx context.Context, userID int64) ([]*models.Key, error) {
			return []*models.Key{}, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(ListKeys)

	capture := rtesting.ServeRequest(engine, "GET", "/api-keys", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.KeyListResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Keys) != 0 {
		t.Errorf("Keys len = %d, want 0", len(resp.Keys))
	}
}

func TestListKeys_StoreError(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnListByUserID: func(ctx context.Context, userID int64) ([]*models.Key, error) {
			return nil, fmt.Errorf("db error")
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(ListKeys)

	capture := rtesting.ServeRequest(engine, "GET", "/api-keys", nil)
	rtesting.AssertStatus(t, capture, 500)
}

// --- GetKey ---

func TestGetKey_Success(t *testing.T) {
	// SetupHandlerTest uses identity ID "1000"; key must match that UserID.
	mk := &vickytest.MockKeys{
		OnGet: func(ctx context.Context, key string) (*models.Key, error) {
			return &models.Key{
				ID:        1,
				UserID:    1000,
				Name:      "Production",
				KeyPrefix: "vky_abc1",
				Scopes:    []string{"search"},
			}, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(GetKey)

	capture := rtesting.ServeRequest(engine, "GET", "/api-keys/1", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.KeyResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Name != "Production" {
		t.Errorf("Name = %q, want %q", resp.Name, "Production")
	}
}

func TestGetKey_NotFound(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnGet: func(ctx context.Context, key string) (*models.Key, error) {
			return nil, fmt.Errorf("not found")
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(GetKey)

	capture := rtesting.ServeRequest(engine, "GET", "/api-keys/99", nil)
	rtesting.AssertStatus(t, capture, 404)
}

func TestGetKey_Forbidden(t *testing.T) {
	// Key belongs to a different user (9999 != 1000).
	mk := &vickytest.MockKeys{
		OnGet: func(ctx context.Context, key string) (*models.Key, error) {
			return &models.Key{
				ID:        1,
				UserID:    9999,
				Name:      "Other User Key",
				KeyPrefix: "vky_abc1",
				Scopes:    []string{"search"},
			}, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(GetKey)

	capture := rtesting.ServeRequest(engine, "GET", "/api-keys/1", nil)
	rtesting.AssertStatus(t, capture, 403)
}

// --- DeleteKey ---

func TestDeleteKey_Success(t *testing.T) {
	deleted := false
	mk := &vickytest.MockKeys{
		OnGet: func(ctx context.Context, key string) (*models.Key, error) {
			return &models.Key{
				ID:        1,
				UserID:    1000,
				Name:      "Production",
				KeyPrefix: "vky_abc1",
				Scopes:    []string{"search"},
			}, nil
		},
		OnDelete: func(ctx context.Context, key string) error {
			deleted = true
			return nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(DeleteKey)

	capture := rtesting.ServeRequest(engine, "DELETE", "/api-keys/1", nil)
	rtesting.AssertStatus(t, capture, 200)

	if !deleted {
		t.Error("expected Delete to be called")
	}
}

func TestDeleteKey_NotFound(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnGet: func(ctx context.Context, key string) (*models.Key, error) {
			return nil, fmt.Errorf("not found")
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(DeleteKey)

	capture := rtesting.ServeRequest(engine, "DELETE", "/api-keys/99", nil)
	rtesting.AssertStatus(t, capture, 404)
}

func TestDeleteKey_Forbidden(t *testing.T) {
	mk := &vickytest.MockKeys{
		OnGet: func(ctx context.Context, key string) (*models.Key, error) {
			return &models.Key{
				ID:        1,
				UserID:    9999,
				Name:      "Other User Key",
				KeyPrefix: "vky_abc1",
				Scopes:    []string{"search"},
			}, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithKeys(mk))
	engine.WithHandlers(DeleteKey)

	capture := rtesting.ServeRequest(engine, "DELETE", "/api-keys/1", nil)
	rtesting.AssertStatus(t, capture, 403)
}
