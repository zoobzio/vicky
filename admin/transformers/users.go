package transformers

import (
	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/models"
)

// UserToAdminResponse transforms a User model to an admin API response.
// Admin responses include all fields without masking.
func UserToAdminResponse(u *models.User) wire.AdminUserResponse {
	return wire.AdminUserResponse{
		ID:          u.ID,
		Login:       u.Login,
		Email:       u.Email,
		Name:        u.Name,
		AvatarURL:   u.AvatarURL,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt,
	}
}

// UsersToAdminList transforms a slice of User models to a paginated admin list response.
func UsersToAdminList(users []*models.User, total, limit, offset int) wire.AdminUserListResponse {
	resp := wire.AdminUserListResponse{
		Users:  make([]wire.AdminUserResponse, len(users)),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	for i, u := range users {
		resp.Users[i] = UserToAdminResponse(u)
	}
	return resp
}

// ApplyAdminUserUpdate applies an AdminUserUpdateRequest to a User model.
// Only updates fields that are present in the request.
func ApplyAdminUserUpdate(req wire.AdminUserUpdateRequest, u *models.User) {
	if req.Name != nil {
		u.Name = req.Name
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Login != nil {
		u.Login = *req.Login
	}
}
