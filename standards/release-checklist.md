# Public Release Checklist

Audit checklist for packages prior to public release. Derived from Vicky standards with sentinel as reference implementation.

## How to Use

For each item:

- **[✓]** Compliant
- **[~]** Partially compliant (note what's missing)
- **[✗]** Non-compliant
- **[N/A]** Not applicable to this package

---

## 1. Go Module

| Item                        | Standard                                          | Check |
| --------------------------- | ------------------------------------------------- | ----- |
| Module path                 | `github.com/zoobzio/[package]`                    | [ ]   |
| Go version                  | `go 1.24` minimum                                 | [ ]   |
| Toolchain directive         | `toolchain go1.25.x`                              | [ ]   |
| No unnecessary dependencies | Layer 0 packages: stdlib only (yaml.v3 exception) | [ ]   |

---

## 2. Package Structure

| Item             | Standard                                  | Check |
| ---------------- | ----------------------------------------- | ----- |
| `api.go` exists  | Public interface entry point              | [ ]   |
| Feature modules  | `[feature].go` files for each feature     | [ ]   |
| 1:1 test mapping | Every `.go` has corresponding `_test.go`* | [ ]   |
| No orphan tests  | No `_test.go` without corresponding `.go` | [ ]   |

\* Exception: Files containing only delegation, re-exports, or trivial wiring (e.g., `api.go` that delegates to internal implementation) may omit corresponding test files if there is no testable logic.

### Optional Directories

| Item        | Standard                 | When Required        | Check |
| ----------- | ------------------------ | -------------------- | ----- |
| `internal/` | Private implementation   | If hiding internals  | [ ]   |
| `pkg/`      | Provider implementations | If multiple backends | [ ]   |

---

## 3. Testing Structure

| Item                      | Standard                   | Check |
| ------------------------- | -------------------------- | ----- |
| `testing/helpers.go`      | Test utilities and helpers | [ ]   |
| `testing/helpers_test.go` | Helpers themselves tested  | [ ]   |
| `testing/benchmarks/`     | Performance tests          | [ ]   |
| `testing/integration/`    | End-to-end tests           | [ ]   |

### Testing Documentation

| File                          | Standard                    | Check |
| ----------------------------- | --------------------------- | ----- |
| `testing/README.md`           | Testing strategy overview   | [ ]   |
| `testing/integration/README.md` | Integration test guidance | [ ]   |
| `testing/benchmarks/README.md`  | Benchmark documentation   | [ ]   |

### Helper Conventions

| Item                              | Standard              | Check |
| --------------------------------- | --------------------- | ----- |
| Helpers call `t.Helper()`         | Clean stack traces    | [ ]   |
| Helpers accept `*testing.T` first | Consistent signature  | [ ]   |
| Helpers are domain-specific       | Not generic utilities | [ ]   |

---

## 4. Repository Files

### Required Files

| File              | Standard                                   | Check |
| ----------------- | ------------------------------------------ | ----- |
| `LICENSE`         | MIT License with correct year              | [ ]   |
| `CONTRIBUTING.md` | Development workflow, links to `make help` | [ ]   |
| `SECURITY.md`     | Vulnerability reporting instructions       | [ ]   |
| `.gitignore`      | Standard Go ignores                        | [ ]   |

### .gitignore Contents

| Item                                    | Check |
| --------------------------------------- | ----- |
| Binaries (_.exe, _.dll, _.so, _.dylib)  | [ ]   |
| Test files (_.test, _.out, coverage.\*) | [ ]   |
| IDE files (.idea/, .vscode/, \*.swp)    | [ ]   |
| OS files (.DS_Store, Thumbs.db)         | [ ]   |
| Build artifacts (dist/, vendor/)        | [ ]   |

### GitHub Templates

| File                                    | Standard                        | Check |
| --------------------------------------- | ------------------------------- | ----- |
| `.github/PULL_REQUEST_TEMPLATE.md`      | Standardised PR format          | [ ]   |
| `.github/ISSUE_TEMPLATE/bug_report.md`  | Bug report template             | [ ]   |
| `.github/ISSUE_TEMPLATE/feature_request.md` | Feature request template    | [ ]   |
| `.github/ISSUE_TEMPLATE/documentation.md` | Documentation issue template  | [ ]   |

---

## 5. Tooling Configuration

### Required Config Files

| File              | Standard               | Check |
| ----------------- | ---------------------- | ----- |
| `.golangci.yml`   | Linter configuration   | [ ]   |
| `.codecov.yml`    | Coverage configuration | [ ]   |
| `.goreleaser.yml` | Release configuration  | [ ]   |

### .golangci.yml - Minimum Linters

| Linter      | Category | Check |
| ----------- | -------- | ----- |
| errcheck    | Required | [ ]   |
| govet       | Required | [ ]   |
| ineffassign | Required | [ ]   |
| staticcheck | Required | [ ]   |
| unused      | Required | [ ]   |

### .golangci.yml - Recommended Linters

| Linter        | Category       | Check |
| ------------- | -------------- | ----- |
| gosec         | Security       | [ ]   |
| noctx         | Security       | [ ]   |
| bodyclose     | Security       | [ ]   |
| sqlclosecheck | Security       | [ ]   |
| errorlint     | Error handling | [ ]   |
| errchkjson    | Error handling | [ ]   |
| wastedassign  | Error handling | [ ]   |
| gocritic      | Best practices | [ ]   |
| revive        | Best practices | [ ]   |
| unconvert     | Best practices | [ ]   |
| dupl          | Best practices | [ ]   |
| goconst       | Best practices | [ ]   |
| misspell      | Best practices | [ ]   |
| prealloc      | Best practices | [ ]   |
| copyloopvar   | Best practices | [ ]   |

### .golangci.yml - Settings

| Setting                          | Standard                                    | Check |
| -------------------------------- | ------------------------------------------- | ----- |
| `version`                        | `"2"` (requires golangci-lint v2.x, note `/v2/` in module path) | [ ]   |
| `run.timeout`                    | 5m                                          | [ ]   |
| `run.tests`                      | true                                        | [ ]   |
| Test file exclusions             | dupl, goconst, govet excluded for \_test.go | [ ]   |
| `errcheck.check-type-assertions` | true                                        | [ ]   |
| `govet.enable-all`               | true                                        | [ ]   |
| `dupl.threshold`                 | 150                                         | [ ]   |
| `goconst.min-len`                | 3                                           | [ ]   |
| `goconst.min-occurrences`        | 3                                           | [ ]   |

### .codecov.yml - Thresholds

| Setting           | Standard | Check |
| ----------------- | -------- | ----- |
| Project target    | 70%      | [ ]   |
| Project threshold | 1%       | [ ]   |
| Patch target      | 80%      | [ ]   |
| Patch threshold   | 0%       | [ ]   |

### .goreleaser.yml

| Setting           | Standard                              | Check |
| ----------------- | ------------------------------------- | ----- |
| `version`         | 2                                     | [ ]   |
| `builds.skip`     | true (libraries don't build binaries) | [ ]   |
| Changelog filters | Exclude docs:, test: prefixes         | [ ]   |

---

## 6. Makefile

### Required Targets

| Target             | Purpose                          | Check |
| ------------------ | -------------------------------- | ----- |
| `test`             | Run all tests with race detector | [ ]   |
| `test-unit`        | Run unit tests only (short mode) | [ ]   |
| `test-integration` | Run integration tests            | [ ]   |
| `test-bench`       | Run benchmarks                   | [ ]   |
| `lint`             | Run linters                      | [ ]   |
| `lint-fix`         | Run linters with auto-fix        | [ ]   |
| `coverage`         | Generate coverage report         | [ ]   |
| `clean`            | Remove generated files           | [ ]   |
| `check`            | Quick validation (test + lint)   | [ ]   |
| `ci`               | Full CI simulation               | [ ]   |
| `install-tools`    | Install dev tools via `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.x.x` | [ ]   |
| `install-hooks`    | Install git hooks                | [ ]   |
| `help`             | Display available commands       | [ ]   |

### Makefile Conventions

| Item                 | Standard                  | Check |
| -------------------- | ------------------------- | ----- |
| `.PHONY` declaration | All targets listed        | [ ]   |
| `.DEFAULT_GOAL`      | Set to `help`             | [ ]   |
| `help` target        | Self-documenting via grep | [ ]   |

---

## 7. CI Workflow

### .github/workflows/ci.yml

| Job         | Purpose                               | Check |
| ----------- | ------------------------------------- | ----- |
| test        | Run tests with matrix (Go 1.24, 1.25) | [ ]   |
| lint        | golangci-lint                         | [ ]   |
| security    | gosec scanner with SARIF upload       | [ ]   |
| coverage    | Upload to Codecov                     | [ ]   |
| benchmark   | Performance tracking                  | [ ]   |
| ci-complete | Aggregated status check               | [ ]   |

### Security Scanning

| Item            | Standard                              | Check |
| --------------- | ------------------------------------- | ----- |
| gosec runs      | Security scanner enabled              | [ ]   |
| SARIF output    | Results in SARIF format               | [ ]   |
| GitHub upload   | Upload to GitHub Code Scanning        | [ ]   |
| CodeQL workflow | `.github/workflows/codeql.yml` exists | [ ]   |

### CI Conventions

| Item                                 | Standard                           | Check |
| ------------------------------------ | ---------------------------------- | ----- |
| Triggers                             | push to main, pull_request to main | [ ]   |
| Go matrix                            | 1.24 and 1.25                      | [ ]   |
| Lint on latest Go only               | 1.25                               | [ ]   |
| Coverage upload on latest Go only    | 1.25                               | [ ]   |
| golangci-lint action                 | `golangci/golangci-lint-action@v7` | [ ]   |
| Uses Makefile targets where possible | Consistency with local dev         | [ ]   |

### .github/workflows/release.yml

| Item         | Standard                               | Check |
| ------------ | -------------------------------------- | ----- |
| Trigger      | `push: tags: ['v*.*.*']`               | [ ]   |
| Validate job | Tests and lints before release         | [ ]   |
| Release job  | Uses `goreleaser/goreleaser-action@v6` | [ ]   |
| Permissions  | `contents: write`                      | [ ]   |

### .github/workflows/codeql.yml

| Item        | Standard                            | Check |
| ----------- | ----------------------------------- | ----- |
| Triggers    | push, pull_request, weekly schedule | [ ]   |
| Language    | `go`                                | [ ]   |
| Permissions | `security-events: write`            | [ ]   |

### .github/workflows/version-preview.yml (Optional)

| Item        | Standard                                  | Check |
| ----------- | ----------------------------------------- | ----- |
| Trigger     | `pull_request`                            | [ ]   |
| Purpose     | Comments next version on PRs              | [ ]   |
| Uses        | Conventional commits to infer version     | [ ]   |

---

## 8. README

The README structure provides rhythm; the library's character dictates voice. These are scaffolding, not script.

### Structure

| Section       | Purpose                                              | Check |
| ------------- | ---------------------------------------------------- | ----- |
| Header        | Title, badges, tagline, supporting sentence          | [ ]   |
| Essence       | The unique hook (section name specific to library)   | [ ]   |
| Install       | go get + requirements                                | [ ]   |
| Quick Start   | Complete runnable example with imports               | [ ]   |
| Capabilities  | Optional — only if distinct modes worth highlighting | [ ]   |
| Why [Name]?   | Concrete benefits, not marketing                     | [ ]   |
| Documentation | Categorised links                                    | [ ]   |
| Contributing  | Brief, link to CONTRIBUTING.md                       | [ ]   |
| License       | One line                                             | [ ]   |

### Essence Section

The essence section name should complete: "[Library] is about \_\_\_."

| Item                      | Standard                                       | Check |
| ------------------------- | ---------------------------------------------- | ----- |
| Specific name             | Not "Overview", "Introduction", "How It Works" | [ ]   |
| Shows the insight         | Minimal code demonstrating core value          | [ ]   |
| Matches library character | Voice derived from what makes it unique        | [ ]   |

### Badges

| Badge                     | Check |
| ------------------------- | ----- |
| CI Status                 | [ ]   |
| Codecov                   | [ ]   |
| Go Report Card            | [ ]   |
| CodeQL                    | [ ]   |
| Go Reference (pkg.go.dev) | [ ]   |
| License                   | [ ]   |
| Go Version                | [ ]   |
| Release                   | [ ]   |

---

## 9. Documentation (docs/)

### Directory Structure

| Path                           | Purpose                | Check |
| ------------------------------ | ---------------------- | ----- |
| `docs/1.overview.md`           | Package overview       | [ ]   |
| `docs/2.learn/`                | Learning materials     | [ ]   |
| `docs/2.learn/1.quickstart.md` | Getting started        | [ ]   |
| `docs/2.learn/2.concepts.md`   | Core concepts          | [ ]   |
| `docs/2.learn/3.architecture.md` | System design        | [ ]   |
| `docs/3.guides/`               | How-to guides          | [ ]   |
| `docs/4.cookbook/`             | Recipes and patterns   | [ ]   |
| `docs/5.reference/`            | API reference          | [ ]   |
| `docs/5.reference/1.api.md`    | Function documentation | [ ]   |

### Frontmatter

| Field         | Check |
| ------------- | ----- |
| `title`       | [ ]   |
| `description` | [ ]   |
| `author`      | [ ]   |
| `published`   | [ ]   |
| `updated`     | [ ]   |
| `tags`        | [ ]   |

### Documentation Conventions

| Item                           | Standard                      | Check |
| ------------------------------ | ----------------------------- | ----- |
| Numbered prefixes              | Control ordering              | [ ]   |
| Clear and direct tone          | Technical but accessible      | [ ]   |
| Example-driven                 | Practical examples throughout | [ ]   |
| Explains "why" alongside "how" | Context for decisions         | [ ]   |

---

## 10. Provider Pattern (if applicable)

Only complete if package supports multiple backends.

| Item                              | Standard                   | Check |
| --------------------------------- | -------------------------- | ----- |
| Core interfaces in root `api.go`  | Provider interface defined | [ ]   |
| Providers in `pkg/[provider]/`    | Each provider isolated     | [ ]   |
| Provider has `provider.go`        | Implementation             | [ ]   |
| Provider has `provider_test.go`   | Unit tests                 | [ ]   |
| No provider-specific Makefile     | Root Makefile handles all  | [ ]   |
| Integration tests cover providers | In `testing/integration/`  | [ ]   |

---

## Summary

| Category              | Items | Compliant | Partial | Non-Compliant |
| --------------------- | ----- | --------- | ------- | ------------- |
| Go Module             |       |           |         |               |
| Package Structure     |       |           |         |               |
| Testing Structure     |       |           |         |               |
| Repository Files      |       |           |         |               |
| Tooling Configuration |       |           |         |               |
| Makefile              |       |           |         |               |
| CI Workflow           |       |           |         |               |
| README                |       |           |         |               |
| Documentation         |       |           |         |               |
| Provider Pattern      |       |           |         |               |
| **TOTAL**             |       |           |         |               |
