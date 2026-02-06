package models

import "time"

// Repository represents a GitHub repository registered by a user.
// Multiple users can register the same GitHub repository independently.
type Repository struct {
	ID            int64     `json:"id" db:"id" constraints:"primarykey" description:"Internal repository ID"`
	GitHubID      int64     `json:"github_id" db:"github_id" constraints:"notnull" description:"GitHub repository ID"`
	UserID        int64     `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning Vicky user"`
	Owner         string    `json:"owner" db:"owner" constraints:"notnull" description:"GitHub org or user" example:"octocat"`
	Name          string    `json:"name" db:"name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	FullName      string    `json:"full_name" db:"full_name" constraints:"notnull" description:"owner/name" example:"octocat/hello-world"`
	Description   *string   `json:"description,omitempty" db:"description" description:"Repository description"`
	DefaultBranch string    `json:"default_branch" db:"default_branch" constraints:"notnull" default:"'main'" description:"Default branch name" example:"main"`
	Private       bool      `json:"private" db:"private" constraints:"notnull" default:"false" description:"Whether repository is private"`
	HTMLURL       string    `json:"html_url" db:"html_url" constraints:"notnull" validate:"url" description:"GitHub URL" example:"https://github.com/octocat/hello-world"`
	CreatedAt     time.Time `json:"created_at" db:"created_at" default:"now()" description:"Registration time"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at" default:"now()" description:"Last sync time"`
}

// Clone returns a deep copy of the Repository.
func (r Repository) Clone() Repository {
	c := r
	if r.Description != nil {
		d := *r.Description
		c.Description = &d
	}
	return c
}
