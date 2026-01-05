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
          go-version: '1.25'
      - uses: golangci/golangci-lint-action@v7

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - uses: securego/gosec@master

  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: make coverage
      - uses: codecov/codecov-action@v4

  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
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

## Additional Workflows

Beyond the main CI workflow, packages include:

### Release Workflow

`.github/workflows/release.yml` triggers on version tags and automates releases via GoReleaser:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'

permissions:
  contents: write

jobs:
  validate:
    # Run tests and linting before release
    steps:
      - run: go mod tidy && git diff --exit-code go.mod go.sum
      - run: make test
      - uses: golangci/golangci-lint-action@v7

  release:
    needs: validate
    steps:
      - uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### CodeQL Workflow

`.github/workflows/codeql.yml` provides GitHub's semantic code analysis:

```yaml
name: CodeQL

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 0 * * 0'  # Weekly

jobs:
  analyze:
    permissions:
      security-events: write
    steps:
      - uses: github/codeql-action/init@v3
        with:
          languages: go
      - uses: github/codeql-action/autobuild@v3
      - uses: github/codeql-action/analyze@v3
```
