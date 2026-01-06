# Makefile

The Makefile serves as the single source of truth for all package operations. CI workflows rely exclusively on Makefile targets.

## Required Targets

Every package must implement these targets:

```makefile
.PHONY: test lint lint-fix coverage clean check ci help test-unit test-integration test-bench install-hooks install-tools

## Testing
test:              ## Run all tests
	go test -race ./...

test-unit:         ## Run unit tests only
	go test -race -short ./...

test-integration:  ## Run integration tests
	go test -race ./testing/integration/...

test-bench:        ## Run benchmarks
	go test -bench=. ./testing/benchmarks/...

## Linting
lint:              ## Run linter
	golangci-lint run

lint-fix:          ## Run linter with auto-fix
	golangci-lint run --fix

## Coverage
coverage:          ## Generate coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Tooling
install-hooks:     ## Install git hooks
	# Implementation specific to package

install-tools:     ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## Maintenance
clean:             ## Remove generated files
	rm -f coverage.out coverage.html

## Workflow
check:             ## Quick validation (test + lint)
	$(MAKE) test
	$(MAKE) lint

ci:                ## Full CI simulation
	$(MAKE) clean
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) coverage

## Help
help:              ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
```

## Rationale

- **Single source of truth**: All operations defined in one place
- **CI consistency**: Workflows call Makefile targets, no duplicated logic
- **Discoverability**: `make help` documents all available operations
- **Portability**: Developers run the same commands locally as CI
