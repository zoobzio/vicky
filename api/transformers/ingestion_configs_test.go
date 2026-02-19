package transformers

import (
	"encoding/json"
	"testing"

	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/api/wire"
)

func TestIngestionConfigToResponse(t *testing.T) {
	c := &models.IngestionConfig{
		ID:              5,
		Language:        models.LanguageGo,
		IncludeDocs:     true,
		ExcludePatterns: []string{"vendor/**"},
		MaxFileSize:     2048,
		LanguageConfig:  json.RawMessage(`{"module_path":"example.com"}`),
	}

	resp := IngestionConfigToResponse(c)

	if resp.ID != 5 {
		t.Errorf("ID = %d, want 5", resp.ID)
	}
	if resp.Language != models.LanguageGo {
		t.Errorf("Language = %q, want %q", resp.Language, models.LanguageGo)
	}
	if !resp.IncludeDocs {
		t.Error("IncludeDocs = false, want true")
	}
	if len(resp.ExcludePatterns) != 1 || resp.ExcludePatterns[0] != "vendor/**" {
		t.Errorf("ExcludePatterns = %v, want [vendor/**]", resp.ExcludePatterns)
	}
	if resp.MaxFileSize != 2048 {
		t.Errorf("MaxFileSize = %d, want 2048", resp.MaxFileSize)
	}
}

func TestApplyIngestionConfigRequest_WithMaxFileSize(t *testing.T) {
	maxSize := int64(512)
	req := wire.IngestionConfigRequest{
		Language:    models.LanguageTypeScript,
		IncludeDocs: true,
		MaxFileSize: &maxSize,
	}
	c := &models.IngestionConfig{}

	ApplyIngestionConfigRequest(req, c)

	if c.Language != models.LanguageTypeScript {
		t.Errorf("Language = %q, want %q", c.Language, models.LanguageTypeScript)
	}
	if c.MaxFileSize != 512 {
		t.Errorf("MaxFileSize = %d, want 512", c.MaxFileSize)
	}
}

func TestApplyIngestionConfigRequest_DefaultMaxFileSize(t *testing.T) {
	req := wire.IngestionConfigRequest{
		Language:    models.LanguageGo,
		MaxFileSize: nil,
	}
	c := &models.IngestionConfig{}

	ApplyIngestionConfigRequest(req, c)

	if c.MaxFileSize != models.DefaultMaxFileSize {
		t.Errorf("MaxFileSize = %d, want %d", c.MaxFileSize, models.DefaultMaxFileSize)
	}
}
