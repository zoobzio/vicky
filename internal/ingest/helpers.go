package ingest

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zoobzio/vicky/models"
)

// Documentation extensions.
var docExtensions = []string{".md", ".mdx", ".markdown"}

// idToKey converts an int64 ID to a string key.
func idToKey(id int64) string {
	return strconv.FormatInt(id, 10)
}

// matchesAnyPattern checks if a path matches any of the given glob patterns.
func matchesAnyPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		// Handle ** prefix patterns
		if strings.HasPrefix(pattern, "**") {
			// Match pattern against any suffix
			suffix := strings.TrimPrefix(pattern, "**")
			suffix = strings.TrimPrefix(suffix, "/")
			if matched, _ := filepath.Match(suffix, filepath.Base(path)); matched {
				return true
			}
		}

		// Handle directory patterns (e.g., "vendor/**")
		if strings.HasSuffix(pattern, "/**") {
			prefix := strings.TrimSuffix(pattern, "/**")
			if strings.HasPrefix(path, prefix+"/") || path == prefix {
				return true
			}
		}

		// Standard glob match
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}

		// Also try matching against base name for simple patterns
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

// containsExt checks if an extension is in the list.
func containsExt(exts []string, ext string) bool {
	for _, e := range exts {
		if e == ext {
			return true
		}
	}
	return false
}

// contentTypeForPath determines the content type based on file extension.
func contentTypeForPath(path string) models.ContentType {
	ext := strings.ToLower(filepath.Ext(path))
	if containsExt(docExtensions, ext) {
		return models.ContentTypeDocs
	}
	return models.ContentTypeCode
}
