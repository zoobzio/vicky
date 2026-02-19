package handlers

import (
	"strconv"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	admincontracts "github.com/zoobzio/vicky/admin/contracts"
	"github.com/zoobzio/vicky/admin/transformers"
	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/api/events"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/stores"
)

// ListJobs returns all jobs with pagination and optional filtering.
var ListJobs = rocco.GET("/admin/jobs", func(req *rocco.Request[rocco.NoBody]) (wire.AdminJobListResponse, error) {
	jobsStore := sum.MustUse[admincontracts.Jobs](req.Context)

	// Parse limit with validation
	limit := 50
	if l := req.Params.Query["limit"]; l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 || parsed > 100 {
			return wire.AdminJobListResponse{}, ErrInvalidLimit
		}
		limit = parsed
	}

	// Parse offset
	offset := 0
	if o := req.Params.Query["offset"]; o != "" {
		parsed, err := strconv.Atoi(o)
		if err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Build filter from query params
	var filter *stores.JobFilter
	userIDFilter := req.Params.Query["user_id"]
	repositoryIDFilter := req.Params.Query["repository_id"]
	versionIDFilter := req.Params.Query["version_id"]
	statusFilter := req.Params.Query["status"]
	ownerFilter := req.Params.Query["owner"]
	repoNameFilter := req.Params.Query["repo_name"]

	if userIDFilter != "" || repositoryIDFilter != "" || versionIDFilter != "" || statusFilter != "" || ownerFilter != "" || repoNameFilter != "" {
		filter = &stores.JobFilter{}

		if userIDFilter != "" {
			parsed, err := strconv.ParseInt(userIDFilter, 10, 64)
			if err != nil {
				return wire.AdminJobListResponse{}, ErrInvalidUserID
			}
			filter.UserID = &parsed
		}

		if repositoryIDFilter != "" {
			parsed, err := strconv.ParseInt(repositoryIDFilter, 10, 64)
			if err != nil {
				return wire.AdminJobListResponse{}, rocco.ErrBadRequest.WithMessage("invalid repository_id parameter")
			}
			filter.RepositoryID = &parsed
		}

		if versionIDFilter != "" {
			parsed, err := strconv.ParseInt(versionIDFilter, 10, 64)
			if err != nil {
				return wire.AdminJobListResponse{}, rocco.ErrBadRequest.WithMessage("invalid version_id parameter")
			}
			filter.VersionID = &parsed
		}

		if statusFilter != "" {
			status := models.JobStatus(statusFilter)
			filter.Status = &status
		}

		if ownerFilter != "" {
			filter.Owner = &ownerFilter
		}

		if repoNameFilter != "" {
			filter.RepoName = &repoNameFilter
		}
	}

	// Execute list query
	jobs, err := jobsStore.List(req.Context, filter, limit, offset)
	if err != nil {
		return wire.AdminJobListResponse{}, err
	}

	// Get total count
	total, err := jobsStore.Count(req.Context, filter)
	if err != nil {
		return wire.AdminJobListResponse{}, err
	}

	return transformers.JobsToAdminList(jobs, total, limit, offset), nil
}).WithSummary("List jobs").
	WithDescription("Returns all jobs with pagination and optional filtering by user_id, repository_id, version_id, status, owner, or repo_name.").
	WithTags("Admin", "Jobs").
	WithQueryParams("limit", "offset", "user_id", "repository_id", "version_id", "status", "owner", "repo_name").
	WithErrors(ErrInvalidLimit, ErrInvalidUserID)

// GetJob returns a single job by ID.
var GetJob = rocco.GET("/admin/jobs/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.AdminJobResponse, error) {
	jobsStore := sum.MustUse[admincontracts.Jobs](req.Context)

	id := req.Params.Path["id"]

	job, err := jobsStore.Get(req.Context, id)
	if err != nil {
		return wire.AdminJobResponse{}, ErrJobNotFound
	}

	return transformers.JobToAdminResponse(job), nil
}).WithSummary("Get job").
	WithDescription("Returns a single job by ID.").
	WithTags("Admin", "Jobs").
	WithPathParams("id").
	WithErrors(ErrJobNotFound)

// CancelJob requests cancellation of a job.
var CancelJob = rocco.POST("/admin/jobs/{id}/cancel", func(req *rocco.Request[rocco.NoBody]) (wire.AdminJobResponse, error) {
	jobsStore := sum.MustUse[admincontracts.Jobs](req.Context)

	id := req.Params.Path["id"]

	// Parse ID to int64
	jobID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return wire.AdminJobResponse{}, ErrJobNotFound
	}

	// Request cancellation
	if err := jobsStore.RequestCancellation(req.Context, jobID); err != nil {
		return wire.AdminJobResponse{}, ErrJobNotCancellable
	}

	// Retrieve updated job
	job, err := jobsStore.Get(req.Context, id)
	if err != nil {
		return wire.AdminJobResponse{}, ErrJobNotFound
	}

	return transformers.JobToAdminResponse(job), nil
}).WithSummary("Cancel job").
	WithDescription("Requests cancellation of a running or pending job. The worker will abort at the next stage boundary.").
	WithTags("Admin", "Jobs").
	WithPathParams("id").
	WithErrors(ErrJobNotFound, ErrJobNotCancellable)

// RetryJob creates a new job for a failed/cancelled job's version.
var RetryJob = rocco.POST("/admin/jobs/{id}/retry", func(req *rocco.Request[rocco.NoBody]) (wire.AdminJobResponse, error) {
	jobsStore := sum.MustUse[admincontracts.Jobs](req.Context)

	id := req.Params.Path["id"]

	// Get original job
	originalJob, err := jobsStore.Get(req.Context, id)
	if err != nil {
		return wire.AdminJobResponse{}, ErrJobNotFound
	}

	// Create new pending job for same version
	newJob := &models.Job{
		VersionID:    originalJob.VersionID,
		RepositoryID: originalJob.RepositoryID,
		UserID:       originalJob.UserID,
		Owner:        originalJob.Owner,
		RepoName:     originalJob.RepoName,
		Tag:          originalJob.Tag,
		Stage:        models.JobStageFetch,
		Status:       models.JobStatusPending,
		Progress:     0,
		ItemsTotal:   0,
		ItemsProcessed: 0,
	}

	// Save new job
	if err := jobsStore.Set(req.Context, "", newJob); err != nil {
		return wire.AdminJobResponse{}, err
	}

	// Emit job creation event
	events.Job.Created.Emit(req.Context, events.JobCreatedEvent{
		Job: newJob,
	})

	return transformers.JobToAdminResponse(newJob), nil
}).WithSummary("Retry job").
	WithDescription("Creates a new pending job for the same version as the specified job. Used to retry failed or cancelled jobs.").
	WithTags("Admin", "Jobs").
	WithPathParams("id").
	WithErrors(ErrJobNotFound).
	WithSuccessStatus(201)

// GetJobStats returns aggregate statistics for jobs.
var GetJobStats = rocco.GET("/admin/jobs/stats", func(req *rocco.Request[rocco.NoBody]) (wire.AdminJobStatsResponse, error) {
	jobsStore := sum.MustUse[admincontracts.Jobs](req.Context)

	// Optional filter by user_id
	var userID *int64
	if u := req.Params.Query["user_id"]; u != "" {
		parsed, err := strconv.ParseInt(u, 10, 64)
		if err != nil {
			return wire.AdminJobStatsResponse{}, ErrInvalidUserID
		}
		userID = &parsed
	}

	// Get counts by status
	counts, err := jobsStore.CountByStatus(req.Context, userID)
	if err != nil {
		return wire.AdminJobStatsResponse{}, err
	}

	// Build response
	total := 0
	for _, count := range counts {
		total += count
	}

	return wire.AdminJobStatsResponse{
		TotalJobs:      total,
		PendingJobs:    counts[models.JobStatusPending],
		RunningJobs:    counts[models.JobStatusRunning],
		CompletedJobs:  counts[models.JobStatusCompleted],
		FailedJobs:     counts[models.JobStatusFailed],
		CancellingJobs: counts[models.JobStatusCancelling],
		CancelledJobs:  counts[models.JobStatusCancelled],
	}, nil
}).WithSummary("Get job statistics").
	WithDescription("Returns aggregate statistics for jobs, optionally filtered by user_id.").
	WithTags("Admin", "Jobs").
	WithQueryParams("user_id").
	WithErrors(ErrInvalidUserID)
