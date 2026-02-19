package wire

import "time"

// AdminUserResponse is the API response for admin user data.
// Unlike user-facing responses, this does NOT mask email or name.
type AdminUserResponse struct {
	ID          int64     `json:"id" description:"GitHub user ID" example:"12345678"`
	Login       string    `json:"login" description:"GitHub username" example:"octocat"`
	Email       string    `json:"email" description:"Email address" example:"octocat@github.com"`
	Name        *string   `json:"name,omitempty" description:"Display name" example:"The Octocat"`
	AvatarURL   *string   `json:"avatar_url,omitempty" description:"GitHub avatar URL"`
	CreatedAt   time.Time `json:"created_at" description:"Account creation time"`
	UpdatedAt   time.Time `json:"updated_at" description:"Last profile sync"`
	LastLoginAt time.Time `json:"last_login_at" description:"Last login time"`
}

// Clone returns a deep copy.
func (r AdminUserResponse) Clone() AdminUserResponse {
	c := r
	if r.Name != nil {
		n := *r.Name
		c.Name = &n
	}
	if r.AvatarURL != nil {
		a := *r.AvatarURL
		c.AvatarURL = &a
	}
	return c
}

// AdminUserListResponse is the API response for listing users.
type AdminUserListResponse struct {
	Users  []AdminUserResponse `json:"users" description:"User array"`
	Total  int                 `json:"total" description:"Total count matching filter"`
	Limit  int                 `json:"limit" description:"Limit used"`
	Offset int                 `json:"offset" description:"Offset used"`
}

// Clone returns a deep copy.
func (r AdminUserListResponse) Clone() AdminUserListResponse {
	c := r
	if r.Users != nil {
		c.Users = make([]AdminUserResponse, len(r.Users))
		for i, u := range r.Users {
			c.Users[i] = u.Clone()
		}
	}
	return c
}

// AdminUserUpdateRequest is the request body for updating user fields.
type AdminUserUpdateRequest struct {
	Name  *string `json:"name,omitempty" description:"New display name" example:"Jane Doe"`
	Email *string `json:"email,omitempty" description:"New email address" example:"user@example.com"`
	Login *string `json:"login,omitempty" description:"New GitHub login" example:"newlogin"`
}

// Clone returns a deep copy.
func (r AdminUserUpdateRequest) Clone() AdminUserUpdateRequest {
	c := r
	if r.Name != nil {
		n := *r.Name
		c.Name = &n
	}
	if r.Email != nil {
		e := *r.Email
		c.Email = &e
	}
	if r.Login != nil {
		l := *r.Login
		c.Login = &l
	}
	return c
}
