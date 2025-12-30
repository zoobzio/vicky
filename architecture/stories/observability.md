# Observability Story

*"How do I observe my application?"*

## The Flow

```
capitan → herald → aperture
```

## The Packages

### capitan - Unified Event Stream

Every package emits signals. Type-safe fields. Async-first.

```go
capitan.Emit(ctx, MySignal,
    MyDurationKey.Field(elapsed),
    MyStatusKey.Field("success"),
)
```

**What it provides:**
- Per-signal worker goroutines (isolation)
- Buffered channels with backpressure
- Hook for local listeners
- Observe for global watchers

### herald - Distributed Events

Extends capitan across process boundaries. Same contracts, distributed.

```go
// Publish local events to external broker
pub := herald.NewPublisher[MyEvent](provider, opts...)
capitan.Hook(ctx, MySignal, func(e *capitan.Event) {
    pub.Publish(ctx, toMyEvent(e))
})

// Subscribe external events to local capitan
sub := herald.NewSubscriber[MyEvent](provider, opts...)
sub.Subscribe(ctx, func(msg herald.Message[MyEvent]) {
    capitan.Emit(ctx, ExternalSignal, ...)
})
```

**What it provides:**
- 11 broker providers
- Bidirectional (publish/subscribe)
- pipz-based reliability

### aperture - OpenTelemetry Bridge

Bridges unified stream to OpenTelemetry. Configure what gets logged, traced, metricked.

```yaml
logs:
  - signal: "my.operation"
    level: info
    attributes: [duration, status]

metrics:
  - signal: "my.operation"
    type: histogram
    field: duration

traces:
  - start_signal: "my.operation.start"
    end_signal: "my.operation.end"
    correlation_key: request_id
```

**What it provides:**
- Name-based matching (hot-reload without recompilation)
- Three independent handlers (logs, metrics, traces)
- User-provided OTEL providers

## The Key Insight

**Observability is built-in, not bolted-on.**

Every package already emits through capitan. You don't instrument - you configure which signals to observe. Emit once, observe everywhere.

```
┌─────────────────────────────────────────────────────────────┐
│                     Your Application                         │
├─────────────────────────────────────────────────────────────┤
│  soy         zyn         rocco        tendo        ...      │
│   │           │            │            │                   │
│   └───────────┴────────────┴────────────┘                   │
│                      │                                       │
│                      ▼                                       │
│               ┌──────────────┐                               │
│               │   capitan    │ ← All signals flow here      │
│               └──────┬───────┘                               │
│                      │                                       │
│         ┌────────────┼────────────┐                         │
│         │            │            │                         │
│         ▼            ▼            ▼                         │
│   ┌──────────┐ ┌──────────┐ ┌──────────┐                   │
│   │  herald  │ │ aperture │ │  custom  │                   │
│   │ (dist)   │ │  (otel)  │ │ (hooks)  │                   │
│   └──────────┘ └──────────┘ └──────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

## Related Stories

- [Distributed](distributed.md) - Uses herald for event coordination
- [Composition](composition.md) - pipz pipelines emit to capitan
