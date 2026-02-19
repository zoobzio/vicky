package wire

import (
	"strings"
	"testing"
)

func TestRegisterRepositoryRequestValidate_Valid(t *testing.T) {
	req := &RegisterRepositoryRequest{
		GitHubID:      123,
		Owner:         "octocat",
		Name:          "hello-world",
		FullName:      "octocat/hello-world",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/octocat/hello-world",
		Config:        IngestionConfigRequest{Language: "go"},
	}
	if err := req.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegisterRepositoryRequestValidate_MissingOwner(t *testing.T) {
	req := &RegisterRepositoryRequest{
		GitHubID:      123,
		Owner:         "",
		Name:          "hello-world",
		FullName:      "octocat/hello-world",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/octocat/hello-world",
		Config:        IngestionConfigRequest{Language: "go"},
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing owner, got nil")
	}
	if !strings.Contains(err.Error(), "owner") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "owner")
	}
}

func TestRegisterRepositoryRequestValidate_MissingName(t *testing.T) {
	req := &RegisterRepositoryRequest{
		GitHubID:      123,
		Owner:         "octocat",
		Name:          "",
		FullName:      "octocat/hello-world",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/octocat/hello-world",
		Config:        IngestionConfigRequest{Language: "go"},
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "name")
	}
}

func TestIngestionConfigRequestValidate_ValidLanguage(t *testing.T) {
	for _, lang := range []string{"go", "typescript"} {
		t.Run(lang, func(t *testing.T) {
			req := &IngestionConfigRequest{Language: "go"}
			if lang == "typescript" {
				req.Language = "typescript"
			}
			if err := req.Validate(); err != nil {
				t.Errorf("unexpected error for language %q: %v", lang, err)
			}
		})
	}
}

func TestIngestionConfigRequestValidate_InvalidLanguage(t *testing.T) {
	req := &IngestionConfigRequest{Language: "python"}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for invalid language, got nil")
	}
	if !strings.Contains(err.Error(), "language") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "language")
	}
}
