package models

import "time"

// VersionStatus represents the ingestion state of a version.
type VersionStatus string

// VersionStatus values.
const (
	VersionStatusPending   VersionStatus = "pending"
	VersionStatusIngesting VersionStatus = "ingesting"
	VersionStatusReady     VersionStatus = "ready"
	VersionStatusFailed    VersionStatus = "failed"
)

// Version represents an ingested snapshot of a repository at a specific tag.
type Version struct {
	ID           int64         `json:"id" db:"id" constraints:"primarykey" description:"Internal version ID"`
	RepositoryID int64         `json:"repository_id" db:"repository_id" constraints:"notnull" references:"repositories(id)" description:"Parent repository"`
	UserID       int64         `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Owner        string        `json:"owner" db:"owner" constraints:"notnull" description:"GitHub org or user" example:"octocat"`
	RepoName     string        `json:"repo_name" db:"repo_name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	Tag          string        `json:"tag" db:"tag" constraints:"notnull" description:"Version tag" example:"v1.0.0"`
	CommitSHA    string        `json:"commit_sha" db:"commit_sha" constraints:"notnull" description:"Git commit SHA" example:"a1b2c3d4e5f6"`
	Status       VersionStatus `json:"status" db:"status" constraints:"notnull" default:"'pending'" description:"Ingestion status"`
	Error        *string       `json:"error,omitempty" db:"error" description:"Ingestion error if failed"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at" default:"now()" description:"Ingestion start time"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at" default:"now()" description:"Last status update"`
}

// Clone returns a deep copy of the Version.
func (v Version) Clone() Version {
	c := v
	if v.Error != nil {
		e := *v.Error
		c.Error = &e
	}
	return c
}
