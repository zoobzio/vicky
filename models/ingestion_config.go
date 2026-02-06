package models

import (
	"encoding/json"
	"time"
)

// Language represents supported languages for SCIP indexing.
type Language string

// Supported languages.
const (
	LanguageGo         Language = "go"
	LanguageTypeScript Language = "typescript"
)

// DefaultExcludePatterns are always applied during ingestion.
var DefaultExcludePatterns = []string{
	".git/**",
	".github/**",
	".vscode/**",
	".idea/**",
	"vendor/**",
	"node_modules/**",
	"dist/**",
	"build/**",
	"*.min.js",
	"*.min.css",
	"*.map",
	"*.lock",
	"package-lock.json",
	"yarn.lock",
	"go.sum",
}

// DefaultMaxFileSize is the default maximum file size to ingest (1MB).
const DefaultMaxFileSize = 1024 * 1024

// IngestionConfig defines per-repo ingestion settings.
type IngestionConfig struct {
	ID              int64           `json:"id" db:"id" constraints:"primarykey" description:"Config ID"`
	RepositoryID    int64           `json:"repository_id" db:"repository_id" constraints:"notnull,unique" references:"repositories(id)" description:"Parent repository"`
	UserID          int64           `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user"`
	Language        Language        `json:"language" db:"language" constraints:"notnull" description:"Primary language for SCIP indexing" example:"go"`
	IncludeDocs     bool            `json:"include_docs" db:"include_docs" constraints:"notnull" default:"true" description:"Include markdown documentation"`
	ExcludePatterns []string        `json:"exclude_patterns" db:"exclude_patterns" description:"Additional glob patterns to exclude"`
	MaxFileSize     int64           `json:"max_file_size" db:"max_file_size" constraints:"notnull" default:"1048576" description:"Maximum file size in bytes"`
	LanguageConfig  json.RawMessage `json:"language_config,omitempty" db:"language_config" description:"Language-specific configuration"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at" default:"now()" description:"Creation time"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at" default:"now()" description:"Last update time"`
}

// AllExcludePatterns returns the combined default and custom exclude patterns.
func (c *IngestionConfig) AllExcludePatterns() []string {
	patterns := make([]string, 0, len(DefaultExcludePatterns)+len(c.ExcludePatterns))
	patterns = append(patterns, DefaultExcludePatterns...)
	patterns = append(patterns, c.ExcludePatterns...)
	return patterns
}

// Clone returns a deep copy of the IngestionConfig.
func (c *IngestionConfig) Clone() *IngestionConfig {
	if c == nil {
		return nil
	}
	clone := *c
	if c.ExcludePatterns != nil {
		clone.ExcludePatterns = make([]string, len(c.ExcludePatterns))
		copy(clone.ExcludePatterns, c.ExcludePatterns)
	}
	if c.LanguageConfig != nil {
		clone.LanguageConfig = make(json.RawMessage, len(c.LanguageConfig))
		copy(clone.LanguageConfig, c.LanguageConfig)
	}
	return &clone
}

// GoConfig holds Go-specific ingestion settings.
type GoConfig struct {
	ModulePath string   `json:"module_path,omitempty"` // Override module path detection
	BuildTags  []string `json:"build_tags,omitempty"`  // Build tags to include
}

// TypeScriptConfig holds TypeScript-specific ingestion settings.
type TypeScriptConfig struct {
	TsConfigPath string `json:"tsconfig_path,omitempty"` // Path to tsconfig.json
	IncludeTests bool   `json:"include_tests,omitempty"` // Include .spec.ts, .test.ts
}

// GetGoConfig parses and returns the Go-specific config.
func (c *IngestionConfig) GetGoConfig() (*GoConfig, error) {
	if c.Language != LanguageGo || c.LanguageConfig == nil {
		return &GoConfig{}, nil
	}
	var cfg GoConfig
	if err := json.Unmarshal(c.LanguageConfig, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetTypeScriptConfig parses and returns the TypeScript-specific config.
func (c *IngestionConfig) GetTypeScriptConfig() (*TypeScriptConfig, error) {
	if c.Language != LanguageTypeScript || c.LanguageConfig == nil {
		return &TypeScriptConfig{}, nil
	}
	var cfg TypeScriptConfig
	if err := json.Unmarshal(c.LanguageConfig, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
