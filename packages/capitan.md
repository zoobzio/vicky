# Capitan

Type-safe event coordination for Go.

[GitHub](https://github.com/zoobzio/capitan)

## Vision

A unified event stream across the framework. Every package we build contributes to the same coordinated system. Type-safe, async-first, production-ready. One stream to observe everything.

## Design Decisions

**Zero dependencies**
Core infrastructure. Standard library only.

**Type-safe fields**
Generic `Key[T]` and `Field` preserve types at compile time. No type assertions, no runtime surprises.

**Per-signal workers**
Each signal gets its own goroutine. Isolation. Slow listeners on one signal don't affect others.

**Buffered channels**
Backpressure control with configurable policies - block or drop. Default buffer size 16.

**Event pooling**
`sync.Pool` reduces allocations on hot paths. Events are recycled.

**Panic recovery**
Per-listener recovery. One bad listener doesn't crash the system.

**Sync mode**
Deterministic testing without async complexity. Same code, predictable execution.

**Three operations**
Emit, Hook, Observe. Covers 90% of use cases. Simple mental model.

## Internal Architecture

```
Capitan Instance
├── Registry: Signal → []*Listener mapping
├── Workers: Per-signal goroutines with buffered channels
├── Observers: Global event watchers (logging, metrics)
├── Config: Per-signal settings (patterns + exact matches)
└── Stats: Runtime metrics (queue depths, emit counts)
```

Workers created lazily on first emission. Single RWMutex protects registry and config. Read-lock hot path for existing workers.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Core types: Signal, Key, Field, Severity |
| `fields.go` | Key implementations (StringKey, IntKey, Float64Key, etc.) |
| `event.go` | Event type with pooling and lifecycle |
| `config.go` | Configuration system with glob pattern support |
| `worker.go` | Worker goroutine loop and event processing |
| `listener.go` | Listener registration and lifecycle |
| `observer.go` | Cross-cutting observers |
| `service.go` | Capitan service and public API |

## Current State / Direction

Stable. Core event system complete. Configuration, rate limiting, and statistics working.

Future considerations:
- Event persistence/replay if patterns warrant
- Distributed coordination if needed

## Framework Context

**Dependencies**: None.

**Role**: Event backbone. All packages in the framework emit through capitan. Unified coordination and observability across everything.
