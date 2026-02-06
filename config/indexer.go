package config

import "github.com/zoobzio/vicky/models"

// Indexer holds gRPC addresses for language-specific SCIP indexer sidecars.
type Indexer struct {
	GoAddr string `env:"VICKY_INDEXER_GO_ADDR"`
	TsAddr string `env:"VICKY_INDEXER_TS_ADDR"`
}

// Validate checks Indexer configuration for required values.
func (c Indexer) Validate() error {
	return nil
}

// Addresses returns a map of language to gRPC address for configured indexers.
func (c Indexer) Addresses() map[models.Language]string {
	addrs := make(map[models.Language]string)
	if c.GoAddr != "" {
		addrs[models.LanguageGo] = c.GoAddr
	}
	if c.TsAddr != "" {
		addrs[models.LanguageTypeScript] = c.TsAddr
	}
	return addrs
}
