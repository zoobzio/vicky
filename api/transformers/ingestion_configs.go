package transformers

import (
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/api/wire"
)

// IngestionConfigToResponse transforms an IngestionConfig model to an API response.
func IngestionConfigToResponse(c *models.IngestionConfig) wire.IngestionConfigResponse {
	return wire.IngestionConfigResponse{
		ID:              c.ID,
		Language:        c.Language,
		IncludeDocs:     c.IncludeDocs,
		ExcludePatterns: c.ExcludePatterns,
		MaxFileSize:     c.MaxFileSize,
		LanguageConfig:  c.LanguageConfig,
	}
}

// ApplyIngestionConfigRequest applies an IngestionConfigRequest to an IngestionConfig model.
func ApplyIngestionConfigRequest(req wire.IngestionConfigRequest, c *models.IngestionConfig) {
	c.Language = req.Language
	c.IncludeDocs = req.IncludeDocs
	c.ExcludePatterns = req.ExcludePatterns
	c.LanguageConfig = req.LanguageConfig

	// Apply max file size with default
	if req.MaxFileSize != nil {
		c.MaxFileSize = *req.MaxFileSize
	} else {
		c.MaxFileSize = models.DefaultMaxFileSize
	}
}
