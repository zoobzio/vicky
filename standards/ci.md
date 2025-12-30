# Continuous Integration

All packages use GitHub Actions for CI with a consistent workflow structure that relies on Makefile targets.

## Workflow Structure

The CI workflow (`/.github/workflows/ci.yml`) consists of these jobs:

```
test → lint → security → coverage → benchmark → ci-complete
```

### Jobs

| Job | Purpose | Makefile Target |
|-----|---------|-----------------|
| test | Run all tests | `make test` |
| lint | Static analysis | `make lint` |
| security | Security scanning | gosec |
| coverage | Coverage reporting | `make coverage` |
| benchmark | Performance tracking | `make test-bench` |
| ci-complete | Aggregated status | - |

## Example Workflow

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.24', '1.25']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - run: make test

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - uses: golangci/golangci-lint-action@v6

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: securego/gosec@master

  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: make coverage
      - uses: codecov/codecov-action@v4

  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: make test-bench

  ci-complete:
    needs: [test, lint, security, coverage, benchmark]
    runs-on: ubuntu-latest
    steps:
      - run: echo "CI complete"
```

## Principles

- **Makefile reliance**: CI calls Makefile targets, no duplicated logic
- **Matrix testing**: Test against Go 1.24 and 1.25
- **Aggregated status**: `ci-complete` job provides single status check for branch protection
