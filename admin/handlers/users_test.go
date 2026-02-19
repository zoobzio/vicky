//go:build testing

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/zoobzio/rocco"
	rtesting "github.com/zoobzio/rocco/testing"
	"github.com/zoobzio/sum"
	sumtest "github.com/zoobzio/sum/testing"
	admincontracts "github.com/zoobzio/vicky/admin/contracts"
	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/stores"
)

// Mock admin stores with List/Count support
type MockAdminUsers struct {
	OnGet    func(ctx context.Context, key string) (*models.User, error)
	OnSet    func(ctx context.Context, key string, user *models.User) error
	OnDelete func(ctx context.Context, key string) error
	OnList   func(ctx context.Context, filter *stores.UserFilter, limit, offset int) ([]*models.User, error)
	OnCount  func(ctx context.Context, filter *stores.UserFilter) (int, error)
}

func (m *MockAdminUsers) Get(ctx context.Context, key string) (*models.User, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.User{}, nil
}

func (m *MockAdminUsers) Set(ctx context.Context, key string, user *models.User) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, user)
	}
	return nil
}

func (m *MockAdminUsers) Delete(ctx context.Context, key string) error {
	if m.OnDelete != nil {
		return m.OnDelete(ctx, key)
	}
	return nil
}

func (m *MockAdminUsers) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockAdminUsers) List(ctx context.Context, filter *stores.UserFilter, limit, offset int) ([]*models.User, error) {
	if m.OnList != nil {
		return m.OnList(ctx, filter, limit, offset)
	}
	return nil, nil
}

func (m *MockAdminUsers) Count(ctx context.Context, filter *stores.UserFilter) (int, error) {
	if m.OnCount != nil {
		return m.OnCount(ctx, filter)
	}
	return 0, nil
}

// setupAdminTest sets up the registry for admin handler tests.
func setupAdminTest(t *testing.T, usersStore *MockAdminUsers) *rocco.Engine {
	t.Helper()
	sum.Reset()
	k := sum.Start()


	// Register mock against the contract
	sum.Register[admincontracts.Users](k, usersStore)

	sum.Freeze(k)
	t.Cleanup(sum.Reset)

	_ = sumtest.TestContext(t)

	identity := rtesting.NewMockIdentity("1000")
	engine := rtesting.TestEngineWithAuth(func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		return identity, nil
	})
	return engine
}

func TestListUsers_DefaultPagination(t *testing.T) {
	users := []*models.User{
		{ID: 1001, Login: "alice", Email: "alice@example.com"},
		{ID: 1002, Login: "bob", Email: "bob@example.com"},
	}

	mu := &MockAdminUsers{
		OnList: func(ctx context.Context, filter *stores.UserFilter, limit, offset int) ([]*models.User, error) {
			if limit != 50 || offset != 0 {
				return nil, fmt.Errorf("expected limit=50 offset=0, got limit=%d offset=%d", limit, offset)
			}
			if filter != nil {
				return nil, fmt.Errorf("expected nil filter")
			}
			return users, nil
		},
		OnCount: func(ctx context.Context, filter *stores.UserFilter) (int, error) {
			return 2, nil
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(ListUsers)

	capture := rtesting.ServeRequest(engine, "GET", "/admin/users", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.AdminUserListResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(resp.Users) != 2 {
		t.Errorf("Users count = %d, want 2", len(resp.Users))
	}
	if resp.Total != 2 {
		t.Errorf("Total = %d, want 2", resp.Total)
	}
	if resp.Limit != 50 {
		t.Errorf("Limit = %d, want 50", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("Offset = %d, want 0", resp.Offset)
	}
}

func TestListUsers_CustomPagination(t *testing.T) {
	users := []*models.User{
		{ID: 1001, Login: "alice", Email: "alice@example.com"},
	}

	mu := &MockAdminUsers{
		OnList: func(ctx context.Context, filter *stores.UserFilter, limit, offset int) ([]*models.User, error) {
			if limit != 10 || offset != 5 {
				return nil, fmt.Errorf("expected limit=10 offset=5, got limit=%d offset=%d", limit, offset)
			}
			return users, nil
		},
		OnCount: func(ctx context.Context, filter *stores.UserFilter) (int, error) {
			return 100, nil
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(ListUsers)

	capture := rtesting.ServeRequest(engine, "GET", "/admin/users?limit=10&offset=5", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.AdminUserListResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if resp.Limit != 10 {
		t.Errorf("Limit = %d, want 10", resp.Limit)
	}
	if resp.Offset != 5 {
		t.Errorf("Offset = %d, want 5", resp.Offset)
	}
}

func TestListUsers_InvalidLimit(t *testing.T) {
	mu := &MockAdminUsers{}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(ListUsers)

	// Limit too high
	capture := rtesting.ServeRequest(engine, "GET", "/admin/users?limit=200", nil)
	rtesting.AssertStatus(t, capture, 400)

	// Limit too low
	capture = rtesting.ServeRequest(engine, "GET", "/admin/users?limit=0", nil)
	rtesting.AssertStatus(t, capture, 400)

	// Invalid limit
	capture = rtesting.ServeRequest(engine, "GET", "/admin/users?limit=abc", nil)
	rtesting.AssertStatus(t, capture, 400)
}

func TestListUsers_WithFilters(t *testing.T) {
	users := []*models.User{
		{ID: 1001, Login: "alice", Email: "alice@example.com"},
	}

	mu := &MockAdminUsers{
		OnList: func(ctx context.Context, filter *stores.UserFilter, limit, offset int) ([]*models.User, error) {
			if filter == nil {
				return nil, fmt.Errorf("expected filter")
			}
			if filter.Login == nil || *filter.Login != "alic" {
				return nil, fmt.Errorf("expected login filter 'alic'")
			}
			if filter.Email == nil || *filter.Email != "example.com" {
				return nil, fmt.Errorf("expected email filter 'example.com'")
			}
			if filter.Name == nil || *filter.Name != "Alice" {
				return nil, fmt.Errorf("expected name filter 'Alice'")
			}
			return users, nil
		},
		OnCount: func(ctx context.Context, filter *stores.UserFilter) (int, error) {
			return 1, nil
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(ListUsers)

	capture := rtesting.ServeRequest(engine, "GET", "/admin/users?login=alic&email=example.com&name=Alice", nil)
	rtesting.AssertStatus(t, capture, 200)
}

func TestGetUser_Found(t *testing.T) {
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  stringPtr("Alice Anderson"),
	}

	mu := &MockAdminUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			if key != "1001" {
				return nil, fmt.Errorf("user not found")
			}
			return user, nil
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(GetUser)

	capture := rtesting.ServeRequest(engine, "GET", "/admin/users/1001", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.AdminUserResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if resp.ID != 1001 {
		t.Errorf("ID = %d, want 1001", resp.ID)
	}
	if resp.Login != "alice" {
		t.Errorf("Login = %q, want %q", resp.Login, "alice")
	}
	if resp.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", resp.Email, "alice@example.com")
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mu := &MockAdminUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			return nil, fmt.Errorf("user not found")
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(GetUser)

	capture := rtesting.ServeRequest(engine, "GET", "/admin/users/9999", nil)
	rtesting.AssertStatus(t, capture, 404)
}

func TestUpdateUser_Success(t *testing.T) {
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  stringPtr("Alice"),
	}

	mu := &MockAdminUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			if key != "1001" {
				return nil, fmt.Errorf("user not found")
			}
			return user, nil
		},
		OnSet: func(ctx context.Context, key string, u *models.User) error {
			// Verify update was applied
			if u.Name == nil || *u.Name != "Alice Anderson" {
				return fmt.Errorf("name not updated")
			}
			if u.Email != "alice@newdomain.com" {
				return fmt.Errorf("email not updated")
			}
			if u.Login != "alice_new" {
				return fmt.Errorf("login not updated")
			}
			return nil
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(UpdateUser)

	newName := "Alice Anderson"
	newEmail := "alice@newdomain.com"
	newLogin := "alice_new"
	body := wire.AdminUserUpdateRequest{
		Name:  &newName,
		Email: &newEmail,
		Login: &newLogin,
	}

	capture := rtesting.ServeRequest(engine, "PATCH", "/admin/users/1001", body)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.AdminUserResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if resp.Name == nil || *resp.Name != "Alice Anderson" {
		t.Errorf("Name = %v, want %q", resp.Name, "Alice Anderson")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	mu := &MockAdminUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			return nil, fmt.Errorf("user not found")
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(UpdateUser)

	newName := "New Name"
	body := wire.AdminUserUpdateRequest{Name: &newName}

	capture := rtesting.ServeRequest(engine, "PATCH", "/admin/users/9999", body)
	rtesting.AssertStatus(t, capture, 404)
}

func TestUpdateUser_PartialUpdate(t *testing.T) {
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  stringPtr("Alice"),
	}

	mu := &MockAdminUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			return user, nil
		},
		OnSet: func(ctx context.Context, key string, u *models.User) error {
			// Only name should be updated
			if u.Name == nil || *u.Name != "Alice Anderson" {
				return fmt.Errorf("name not updated")
			}
			if u.Email != "alice@example.com" {
				return fmt.Errorf("email should not change")
			}
			if u.Login != "alice" {
				return fmt.Errorf("login should not change")
			}
			return nil
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(UpdateUser)

	// Only update name
	newName := "Alice Anderson"
	body := wire.AdminUserUpdateRequest{Name: &newName}

	capture := rtesting.ServeRequest(engine, "PATCH", "/admin/users/1001", body)
	rtesting.AssertStatus(t, capture, 200)
}

func TestDeleteUser_Success(t *testing.T) {
	deleted := false

	mu := &MockAdminUsers{
		OnDelete: func(ctx context.Context, key string) error {
			if key != "1001" {
				return fmt.Errorf("user not found")
			}
			deleted = true
			return nil
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(DeleteUser)

	capture := rtesting.ServeRequest(engine, "DELETE", "/admin/users/1001", nil)
	rtesting.AssertStatus(t, capture, 204)

	if !deleted {
		t.Error("Delete was not called")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mu := &MockAdminUsers{
		OnDelete: func(ctx context.Context, key string) error {
			return fmt.Errorf("user not found")
		},
	}

	engine := setupAdminTest(t, mu)
	engine.WithHandlers(DeleteUser)

	capture := rtesting.ServeRequest(engine, "DELETE", "/admin/users/9999", nil)
	rtesting.AssertStatus(t, capture, 404)
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}
