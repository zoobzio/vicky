package models

import "testing"

func TestRepositoryClone(t *testing.T) {
	desc := "A test repository"
	orig := Repository{
		ID:          1,
		Description: &desc,
	}
	clone := orig.Clone()

	// Modify clone pointer
	*clone.Description = "CHANGED"

	if *orig.Description != "A test repository" {
		t.Error("Clone did not isolate Description pointer")
	}
}

func TestRepositoryClone_NilDescription(t *testing.T) {
	orig := Repository{ID: 1, Description: nil}
	clone := orig.Clone()

	if clone.Description != nil {
		t.Error("Clone should preserve nil Description")
	}
	if clone.ID != 1 {
		t.Errorf("ID = %d, want 1", clone.ID)
	}
}
