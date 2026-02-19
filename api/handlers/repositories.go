package handlers

import (
	"strconv"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/wire"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/api/transformers"
)

// ListRepositories returns all repositories for the authenticated user.
var ListRepositories = rocco.GET("/repositories", func(req *rocco.Request[rocco.NoBody]) (wire.RepositoryListResponse, error) {
	repos := sum.MustUse[contracts.Repositories](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.RepositoryListResponse{}, err
	}

	list, err := repos.ListByUserID(req.Context, userID)
	if err != nil {
		return wire.RepositoryListResponse{}, err
	}

	return transformers.RepositoriesToList(list), nil
}).WithSummary("List repositories").
	WithDescription("Returns all repositories registered by the authenticated user.").
	WithTags("Repositories").
	WithAuthentication()

// RegisterRepository registers a new repository for ingestion.
var RegisterRepository = rocco.POST("/repositories", func(req *rocco.Request[wire.RegisterRepositoryRequest]) (wire.RepositoryResponse, error) {
	repos := sum.MustUse[contracts.Repositories](req.Context)
	configs := sum.MustUse[contracts.IngestionConfigs](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.RepositoryResponse{}, err
	}

	// Create repository
	repo := &models.Repository{UserID: userID}
	transformers.ApplyRepositoryRegistration(req.Body, repo)

	if err := repos.Set(req.Context, "", repo); err != nil {
		return wire.RepositoryResponse{}, err
	}

	// Create ingestion config
	config := &models.IngestionConfig{
		RepositoryID: repo.ID,
		UserID:       userID,
	}
	transformers.ApplyIngestionConfigRequest(req.Body.Config, config)

	if err := configs.Set(req.Context, "", config); err != nil {
		return wire.RepositoryResponse{}, err
	}

	return transformers.RepositoryToResponse(repo), nil
}).WithSummary("Register repository").
	WithDescription("Registers a GitHub repository for ingestion with configuration.").
	WithTags("Repositories").
	WithAuthentication().
	WithSuccessStatus(201)

// GetRepository returns a specific repository.
var GetRepository = rocco.GET("/repositories/{owner}/{repo}", func(req *rocco.Request[rocco.NoBody]) (wire.RepositoryResponse, error) {
	repos := sum.MustUse[contracts.Repositories](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.RepositoryResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]

	repo, err := repos.GetByUserOwnerAndName(req.Context, userID, owner, repoName)
	if err != nil {
		return wire.RepositoryResponse{}, ErrRepositoryNotFound
	}

	return transformers.RepositoryToResponse(repo), nil
}).WithPathParams("owner", "repo").
	WithSummary("Get repository").
	WithDescription("Returns a specific registered repository.").
	WithTags("Repositories").
	WithErrors(ErrRepositoryNotFound).
	WithAuthentication()
