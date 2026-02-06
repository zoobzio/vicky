package models

import (
	"testing"
	"time"
)

func TestJobClone(t *testing.T) {
	errMsg := "failed"
	started := time.Now()
	completed := started.Add(time.Minute)

	orig := &Job{
		ID:          1,
		Error:       &errMsg,
		StartedAt:   &started,
		CompletedAt: &completed,
	}
	clone := orig.Clone()

	// Modify clone pointers
	*clone.Error = "CHANGED"
	*clone.StartedAt = time.Time{}
	*clone.CompletedAt = time.Time{}

	if *orig.Error != "failed" {
		t.Error("Clone did not isolate Error pointer")
	}
	if orig.StartedAt.IsZero() {
		t.Error("Clone did not isolate StartedAt pointer")
	}
	if orig.CompletedAt.IsZero() {
		t.Error("Clone did not isolate CompletedAt pointer")
	}
}

func TestJobClone_NilPointers(t *testing.T) {
	orig := &Job{ID: 1}
	clone := orig.Clone()

	if clone.Error != nil || clone.StartedAt != nil || clone.CompletedAt != nil {
		t.Error("Clone should preserve nil pointers")
	}
}

func TestJobClone_Nil(t *testing.T) {
	var j *Job
	if j.Clone() != nil {
		t.Error("Clone of nil should return nil")
	}
}
