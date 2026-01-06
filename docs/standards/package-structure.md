# Package Structure

All packages follow a consistent directory and file organisation.

## Directory Layout

```
/
├── api.go               # Public interface entry point
├── [feature].go         # Feature modules
├── [feature]_test.go    # Corresponding unit tests
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── .codecov.yml
├── .goreleaser.yml
├── .github/
│   └── workflows/
│       └── ci.yml
├── docs/
├── testing/
│   ├── helpers.go
│   ├── helpers_test.go
│   ├── benchmarks/
│   └── integration/
├── internal/            # Optional: private implementation
├── [provider-a]/        # Optional: provider implementations
├── [provider-b]/
└── [provider-c]/
```

## Entry Point

Every package exposes its public interface through `api.go`. This file serves as the primary entry point and should contain or re-export all public types, functions, and interfaces.

## Module Naming

```
github.com/zoobzio/[package-name]
```

## internal/

Use `internal/` when a package requires private implementation details that should not be importable by external consumers. This is discretionary based on package needs.

## Providers

Provider implementations live in root-level directories (e.g., `redis/`, `postgres/`, `s3/`). See [providers](./providers.md) for details.
