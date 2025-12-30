# Provider Abstraction

*Single interface, multiple backends.*

## The Problem

Applications need to interact with different backends - databases, message brokers, LLM APIs, storage systems. Without abstraction, changing backends means changing application code.

## The Pattern

Define a minimal interface. Implement it for each backend. Application code uses the interface, never the concrete implementation.

```go
// Interface in core package
type Provider interface {
    Publish(ctx context.Context, msg Message) error
    Subscribe(ctx context.Context, handler func(Message)) error
    Ping(ctx context.Context) error
    Close() error
}

// Implementation in pkg/ subpackage
// pkg/kafka/provider.go
// pkg/nats/provider.go
// pkg/redis/provider.go
```

## Where It's Used

| Package | Interface | Providers |
|---------|-----------|-----------|
| herald | Provider (4 methods) | 11 message brokers |
| grub | Provider | 9 storage backends |
| flux | Watcher | 8 configuration sources |
| astql | Renderer | 4 SQL dialects |
| vectql | Renderer | 4 vector databases |
| zyn | Provider | 3 LLM APIs |
| tendo | Backend | CPU, CUDA |

## Implementation Details

### Minimal Interface

Keep interfaces small. herald's Provider has just 4 methods:
- Publish
- Subscribe
- Ping
- Close

More methods = harder to implement = fewer providers.

### Optional Capabilities

Some backends support features others don't. Use type assertion:

```go
type Lifecycle interface {
    Connect(ctx context.Context) error
    Close() error
    Health(ctx context.Context) error
}

// Check if provider supports lifecycle
if lc, ok := provider.(Lifecycle); ok {
    lc.Connect(ctx)
}
```

### Isolated Dependencies

Each provider lives in its own subpackage. Dependencies isolated:

```
herald/
├── api.go            # Interface definition
├── publisher.go      # Uses interface
├── pkg/
│   ├── kafka/        # Only imports kafka client
│   ├── nats/         # Only imports nats client
│   └── redis/        # Only imports redis client
```

Users only pull dependencies they use.

## Why This Pattern

**No vendor lock-in.** Switch from Redis to Kafka by changing one line.

**Easy to extend.** New backend = new implementation. No core changes.

**Testable.** Mock implementations for unit tests.

**Clear contracts.** Interface documents what backends must provide.

## Trade-offs

**Lowest common denominator.** Interface can only include features all backends support. Advanced features require provider-specific code.

**Configuration complexity.** Each provider has different configuration. Provider-specific options needed.

**Testing burden.** Each provider needs testing. Integration tests with real backends required.

## Alternatives Considered

**Driver pattern** (database/sql style): More complex, better for databases specifically. Provider pattern simpler for general use.

**Plugin system**: Dynamic loading adds complexity. Static compilation preferred for reliability.

## Related Patterns

- [Layered Architecture](layered-architecture.md) - Providers live at specific layers
- [Signal Emission](signal-emission.md) - Providers emit to capitan
