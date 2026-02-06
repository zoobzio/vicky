package transformers

import (
	"github.com/zoobzio/vicky/wire"
	"github.com/zoobzio/vicky/models"
)

// VersionToResponse transforms a Version model to an API response.
func VersionToResponse(v *models.Version) wire.VersionResponse {
	return wire.VersionResponse{
		ID:           v.ID,
		RepositoryID: v.RepositoryID,
		Owner:        v.Owner,
		RepoName:     v.RepoName,
		Tag:          v.Tag,
		CommitSHA:    v.CommitSHA,
		Status:       v.Status,
		Error:        v.Error,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

// VersionsToList transforms a slice of Version models to an API list response.
func VersionsToList(versions []*models.Version) wire.VersionListResponse {
	resp := wire.VersionListResponse{
		Versions: make([]wire.VersionResponse, len(versions)),
	}
	for i, v := range versions {
		resp.Versions[i] = VersionToResponse(v)
	}
	return resp
}

// ApplyIngestRequest applies an IngestRequest to a Version model.
func ApplyIngestRequest(req wire.IngestRequest, v *models.Version) {
	v.CommitSHA = req.CommitSHA
}
