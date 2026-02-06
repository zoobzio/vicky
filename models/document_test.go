package models

import "testing"

func TestDocumentClone(t *testing.T) {
	orig := Document{
		ID:     1,
		Vector: []float32{0.1, 0.2, 0.3},
	}
	clone := orig.Clone()

	clone.Vector[0] = 9.9

	if orig.Vector[0] != 0.1 {
		t.Error("Clone did not isolate Vector")
	}
}

func TestDocumentClone_NilVector(t *testing.T) {
	orig := Document{ID: 1}
	clone := orig.Clone()

	if clone.Vector != nil {
		t.Error("Clone should preserve nil Vector")
	}
	if clone.ID != 1 {
		t.Errorf("ID = %d, want 1", clone.ID)
	}
}
