package models

import "testing"

func TestChunkClone(t *testing.T) {
	sym := "NewClient"
	orig := Chunk{
		ID:      1,
		Symbol:  &sym,
		Context: []string{"type UserService", "method GetUser"},
		Vector:  []float32{0.5, 0.6},
	}
	clone := orig.Clone()

	// Modify clone fields
	*clone.Symbol = "CHANGED"
	clone.Context[0] = "CHANGED"
	clone.Vector[0] = 9.9

	if *orig.Symbol != "NewClient" {
		t.Error("Clone did not isolate Symbol")
	}
	if orig.Context[0] != "type UserService" {
		t.Error("Clone did not isolate Context")
	}
	if orig.Vector[0] != 0.5 {
		t.Error("Clone did not isolate Vector")
	}
}

func TestChunkClone_NilPointers(t *testing.T) {
	orig := Chunk{ID: 1}
	clone := orig.Clone()

	if clone.Symbol != nil || clone.Context != nil || clone.Vector != nil {
		t.Error("Clone should preserve nil fields")
	}
}
