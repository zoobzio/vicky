# Tooling

All packages use a consistent set of configuration files for linting, coverage, and releases.

## Required Configuration Files

```
/
├── .golangci.yml        # Linter configuration
├── .codecov.yml         # Coverage configuration
└── .goreleaser.yml      # Release configuration
```

## golangci-lint

`.golangci.yml` configures static analysis. All packages should use a consistent linter configuration.

### Minimum Required

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - ineffassign
    - staticcheck
    - unused

linters-settings:
  errcheck:
    check-type-assertions: true
```

### Recommended Full Configuration

Production packages should enable additional linters for security and quality:

```yaml
version: "2"

run:
  timeout: 5m
  tests: true

linters:
  enable:
    # Required (minimum)
    - errcheck
    - govet
    - ineffassign
    - staticcheck
    - unused

    # Security
    - gosec
    - noctx
    - bodyclose
    - sqlclosecheck

    # Error handling
    - errorlint
    - errchkjson
    - wastedassign

    # Best practices
    - gocritic
    - revive
    - unconvert
    - dupl
    - goconst
    - misspell
    - prealloc
    - copyloopvar

  exclusions:
    rules:
      - path: _test\.go
        linters:
          - dupl
          - goconst
          - govet

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true

  govet:
    enable-all: true

  dupl:
    threshold: 150

  goconst:
    min-len: 3
    min-occurrences: 3
```

### Linter Categories

| Category | Linters                                              | Purpose                 |
| -------- | ---------------------------------------------------- | ----------------------- |
| Required | errcheck, govet, ineffassign, staticcheck, unused    | Core correctness        |
| Security | gosec, noctx, bodyclose, sqlclosecheck               | Vulnerability detection |
| Errors   | errorlint, errchkjson, wastedassign                  | Error handling          |
| Quality  | gocritic, revive, unconvert, dupl, goconst, misspell | Code quality            |

## Codecov

`.codecov.yml` configures coverage reporting. See [testing.md](testing.md#coverage) for threshold rationale.

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

## GoReleaser

`.goreleaser.yml` is present but skips binary builds (all packages are libraries):

```yaml
version: 2

builds:
  - skip: true

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
```

## Development Tools

Install development tools via Makefile:

```bash
make install-tools
```

This installs:

- `golangci-lint` - Linter
- Any package-specific tools
