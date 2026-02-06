.PHONY: test lint lint-fix coverage clean check ci help test-unit test-integration test-bench install-hooks install-tools build run migrate

## Build
build:            ## Build the vicky binary
	go build -o bin/vicky ./cmd/vicky

run:              ## Run the server locally
	go run ./cmd/vicky

## Testing
test:             ## Run all tests
	go test -tags testing -race ./...

test-unit:        ## Run unit tests only
	go test -tags testing -race -short ./...

test-integration: ## Run integration tests
	go test -tags testing -race ./testing/integration/...

test-bench:       ## Run benchmarks
	go test -tags testing -bench=. -benchmem ./testing/benchmarks/...

## Linting
lint:             ## Run linter
	golangci-lint run

lint-fix:         ## Run linter with auto-fix
	golangci-lint run --fix

## Coverage
coverage:         ## Generate coverage report (unit + integration)
	@go test -tags testing -coverprofile=coverage-unit.out -covermode=atomic ./...
	@go test -tags testing -coverprofile=coverage-integration.out -covermode=atomic ./testing/integration/... 2>/dev/null || true
	@echo "mode: atomic" > coverage.out
	@tail -n +2 coverage-unit.out >> coverage.out
	@tail -n +2 coverage-integration.out >> coverage.out 2>/dev/null || true
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | tail -1

## Database
migrate:          ## Run database migrations
	@echo "Running migrations..."
	@for f in migrations/*.up.sql; do \
		echo "Applying $$f"; \
		psql "$$VICKY_DB_DSN" -f "$$f"; \
	done

migrate-down:     ## Rollback database migrations
	@echo "Rolling back migrations..."
	@for f in $$(ls -r migrations/*.down.sql); do \
		echo "Rolling back $$f"; \
		psql "$$VICKY_DB_DSN" -f "$$f"; \
	done

## Docker
docker-build:     ## Build Docker image
	docker build -t vicky:latest .

docker-up:        ## Start local development environment
	docker-compose up -d

docker-down:      ## Stop local development environment
	docker-compose down

docker-logs:      ## View container logs
	docker-compose logs -f

## Tooling
install-hooks:    ## Install git hooks
	@echo "No hooks configured"

install-tools:    ## Install development tools
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2

## Maintenance
clean:            ## Remove generated files
	rm -f coverage.out coverage-unit.out coverage-integration.out coverage.html
	rm -rf bin/

## Workflow
check:            ## Quick validation (test + lint)
	$(MAKE) test
	$(MAKE) lint

ci:               ## Full CI simulation
	$(MAKE) clean
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) coverage

## Help
help:             ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
