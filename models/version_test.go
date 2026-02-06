package models

import "testing"

func TestVersionClone(t *testing.T) {
	errMsg := "something went wrong"
	orig := Version{
		ID:    1,
		Error: &errMsg,
	}
	clone := orig.Clone()

	// Modify clone pointer
	*clone.Error = "CHANGED"

	if *orig.Error != "something went wrong" {
		t.Error("Clone did not isolate Error pointer")
	}
}

func TestVersionClone_NilError(t *testing.T) {
	orig := Version{ID: 1, Error: nil}
	clone := orig.Clone()

	if clone.Error != nil {
		t.Error("Clone should preserve nil Error")
	}
	if clone.ID != 1 {
		t.Errorf("ID = %d, want 1", clone.ID)
	}
}
