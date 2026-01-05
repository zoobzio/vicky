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

Common helper patterns:

```go
// Assertions - wrap t.Helper() for clean stack traces
func AssertEqual(t *testing.T, got, want any)
func AssertNoError(t *testing.T, err error)
func AssertContains(t *testing.T, slice []string, item string)

// Setup/teardown - manage test state
func ResetState(t *testing.T)
func WithTimeout(t *testing.T, d time.Duration) context.Context

// Fixtures - create test data
func NewTestUser(t *testing.T) User
func NewTestConfig(t *testing.T, opts ...Option) Config
```

Helpers should:
- Call `t.Helper()` for accurate line numbers in failures
- Accept `*testing.T` as first parameter
- Be specific to package domain, not generic utilities

### Benchmarks

Performance tests live in `testing/benchmarks/`. These measure execution time, memory allocation, and other performance characteristics.

### Integration

End-to-end tests live in `testing/integration/`. These test complete workflows, external dependencies, and provider implementations when applicable.

## Provider Tests

Packages with providers in `pkg/` follow additional testing patterns:

- Each provider has its own unit tests (colocated `_test.go` files)
- Integration tests should cover providers where applicable
- Provider tests may require external dependencies (databases, message queues, etc.)

## Coverage

Coverage thresholds ensure code quality without being punitive:

| Metric | Target | Threshold |
|--------|--------|-----------|
| Project | 70% | 1% drop allowed |
| Patch (new code) | 80% | 0% drop allowed |

### Rationale

- **70% project floor**: Achievable for most packages; forces coverage of core paths
- **80% patch target**: New code should be well-tested; prevents coverage erosion
- **1% threshold**: Allows minor fluctuations without blocking PRs

Configure in `.codecov.yml`:

```yaml
coverage:
  status:
    project:
      default:
        target: 70%
        threshold: 1%
    patch:
      default:
        target: 80%
```
