package transformers

import (
	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/models"
)

// UserToResponse transforms a User model to an API response.
func UserToResponse(u *models.User) wire.UserResponse {
	return wire.UserResponse{
		ID:        u.ID,
		Login:     u.Login,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
	}
}

// ApplyUserUpdate applies a UserUpdateRequest to a User model.
func ApplyUserUpdate(req wire.UserUpdateRequest, u *models.User) {
	if req.Name != nil {
		u.Name = req.Name
	}
}
