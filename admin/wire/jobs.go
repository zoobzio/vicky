package wire

import "time"

// AdminJobResponse is the API response for admin job data.
type AdminJobResponse struct {
	ID             int64      `json:"id" description:"Job ID" example:"12345"`
	VersionID      int64      `json:"version_id" description:"Parent version ID"`
	RepositoryID   int64      `json:"repository_id" description:"Parent repository ID"`
	UserID         int64      `json:"user_id" description:"Owning user ID"`
	Owner          string     `json:"owner" description:"Repository owner" example:"octocat"`
	RepoName       string     `json:"repo_name" description:"Repository name" example:"hello-world"`
	Tag            string     `json:"tag" description:"Version tag" example:"v1.0.0"`
	Stage          string     `json:"stage" description:"Current processing stage" example:"embed"`
	Status         string     `json:"status" description:"Job status" example:"running"`
	Progress       int        `json:"progress" description:"Percentage completion 0-100" example:"45"`
	Error          *string    `json:"error,omitempty" description:"Error message if failed"`
	ItemsTotal     int        `json:"items_total" description:"Total items to process"`
	ItemsProcessed int        `json:"items_processed" description:"Items processed so far"`
	CreatedAt      time.Time  `json:"created_at" description:"Job creation time"`
	StartedAt      *time.Time `json:"started_at,omitempty" description:"Processing start time"`
	CompletedAt    *time.Time `json:"completed_at,omitempty" description:"Processing completion time"`
	UpdatedAt      time.Time  `json:"updated_at" description:"Last update time"`
}

// Clone returns a deep copy.
func (r AdminJobResponse) Clone() AdminJobResponse {
	c := r
	if r.Error != nil {
		e := *r.Error
		c.Error = &e
	}
	if r.StartedAt != nil {
		s := *r.StartedAt
		c.StartedAt = &s
	}
	if r.CompletedAt != nil {
		comp := *r.CompletedAt
		c.CompletedAt = &comp
	}
	return c
}

// AdminJobListResponse is the API response for listing jobs.
type AdminJobListResponse struct {
	Jobs   []AdminJobResponse `json:"jobs" description:"Job array"`
	Total  int                `json:"total" description:"Total count matching filter"`
	Limit  int                `json:"limit" description:"Limit used"`
	Offset int                `json:"offset" description:"Offset used"`
}

// Clone returns a deep copy.
func (r AdminJobListResponse) Clone() AdminJobListResponse {
	c := r
	if r.Jobs != nil {
		c.Jobs = make([]AdminJobResponse, len(r.Jobs))
		for i, j := range r.Jobs {
			c.Jobs[i] = j.Clone()
		}
	}
	return c
}

// AdminJobStatsResponse provides aggregate job statistics.
type AdminJobStatsResponse struct {
	TotalJobs      int `json:"total_jobs" description:"Total number of jobs"`
	PendingJobs    int `json:"pending_jobs" description:"Jobs waiting to run"`
	RunningJobs    int `json:"running_jobs" description:"Jobs currently running"`
	CompletedJobs  int `json:"completed_jobs" description:"Successfully completed jobs"`
	FailedJobs     int `json:"failed_jobs" description:"Failed jobs"`
	CancellingJobs int `json:"cancelling_jobs" description:"Jobs marked for cancellation"`
	CancelledJobs  int `json:"cancelled_jobs" description:"Cancelled jobs"`
}

// Clone returns a deep copy.
func (r AdminJobStatsResponse) Clone() AdminJobStatsResponse {
	return r
}
