# Pipz

Composable operations for Go.

[GitHub](https://github.com/zoobzio/pipz)

## Vision

The `Chainable[T]` interface is the product. Everything else - processors, connectors, patterns - exists to demonstrate and leverage that interface.

Any type implementing `Chainable[T]` composes with any other `Chainable[T]`. Framework primitives have no privileged position. User code satisfying the interface is indistinguishable from framework code. The API surface is unbounded.

## Usage Capacities

Pipz serves three distinct roles:

**Direct pipelines**
Users compose pipz primitives directly over their own types. Build data pipelines, request processing, ETL workflows.

**Package extensibility**
Packages like zyn and herald manage internal pipz pipelines but expose insertion points. The package defines T, users inject `Chainable[T]` implementations. Users extend package behaviour without forking.

**Downstream foundations**
Packages like cogito are built *on top of* pipz where T is a known quantity (`*Thought`). All pipz connectors work via type aliasing. Any user type implementing `Chainable[*Thought]` becomes a first-class primitive in cogito's ecosystem. Hundreds of sub-packages could contribute to cogito's API surface simply by satisfying the core interface.

## Design Decisions

**Uniform interface**
Everything implements `Chainable[T]`. Processors and connectors compose freely at any depth.

**Processors are leaves, connectors are branches**
Processors wrap functions - they're leaf nodes with no children. Connectors orchestrate other Chainables - they're branch nodes containing children. Processors are value types (stateless, safe to share). Connectors are pointer types (stateful, manage children).

**Generics throughout**
Compile-time type safety. No `interface{}`, no reflection.

**Context everywhere**
Timeout and cancellation supported at every level.

**Fail-fast**
First error stops pipeline, bubbles up with full path.

**Path tracking**
Errors carry the exact failure location through nested pipelines. UUIDs enable correlation.

**Panic recovery**
Automatic recovery with security-focused message sanitisation.

## Processors

Leaf nodes. Wrap functions with identity metadata. Value types.

| Processor | Purpose |
|-----------|---------|
| `Transform` | Pure transformation (cannot fail) |
| `Apply` | Fallible transformation |
| `Effect` | Side effects (data passes through) |
| `Mutate` | Conditional transformation |
| `Enrich` | Optional enhancement (errors ignored) |

## Connectors

Branch nodes. Orchestrate child Chainables. Pointer types with state.

| Category | Connectors |
|----------|------------|
| Sequential | `Sequence` |
| Parallel | `Concurrent`, `Race`, `Contest`, `WorkerPool` |
| Error handling | `Retry`, `Backoff`, `Fallback`, `Handle` |
| Flow control | `Switch`, `Filter`, `Timeout` |
| Resource protection | `RateLimiter`, `CircuitBreaker` |
| Background | `Scaffold` |

## Internal Architecture

```
Chainable[T] interface
├── Process(ctx, T) → (T, error)
├── Identity() → name, description, UUID
├── Schema() → introspectable structure
└── Close() → cleanup
```

Processors implement Chainable by wrapping a function. Connectors implement Chainable by delegating to child Chainables. Both satisfy the same interface. Nest freely.

## Code Organisation

| Category | Files |
|----------|-------|
| Core | `api.go`, `error.go`, `signals.go`, `schema.go` |
| Processors | `transform.go`, `apply.go`, `effect.go`, `mutate.go`, `enrich.go` |
| Sequential | `sequence.go` |
| Parallel | `concurrent.go`, `race.go`, `contest.go`, `workerpool.go` |
| Error handling | `retry.go`, `backoff.go`, `fallback.go`, `handle.go` |
| Flow control | `switch.go`, `filter.go`, `timeout.go` |
| Resource protection | `circuitbreaker.go`, `ratelimiter.go` |
| Background | `scaffold.go` |

## Current State / Direction

Stable. All major patterns implemented.

Future considerations:
- Additional connectors as patterns emerge
- Composition syntax sugar

## Framework Context

**Dependencies**: clockz (time abstraction), capitan (signals), uuid.

**Role**: The composability foundation. Defines the interface that enables unbounded extensibility across the framework.
