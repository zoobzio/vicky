package capacitors

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// DBWatcher watches a config domain in the database using LISTEN/NOTIFY.
// Implements flux.Watcher interface.
type DBWatcher struct {
	db      *sql.DB
	domain  string
	channel string
}

// NewDBWatcher creates a watcher for a specific config domain.
// The domain corresponds to a row in the configs table.
func NewDBWatcher(db *sql.DB, domain string) *DBWatcher {
	return &DBWatcher{
		db:      db,
		domain:  domain,
		channel: "config_" + domain,
	}
}

// Watch starts watching for config changes.
// Emits the current value immediately, then emits on each change.
// Closes the channel when context is cancelled.
func (w *DBWatcher) Watch(ctx context.Context) (<-chan []byte, error) {
	out := make(chan []byte)

	// Set up LISTEN
	listener := pq.NewListener(
		dsn(w.db),
		10*time.Second,
		time.Minute,
		func(_ pq.ListenerEventType, _ error) {
			// Connection event handler - could log here
		},
	)

	if err := listener.Listen(w.channel); err != nil {
		return nil, err
	}

	go func() {
		defer close(out)
		defer func() { _ = listener.Close() }()

		// Emit initial value
		data, err := w.fetch(ctx)
		if err == nil && data != nil {
			select {
			case out <- data:
			case <-ctx.Done():
				return
			}
		}

		// Watch for changes
		for {
			select {
			case <-ctx.Done():
				return
			case <-listener.Notify:
				// Config changed, fetch new value
				data, err := w.fetch(ctx)
				if err == nil && data != nil {
					select {
					case out <- data:
					case <-ctx.Done():
						return
					}
				}
			case <-time.After(90 * time.Second):
				// Ping to keep connection alive
				_ = listener.Ping()
			}
		}
	}()

	return out, nil
}

// fetch retrieves the current config value from the database.
func (w *DBWatcher) fetch(ctx context.Context) ([]byte, error) {
	var data json.RawMessage
	err := w.db.QueryRowContext(ctx,
		"SELECT data FROM configs WHERE domain = $1",
		w.domain,
	).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

// dsn extracts the connection string from a *sql.DB.
// This is a workaround since pq.NewListener needs the DSN.
func dsn(_ *sql.DB) string {
	// In practice, you'd pass the DSN directly or store it.
	// For now, this is a placeholder that needs proper implementation.
	// The caller should provide the DSN separately.
	return ""
}

// DBWatcherWithDSN creates a watcher with an explicit DSN for the listener.
type DBWatcherWithDSN struct {
	db      *sql.DB
	dsn     string
	domain  string
	channel string
}

// NewDBWatcherWithDSN creates a watcher with explicit DSN.
func NewDBWatcherWithDSN(db *sql.DB, dsn, domain string) *DBWatcherWithDSN {
	return &DBWatcherWithDSN{
		db:      db,
		dsn:     dsn,
		domain:  domain,
		channel: "config_" + domain,
	}
}

// Watch starts watching for config changes.
func (w *DBWatcherWithDSN) Watch(ctx context.Context) (<-chan []byte, error) {
	out := make(chan []byte)

	listener := pq.NewListener(
		w.dsn,
		10*time.Second,
		time.Minute,
		nil,
	)

	if err := listener.Listen(w.channel); err != nil {
		return nil, err
	}

	go func() {
		defer close(out)
		defer func() { _ = listener.Close() }()

		// Emit initial value
		data, err := w.fetch(ctx)
		if err == nil && data != nil {
			select {
			case out <- data:
			case <-ctx.Done():
				return
			}
		}

		// Watch for changes
		for {
			select {
			case <-ctx.Done():
				return
			case <-listener.Notify:
				data, err := w.fetch(ctx)
				if err == nil && data != nil {
					select {
					case out <- data:
					case <-ctx.Done():
						return
					}
				}
			case <-time.After(90 * time.Second):
				_ = listener.Ping()
			}
		}
	}()

	return out, nil
}

func (w *DBWatcherWithDSN) fetch(ctx context.Context) ([]byte, error) {
	var data json.RawMessage
	err := w.db.QueryRowContext(ctx,
		"SELECT data FROM configs WHERE domain = $1",
		w.domain,
	).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
