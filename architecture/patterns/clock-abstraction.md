# Clock Abstraction

*Inject time dependency for testability.*

## The Problem

Time-dependent code is hard to test:
- Real sleeps make tests slow
- Non-deterministic timing creates flakiness
- Testing timeout logic requires waiting

## The Pattern

Abstract time behind an interface. Inject real clock in production, fake clock in tests.

```go
// Interface mirrors time package
type Clock interface {
    Now() time.Time
    After(d time.Duration) <-chan time.Time
    Sleep(d time.Duration)
    NewTimer(d time.Duration) Timer
    NewTicker(d time.Duration) Ticker
    AfterFunc(d time.Duration, f func()) Timer
}

// Production
clock := clockz.Real()

// Test
clock := clockz.NewFake(time.Now())
clock.Advance(5 * time.Minute) // Instant
```

## Where It's Used

| Package | Time-Dependent Features |
|---------|-------------------------|
| clockz | Provides the abstraction |
| streamz | Windowing, batching, throttling |
| pipz | Timeouts, backoff |
| flux | Debouncing |
| tendo | Profiling (timing) |

## Implementation Details

### Real Clock

Zero overhead. Singleton. Delegates to time package.

```go
var realClock = &RealClock{}

func Real() Clock {
    return realClock
}

func (c *RealClock) Now() time.Time {
    return time.Now()
}
```

### Fake Clock

Explicit time control. Deterministic.

```go
fake := clockz.NewFake(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

// Advance fires waiters synchronously
fake.Advance(5 * time.Minute)

// Jump to specific time
fake.SetTime(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC))

// Wait for pending operations
fake.BlockUntilReady()
```

### Sorted Waiter Processing

Multiple timers firing in same Advance execute in order:

```go
fake.NewTimer(1 * time.Second)  // Fires first
fake.NewTimer(2 * time.Second)  // Fires second
fake.NewTimer(3 * time.Second)  // Fires third

fake.Advance(5 * time.Second)
// All three fire in order
```

### Context Integration

```go
ctx, cancel := clock.WithTimeout(ctx, 5*time.Second)
defer cancel()

// In test:
fake.Advance(6 * time.Second)
// Context is now cancelled
```

## Why This Pattern

**Fast tests.** No real waiting. Advance is instant.

**Deterministic.** No race conditions. Same input = same output.

**Full coverage.** Test timeout paths, expiration, scheduling.

**No code changes.** Same production code, different clock injection.

## Usage Pattern

```go
// Accept clock as dependency
type Service struct {
    clock clockz.Clock
}

func NewService(clock clockz.Clock) *Service {
    return &Service{clock: clock}
}

// Use clock instead of time package
func (s *Service) Process(ctx context.Context) {
    start := s.clock.Now()
    // ...
    elapsed := s.clock.Now().Sub(start)
}

// Production
service := NewService(clockz.Real())

// Test
fake := clockz.NewFake(time.Now())
service := NewService(fake)
fake.Advance(5 * time.Second)
```

## Trade-offs

**API surface.** Must mirror time package. Additional interface methods as needed.

**Discipline required.** Must consistently use clock, never time package directly.

**Complexity.** Fake clock implementation is non-trivial (waiter management, ordering).

## Test Example

```go
func TestTimeout(t *testing.T) {
    fake := clockz.NewFake(time.Now())
    service := NewService(fake)

    ctx, cancel := fake.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Start operation
    go service.LongOperation(ctx)

    // Advance past timeout
    fake.Advance(6 * time.Second)

    // Verify timeout was handled
    assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}
```

## Related Patterns

- [Chainable Composition](chainable-composition.md) - Timeouts use clock
- [Signal Emission](signal-emission.md) - Sync mode for deterministic testing
