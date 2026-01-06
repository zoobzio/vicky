# Chainable Composition

*`pipz.Chainable[T]` as universal interface.*

## The Problem

Operations need to compose. Retry wraps timeout wraps transform. But how do you compose heterogeneous operations into uniform pipelines?

## The Pattern

Define a universal interface. Everything implements it. Compose freely at any depth.

```go
type Chainable[T any] interface {
    Process(ctx context.Context, input T) (T, error)
    Identity() (name, description string, id uuid.UUID)
    Schema() any
    Close() error
}
```

## Where It's Used

| Package | What Implements Chainable |
|---------|---------------------------|
| pipz | All processors and connectors |
| tendo | All tensor operations |
| cogito | All reasoning primitives |
| zyn | Service[T] |
| flux | Processing pipeline |
| flume | Built pipelines |
| herald | Publisher/Subscriber pipelines |
| ago | All coordination primitives |

## Implementation Details

### Processors vs Connectors (pipz)

**Processors** - Immutable values. Safe to share.
```go
transform := pipz.Transform(func(x int) int { return x * 2 })
```

**Connectors** - Stateful containers. Wrap other Chainables.
```go
retry := pipz.NewRetry(transform, pipz.WithMaxRetries(3))
sequence := pipz.NewSequence(validate, transform, persist)
```

Both implement `Chainable[T]`.

### Composition Depth

Nest freely:
```go
pipeline := pipz.NewSequence(
    pipz.NewRetry(
        pipz.NewTimeout(
            pipz.Apply(callExternalAPI),
            5*time.Second,
        ),
        pipz.WithMaxRetries(3),
    ),
    pipz.NewCircuitBreaker(
        pipz.Apply(processResult),
        pipz.WithFailureThreshold(5),
    ),
)
```

### Path Tracking

Errors carry their path through the chain:
```go
// Error includes: "sequence[0].retry[2].timeout.apply: context deadline exceeded"
```

UUIDs enable correlation across distributed systems.

### Schema Introspection

```go
schema := pipeline.Schema()
// Returns tree of operations, parameters, nested schemas
// Useful for documentation, debugging, visualisation
```

## Why This Pattern

**Uniform composition.** Any Chainable composes with any other. No special cases.

**Reliability patterns built in.** Retry, timeout, circuit breaker apply to anything.

**Introspectable.** Identity and Schema enable debugging, logging, visualisation.

**Nested error handling.** Errors bubble up with full context.

## The Four Methods

| Method | Purpose | When Called |
|--------|---------|-------------|
| `Process` | Execute the operation | Every invocation |
| `Identity` | Return name, description, UUID | Debugging, logging, error messages |
| `Schema` | Return introspectable structure | Documentation, visualisation |
| `Close` | Cleanup resources | Shutdown |

## Pattern Variations

### tendo Operations

```go
// Tensor operations are Chainable
result, err := tendo.MatMul(weights).Process(ctx, input)

// Compose into neural network
network := pipz.NewSequence[*tendo.Tensor](
    tendo.MatMul(w1),
    tendo.ReLU(),
    tendo.MatMul(w2),
    tendo.Softmax(-1),
)
```

### cogito Primitives

```go
// Reasoning primitives are Chainable
chain := pipz.NewSequence[*cogito.Thought](
    cogito.Analyze("input"),
    cogito.Decide("route", options),
    cogito.Checkpoint("save"),
)
```

## Trade-offs

**Interface overhead.** Every operation needs 4 methods. Most are trivial but required.

**Generic constraints.** All operations in a chain must have same T. Type conversion needed at boundaries.

**Learning curve.** Understanding composition patterns takes time.

## Related Patterns

- [Signal Emission](signal-emission.md) - Operations emit to capitan
- [Provider Abstraction](provider-abstraction.md) - Providers can wrap Chainables
