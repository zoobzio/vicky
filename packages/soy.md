# Soy

Type-safe database access for Go.

[GitHub](https://github.com/zoobzio/soy)

## Vision

Schema-validated queries with compile-time safety. Struct tags as schema definition. Reflection once at initialisation, zero on hot path. Multi-database through astql providers. Type-safe results via generics.

## Design Decisions

**Layered architecture**
Soy (ergonomic API) → ASTQL (validation) → sqlx (execution). Each layer has single responsibility. Validation at boundaries.

**Generic types**
`Soy[T]` flows type safety through entire chain. Select returns `*T`. Query returns `[]*T`. No `interface{}` gymnastics.

**Metadata caching**
Sentinel inspects struct once in New(). Cached forever. Zero reflection on execution hot path.

**Named parameters**
All values via `:param` syntax. SQL injection prevention built in. Prepared statement reuse.

**Required WHERE**
UPDATE and DELETE refuse execution without WHERE clause. Checked at Exec() time. Prevents accidental data loss.

**Error deferral**
Builders chain fluently without error checking at each step. Error captured internally, reported at Exec().

## Schema Flow

```
Struct Tags → Sentinel Metadata → DBML → ASTQL Instance
                                            ↓
                                      Validation Engine
```

Tags registered: `db`, `type`, `constraints`, `default`, `check`, `index`, `references`.

## Builders

| Builder | Purpose | Returns |
|---------|---------|---------|
| Select | Single record (adds LIMIT 1) | `*T` |
| Query | Multiple records | `[]*T` |
| Create | INSERT with RETURNING | `*T` |
| Update | UPDATE with required WHERE | `*T` |
| Delete | DELETE with required WHERE | `int64` (rows affected) |
| Aggregate | COUNT/SUM/AVG/MIN/MAX | varies |
| Compound | UNION/INTERSECT/EXCEPT | `[]*T` |

## Internal Architecture

```
Soy[T]
├── metadata (sentinel.Metadata, cached)
├── instance (astql.ASTQL, validates fields)
├── db (sqlx.DB)
└── renderer (astql provider)
        ↓
Builder.Exec(ctx, params)
        ↓
ASTQL Render → SQL + RequiredParams
        ↓
sqlx.NamedQueryContext → rows
        ↓
StructScan → *T or []*T
```

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Soy[T], New(), entry points |
| `select.go` | Single-record queries |
| `query.go` | Multi-record queries |
| `create.go` | INSERT operations |
| `update.go` | UPDATE with WHERE enforcement |
| `delete.go` | DELETE with WHERE enforcement |
| `aggregate.go` | COUNT/SUM/AVG/MIN/MAX |
| `compound.go` | UNION/INTERSECT/EXCEPT |
| `where.go` | Shared WHERE builder logic |
| `batch.go` | Bulk operations |
| `dbml.go` | Struct tags → DBML conversion |
| `events.go` | Capitan signals |

## Current State / Direction

Stable. Core operations complete. Window functions and CASE expressions working.

Future considerations:
- Additional expression types
- Schema migration helpers

## Framework Context

**Dependencies**: astql (SQL validation), sentinel (metadata), dbml (schema), capitan (events), sqlx (execution).

**Role**: Database layer. Type-safe queries built on astql's schema validation. Sentinel provides type intelligence for metadata caching. Used by cogito for SoyMemory persistence.
