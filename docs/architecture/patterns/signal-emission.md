# Signal Emission

*capitan signals for observability.*

## The Problem

Observability is critical. But instrumenting every package with logging, metrics, and tracing creates coupling and boilerplate. Different backends need different code.

## The Pattern

Emit structured events to a unified stream. Observers decide what to do with them.

```go
// Package emits
capitan.Emit(ctx, MyOperationSignal,
    DurationKey.Field(elapsed),
    StatusKey.Field("success"),
    RecordCountKey.Field(count),
)

// Observer handles
capitan.Hook(ctx, MyOperationSignal, func(e *capitan.Event) {
    log.Info(e.Fields...)
})

// Or global observer for everything
capitan.Observe(ctx, func(e *capitan.Event) {
    metrics.Record(e)
})
```

## Where It's Used

**Emitters:** All packages except schema languages (dbml, ddml, vdml, openapi).

**Consumers:**
- aperture - OpenTelemetry bridge
- herald - Distributed messaging
- Custom observers - Application-specific

## Implementation Details

### Type-Safe Fields

```go
// Define typed keys
var DurationKey = capitan.Float64Key("duration_ms")
var StatusKey = capitan.StringKey("status")
var CountKey = capitan.IntKey("count")

// Use in emissions
capitan.Emit(ctx, Signal,
    DurationKey.Field(42.5),  // Compile-time type checking
    StatusKey.Field("ok"),
)
```

No `interface{}`. No type assertions at runtime.

### Per-Signal Workers

Each signal gets its own goroutine:
- Isolation between signals
- Slow listener on one signal doesn't affect others
- Buffered channels for backpressure

### Sync Mode

For testing:
```go
cap := capitan.New(capitan.WithSyncMode(true))
// Events processed synchronously
// Deterministic test execution
```

### Panic Recovery

```go
capitan.Hook(ctx, Signal, func(e *capitan.Event) {
    panic("listener bug")
})
// Panic recovered. Other listeners continue.
// System stays up.
```

## Why This Pattern

**Observability without instrumentation.** Packages emit, observers decide. No logging/metrics/tracing code in business logic.

**Unified stream.** One event system. One vocabulary. One place to hook.

**Decoupled backends.** Switch from stdout to OpenTelemetry by changing observers, not emitters.

**Type safety.** Compile-time checking on field types.

## Signal Design

Good signals:
- Named for the operation, not the outcome (`order.placed`, not `order.placed.success`)
- Include relevant context (IDs, durations, counts)
- Use typed field keys
- Emit at consistent points (start/end, or just completion)

## Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Framework Packages                        │
│  soy   zyn   rocco   tendo   herald   flux   cogito   ...   │
│   │     │      │       │        │       │       │           │
│   └─────┴──────┴───────┴────────┴───────┴───────┘           │
│                         │                                    │
│                         ▼                                    │
│                 ┌──────────────┐                             │
│                 │   capitan    │  ← All signals flow here   │
│                 │   (events)   │                             │
│                 └──────┬───────┘                             │
│                        │                                     │
│         ┌──────────────┼──────────────┐                     │
│         │              │              │                     │
│         ▼              ▼              ▼                     │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐               │
│   │ aperture │   │  herald  │   │  custom  │               │
│   │  (otel)  │   │  (dist)  │   │(your app)│               │
│   └──────────┘   └──────────┘   └──────────┘               │
└─────────────────────────────────────────────────────────────┘
```

## Trade-offs

**Event volume.** High-throughput systems generate many events. Backpressure policies needed.

**Schema evolution.** Adding/removing fields affects observers. Versioning strategy needed.

**Debugging complexity.** Async events harder to trace than synchronous calls.

## Related Patterns

- [Chainable Composition](chainable-composition.md) - Operations emit signals
- [Provider Abstraction](provider-abstraction.md) - aperture bridges to OTEL
