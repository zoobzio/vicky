package wire

import "time"

// AdminRepositoryResponse is the API response for admin repository data.
// Unlike user-facing responses, this includes all fields without masking.
type AdminRepositoryResponse struct {
	ID            int64     `json:"id" description:"Internal repository ID" example:"1"`
	GitHubID      int64     `json:"github_id" description:"GitHub repository ID" example:"123456789"`
	UserID        int64     `json:"user_id" description:"Owning Vicky user ID" example:"12345678"`
	Owner         string    `json:"owner" description:"GitHub org or user" example:"octocat"`
	Name          string    `json:"name" description:"Repository name" example:"hello-world"`
	FullName      string    `json:"full_name" description:"Full repository name (owner/name)" example:"octocat/hello-world"`
	Description   *string   `json:"description,omitempty" description:"Repository description"`
	DefaultBranch string    `json:"default_branch" description:"Default branch name" example:"main"`
	Private       bool      `json:"private" description:"Whether repository is private"`
	HTMLURL       string    `json:"html_url" description:"GitHub URL" example:"https://github.com/octocat/hello-world"`
	CreatedAt     time.Time `json:"created_at" description:"Registration time"`
	UpdatedAt     time.Time `json:"updated_at" description:"Last sync time"`
}

// Clone returns a deep copy.
func (r AdminRepositoryResponse) Clone() AdminRepositoryResponse {
	c := r
	if r.Description != nil {
		d := *r.Description
		c.Description = &d
	}
	return c
}

// AdminRepositoryListResponse is the API response for listing repositories.
type AdminRepositoryListResponse struct {
	Repositories []AdminRepositoryResponse `json:"repositories" description:"Repository array"`
	Total        int                       `json:"total" description:"Total count matching filter"`
	Limit        int                       `json:"limit" description:"Limit used"`
	Offset       int                       `json:"offset" description:"Offset used"`
}

// Clone returns a deep copy.
func (r AdminRepositoryListResponse) Clone() AdminRepositoryListResponse {
	c := r
	if r.Repositories != nil {
		c.Repositories = make([]AdminRepositoryResponse, len(r.Repositories))
		for i, repo := range r.Repositories {
			c.Repositories[i] = repo.Clone()
		}
	}
	return c
}

// AdminRepositoryUpdateRequest is the request body for updating repository fields.
type AdminRepositoryUpdateRequest struct {
	Owner         *string `json:"owner,omitempty" description:"New owner" example:"newowner"`
	Name          *string `json:"name,omitempty" description:"New name" example:"new-repo"`
	FullName      *string `json:"full_name,omitempty" description:"New full name" example:"newowner/new-repo"`
	Description   *string `json:"description,omitempty" description:"New description"`
	DefaultBranch *string `json:"default_branch,omitempty" description:"New default branch" example:"main"`
	Private       *bool   `json:"private,omitempty" description:"New private status"`
	HTMLURL       *string `json:"html_url,omitempty" description:"New GitHub URL"`
}

// Clone returns a deep copy.
func (r AdminRepositoryUpdateRequest) Clone() AdminRepositoryUpdateRequest {
	c := r
	if r.Owner != nil {
		o := *r.Owner
		c.Owner = &o
	}
	if r.Name != nil {
		n := *r.Name
		c.Name = &n
	}
	if r.FullName != nil {
		f := *r.FullName
		c.FullName = &f
	}
	if r.Description != nil {
		d := *r.Description
		c.Description = &d
	}
	if r.DefaultBranch != nil {
		b := *r.DefaultBranch
		c.DefaultBranch = &b
	}
	if r.Private != nil {
		p := *r.Private
		c.Private = &p
	}
	if r.HTMLURL != nil {
		u := *r.HTMLURL
		c.HTMLURL = &u
	}
	return c
}
