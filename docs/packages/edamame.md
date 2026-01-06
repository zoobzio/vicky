# Edamame

Capability-driven query factory for Go.

[GitHub](https://github.com/zoobzio/edamame)

## Vision

Named, introspectable database operations. Specs are pure data (JSON-serialisable), safe for LLM consumption. Bridge between scattered SQL strings and heavyweight ORMs. Semantic layer over soy adding capability management and introspection.

## Design Decisions

**Spec-to-builder pattern**
Pure data specs converted to soy builders on demand. Specs are introspectable, exportable for LLM consumption, cacheable after rendering. No code injection risk.

**Generic Factory[T]**
Wraps `soy.Soy[T]`. Type-safe results throughout. Query returns `[]*T`, Select returns `*T`.

**Thread-safe registry**
RWMutex for all capability registries. Concurrent query execution (read-lock). Serialised registration (write-lock). No performance penalty for read-heavy workloads.

**Default CRUD**
Every factory auto-registers query, select, delete, count. Derived from struct metadata. Can be overridden.

**SQL caching**
Rendered SQL cached by "type:name" key. Avoids re-rendering identical specs on every execution.

**Recursive conditions**
ConditionSpec supports nested Logic:"OR"/"AND" groups. Enables complex WHERE clauses. MaxConditionDepth prevents abuse (default 10 levels).

**Capitan events**
FactoryCreated, CapabilityAdded, CapabilityRemoved, CapabilityNotFound. Decouples observability.

## Internal Architecture

```
Application
    ↓
Factory[T] (capability registry, SQL cache)
    ↓
Dispatch (ExecQuery, ExecSelect, ExecUpdate, ExecDelete)
    ↓
Convert (spec → soy builder)
    ↓
Soy (query building, execution)
    ↓
ASTQL → sqlx
```

## Capability Types

| Type | Purpose | Returns |
|------|---------|---------|
| QueryCapability | Multi-record SELECT | `[]*T` |
| SelectCapability | Single-record SELECT | `*T` |
| UpdateCapability | UPDATE operations | `*T` |
| DeleteCapability | DELETE operations | `int64` |
| AggregateCapability | COUNT/SUM/AVG/MIN/MAX | `float64` |

Each capability has Name, Description, Spec, Params, Tags.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `factory.go` | Factory definition, registration, rendering, prepared statements |
| `dispatch.go` | Execution methods, transaction variants |
| `convert.go` | Spec → soy builder conversion, condition handling |
| `spec.go` | Spec types, FactorySpec generation, JSON export |
| `capability.go` | Capability type definitions |
| `events.go` | Capitan signals |

## LLM Integration

```go
json, _ := factory.SpecJSON()
// {"table":"users", "queries":[...], "params":[...]}
// LLM selects capability and provides params
// Execute: factory.ExecQuery(ctx, name, params)
```

Specs are pure data. No code. Safe to expose to LLMs.

## Current State / Direction

Stable. Core capability pattern complete.

Future considerations:
- Compound queries (UNION/INTERSECT/EXCEPT)
- Pre-compiled specs at factory creation

## Framework Context

**Dependencies**: soy (query building), astql (validation), capitan (events).

**Role**: Semantic layer over soy. Named capabilities for discoverability. JSON-serialisable specs for LLM integration. Introspectable database operations.
