# Integration Tests

End-to-end tests for vicky.

## Prerequisites

- PostgreSQL (pgvector)
- MinIO (S3-compatible object storage)
- Redis

## Running

```bash
make test-integration
# or
go test -v -race -tags testing ./testing/integration/...
```

## Test Isolation

Each test should set up and tear down its own state.
