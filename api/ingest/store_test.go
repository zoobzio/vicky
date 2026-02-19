//go:build testing

package ingest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	vickytest "github.com/zoobzio/vicky/testing"

	"github.com/zoobzio/vicky/models"
)

func TestStoreStage(t *testing.T) {
	var calledID int64
	var calledStatus models.VersionStatus
	var calledErr *string

	mv := &vickytest.MockVersions{
		OnUpdateStatus: func(ctx context.Context, id int64, status models.VersionStatus, versionErr *string) (*models.Version, error) {
			calledID = id
			calledStatus = status
			calledErr = versionErr
			return &models.Version{ID: id, Status: status}, nil
		},
	}

	ctx := vickytest.SetupRegistry(t, vickytest.WithVersions(mv))
	job := vickytest.NewJob(t)

	result, err := storeStage(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Stage != models.JobStageStore {
		t.Errorf("Stage = %q, want %q", result.Stage, models.JobStageStore)
	}
	if calledID != job.VersionID {
		t.Errorf("UpdateStatus called with id=%d, want %d", calledID, job.VersionID)
	}
	if calledStatus != models.VersionStatusReady {
		t.Errorf("UpdateStatus called with status=%q, want %q", calledStatus, models.VersionStatusReady)
	}
	if calledErr != nil {
		t.Errorf("UpdateStatus called with non-nil error: %v", calledErr)
	}
}

func TestStoreStage_UpdateStatusError(t *testing.T) {
	mv := &vickytest.MockVersions{
		OnUpdateStatus: func(ctx context.Context, id int64, status models.VersionStatus, versionErr *string) (*models.Version, error) {
			return nil, fmt.Errorf("db down")
		},
	}

	ctx := vickytest.SetupRegistry(t, vickytest.WithVersions(mv))
	job := vickytest.NewJob(t)

	result, err := storeStage(ctx, job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "update version status") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "update version status")
	}
	if result == nil {
		t.Error("expected job to be returned even on error")
	}
}
