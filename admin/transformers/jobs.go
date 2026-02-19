package transformers

import (
	"github.com/zoobzio/vicky/admin/wire"
	"github.com/zoobzio/vicky/models"
)

// JobToAdminResponse transforms a Job model to an admin API response.
func JobToAdminResponse(j *models.Job) wire.AdminJobResponse {
	return wire.AdminJobResponse{
		ID:             j.ID,
		VersionID:      j.VersionID,
		RepositoryID:   j.RepositoryID,
		UserID:         j.UserID,
		Owner:          j.Owner,
		RepoName:       j.RepoName,
		Tag:            j.Tag,
		Stage:          string(j.Stage),
		Status:         string(j.Status),
		Progress:       j.Progress,
		Error:          j.Error,
		ItemsTotal:     j.ItemsTotal,
		ItemsProcessed: j.ItemsProcessed,
		CreatedAt:      j.CreatedAt,
		StartedAt:      j.StartedAt,
		CompletedAt:    j.CompletedAt,
		UpdatedAt:      j.UpdatedAt,
	}
}

// JobsToAdminList transforms a slice of Job models to a paginated admin list response.
func JobsToAdminList(jobs []*models.Job, total, limit, offset int) wire.AdminJobListResponse {
	resp := wire.AdminJobListResponse{
		Jobs:   make([]wire.AdminJobResponse, len(jobs)),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	for i, j := range jobs {
		resp.Jobs[i] = JobToAdminResponse(j)
	}
	return resp
}
