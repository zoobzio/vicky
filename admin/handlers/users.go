package handlers

import (
	"strconv"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	admincontracts "github.com/zoobzio/vicky/admin/contracts"
	"github.com/zoobzio/vicky/admin/transformers"
	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/stores"
)

// ListUsers returns all users with pagination and optional filtering.
var ListUsers = rocco.GET("/admin/users", func(req *rocco.Request[rocco.NoBody]) (wire.AdminUserListResponse, error) {
	usersStore := sum.MustUse[admincontracts.Users](req.Context)

	// Parse limit with validation
	limit := 50
	if l := req.Params.Query["limit"]; l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 || parsed > 100 {
			return wire.AdminUserListResponse{}, ErrInvalidLimit
		}
		limit = parsed
	}

	// Parse offset
	offset := 0
	if o := req.Params.Query["offset"]; o != "" {
		parsed, err := strconv.Atoi(o)
		if err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Build filter from query params
	var filter *stores.UserFilter
	loginFilter := req.Params.Query["login"]
	emailFilter := req.Params.Query["email"]
	nameFilter := req.Params.Query["name"]

	if loginFilter != "" || emailFilter != "" || nameFilter != "" {
		filter = &stores.UserFilter{}
		if loginFilter != "" {
			filter.Login = &loginFilter
		}
		if emailFilter != "" {
			filter.Email = &emailFilter
		}
		if nameFilter != "" {
			filter.Name = &nameFilter
		}
	}

	// Execute list query
	users, err := usersStore.List(req.Context, filter, limit, offset)
	if err != nil {
		return wire.AdminUserListResponse{}, err
	}

	// Get total count
	total, err := usersStore.Count(req.Context, filter)
	if err != nil {
		return wire.AdminUserListResponse{}, err
	}

	return transformers.UsersToAdminList(users, total, limit, offset), nil
}).WithSummary("List users").
	WithDescription("Returns all users with pagination and optional filtering by login, email, or name.").
	WithTags("Admin", "Users").
	WithQueryParams("limit", "offset", "login", "email", "name").
	WithErrors(ErrInvalidLimit)

// GetUser returns a single user by ID.
var GetUser = rocco.GET("/admin/users/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.AdminUserResponse, error) {
	usersStore := sum.MustUse[admincontracts.Users](req.Context)

	id := req.Params.Path["id"]

	user, err := usersStore.Get(req.Context, id)
	if err != nil {
		return wire.AdminUserResponse{}, ErrUserNotFound
	}

	return transformers.UserToAdminResponse(user), nil
}).WithSummary("Get user").
	WithDescription("Returns a single user by ID.").
	WithTags("Admin", "Users").
	WithPathParams("id").
	WithErrors(ErrUserNotFound)

// UpdateUser updates user fields.
var UpdateUser = rocco.PATCH("/admin/users/{id}", func(req *rocco.Request[wire.AdminUserUpdateRequest]) (wire.AdminUserResponse, error) {
	usersStore := sum.MustUse[admincontracts.Users](req.Context)

	id := req.Params.Path["id"]

	// Get existing user
	user, err := usersStore.Get(req.Context, id)
	if err != nil {
		return wire.AdminUserResponse{}, ErrUserNotFound
	}

	// Apply update
	transformers.ApplyAdminUserUpdate(req.Body, user)

	// Save updated user
	if err := usersStore.Set(req.Context, id, user); err != nil {
		return wire.AdminUserResponse{}, err
	}

	return transformers.UserToAdminResponse(user), nil
}).WithSummary("Update user").
	WithDescription("Updates user fields. Only provided fields are updated.").
	WithTags("Admin", "Users").
	WithPathParams("id").
	WithErrors(ErrUserNotFound)

// DeleteUser hard deletes a user with cascade.
var DeleteUser = rocco.DELETE("/admin/users/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	usersStore := sum.MustUse[admincontracts.Users](req.Context)

	id := req.Params.Path["id"]

	// Delete will cascade via database constraints
	if err := usersStore.Delete(req.Context, id); err != nil {
		return rocco.NoBody{}, ErrUserNotFound
	}

	return rocco.NoBody{}, nil
}).WithSummary("Delete user").
	WithDescription("Hard deletes a user with cascade. Removes all user data including repositories, versions, documents, chunks, symbols, SCIP data, ingestion configs, and jobs.").
	WithTags("Admin", "Users").
	WithPathParams("id").
	WithErrors(ErrUserNotFound).
	WithSuccessStatus(204)
