package models

import "time"

// JobStage represents the processing stage of an ingestion job.
type JobStage string

// JobStage values.
const (
	JobStageFetch JobStage = "fetch"
	JobStageParse JobStage = "parse"
	JobStageChunk JobStage = "chunk"
	JobStageEmbed JobStage = "embed"
	JobStageStore JobStage = "store"
)

// JobStatus represents the overall status of a job.
type JobStatus string

// JobStatus values.
const (
	JobStatusPending    JobStatus = "pending"
	JobStatusRunning    JobStatus = "running"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelling JobStatus = "cancelling"
	JobStatusCancelled  JobStatus = "cancelled"
)

// Job tracks an async ingestion pipeline execution.
type Job struct {
	ID           int64      `json:"id" db:"id" constraints:"primarykey" description:"Job ID"`
	VersionID    int64      `json:"version_id" db:"version_id" constraints:"notnull" references:"versions(id)" description:"Parent version"`
	RepositoryID int64      `json:"repository_id" db:"repository_id" constraints:"notnull" references:"repositories(id)" description:"Parent repository"`
	UserID       int64      `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Owner        string     `json:"owner" db:"owner" constraints:"notnull" description:"Repository owner" example:"octocat"`
	RepoName       string     `json:"repo_name" db:"repo_name" constraints:"notnull" description:"Repository name" example:"hello-world"`
	Tag            string     `json:"tag" db:"tag" constraints:"notnull" description:"Version tag" example:"v1.0.0"`
	Stage          JobStage   `json:"stage" db:"stage" constraints:"notnull" default:"'fetch'" description:"Current processing stage"`
	Status         JobStatus  `json:"status" db:"status" constraints:"notnull" default:"'pending'" description:"Overall job status"`
	Progress       int        `json:"progress" db:"progress" constraints:"notnull" default:"0" description:"Percentage completion 0-100"`
	Error          *string    `json:"error,omitempty" db:"error" description:"Error message if failed"`
	ItemsTotal     int        `json:"items_total" db:"items_total" constraints:"notnull" default:"0" description:"Total items to process"`
	ItemsProcessed int        `json:"items_processed" db:"items_processed" constraints:"notnull" default:"0" description:"Items processed so far"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at" default:"now()" description:"Job creation time"`
	StartedAt      *time.Time `json:"started_at,omitempty" db:"started_at" description:"Processing start time"`
	CompletedAt    *time.Time `json:"completed_at,omitempty" db:"completed_at" description:"Processing completion time"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at" default:"now()" description:"Last update time"`
}

// Clone returns a deep copy of the Job.
// Pointer receiver implements pipz.Cloner[*Job] for use in pipelines.
func (j *Job) Clone() *Job {
	if j == nil {
		return nil
	}
	c := *j
	if j.Error != nil {
		e := *j.Error
		c.Error = &e
	}
	if j.StartedAt != nil {
		s := *j.StartedAt
		c.StartedAt = &s
	}
	if j.CompletedAt != nil {
		comp := *j.CompletedAt
		c.CompletedAt = &comp
	}
	return &c
}
