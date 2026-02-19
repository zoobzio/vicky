package transformers

import (
	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/models"
)

// RepositoryToAdminResponse transforms a Repository model to an admin API response.
// Admin responses include all fields without masking.
func RepositoryToAdminResponse(r *models.Repository) wire.AdminRepositoryResponse {
	return wire.AdminRepositoryResponse{
		ID:            r.ID,
		GitHubID:      r.GitHubID,
		UserID:        r.UserID,
		Owner:         r.Owner,
		Name:          r.Name,
		FullName:      r.FullName,
		Description:   r.Description,
		DefaultBranch: r.DefaultBranch,
		Private:       r.Private,
		HTMLURL:       r.HTMLURL,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

// RepositoriesToAdminList transforms a slice of Repository models to a paginated admin list response.
func RepositoriesToAdminList(repos []*models.Repository, total, limit, offset int) wire.AdminRepositoryListResponse {
	resp := wire.AdminRepositoryListResponse{
		Repositories: make([]wire.AdminRepositoryResponse, len(repos)),
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	}
	for i, r := range repos {
		resp.Repositories[i] = RepositoryToAdminResponse(r)
	}
	return resp
}

// ApplyAdminRepositoryUpdate applies an AdminRepositoryUpdateRequest to a Repository model.
// Only updates fields that are present in the request.
func ApplyAdminRepositoryUpdate(req wire.AdminRepositoryUpdateRequest, r *models.Repository) {
	if req.Owner != nil {
		r.Owner = *req.Owner
	}
	if req.Name != nil {
		r.Name = *req.Name
	}
	if req.FullName != nil {
		r.FullName = *req.FullName
	}
	if req.Description != nil {
		r.Description = req.Description
	}
	if req.DefaultBranch != nil {
		r.DefaultBranch = *req.DefaultBranch
	}
	if req.Private != nil {
		r.Private = *req.Private
	}
	if req.HTMLURL != nil {
		r.HTMLURL = *req.HTMLURL
	}
}
