package wire

import (
	"encoding/json"

	"github.com/zoobzio/check"
	"github.com/zoobzio/vicky/models"
)

// RepositoryResponse is the API response for repository data.
type RepositoryResponse struct {
	ID            int64   `json:"id" description:"Repository ID"`
	GitHubID      int64   `json:"github_id" description:"GitHub repository ID" example:"123456789"`
	Owner         string  `json:"owner" description:"Repository owner" example:"octocat"`
	Name          string  `json:"name" description:"Repository name" example:"hello-world"`
	FullName      string  `json:"full_name" description:"Full repository name" example:"octocat/hello-world"`
	Description   *string `json:"description,omitempty" description:"Repository description"`
	DefaultBranch string  `json:"default_branch" description:"Default branch" example:"main"`
	Private       bool    `json:"private" description:"Whether repository is private" example:"false"`
	HTMLURL       string  `json:"html_url" description:"GitHub URL" example:"https://github.com/octocat/hello-world"`
}

// RepositoryListResponse is the API response for listing repositories.
type RepositoryListResponse struct {
	Repositories []RepositoryResponse `json:"repositories" description:"List of registered repositories"`
}

// RegisterRepositoryRequest is the request body for registering a repository.
type RegisterRepositoryRequest struct {
	GitHubID      int64                   `json:"github_id" description:"GitHub repository ID" example:"123456789" validate:"required"`
	Owner         string                  `json:"owner" description:"Repository owner" example:"octocat" validate:"required,max=255"`
	Name          string                  `json:"name" description:"Repository name" example:"hello-world" validate:"required,max=255"`
	FullName      string                  `json:"full_name" description:"Full repository name" example:"octocat/hello-world" validate:"required,max=512"`
	Description   *string                 `json:"description,omitempty" description:"Repository description" validate:"omitempty,max=1000"`
	DefaultBranch string                  `json:"default_branch" description:"Default branch" example:"main" validate:"required,max=255"`
	Private       bool                    `json:"private" description:"Whether repository is private" example:"false"`
	HTMLURL       string                  `json:"html_url" description:"GitHub URL" example:"https://github.com/octocat/hello-world" validate:"required,url"`
	Config        IngestionConfigRequest  `json:"config" description:"Ingestion configuration"`
}

// IngestionConfigRequest is the request body for ingestion configuration.
type IngestionConfigRequest struct {
	Language        models.Language `json:"language" description:"Primary language for SCIP indexing" example:"go" validate:"required,oneof=go typescript"`
	IncludeDocs     bool            `json:"include_docs" description:"Include markdown documentation" example:"true"`
	ExcludePatterns []string        `json:"exclude_patterns,omitempty" description:"Additional glob patterns to exclude"`
	MaxFileSize     *int64          `json:"max_file_size,omitempty" description:"Maximum file size in bytes (default 1MB)"`
	LanguageConfig  json.RawMessage `json:"language_config,omitempty" description:"Language-specific configuration"`
}

// IngestionConfigResponse is the API response for ingestion configuration.
type IngestionConfigResponse struct {
	ID              int64           `json:"id" description:"Config ID"`
	Language        models.Language `json:"language" description:"Primary language for SCIP indexing" example:"go"`
	IncludeDocs     bool            `json:"include_docs" description:"Include markdown documentation"`
	ExcludePatterns []string        `json:"exclude_patterns" description:"Additional glob patterns to exclude"`
	MaxFileSize     int64           `json:"max_file_size" description:"Maximum file size in bytes"`
	LanguageConfig  json.RawMessage `json:"language_config,omitempty" description:"Language-specific configuration"`
}

// Clone returns a deep copy of the RepositoryResponse.
func (r RepositoryResponse) Clone() RepositoryResponse {
	c := r
	if r.Description != nil {
		d := *r.Description
		c.Description = &d
	}
	return c
}

// Clone returns a deep copy of the RepositoryListResponse.
func (r RepositoryListResponse) Clone() RepositoryListResponse {
	c := r
	if r.Repositories != nil {
		c.Repositories = make([]RepositoryResponse, len(r.Repositories))
		for idx, repo := range r.Repositories {
			c.Repositories[idx] = repo.Clone()
		}
	}
	return c
}

// Clone returns a deep copy of the RegisterRepositoryRequest.
func (r RegisterRepositoryRequest) Clone() RegisterRepositoryRequest {
	c := r
	if r.Description != nil {
		d := *r.Description
		c.Description = &d
	}
	c.Config = r.Config.Clone()
	return c
}

// Clone returns a deep copy of the IngestionConfigRequest.
func (r IngestionConfigRequest) Clone() IngestionConfigRequest {
	c := r
	if r.ExcludePatterns != nil {
		c.ExcludePatterns = make([]string, len(r.ExcludePatterns))
		copy(c.ExcludePatterns, r.ExcludePatterns)
	}
	if r.MaxFileSize != nil {
		m := *r.MaxFileSize
		c.MaxFileSize = &m
	}
	if r.LanguageConfig != nil {
		c.LanguageConfig = make(json.RawMessage, len(r.LanguageConfig))
		copy(c.LanguageConfig, r.LanguageConfig)
	}
	return c
}

// Clone returns a deep copy of the IngestionConfigResponse.
func (r IngestionConfigResponse) Clone() IngestionConfigResponse {
	c := r
	if r.ExcludePatterns != nil {
		c.ExcludePatterns = make([]string, len(r.ExcludePatterns))
		copy(c.ExcludePatterns, r.ExcludePatterns)
	}
	if r.LanguageConfig != nil {
		c.LanguageConfig = make(json.RawMessage, len(r.LanguageConfig))
		copy(c.LanguageConfig, r.LanguageConfig)
	}
	return c
}

// Validate validates the RegisterRepositoryRequest.
func (r *RegisterRepositoryRequest) Validate() error {
	if err := check.All(
		check.Int(r.GitHubID, "github_id").Positive().V(),
		check.Str(r.Owner, "owner").Required().MaxLen(255).V(),
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
		check.Str(r.FullName, "full_name").Required().MaxLen(512).V(),
		check.OptStr(r.Description, "description").MaxLen(1000).V(),
		check.Str(r.DefaultBranch, "default_branch").Required().MaxLen(255).V(),
		check.Str(r.HTMLURL, "html_url").Required().URL().V(),
	).Err(); err != nil {
		return err
	}
	return r.Config.Validate()
}

// Validate validates the IngestionConfigRequest.
func (r *IngestionConfigRequest) Validate() error {
	return check.All(
		check.Str(string(r.Language), "language").Required().OneOf([]string{"go", "typescript"}).V(),
	).Err()
}
