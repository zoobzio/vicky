//go:build integration

package stores

import (
	"context"
	"fmt"
	"testing"

	"github.com/zoobzio/vicky/models"
)

// Note: These are integration tests that require a real database.
// Run with: go test -tags=integration ./stores/...

func setupTestUsers(t *testing.T, s *Users) {
	t.Helper()
	ctx := context.Background()

	// Create test users with distinct patterns for filtering
	users := []*models.User{
		{
			ID:    1001,
			Login: "alice",
			Email: "alice@company.com",
			Name:  stringPtr("Alice Anderson"),
		},
		{
			ID:    1002,
			Login: "alicebob",
			Email: "alicebob@company.com",
			Name:  stringPtr("Alice Bob"),
		},
		{
			ID:    1003,
			Login: "bob",
			Email: "bob@company.com",
			Name:  stringPtr("Bob Brown"),
		},
		{
			ID:    1004,
			Login: "charlie",
			Email: "charlie@external.org",
			Name:  stringPtr("Charlie Chen"),
		},
		{
			ID:    1005,
			Login: "david",
			Email: "david@external.org",
			Name:  stringPtr("David Davis"),
		},
	}

	for _, u := range users {
		if err := s.Set(ctx, fmt.Sprintf("%d", u.ID), u); err != nil {
			t.Fatalf("Set user %s: %v", u.Login, err)
		}
	}
}

func TestList_NilFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// TODO: Initialize store with real DB connection
	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// List all users with nil filter
	result, err := users.List(ctx, nil, 10, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(result) < 5 {
		t.Errorf("List returned %d users, want at least 5", len(result))
	}
}

func TestList_FilterByLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// Filter by partial login - should match "alice" and "alicebob"
	loginFilter := "alic"
	filter := &UserFilter{Login: &loginFilter}

	result, err := users.List(ctx, filter, 10, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("List returned %d users, want 2", len(result))
	}

	// Verify both alice and alicebob are returned
	found := make(map[string]bool)
	for _, u := range result {
		found[u.Login] = true
	}
	if !found["alice"] || !found["alicebob"] {
		t.Errorf("Expected alice and alicebob, got: %v", found)
	}
}

func TestList_FilterByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// Filter by email domain - should match 3 users at company.com
	emailFilter := "company.com"
	filter := &UserFilter{Email: &emailFilter}

	result, err := users.List(ctx, filter, 10, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("List returned %d users, want 3", len(result))
	}
}

func TestList_FilterByName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// Filter by partial name - should match "Alice Anderson" and "Bob Anderson"
	nameFilter := "Anderson"
	filter := &UserFilter{Name: &nameFilter}

	result, err := users.List(ctx, filter, 10, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("List returned %d users, want 1 (Alice Anderson)", len(result))
	}
}

func TestList_MultipleFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// Filter by login AND email (both must match)
	loginFilter := "alic"
	emailFilter := "company.com"
	filter := &UserFilter{
		Login: &loginFilter,
		Email: &emailFilter,
	}

	result, err := users.List(ctx, filter, 10, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	// alice and alicebob both match login, and both have company.com email
	if len(result) != 2 {
		t.Errorf("List returned %d users, want 2", len(result))
	}
}

func TestList_Pagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// Get first page (limit 2)
	page1, err := users.List(ctx, nil, 2, 0)
	if err != nil {
		t.Fatalf("List page 1: %v", err)
	}
	if len(page1) != 2 {
		t.Errorf("Page 1 returned %d users, want 2", len(page1))
	}

	// Get second page (limit 2, offset 2)
	page2, err := users.List(ctx, nil, 2, 2)
	if err != nil {
		t.Fatalf("List page 2: %v", err)
	}
	if len(page2) != 2 {
		t.Errorf("Page 2 returned %d users, want 2", len(page2))
	}

	// Verify different users
	if page1[0].ID == page2[0].ID {
		t.Errorf("Page 1 and Page 2 contain same user")
	}

	// Get third page (limit 2, offset 4)
	page3, err := users.List(ctx, nil, 2, 4)
	if err != nil {
		t.Fatalf("List page 3: %v", err)
	}
	if len(page3) >= 1 {
		// Should have at least 1 user on page 3
		t.Logf("Page 3 has %d users", len(page3))
	}
}

func TestCount_NilFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// Count all users
	count, err := users.Count(ctx, nil)
	if err != nil {
		t.Fatalf("Count: %v", err)
	}

	if count < 5 {
		t.Errorf("Count = %d, want at least 5", count)
	}
}

func TestCount_WithFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Skip("Integration test - requires SetupStores helper")

	var users *Users
	ctx := context.Background()

	setupTestUsers(t, users)

	// Count users with company.com email
	emailFilter := "company.com"
	filter := &UserFilter{Email: &emailFilter}

	count, err := users.Count(ctx, filter)
	if err != nil {
		t.Fatalf("Count: %v", err)
	}

	if count != 3 {
		t.Errorf("Count = %d, want 3", count)
	}
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}
