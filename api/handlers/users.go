package handlers

import (
	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/api/transformers"
)

// GetMe returns the authenticated user's profile.
var GetMe = rocco.GET("/me", func(req *rocco.Request[rocco.NoBody]) (wire.UserResponse, error) {
	users := sum.MustUse[contracts.Users](req.Context)

	user, err := users.Get(req.Context, req.Identity.ID())
	if err != nil {
		return wire.UserResponse{}, err
	}

	return transformers.UserToResponse(user), nil
}).WithSummary("Get current user").
	WithDescription("Returns the authenticated user's profile.").
	WithTags("Users").
	WithAuthentication()

// UpdateMe updates the authenticated user's profile.
var UpdateMe = rocco.PATCH("/me", func(req *rocco.Request[wire.UserUpdateRequest]) (wire.UserResponse, error) {
	users := sum.MustUse[contracts.Users](req.Context)

	user, err := users.Get(req.Context, req.Identity.ID())
	if err != nil {
		return wire.UserResponse{}, err
	}

	transformers.ApplyUserUpdate(req.Body, user)

	if err := users.Set(req.Context, req.Identity.ID(), user); err != nil {
		return wire.UserResponse{}, err
	}

	return transformers.UserToResponse(user), nil
}).WithSummary("Update current user").
	WithDescription("Updates the authenticated user's profile.").
	WithTags("Users").
	WithAuthentication()
