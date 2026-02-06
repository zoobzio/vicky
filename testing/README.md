# Testing

Testing infrastructure for vicky.

## Running Tests

```bash
make test              # All tests with race detector
make test-unit         # Unit tests only (short mode)
make test-integration  # Integration tests
make test-bench        # Benchmarks
make coverage          # Combined coverage report
```

## Structure

- `helpers.go` — Domain-specific test utilities (fixtures, registry setup)
- `mocks.go` — Function-field mock implementations of contracts
- `benchmarks/` — Performance tests
- `integration/` — End-to-end tests

## Coverage

Target: 70% project, 80% new code (patch).
