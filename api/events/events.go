// Package events provides bidirectional typed events for vicky's event stream.
//
// Each event domain exposes a namespace with Event[T] instances that support
// both emitting and listening. Signals and event structs are also exported
// for direct access when needed.
//
// # Emitting Events
//
//	events.Ingest.Completed.Emit(ctx, events.IngestCompletedEvent{
//	    RepositoryID: repoID,
//	    VersionTag:   "v1.2.0",
//	    ChunkCount:   1542,
//	    Duration:     elapsed,
//	})
//
// # Listening to Events
//
//	listener := events.Ingest.Completed.Listen(func(ctx context.Context, e events.IngestCompletedEvent) {
//	    log.Printf("Ingestion completed: %d chunks in %v", e.ChunkCount, e.Duration)
//	})
//	defer listener.Close()
//
// # Organization
//
// Events are organized by domain:
//   - Ingest: Pipeline lifecycle (fetch, parse, chunk, embed, store)
//   - Search: Vector similarity and symbol search
//   - Repository: Repository and version lifecycle
//   - User: Authentication and profile events
//   - Embedding: Embedding provider and request events
//   - Config: Hot-reload configuration events
//
// # Direct Signal Access
//
// Signals are exported for direct use with capitan when needed:
//
//	capitan.Hook(events.IngestCompletedSignal, func(ctx context.Context, e *capitan.Event) {
//	    // low-level access
//	})
package events
