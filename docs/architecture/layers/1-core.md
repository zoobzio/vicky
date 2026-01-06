# Layer 1: Core Abstractions

Build on primitives to provide patterns used everywhere.

## Rule

Packages in this layer may depend on Layer 0 only.

## Packages

| Package | Role | Framework Deps |
|---------|------|----------------|
| [pipz](../../packages/pipz.md) | Composable operations | clockz, capitan |
| [atom](../../packages/atom.md) | Type-segregated decomposition | sentinel |
| [streamz](../../packages/streamz.md) | Channel stream processing | clockz |

## Package Details

### pipz

Composable operations for Go. A vocabulary around T.

**Key capabilities:**
- `Chainable[T]` universal interface
- Processors: Transform, Apply, Effect, Mutate, Enrich
- Connectors: Sequence, Concurrent, Race, Retry, Backoff, CircuitBreaker, RateLimiter
- Path tracking through nested pipelines
- Panic recovery with sanitised messages

**Framework dependencies:**
- clockz: Time-based operations (timeouts, backoff)
- capitan: Signal emission for observability

**Why Layer 1:** Defines the composability vocabulary. Used by flux, flume, herald, zyn, cogito, tendo.

### atom

Type-segregated atomic value storage. Decompose structs to typed maps.

**Key capabilities:**
- 21 typed tables (strings, ints, floats, bools, times, bytes + pointers + slices)
- `Use[T]()` registry with caching
- `Atomize(*T)` / `Deatomize(*Atom)` bidirectional conversion
- Flatten/Unflatten for JSON interop
- Nested struct support with cycle handling

**Framework dependencies:**
- sentinel: Type metadata extraction

**Why Layer 1:** Type-safe struct decomposition. Foundation for grub's provider-agnostic storage.

### streamz

Type-safe stream processing for Go channels.

**Key capabilities:**
- `Result[T]` pattern (value or error in one channel)
- Transformation: Filter, Mapper, AsyncMapper, Tap, Sample
- Routing: FanIn, FanOut, Partition, Switch, DeadLetterQueue
- Batching: Batcher, Buffer
- Windowing: TumblingWindow, SlidingWindow, SessionWindow
- Flow Control: Throttle, Debounce
- Single-goroutine per processor (predictable, no leaks)

**Framework dependencies:**
- clockz: Deterministic time-based testing

**Why Layer 1:** Channel-based complement to pipz. Real-time data pipelines.

## Architectural Significance

Layer 1 establishes the core patterns:

| Pattern | Package | Used By |
|---------|---------|---------|
| `Chainable[T]` composition | pipz | cogito, zyn, flux, flume, herald, tendo |
| Type-segregated storage | atom | grub |
| `Result[T]` streaming | streamz | (stream pipelines) |

These patterns propagate through all higher layers.
