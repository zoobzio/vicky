package wire

import (
	"time"

	"github.com/zoobzio/check"
	"github.com/zoobzio/vicky/models"
)

// VersionResponse is the API response for version data.
type VersionResponse struct {
	ID           int64                `json:"id" description:"Version ID"`
	RepositoryID int64                `json:"repository_id" description:"Parent repository ID"`
	Owner        string               `json:"owner" description:"Repository owner" example:"octocat"`
	RepoName     string               `json:"repo_name" description:"Repository name" example:"hello-world"`
	Tag          string               `json:"tag" description:"Version tag" example:"v1.0.0"`
	CommitSHA    string               `json:"commit_sha" description:"Git commit SHA" example:"abc123def456"`
	Status       models.VersionStatus `json:"status" description:"Ingestion status" example:"ready"`
	Error        *string              `json:"error,omitempty" description:"Error message if failed"`
	CreatedAt    time.Time            `json:"created_at" description:"Creation timestamp"`
	UpdatedAt    time.Time            `json:"updated_at" description:"Last update timestamp"`
}

// VersionListResponse is the API response for listing versions.
type VersionListResponse struct {
	Versions []VersionResponse `json:"versions" description:"List of versions"`
}

// IngestRequest is the request body for triggering ingestion.
type IngestRequest struct {
	CommitSHA string `json:"commit_sha" description:"Git commit SHA to ingest" example:"abc123def456" validate:"required,len=40"`
}

// Clone returns a deep copy of the VersionResponse.
func (v VersionResponse) Clone() VersionResponse {
	c := v
	if v.Error != nil {
		e := *v.Error
		c.Error = &e
	}
	return c
}

// Clone returns a deep copy of the VersionListResponse.
func (v VersionListResponse) Clone() VersionListResponse {
	c := v
	if v.Versions != nil {
		c.Versions = make([]VersionResponse, len(v.Versions))
		for idx, ver := range v.Versions {
			c.Versions[idx] = ver.Clone()
		}
	}
	return c
}

// Clone returns a deep copy of the IngestRequest.
func (r IngestRequest) Clone() IngestRequest { return r }

// Validate validates the IngestRequest.
func (r *IngestRequest) Validate() error {
	return check.All(
		check.Str(r.CommitSHA, "commit_sha").Required().Len(40).V(),
	).Err()
}
