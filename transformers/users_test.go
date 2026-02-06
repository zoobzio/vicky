package transformers

import (
	"testing"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/wire"
)

func TestUserToResponse(t *testing.T) {
	name := "The Octocat"
	avatar := "https://example.com/avatar.png"
	u := &models.User{
		ID:        42,
		Login:     "octocat",
		Email:     "octo@example.com",
		Name:      &name,
		AvatarURL: &avatar,
	}

	resp := UserToResponse(u)

	if resp.ID != 42 {
		t.Errorf("ID = %d, want 42", resp.ID)
	}
	if resp.Login != "octocat" {
		t.Errorf("Login = %q, want %q", resp.Login, "octocat")
	}
	if resp.Email != "octo@example.com" {
		t.Errorf("Email = %q, want %q", resp.Email, "octo@example.com")
	}
	if resp.Name == nil || *resp.Name != "The Octocat" {
		t.Errorf("Name = %v, want %q", resp.Name, "The Octocat")
	}
	if resp.AvatarURL == nil || *resp.AvatarURL != avatar {
		t.Errorf("AvatarURL = %v, want %q", resp.AvatarURL, avatar)
	}
}

func TestApplyUserUpdate(t *testing.T) {
	name := "New Name"
	req := wire.UserUpdateRequest{Name: &name}
	u := &models.User{ID: 1, Login: "octocat"}

	ApplyUserUpdate(req, u)

	if u.Name == nil || *u.Name != "New Name" {
		t.Errorf("Name = %v, want %q", u.Name, "New Name")
	}
}

func TestApplyUserUpdate_NilName(t *testing.T) {
	existing := "Existing"
	req := wire.UserUpdateRequest{Name: nil}
	u := &models.User{ID: 1, Name: &existing}

	ApplyUserUpdate(req, u)

	if u.Name == nil || *u.Name != "Existing" {
		t.Errorf("Name = %v, want %q (unchanged)", u.Name, "Existing")
	}
}
