package wire

import (
	"strings"
	"testing"
)

func TestIngestRequestValidate_Valid(t *testing.T) {
	req := &IngestRequest{CommitSHA: strings.Repeat("a", 40)}
	if err := req.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestIngestRequestValidate_Empty(t *testing.T) {
	req := &IngestRequest{CommitSHA: ""}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty commit SHA, got nil")
	}
	if !strings.Contains(err.Error(), "commit_sha") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "commit_sha")
	}
}

func TestIngestRequestValidate_WrongLength(t *testing.T) {
	req := &IngestRequest{CommitSHA: "tooshort"}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for wrong length, got nil")
	}
	if !strings.Contains(err.Error(), "commit_sha") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "commit_sha")
	}
}
