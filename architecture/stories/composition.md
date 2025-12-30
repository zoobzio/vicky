# Composition Story

*"How do I build reliable pipelines?"*

## The Flow

```
pipz → flume / streamz
```

## The Packages

### pipz - Composable Operations

A vocabulary around T. Reliability patterns that compose.

```go
pipeline := pipz.NewSequence[Order](
    pipz.Apply(validateOrder),
    pipz.NewRetry(
        pipz.Apply(chargePayment),
        pipz.WithMaxRetries(3),
        pipz.WithBackoff(pipz.ExponentialBackoff(100*time.Millisecond)),
    ),
    pipz.Apply(sendConfirmation),
)

result, err := pipeline.Process(ctx, order)
```

**Processors** (immutable values):
| Processor | Purpose |
|-----------|---------|
| Transform | Pure transformation (cannot fail) |
| Apply | Fallible transformation |
| Effect | Side effects (data passes through) |
| Mutate | Conditional transformation |
| Enrich | Optional enhancement |

**Connectors** (stateful containers):
| Category | Connectors |
|----------|------------|
| Sequential | Sequence |
| Parallel | Concurrent, Race, Contest, WorkerPool |
| Error handling | Retry, Backoff, Fallback, Handle |
| Flow control | Switch, Filter, Timeout |
| Resource protection | RateLimiter, CircuitBreaker |

### flume - Schema-Driven Pipelines

Define pipelines in YAML. Hot-reload without restarts.

```yaml
root:
  type: sequence
  steps:
    - ref: validate_order
    - type: retry
      max_retries: 3
      backoff: exponential
      processor:
        ref: charge_payment
    - ref: send_confirmation
```

```go
factory := flume.NewFactory[Order]()
factory.Add("validate_order", pipz.Apply(validateOrder))
factory.Add("charge_payment", pipz.Apply(chargePayment))
factory.Add("send_confirmation", pipz.Apply(sendConfirmation))

binding := factory.Bind(identity, schema)
// Later: binding.Update(newSchema) for hot-reload
```

### streamz - Channel Streaming

Same patterns for continuous data.

```go
out := streamz.Filter(
    streamz.Mapper(input, transformItem),
    isValid,
)

// With windowing
windows := streamz.TumblingWindow(input, 5*time.Minute, clock)

// With routing
success, errors := streamz.DeadLetterQueue(input, process)
```

**Processors:**
| Category | Processors |
|----------|------------|
| Transformation | Filter, Mapper, AsyncMapper, Tap, Sample |
| Routing | FanIn, FanOut, Partition, Switch, DeadLetterQueue |
| Batching | Batcher, Buffer |
| Windowing | TumblingWindow, SlidingWindow, SessionWindow |
| Flow Control | Throttle, Debounce |

## The Key Insight

**One vocabulary. Same patterns whether request/response or streaming.**

pipz defines the patterns. streamz applies them to channels. flume makes them configurable. All use `Chainable[T]` or `Result[T]`.

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Need                         │
├───────────────────────┬─────────────────────────────────────┤
│   Request/Response    │          Streaming                   │
│         │             │              │                       │
│         ▼             │              ▼                       │
│   ┌──────────┐        │       ┌──────────┐                  │
│   │   pipz   │        │       │  streamz │                  │
│   │Chainable[T]       │       │Result[T] │                  │
│   └────┬─────┘        │       └──────────┘                  │
│        │              │                                      │
│        ▼              │                                      │
│   ┌──────────┐        │                                      │
│   │  flume   │ (optional schema-driven)                     │
│   └──────────┘        │                                      │
└───────────────────────┴─────────────────────────────────────┘
```

## Pattern Examples

**Retry with exponential backoff:**
```go
pipz.NewRetry(processor,
    pipz.WithMaxRetries(5),
    pipz.WithBackoff(pipz.ExponentialBackoff(100*time.Millisecond)),
)
```

**Circuit breaker:**
```go
pipz.NewCircuitBreaker(processor,
    pipz.WithFailureThreshold(5),
    pipz.WithResetTimeout(30*time.Second),
)
```

**Rate limiting:**
```go
pipz.NewRateLimiter(processor,
    pipz.WithRate(100),
    pipz.WithBurst(10),
)
```

## Related Stories

- [Intelligence](intelligence.md) - zyn/cogito use pipz for reliability
- [Configuration](configuration.md) - flume enables runtime configuration
- [Numerical](numerical.md) - tendo operations are Chainable
