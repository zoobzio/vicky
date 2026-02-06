package models

import "testing"

func TestUserClone(t *testing.T) {
	name := "The Octocat"
	avatar := "https://example.com/avatar.png"

	orig := User{
		ID:        1,
		Name:      &name,
		AvatarURL: &avatar,
	}
	clone := orig.Clone()

	*clone.Name = "CHANGED"
	*clone.AvatarURL = "CHANGED"

	if *orig.Name != "The Octocat" {
		t.Error("Clone did not isolate Name")
	}
	if *orig.AvatarURL != "https://example.com/avatar.png" {
		t.Error("Clone did not isolate AvatarURL")
	}
}

func TestUserClone_NilPointers(t *testing.T) {
	orig := User{ID: 1}
	clone := orig.Clone()

	if clone.Name != nil || clone.AvatarURL != nil {
		t.Error("Clone should preserve nil pointers")
	}
}
