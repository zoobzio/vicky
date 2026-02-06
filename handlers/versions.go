package handlers

import (
	"strconv"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/contracts"
	"github.com/zoobzio/vicky/events"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/transformers"
	"github.com/zoobzio/vicky/wire"
)

// ListVersions returns all versions for a repository.
var ListVersions = rocco.GET("/repositories/{owner}/{repo}/versions", func(req *rocco.Request[rocco.NoBody]) (wire.VersionListResponse, error) {
	versions := sum.MustUse[contracts.Versions](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.VersionListResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]

	list, err := versions.ListByUserAndRepo(req.Context, userID, owner, repoName)
	if err != nil {
		return wire.VersionListResponse{}, err
	}

	return transformers.VersionsToList(list), nil
}).WithPathParams("owner", "repo").
	WithSummary("List versions").
	WithDescription("Returns all ingested versions for a repository.").
	WithTags("Versions").
	WithAuthentication()

// GetVersion returns a specific version.
var GetVersion = rocco.GET("/repositories/{owner}/{repo}/versions/{tag}", func(req *rocco.Request[rocco.NoBody]) (wire.VersionResponse, error) {
	versions := sum.MustUse[contracts.Versions](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.VersionResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]

	version, err := versions.GetByUserRepoAndTag(req.Context, userID, owner, repoName, tag)
	if err != nil {
		return wire.VersionResponse{}, ErrVersionNotFound
	}

	return transformers.VersionToResponse(version), nil
}).WithPathParams("owner", "repo", "tag").
	WithSummary("Get version").
	WithDescription("Returns a specific version's ingestion status.").
	WithTags("Versions").
	WithErrors(ErrVersionNotFound).
	WithAuthentication()

// TriggerIngest initiates ingestion for a version.
var TriggerIngest = rocco.POST("/repositories/{owner}/{repo}/versions/{tag}", func(req *rocco.Request[wire.IngestRequest]) (wire.VersionResponse, error) {
	versions := sum.MustUse[contracts.Versions](req.Context)
	repos := sum.MustUse[contracts.Repositories](req.Context)
	jobs := sum.MustUse[contracts.Jobs](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.VersionResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]

	// Find the repository
	repo, err := repos.GetByUserOwnerAndName(req.Context, userID, owner, repoName)
	if err != nil {
		return wire.VersionResponse{}, ErrRepositoryNotFound
	}

	// Create pending version
	version := &models.Version{
		RepositoryID: repo.ID,
		UserID:       userID,
		Owner:        owner,
		RepoName:     repoName,
		Tag:          tag,
		Status:       models.VersionStatusPending,
	}
	transformers.ApplyIngestRequest(req.Body, version)

	if err := versions.Set(req.Context, "", version); err != nil {
		return wire.VersionResponse{}, err
	}

	// Create pending job for async processing
	job := &models.Job{
		VersionID:    version.ID,
		RepositoryID: repo.ID,
		UserID:       userID,
		Owner:        owner,
		RepoName:     repoName,
		Tag:          tag,
		Stage:        models.JobStageFetch,
		Status:       models.JobStatusPending,
	}

	if err := jobs.Set(req.Context, "", job); err != nil {
		return wire.VersionResponse{}, err
	}

	// Emit event to trigger async ingestion via worker pool
	events.Job.Created.Emit(req.Context, events.JobCreatedEvent{Job: job})

	return transformers.VersionToResponse(version), nil
}).WithPathParams("owner", "repo", "tag").
	WithSummary("Trigger ingestion").
	WithDescription("Initiates ingestion for a repository version.").
	WithTags("Versions").
	WithErrors(ErrRepositoryNotFound).
	WithAuthentication().
	WithSuccessStatus(202)
