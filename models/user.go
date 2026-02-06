package models

import (
	"context"
	"time"

	"github.com/zoobzio/sum"
)

// User represents an authenticated GitHub user.
// GitHub is the sole identity provider - the GitHub user ID is the primary key.
type User struct {
	ID          int64      `json:"id" db:"id" constraints:"primarykey" description:"GitHub user ID" example:"12345678"`
	Login       string     `json:"login" db:"login" constraints:"notnull,unique" description:"GitHub username" example:"octocat"`
	Email       string     `json:"email" db:"email" constraints:"notnull" validate:"email" description:"GitHub email" example:"octocat@github.com"`
	Name        *string    `json:"name,omitempty" db:"name" description:"Display name" example:"The Octocat"`
	AvatarURL   *string    `json:"avatar_url,omitempty" db:"avatar_url" validate:"url" description:"GitHub avatar URL"`
	AccessToken string     `json:"-" db:"access_token" constraints:"notnull" store.encrypt:"aes" load.decrypt:"aes"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" default:"now()" description:"Account creation time"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" default:"now()" description:"Last profile sync"`
	LastLoginAt time.Time  `json:"last_login_at" db:"last_login_at" default:"now()" description:"Last login time"`
}

// BeforeSave encrypts sensitive fields before persistence.
func (u *User) BeforeSave(ctx context.Context) error {
	b := sum.MustUse[*sum.Boundary[User]](ctx)
	stored, err := b.Store(ctx, *u)
	if err != nil {
		return err
	}
	*u = stored
	return nil
}

// AfterLoad decrypts sensitive fields after loading from storage.
func (u *User) AfterLoad(ctx context.Context) error {
	b := sum.MustUse[*sum.Boundary[User]](ctx)
	loaded, err := b.Load(ctx, *u)
	if err != nil {
		return err
	}
	*u = loaded
	return nil
}

// Clone returns a deep copy of the User.
func (u User) Clone() User {
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
