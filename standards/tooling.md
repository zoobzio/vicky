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

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused

linters-settings:
  errcheck:
    check-type-assertions: true
```

## Codecov

`.codecov.yml` configures coverage reporting:

```yaml
coverage:
  status:
    project:
      default:
        target: auto
        threshold: 1%
    patch:
      default:
        target: auto
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
      - '^docs:'
      - '^test:'
```

## Development Tools

Install development tools via Makefile:

```bash
make install-tools
```

This installs:

- `golangci-lint` - Linter
- Any package-specific tools
