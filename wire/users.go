package wire

import (
	"context"

	"github.com/zoobzio/check"
	"github.com/zoobzio/sum"
)

// UserResponse is the API response for user data.
type UserResponse struct {
	ID        int64   `json:"id" description:"GitHub user ID" example:"12345678"`
	Login     string  `json:"login" description:"GitHub username" example:"octocat"`
	Email     string  `json:"email" description:"Email address" example:"octocat@github.com" send.mask:"email"`
	Name      *string `json:"name,omitempty" description:"Display name" example:"The Octocat" send.mask:"name"`
	AvatarURL *string `json:"avatar_url,omitempty" description:"Profile image URL" example:"https://avatars.githubusercontent.com/u/583231"`
}

// OnSend applies boundary masking before the response is marshaled.
func (u *UserResponse) OnSend(ctx context.Context) error {
	b := sum.MustUse[*sum.Boundary[UserResponse]](ctx)
	masked, err := b.Send(ctx, *u)
	if err != nil {
		return err
	}
	*u = masked
	return nil
}

// UserUpdateRequest is the request body for updating user profile.
type UserUpdateRequest struct {
	Name *string `json:"name,omitempty" description:"New display name" example:"The Octocat" validate:"omitempty,max=255"`
}

// Clone returns a deep copy of the UserResponse.
func (u UserResponse) Clone() UserResponse {
	c := u
	if u.Name != nil {
		n := *u.Name
		c.Name = &n
	}
	if u.AvatarURL != nil {
		a := *u.AvatarURL
		c.AvatarURL = &a
	}
	return c
}

// Clone returns a deep copy of the UserUpdateRequest.
func (r UserUpdateRequest) Clone() UserUpdateRequest {
	c := r
	if r.Name != nil {
		n := *r.Name
		c.Name = &n
	}
	return c
}

// Validate validates the UserUpdateRequest.
func (r *UserUpdateRequest) Validate() error {
	return check.All(
		check.OptStr(r.Name, "name").MaxLen(255).V(),
	).Err()
}
