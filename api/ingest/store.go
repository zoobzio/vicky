package ingest

import (
	"context"
	"fmt"

	"github.com/zoobzio/pipz"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/api/events"
	"github.com/zoobzio/vicky/models"
)

// Stage identity.
var StoreStageID = pipz.NewIdentity("store", "Finalizes ingestion and marks version ready")

// storeStage finalizes ingestion by marking the version as ready.
func storeStage(ctx context.Context, job *models.Job) (*models.Job, error) {
	events.Ingest.Store.Started.Emit(ctx, events.StoreEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
	})

	// Update job stage
	job.Stage = models.JobStageStore

	// Resolve dependencies
	versions := sum.MustUse[contracts.Versions](ctx)

	// Mark version as ready
	if _, err := versions.UpdateStatus(ctx, job.VersionID, models.VersionStatusReady, nil); err != nil {
		return job, fmt.Errorf("update version status: %w", err)
	}

	events.Ingest.Store.Completed.Emit(ctx, events.StoreEvent{
		RepositoryID: job.RepositoryID,
		VersionID:    job.VersionID,
	})

	return job, nil
}
