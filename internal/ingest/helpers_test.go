//go:build testing

package ingest

import (
	"math"
	"testing"

	"github.com/zoobzio/vicky/models"
)

func TestIdToKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   int64
		want string
	}{
		{"zero", 0, "0"},
		{"positive", 1, "1"},
		{"negative", -1, "-1"},
		{"large", 12345678, "12345678"},
		{"max int64", math.MaxInt64, "9223372036854775807"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := idToKey(tt.id)
			if got != tt.want {
				t.Errorf("idToKey(%d) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestMatchesAnyPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{"vendor dir", "vendor/pkg/foo.go", []string{"vendor/**"}, true},
		{"no match vendor", "src/main.go", []string{"vendor/**"}, false},
		{"minified js", "src/foo.min.js", []string{"*.min.js"}, true},
		{"non-minified js", "src/app.js", []string{"*.min.js"}, false},
		{"glob md root", "README.md", []string{"**/*.md"}, true},
		{"glob md nested", "docs/guide.md", []string{"**/*.md"}, true},
		{"empty patterns", "main.go", []string{}, false},
		{"github dir", ".github/workflows/ci.yml", []string{".github/**"}, true},
		{"node_modules", "node_modules/pkg/index.js", []string{"node_modules/**"}, true},
		{"multiple patterns first", "vendor/x.go", []string{"vendor/**", "*.test.js"}, true},
		{"multiple patterns second", "foo.test.js", []string{"vendor/**", "*.test.js"}, true},
		{"multiple patterns none", "src/app.go", []string{"vendor/**", "*.test.js"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := matchesAnyPattern(tt.path, tt.patterns)
			if got != tt.want {
				t.Errorf("matchesAnyPattern(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}

func TestContainsExt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		exts []string
		ext  string
		want bool
	}{
		{"found", []string{".go", ".mod"}, ".go", true},
		{"not found", []string{".go", ".mod"}, ".rs", false},
		{"empty slice", []string{}, ".go", false},
		{"exact match", []string{".md"}, ".md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := containsExt(tt.exts, tt.ext)
			if got != tt.want {
				t.Errorf("containsExt(%v, %q) = %v, want %v", tt.exts, tt.ext, got, tt.want)
			}
		})
	}
}

func TestContentTypeForPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want models.ContentType
	}{
		{"go file", "main.go", models.ContentTypeCode},
		{"markdown", "README.md", models.ContentTypeDocs},
		{"mdx", "docs/guide.mdx", models.ContentTypeDocs},
		{"markdown long ext", "notes.markdown", models.ContentTypeDocs},
		{"typescript", "app.ts", models.ContentTypeCode},
		{"makefile", "Makefile", models.ContentTypeCode},
		{"nested code", "src/pkg/handler.go", models.ContentTypeCode},
		{"nested docs", "docs/api/overview.md", models.ContentTypeDocs},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := contentTypeForPath(tt.path)
			if got != tt.want {
				t.Errorf("contentTypeForPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
