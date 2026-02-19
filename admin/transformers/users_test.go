//go:build testing

package transformers

import (
	"testing"
	"time"

	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/models"
)

func TestUserToAdminResponse(t *testing.T) {
	now := time.Now()
	name := "Alice Anderson"
	avatarURL := "https://example.com/avatar.png"

	user := &models.User{
		ID:          1001,
		Login:       "alice",
		Email:       "alice@example.com",
		Name:        &name,
		AvatarURL:   &avatarURL,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: now,
	}

	resp := UserToAdminResponse(user)

	// Verify all fields are mapped correctly
	if resp.ID != 1001 {
		t.Errorf("ID = %d, want 1001", resp.ID)
	}
	if resp.Login != "alice" {
		t.Errorf("Login = %q, want %q", resp.Login, "alice")
	}
	if resp.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", resp.Email, "alice@example.com")
	}
	if resp.Name == nil || *resp.Name != "Alice Anderson" {
		t.Errorf("Name = %v, want %q", resp.Name, "Alice Anderson")
	}
	if resp.AvatarURL == nil || *resp.AvatarURL != avatarURL {
		t.Errorf("AvatarURL = %v, want %q", resp.AvatarURL, avatarURL)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", resp.CreatedAt, now)
	}
	if !resp.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", resp.UpdatedAt, now)
	}
	if !resp.LastLoginAt.Equal(now) {
		t.Errorf("LastLoginAt = %v, want %v", resp.LastLoginAt, now)
	}
}

func TestUserToAdminResponse_NilOptionalFields(t *testing.T) {
	now := time.Now()

	user := &models.User{
		ID:          1001,
		Login:       "alice",
		Email:       "alice@example.com",
		Name:        nil,
		AvatarURL:   nil,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: now,
	}

	resp := UserToAdminResponse(user)

	if resp.Name != nil {
		t.Errorf("Name = %v, want nil", resp.Name)
	}
	if resp.AvatarURL != nil {
		t.Errorf("AvatarURL = %v, want nil", resp.AvatarURL)
	}
}

func TestUsersToAdminList(t *testing.T) {
	now := time.Now()
	name1 := "Alice"
	name2 := "Bob"

	users := []*models.User{
		{
			ID:          1001,
			Login:       "alice",
			Email:       "alice@example.com",
			Name:        &name1,
			CreatedAt:   now,
			UpdatedAt:   now,
			LastLoginAt: now,
		},
		{
			ID:          1002,
			Login:       "bob",
			Email:       "bob@example.com",
			Name:        &name2,
			CreatedAt:   now,
			UpdatedAt:   now,
			LastLoginAt: now,
		},
	}

	resp := UsersToAdminList(users, 100, 10, 20)

	// Verify pagination metadata
	if resp.Total != 100 {
		t.Errorf("Total = %d, want 100", resp.Total)
	}
	if resp.Limit != 10 {
		t.Errorf("Limit = %d, want 10", resp.Limit)
	}
	if resp.Offset != 20 {
		t.Errorf("Offset = %d, want 20", resp.Offset)
	}

	// Verify user array
	if len(resp.Users) != 2 {
		t.Fatalf("Users count = %d, want 2", len(resp.Users))
	}

	// Verify first user
	if resp.Users[0].ID != 1001 {
		t.Errorf("Users[0].ID = %d, want 1001", resp.Users[0].ID)
	}
	if resp.Users[0].Login != "alice" {
		t.Errorf("Users[0].Login = %q, want %q", resp.Users[0].Login, "alice")
	}

	// Verify second user
	if resp.Users[1].ID != 1002 {
		t.Errorf("Users[1].ID = %d, want 1002", resp.Users[1].ID)
	}
	if resp.Users[1].Login != "bob" {
		t.Errorf("Users[1].Login = %q, want %q", resp.Users[1].Login, "bob")
	}
}

func TestUsersToAdminList_EmptyArray(t *testing.T) {
	users := []*models.User{}

	resp := UsersToAdminList(users, 0, 10, 0)

	if len(resp.Users) != 0 {
		t.Errorf("Users count = %d, want 0", len(resp.Users))
	}
	if resp.Total != 0 {
		t.Errorf("Total = %d, want 0", resp.Total)
	}
}

func TestApplyAdminUserUpdate_AllFields(t *testing.T) {
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  stringPtr("Alice"),
	}

	newName := "Alice Anderson"
	newEmail := "alice@newdomain.com"
	newLogin := "alice_new"

	req := wire.AdminUserUpdateRequest{
		Name:  &newName,
		Email: &newEmail,
		Login: &newLogin,
	}

	ApplyAdminUserUpdate(req, user)

	// Verify all fields were updated
	if user.Name == nil || *user.Name != "Alice Anderson" {
		t.Errorf("Name = %v, want %q", user.Name, "Alice Anderson")
	}
	if user.Email != "alice@newdomain.com" {
		t.Errorf("Email = %q, want %q", user.Email, "alice@newdomain.com")
	}
	if user.Login != "alice_new" {
		t.Errorf("Login = %q, want %q", user.Login, "alice_new")
	}
}

func TestApplyAdminUserUpdate_PartialUpdate(t *testing.T) {
	originalName := "Alice"
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  &originalName,
	}

	newEmail := "alice@newdomain.com"

	req := wire.AdminUserUpdateRequest{
		Email: &newEmail,
		// Name and Login not provided
	}

	ApplyAdminUserUpdate(req, user)

	// Only email should be updated
	if user.Email != "alice@newdomain.com" {
		t.Errorf("Email = %q, want %q", user.Email, "alice@newdomain.com")
	}

	// Name should remain unchanged
	if user.Name == nil || *user.Name != "Alice" {
		t.Errorf("Name = %v, want %q", user.Name, "Alice")
	}

	// Login should remain unchanged
	if user.Login != "alice" {
		t.Errorf("Login = %q, want %q", user.Login, "alice")
	}
}

func TestApplyAdminUserUpdate_NoFields(t *testing.T) {
	originalName := "Alice"
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  &originalName,
	}

	req := wire.AdminUserUpdateRequest{}

	ApplyAdminUserUpdate(req, user)

	// Nothing should change
	if user.Login != "alice" {
		t.Errorf("Login = %q, want %q", user.Login, "alice")
	}
	if user.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "alice@example.com")
	}
	if user.Name == nil || *user.Name != "Alice" {
		t.Errorf("Name = %v, want %q", user.Name, "Alice")
	}
}

func TestApplyAdminUserUpdate_SetNameToNil(t *testing.T) {
	originalName := "Alice"
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  &originalName,
	}

	// Explicitly set name to nil (empty string pointer means remove name)
	emptyName := ""
	req := wire.AdminUserUpdateRequest{
		Name: &emptyName,
	}

	ApplyAdminUserUpdate(req, user)

	// Name should be updated to empty string pointer
	if user.Name == nil {
		t.Error("Name should not be nil, should be pointer to empty string")
	} else if *user.Name != "" {
		t.Errorf("Name = %q, want empty string", *user.Name)
	}
}

func TestApplyAdminUserUpdate_UpdateOnlyLogin(t *testing.T) {
	originalName := "Alice"
	user := &models.User{
		ID:    1001,
		Login: "alice",
		Email: "alice@example.com",
		Name:  &originalName,
	}

	newLogin := "alice_updated"
	req := wire.AdminUserUpdateRequest{
		Login: &newLogin,
	}

	ApplyAdminUserUpdate(req, user)

	// Only login should be updated
	if user.Login != "alice_updated" {
		t.Errorf("Login = %q, want %q", user.Login, "alice_updated")
	}

	// Email and Name should remain unchanged
	if user.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "alice@example.com")
	}
	if user.Name == nil || *user.Name != "Alice" {
		t.Errorf("Name = %v, want %q", user.Name, "Alice")
	}
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}
