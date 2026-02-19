package handlers

import (
	"strconv"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	admincontracts "github.com/zoobzio/vicky/admin/contracts"
	"github.com/zoobzio/vicky/admin/transformers"
	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/stores"
)

// ListRepositories returns all repositories with pagination and optional filtering.
var ListRepositories = rocco.GET("/admin/repositories", func(req *rocco.Request[rocco.NoBody]) (wire.AdminRepositoryListResponse, error) {
	reposStore := sum.MustUse[admincontracts.Repositories](req.Context)

	// Parse limit with validation
	limit := 50
	if l := req.Params.Query["limit"]; l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 || parsed > 100 {
			return wire.AdminRepositoryListResponse{}, ErrInvalidLimit
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
	var filter *stores.RepositoryFilter
	ownerFilter := req.Params.Query["owner"]
	nameFilter := req.Params.Query["name"]
	userIDFilter := req.Params.Query["user_id"]

	if ownerFilter != "" || nameFilter != "" || userIDFilter != "" {
		filter = &stores.RepositoryFilter{}
		if ownerFilter != "" {
			filter.Owner = &ownerFilter
		}
		if nameFilter != "" {
			filter.Name = &nameFilter
		}
		if userIDFilter != "" {
			parsed, err := strconv.ParseInt(userIDFilter, 10, 64)
			if err != nil {
				return wire.AdminRepositoryListResponse{}, ErrInvalidUserID
			}
			filter.UserID = &parsed
		}
	}

	// Execute list query
	repos, err := reposStore.List(req.Context, filter, limit, offset)
	if err != nil {
		return wire.AdminRepositoryListResponse{}, err
	}

	// Get total count
	total, err := reposStore.Count(req.Context, filter)
	if err != nil {
		return wire.AdminRepositoryListResponse{}, err
	}

	return transformers.RepositoriesToAdminList(repos, total, limit, offset), nil
}).WithSummary("List repositories").
	WithDescription("Returns all repositories with pagination and optional filtering by owner, name, or user_id.").
	WithTags("Admin", "Repositories").
	WithQueryParams("limit", "offset", "owner", "name", "user_id").
	WithErrors(ErrInvalidLimit, ErrInvalidUserID)

// GetRepository returns a single repository by ID.
var GetRepository = rocco.GET("/admin/repositories/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.AdminRepositoryResponse, error) {
	reposStore := sum.MustUse[admincontracts.Repositories](req.Context)

	id := req.Params.Path["id"]

	repo, err := reposStore.Get(req.Context, id)
	if err != nil {
		return wire.AdminRepositoryResponse{}, ErrRepositoryNotFound
	}

	return transformers.RepositoryToAdminResponse(repo), nil
}).WithSummary("Get repository").
	WithDescription("Returns a single repository by ID.").
	WithTags("Admin", "Repositories").
	WithPathParams("id").
	WithErrors(ErrRepositoryNotFound)

// UpdateRepository updates repository fields.
var UpdateRepository = rocco.PATCH("/admin/repositories/{id}", func(req *rocco.Request[wire.AdminRepositoryUpdateRequest]) (wire.AdminRepositoryResponse, error) {
	reposStore := sum.MustUse[admincontracts.Repositories](req.Context)

	id := req.Params.Path["id"]

	// Get existing repository
	repo, err := reposStore.Get(req.Context, id)
	if err != nil {
		return wire.AdminRepositoryResponse{}, ErrRepositoryNotFound
	}

	// Apply update
	transformers.ApplyAdminRepositoryUpdate(req.Body, repo)

	// Save updated repository
	if err := reposStore.Set(req.Context, id, repo); err != nil {
		return wire.AdminRepositoryResponse{}, err
	}

	return transformers.RepositoryToAdminResponse(repo), nil
}).WithSummary("Update repository").
	WithDescription("Updates repository fields. Only provided fields are updated.").
	WithTags("Admin", "Repositories").
	WithPathParams("id").
	WithErrors(ErrRepositoryNotFound)

// DeleteRepository hard deletes a repository with cascade.
var DeleteRepository = rocco.DELETE("/admin/repositories/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	reposStore := sum.MustUse[admincontracts.Repositories](req.Context)

	id := req.Params.Path["id"]

	// Delete will cascade via database constraints
	if err := reposStore.Delete(req.Context, id); err != nil {
		return rocco.NoBody{}, ErrRepositoryNotFound
	}

	return rocco.NoBody{}, nil
}).WithSummary("Delete repository").
	WithDescription("Hard deletes a repository with cascade. Removes all repository data including versions, documents, chunks, symbols, and SCIP data.").
	WithTags("Admin", "Repositories").
	WithPathParams("id").
	WithErrors(ErrRepositoryNotFound).
	WithSuccessStatus(204)
