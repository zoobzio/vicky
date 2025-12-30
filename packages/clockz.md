# Clockz

Deterministic time control for Go.

[GitHub](https://github.com/zoobzio/clockz)

## Vision

Control time in tests. Time-dependent code is notoriously difficult to test. Real sleeps make tests slow. Non-deterministic timing creates flakiness. clockz eliminates both. One interface, two implementations. Production uses real time. Tests control time explicitly.

## Design Decisions

**Interface-based abstraction**
Single `Clock` interface mirrors the `time` package. Dependency injection enables zero-code-change test swaps. Accept a clock, test with fake time.

**Separated Timer/Ticker interfaces**
Return interfaces, not concrete types. Channel via `C()` method, not field. Consistent patterns, mockable internals.

**Synchronous callback execution**
`AfterFunc` callbacks execute synchronously during `Advance()`, sorted by target time. Deterministic ordering. No goroutine scheduling uncertainty. Tests observe results immediately.

**Sorted waiter processing**
Multiple timers firing in same `Advance()` execute in predictable sequence. Prevents subtle ordering bugs. Determinism over performance.

**BlockUntilReady() primitive**
Separates clock advancement from operation delivery. Handles non-blocking channel sends. Ensures tests don't race against timer buffering.

**No backward time movement**
Time progresses monotonically. Panics if time moves backward during advance. Forces linear thinking about time. Matches `time.Time` semantics.

**Context integration**
`WithTimeout()` and `WithDeadline()` bridge clock abstraction into context system. Three cancellation paths: manual, deadline, parent. Correct error codes.

## The Vocabulary

### Clock Interface

| Method | Purpose |
|--------|---------|
| `Now()` | Current time |
| `After(d)` | Channel that receives after duration |
| `Sleep(d)` | Block for duration |
| `NewTimer(d)` | Resettable timer |
| `NewTicker(d)` | Repeating ticker |
| `AfterFunc(d, f)` | Callback after duration |
| `WithTimeout(ctx, d)` | Context with timeout |
| `WithDeadline(ctx, t)` | Context with deadline |

### FakeClock Additions

| Method | Purpose |
|--------|---------|
| `Advance(d)` | Move time forward, fire waiters |
| `SetTime(t)` | Jump to specific time |
| `BlockUntilReady()` | Wait for pending operations |

## Internal Architecture

```
Clock Interface
├── RealClock (singleton, zero overhead)
│   └── Delegates to time package
│
└── FakeClock
    ├── time (current fake time)
    ├── waiters (sorted by target time)
    │   ├── Timer waiters
    │   ├── Ticker waiters
    │   └── AfterFunc callbacks
    ├── contextWaiters (deadline contexts)
    └── pendingSends (queued channel deliveries)
```

`Advance()` sorts waiters, fires those past target, queues sends. `BlockUntilReady()` drains pending sends. Deterministic. Thread-safe.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Public interfaces: Clock, Timer, Ticker |
| `real.go` | RealClock singleton, delegates to stdlib |
| `fake.go` | FakeClock implementation, waiter management, context integration |

## Current State / Direction

Stable. Core clock abstraction complete with context integration.

Future considerations:
- Additional helper patterns as usage emerges
- Performance optimisations for high timer counts

## Framework Context

**Dependencies**: None. Standard library only.

**Role**: Time abstraction foundation. Enables deterministic testing of time-dependent code. Used by streamz, pipz, and any package with timing logic.
