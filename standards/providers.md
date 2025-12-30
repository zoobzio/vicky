# Providers

Packages that support pluggable implementations use a consistent provider pattern.

## When to Use Providers

Use the `pkg/` structure when a package:

- Supports multiple backends (databases, message queues, LLM services)
- Has pluggable implementations of a core interface
- Separates core logic from external dependencies

## Directory Structure

```
/
├── api.go                    # Core interface definitions
├── [core].go                 # Core implementation
├── [core]_test.go
└── pkg/
    ├── provider-a/
    │   ├── provider.go
    │   └── provider_test.go
    ├── provider-b/
    │   ├── provider.go
    │   └── provider_test.go
    └── provider-c/
        ├── provider.go
        └── provider_test.go
```

## Examples

| Package | Providers |
|---------|-----------|
| herald | amqp, bolt, firestore, jetstream, kafka, nats, redis, etc. |
| zyn | anthropic, gemini, openai |
| astql | postgres, mysql, mssql, sqlite |

## Testing

### Unit Tests

Each provider maintains its own unit tests:

```
pkg/postgres/
├── provider.go
└── provider_test.go     # Unit tests for postgres provider
```

### Integration Tests

Integration tests in `testing/integration/` should cover providers where applicable. These tests may require external dependencies (Docker, testcontainers, etc.).

## Build

- Providers do not have their own Makefile
- The root Makefile handles all operations
- Provider tests run as part of `make test`
