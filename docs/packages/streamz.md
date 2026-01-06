# Streamz

Type-safe stream processing for Go channels.

[GitHub](https://github.com/zoobzio/streamz)

## Vision

Composable channel primitives. Filter, transform, batch, window, route - each processor takes a channel and returns a channel. Errors flow through with context. Time-based logic tests deterministically.

The package is stable and complete. Its role in the broader framework remains to be determined - no other packages currently depend on it.

## Design Decisions

**Result[T] pattern**
Single channel carries both success and error. No dual-channel deadlocks. Cleaner composition. Backpressure on one path.

**Metadata attachment**
Results carry optional context - window boundaries, partition indices, processor names, timestamps. Flows through automatically. Zero overhead when unused.

**Clock abstraction**
Inject clock dependency. Real clock in production, fake clock in tests. Deterministic verification of time-based logic. No flaky tests.

**Single-goroutine per processor**
One goroutine with select loop, not worker pools. No goroutine leaks. No shared state races. Predictable memory. AsyncMapper is the exception when parallelism needed.

**Panic recovery throughout**
User predicates and mappers wrapped in deferred recovery. Panics become errors. Pipelines survive bad user code.

**Structured errors**
StreamError captures original item, processor name, timestamp, wrapped error. Enables retry logic, DLQ patterns, latency analysis.

**Non-blocking drops**
Dead letter queue drops rather than deadlocks. Atomic counter tracks losses. Asymmetric consumption supported.

## Processors

### Transformation

| Processor | Purpose |
|-----------|---------|
| `Filter` | Predicate-based pass/drop |
| `Mapper` | Synchronous transformation |
| `AsyncMapper` | Concurrent transformation, ordered or unordered |
| `Tap` | Side effects, 100% throughput |
| `Sample` | Probabilistic sampling |

### Routing

| Processor | Purpose |
|-----------|---------|
| `FanIn` | Merge N channels to 1 |
| `FanOut` | Duplicate 1 channel to N |
| `Partition` | Hash or round-robin sharding |
| `Switch` | Predicate-based routing |
| `DeadLetterQueue` | Success/failure split |

### Batching

| Processor | Purpose |
|-----------|---------|
| `Batcher` | Size or time-triggered batches |
| `Buffer` | Buffered channel wrapper |

### Windowing

| Processor | Purpose |
|-----------|---------|
| `TumblingWindow` | Fixed non-overlapping windows |
| `SlidingWindow` | Overlapping windows |
| `SessionWindow` | Activity-gap based sessions |
| `WindowCollector` | Groups by window metadata |

### Flow Control

| Processor | Purpose |
|-----------|---------|
| `Throttle` | Leading-edge rate limiting |
| `Debounce` | Trailing-edge quiet period |

## Internal Architecture

```
Input Channel
     │
     ▼
Result[T] (value or StreamError)
     │
     ▼
Processor (single goroutine, select loop)
├── Transform/filter item
├── Handle errors (pass through or route)
├── Attach metadata (window, partition)
└── Recover from panics
     │
     ▼
Output Channel (Result[Out])
     │
     ▼
Next Processor...
```

Each processor: `Process(ctx, in <-chan Result[T]) <-chan Result[Out]`

Connect any to any. Errors flow through with context.

## Code Organisation

| Category | Files |
|----------|-------|
| Core | `api.go`, `result.go`, `error.go`, `clock.go` |
| Transformation | `filter.go`, `mapper.go`, `async_mapper.go`, `tap.go`, `sample.go` |
| Routing | `fanin.go`, `fanout.go`, `partition.go`, `switch.go`, `dlq.go` |
| Batching | `batcher.go`, `buffer.go` |
| Windowing | `window_tumbling.go`, `window_sliding.go`, `window_session.go` |
| Flow Control | `throttle.go`, `debounce.go` |

## Current State / Direction

Stable. All 24 processors complete with race detection, panic recovery, and edge case handling.

Future considerations:
- Additional windowing strategies as patterns emerge
- Metrics integration
- Backpressure signaling

## Framework Context

**Dependencies**: clockz (time abstraction).

**Role**: Channel-based stream processing primitives. Currently standalone - no framework packages depend on streamz yet. May find its place as real-time processing needs emerge.
