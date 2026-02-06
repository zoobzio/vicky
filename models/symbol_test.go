package models

import "testing"

func TestSymbolClone(t *testing.T) {
	sig := "func Foo() error"
	doc := "Foo does things."
	parentID := int64(42)

	orig := Symbol{
		ID:        1,
		Signature: &sig,
		Doc:       &doc,
		ParentID:  &parentID,
		Vector:    []float32{0.1, 0.2, 0.3},
	}
	clone := orig.Clone()

	// Modify clone fields
	*clone.Signature = "CHANGED"
	*clone.Doc = "CHANGED"
	*clone.ParentID = 999
	clone.Vector[0] = 9.9

	if *orig.Signature != "func Foo() error" {
		t.Error("Clone did not isolate Signature")
	}
	if *orig.Doc != "Foo does things." {
		t.Error("Clone did not isolate Doc")
	}
	if *orig.ParentID != 42 {
		t.Error("Clone did not isolate ParentID")
	}
	if orig.Vector[0] != 0.1 {
		t.Error("Clone did not isolate Vector")
	}
}

func TestSymbolClone_NilPointers(t *testing.T) {
	orig := Symbol{ID: 1}
	clone := orig.Clone()

	if clone.Signature != nil || clone.Doc != nil || clone.ParentID != nil || clone.Vector != nil {
		t.Error("Clone should preserve nil pointers")
	}
}
