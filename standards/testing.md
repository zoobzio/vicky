# Testing

Packages follow a structured approach to testing with clear separation between unit tests, integration tests, benchmarks, and test utilities.

## Unit Tests

Unit tests maintain a strict 1:1 relationship with source files:

```
/
├── api.go
├── api_test.go          # Tests for api.go
├── cache.go
├── cache_test.go        # Tests for cache.go
├── metadata.go
├── metadata_test.go     # Tests for metadata.go
```

### Rules

- Every `.go` file must have a corresponding `_test.go` file
- No `_test.go` file should exist without a corresponding `.go` file
- Tests are colocated with source in the root package

## testing/ Directory

The `testing/` directory contains test infrastructure:

```
testing/
├── helpers.go           # Test utilities and helpers
├── helpers_test.go      # Tests for the helpers themselves
├── benchmarks/          # Performance tests
└── integration/         # End-to-end tests
```

### Helpers

Test helpers provide shared utilities for testing across the package. Helpers must themselves be tested to ensure reliability.

### Benchmarks

Performance tests live in `testing/benchmarks/`. These measure execution time, memory allocation, and other performance characteristics.

### Integration

End-to-end tests live in `testing/integration/`. These test complete workflows, external dependencies, and provider implementations when applicable.

## Provider Tests

Packages with providers in `pkg/` follow additional testing patterns:

- Each provider has its own unit tests (colocated `_test.go` files)
- Integration tests should cover providers where applicable
- Provider tests may require external dependencies (databases, message queues, etc.)
