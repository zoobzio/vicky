// Package capacitors provides hot-reloadable runtime configuration via flux.
//
// Capacitors are backed by database storage with LISTEN/NOTIFY for real-time updates.
// Use config/ for static bootstrap configuration that requires restart to change.
//
// Usage:
//
//	// At startup, after database is connected
//	err := capacitors.Init(ctx, db, dsn, ap)
//
//	// Access current config anywhere
//	cfg := capacitors.EmbeddingCurrent()
package capacitors

import (
	"context"
	"database/sql"

	"github.com/zoobzio/aperture"
)

// Init initializes all capacitors with database-backed watchers.
// Should be called after database connection is established.
func Init(ctx context.Context, db *sql.DB, dsn string, ap *aperture.Aperture) error {
	// Pipeline stage capacitors
	fetchWatcher := NewDBWatcherWithDSN(db, dsn, DomainFetch)
	if err := InitFetch(ctx, fetchWatcher); err != nil {
		return err
	}

	parseWatcher := NewDBWatcherWithDSN(db, dsn, DomainParse)
	if err := InitParse(ctx, parseWatcher); err != nil {
		return err
	}

	chunkWatcher := NewDBWatcherWithDSN(db, dsn, DomainChunk)
	if err := InitChunk(ctx, chunkWatcher); err != nil {
		return err
	}

	embeddingWatcher := NewDBWatcherWithDSN(db, dsn, DomainEmbedding)
	if err := InitEmbedding(ctx, embeddingWatcher); err != nil {
		return err
	}

	// System capacitors
	eventsWatcher := NewDBWatcherWithDSN(db, dsn, DomainEvents)
	if err := InitEvents(ctx, eventsWatcher); err != nil {
		return err
	}

	observabilityWatcher := NewDBWatcherWithDSN(db, dsn, DomainObservability)
	if err := InitObservability(ctx, observabilityWatcher, ap); err != nil {
		return err
	}

	return nil
}

// Domain constants for config table.
const (
	// Pipeline stages
	DomainFetch     = "fetch"
	DomainParse     = "parse"
	DomainChunk     = "chunk"
	DomainEmbedding = "embedding"

	// System
	DomainEvents        = "events"
	DomainObservability = "observability"
)
