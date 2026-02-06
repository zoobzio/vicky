package transformers

import (
	"github.com/zoobzio/vicky/wire"
	"github.com/zoobzio/vicky/models"
)

// RepositoryToResponse transforms a Repository model to an API response.
func RepositoryToResponse(r *models.Repository) wire.RepositoryResponse {
	return wire.RepositoryResponse{
		ID:            r.ID,
		GitHubID:      r.GitHubID,
		Owner:         r.Owner,
		Name:          r.Name,
		FullName:      r.FullName,
		Description:   r.Description,
		DefaultBranch: r.DefaultBranch,
		Private:       r.Private,
		HTMLURL:       r.HTMLURL,
	}
}

// RepositoriesToList transforms a slice of Repository models to an API list response.
func RepositoriesToList(repos []*models.Repository) wire.RepositoryListResponse {
	resp := wire.RepositoryListResponse{
		Repositories: make([]wire.RepositoryResponse, len(repos)),
	}
	for i, r := range repos {
		resp.Repositories[i] = RepositoryToResponse(r)
	}
	return resp
}

// ApplyRepositoryRegistration applies a RegisterRepositoryRequest to a Repository model.
func ApplyRepositoryRegistration(req wire.RegisterRepositoryRequest, r *models.Repository) {
	r.GitHubID = req.GitHubID
	r.Owner = req.Owner
	r.Name = req.Name
	r.FullName = req.FullName
	r.Description = req.Description
	r.DefaultBranch = req.DefaultBranch
	r.Private = req.Private
	r.HTMLURL = req.HTMLURL
}
