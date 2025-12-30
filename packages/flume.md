# Flume

Dynamic pipeline factory for Go.

[GitHub](https://github.com/zoobzio/flume)

## Vision

Schema-driven pipeline construction for pipz. Define pipelines in YAML/JSON instead of code. Hot-reload without restarts. Configuration as the interface between pipeline design and execution.

## Design Decisions

**Schema-driven construction**
YAML/JSON schemas define pipeline structure. Processors referenced by name. Connectors configured declaratively. No code changes for pipeline modifications.

**Hot-reloading via Bindings**
`Binding.Update()` swaps pipeline at runtime. Previous versions retained in history. `Rollback()` reverts instantly. Zero-downtime pipeline evolution.

**Component registries**
Factory holds processors, predicates, conditions, reducers, error handlers, channels. Schemas reference by name. Separation between component implementation and composition.

**Two-phase validation**
`ValidateSchema()` checks all references exist. `ValidateSchemaStructure()` validates syntax without factory. Enables CI/CD linting of schemas before deployment.

**Cycle detection**
Prevents circular processor references in schemas. Validated at build time. Mutually exclusive branches (filter then/else, switch routes) allow same processor in different branches.

**Managed identities**
`Factory.Identity()` caches by name. Consistent references across registrations. Panics if same name registered with different description.

**Channel integration**
`stream` nodes push to registered channels. Terminal endpoints for synchronous pipelines. Optional timeout on channel writes.

## Connector Types

| Category | Connectors |
|----------|------------|
| Composition | sequence, concurrent, race, fallback |
| Resilience | retry, timeout, circuit-breaker, rate-limit |
| Routing | filter (then/else), switch (routes/default), contest |
| Structure | handle, scaffold, worker-pool |

14 connector types mapping to pipz primitives.

## Registries

| Registry | Purpose | Schema Reference |
|----------|---------|------------------|
| processors | pipz.Chainable components | `ref: name` |
| predicates | Boolean functions | `predicate: name` |
| conditions | String-returning functions | `condition: name` |
| reducers | Concurrent merge functions | `reducer: name` |
| errorHandlers | Error pipelines | `error_handler: name` |
| channels | Output channels | `stream: name` |

## Internal Architecture

```
Factory[T]
    ↓
Register: Add(), AddPredicate(), AddCondition(), ...
    ↓
Schema (YAML/JSON)
    ↓
ValidateSchema() → Build() → pipz.Chainable[T]
    OR
Bind(identity, schema) → Binding[T]
    ↓
Binding.Update(newSchema) → hot-reload
Binding.Rollback() → revert
```

## Code Organisation

| File | Responsibility |
|------|----------------|
| `factory.go` | Factory[T], registration, Build() |
| `binding.go` | Binding[T], Update, Rollback, history |
| `schema.go` | Schema, Node type definitions |
| `builders.go` | Connector builders (14 types) |
| `validation.go` | ValidateSchema, ValidateSchemaStructure |
| `loader.go` | BuildFromFile, BuildFromYAML, BuildFromJSON |
| `spec.go` | Schema generation/export |
| `observability.go` | Capitan signals |

## Current State / Direction

Stable. Core schema-driven construction complete.

Future considerations:
- Schema composition/inheritance
- Visual schema editor integration

## Framework Context

**Dependencies**: pipz (pipeline primitives), capitan (events), yaml.v3.

**Role**: Configuration layer over pipz. Define pipelines declaratively. Hot-reload for operational flexibility. Bridge between pipeline design and runtime execution.
