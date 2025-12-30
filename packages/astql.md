# Astql

Type-safe SQL query builder for Go.

[GitHub](https://github.com/zoobzio/astql)

## Vision

Schema-validated SQL generation. DBML schema as allowlist for all identifiers. Validate at construction, not execution. Multi-dialect support from a single AST. Security through defense in depth.

## Design Decisions

**Three-stage pipeline**
Validation → AST → Rendering. Schema validation happens first. Invalid identifiers never enter the system. Providers are stateless renderers that never see raw user input.

**Instance-based API**
All identifiers created through ASTQL instance methods (T, F, P, C). Types are internal, cannot be constructed directly. All construction goes through validation.

**Multi-layer security**
1. Schema validation (table/field allowlist via DBML)
2. Identifier validation (alphanumeric + underscore only)
3. SQL keyword blocking
4. Alias restrictions (single lowercase letter)
5. Quoted identifiers (dialect-specific escaping)
6. Parameterised queries (values never interpolated)
7. Subquery depth limits (max 3 levels)

**Provider pattern**
Single AST, multiple dialect-specific renderers. Each provider handles its own syntax. Easy to add new dialects without changing core logic.

**Panicking vs error variants**
T/TryT, F/TryF, P/TryP. Panicking variants for validated schemas. Error-returning variants for runtime input. Try variants mandatory for user-provided input.

## Providers

| Provider | Quoting | Placeholder | RETURNING | Upsert |
|----------|---------|-------------|-----------|--------|
| PostgreSQL | `"name"` | `:name` | Native | ON CONFLICT |
| MySQL | `` `name` `` | `:name` | Unsupported | ON DUPLICATE KEY |
| SQLite | `"name"` | `:name` | Native | ON CONFLICT |
| SQL Server | `[name]` | `@name` | OUTPUT | Unsupported |

Unsupported features raise errors rather than silently falling back.

## Internal Architecture

```
ASTQL Instance (holds DBML schema)
        ↓
    T(), F(), P(), C() → Validated identifiers
        ↓
    Builder (fluent API)
        ↓
    AST (intermediate representation)
        ↓
    Provider.Render() → SQL + RequiredParams
```

## Code Organisation

| File | Responsibility |
|------|----------------|
| `instance.go` | ASTQL instance, validation logic |
| `builder.go` | Fluent query builder |
| `expressions.go` | Aggregates, conditions, case, window functions |
| `renderer.go` | Renderer interface |
| `internal/types/` | AST structures (table, field, param, condition, operator) |
| `pkg/postgres/` | PostgreSQL renderer |
| `pkg/mysql/` | MySQL renderer |
| `pkg/sqlite/` | SQLite renderer |
| `pkg/mssql/` | SQL Server renderer |

## Current State / Direction

Stable. Four providers complete. Core security model established.

Future considerations:
- Additional providers as needs arise
- Extended expression support

## Framework Context

**Dependencies**: dbml (schema parsing). Database drivers for testing only.

**Role**: SQL generation layer. Schema-validated, injection-resistant query building. Generates SQL, doesn't execute it. Used by soy for higher-level database operations.
