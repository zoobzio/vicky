//go:build testing

package handlers

import (
	"context"
	"fmt"
	"testing"

	rtesting "github.com/zoobzio/rocco/testing"
	vickytest "github.com/zoobzio/vicky/testing"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/api/wire"
)

func TestGetMe(t *testing.T) {
	user := vickytest.NewUser(t)
	mu := &vickytest.MockUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			return user, nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithUsers(mu))
	engine.WithHandlers(GetMe)

	capture := rtesting.ServeRequest(engine, "GET", "/me", nil)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.UserResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Login != "testuser" {
		t.Errorf("Login = %q, want %q", resp.Login, "testuser")
	}
}

func TestGetMe_UserNotFound(t *testing.T) {
	mu := &vickytest.MockUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			return nil, fmt.Errorf("user not found")
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithUsers(mu))
	engine.WithHandlers(GetMe)

	capture := rtesting.ServeRequest(engine, "GET", "/me", nil)
	rtesting.AssertStatus(t, capture, 500)
}

func TestUpdateMe(t *testing.T) {
	user := vickytest.NewUser(t)
	mu := &vickytest.MockUsers{
		OnGet: func(ctx context.Context, key string) (*models.User, error) {
			return user, nil
		},
		OnSet: func(ctx context.Context, key string, u *models.User) error {
			return nil
		},
	}

	engine := vickytest.SetupHandlerTest(t, vickytest.WithUsers(mu))
	engine.WithHandlers(UpdateMe)

	name := "New Name"
	body := wire.UserUpdateRequest{Name: &name}
	capture := rtesting.ServeRequest(engine, "PATCH", "/me", body)
	rtesting.AssertStatus(t, capture, 200)

	var resp wire.UserResponse
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Name == nil || *resp.Name != "New Name" {
		t.Errorf("Name = %v, want %q", resp.Name, "New Name")
	}
}
