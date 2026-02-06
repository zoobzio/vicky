package models

import (
	"encoding/json"
	"testing"
)

func TestAllExcludePatterns(t *testing.T) {
	c := &IngestionConfig{
		ExcludePatterns: []string{"custom/**", "*.tmp"},
	}
	patterns := c.AllExcludePatterns()

	want := len(DefaultExcludePatterns) + 2
	if len(patterns) != want {
		t.Fatalf("len(patterns) = %d, want %d", len(patterns), want)
	}
	// Defaults come first
	if patterns[0] != DefaultExcludePatterns[0] {
		t.Errorf("first pattern = %q, want %q", patterns[0], DefaultExcludePatterns[0])
	}
	// Custom patterns appended
	if patterns[len(patterns)-2] != "custom/**" {
		t.Errorf("second-to-last = %q, want %q", patterns[len(patterns)-2], "custom/**")
	}
}

func TestAllExcludePatterns_NilCustom(t *testing.T) {
	c := &IngestionConfig{}
	patterns := c.AllExcludePatterns()

	if len(patterns) != len(DefaultExcludePatterns) {
		t.Errorf("len(patterns) = %d, want %d", len(patterns), len(DefaultExcludePatterns))
	}
}

func TestIngestionConfigClone(t *testing.T) {
	orig := &IngestionConfig{
		ID:              1,
		ExcludePatterns: []string{"vendor/**"},
		LanguageConfig:  json.RawMessage(`{"module_path":"example.com/foo"}`),
	}
	clone := orig.Clone()

	// Modify clone slices
	clone.ExcludePatterns[0] = "CHANGED"
	clone.LanguageConfig[0] = 'X'

	// Original should be unaffected
	if orig.ExcludePatterns[0] != "vendor/**" {
		t.Error("Clone did not isolate ExcludePatterns")
	}
	if orig.LanguageConfig[0] != '{' {
		t.Error("Clone did not isolate LanguageConfig")
	}
}

func TestIngestionConfigClone_Nil(t *testing.T) {
	var c *IngestionConfig
	if c.Clone() != nil {
		t.Error("Clone of nil should return nil")
	}
}

func TestGetGoConfig(t *testing.T) {
	c := &IngestionConfig{
		Language:       LanguageGo,
		LanguageConfig: json.RawMessage(`{"module_path":"example.com/test","build_tags":["integration"]}`),
	}
	cfg, err := c.GetGoConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ModulePath != "example.com/test" {
		t.Errorf("ModulePath = %q, want %q", cfg.ModulePath, "example.com/test")
	}
	if len(cfg.BuildTags) != 1 || cfg.BuildTags[0] != "integration" {
		t.Errorf("BuildTags = %v, want [integration]", cfg.BuildTags)
	}
}

func TestGetGoConfig_WrongLanguage(t *testing.T) {
	c := &IngestionConfig{
		Language:       LanguageTypeScript,
		LanguageConfig: json.RawMessage(`{"module_path":"example.com/test"}`),
	}
	cfg, err := c.GetGoConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ModulePath != "" {
		t.Errorf("expected empty GoConfig for non-Go language, got ModulePath=%q", cfg.ModulePath)
	}
}

func TestGetTypeScriptConfig(t *testing.T) {
	c := &IngestionConfig{
		Language:       LanguageTypeScript,
		LanguageConfig: json.RawMessage(`{"tsconfig_path":"tsconfig.json","include_tests":true}`),
	}
	cfg, err := c.GetTypeScriptConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TsConfigPath != "tsconfig.json" {
		t.Errorf("TsConfigPath = %q, want %q", cfg.TsConfigPath, "tsconfig.json")
	}
	if !cfg.IncludeTests {
		t.Error("IncludeTests = false, want true")
	}
}

func TestGetTypeScriptConfig_NilConfig(t *testing.T) {
	c := &IngestionConfig{
		Language:       LanguageTypeScript,
		LanguageConfig: nil,
	}
	cfg, err := c.GetTypeScriptConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TsConfigPath != "" {
		t.Errorf("expected empty TypeScriptConfig for nil config, got TsConfigPath=%q", cfg.TsConfigPath)
	}
}
