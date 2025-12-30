# Schema-Driven Validation

*Schema as allowlist. Validate at construction.*

## The Problem

User input is dangerous. SQL injection, invalid data, malformed requests. Runtime validation is error-prone and often forgotten.

## The Pattern

Define schema upfront. Validate every operation against schema at construction time. Invalid data never enters the system.

```go
// Schema defines what's allowed
schema := dbml.NewProject("myapp").
    AddTable("users").
        AddColumn("id").WithType("uuid").Done().
        AddColumn("email").WithType("varchar(255)").Done().
    Done()

// ASTQL enforces schema
instance := astql.New(schema)

// This works - "users" and "email" are in schema
query := instance.Query().
    Select(instance.F("email")).
    From(instance.T("users"))

// This panics - "passwords" not in schema
query := instance.Query().
    Select(instance.F("passwords")) // PANIC
```

## Where It's Used

| Schema | Validator | Consumer |
|--------|-----------|----------|
| dbml | astql | soy |
| vdml | vectql | (future vector storage) |
| openapi | rocco | API clients |
| flume schema | Factory | pipz pipelines |

## Implementation Details

### Validation Layers

astql implements 7 validation layers:

1. **Schema validation** - Table/field allowlist via DBML
2. **Identifier validation** - Alphanumeric + underscore only
3. **SQL keyword blocking** - Reject keywords as identifiers
4. **Alias restrictions** - Single lowercase letter
5. **Quoted identifiers** - Dialect-specific escaping
6. **Parameterised queries** - Values never interpolated
7. **Subquery depth limits** - Max 3 levels

### Fail-Fast vs Try Variants

Two API styles:

```go
// Fail-fast (panics) - for validated schemas
table := instance.T("users")
field := instance.F("email")

// Try variants (returns error) - for runtime input
table, err := instance.TryT(userInput)
field, err := instance.TryF(fieldInput)
```

Use fail-fast when schema is known at compile time. Use Try variants for user-provided input.

### Two-Phase Validation (flume)

```go
// Phase 1: Syntax validation (no factory needed)
err := flume.ValidateSchemaStructure(schema)

// Phase 2: Reference validation (factory required)
err := factory.ValidateSchema(schema)
```

Enables CI/CD linting before deployment.

## Why This Pattern

**Security through defense in depth.** Multiple validation layers. Even if one fails, others catch it.

**Early errors.** Invalid data rejected at construction, not execution. Clear stack traces.

**Auditable.** Schema documents allowed operations. Easy to review.

**Injection-resistant.** Identifiers validated. Values parameterised. No string interpolation.

## Trade-offs

**Schema maintenance.** Schema must stay in sync with reality. Migration tooling helps.

**Flexibility cost.** Dynamic queries harder. Must use Try variants with error handling.

**Performance.** Validation has cost. Registry caching mitigates for repeated operations.

## Security Model

```
┌─────────────────────────────────────────────────────────────┐
│                     User Input                               │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │   Schema Allowlist      │  ← Only schema items pass
              └────────────┬────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │   Identifier Rules      │  ← No special characters
              └────────────┬────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │   Parameterisation      │  ← Values never interpolated
              └────────────┬────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │     Safe Operation      │
              └─────────────────────────┘
```

## Related Patterns

- [Layered Architecture](layered-architecture.md) - Schema languages at Layer 0
- [Registry Caching](registry-caching.md) - Schema metadata cached
