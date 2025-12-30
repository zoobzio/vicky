# Layer 3: Data Access

Storage abstraction with validation.

## Rule

Packages in this layer may depend on Layers 0-2.

## Packages

| Package | Role | Framework Deps |
|---------|------|----------------|
| [soy](../../packages/soy.md) | Type-safe SQL access | astql, sentinel, dbml, capitan |
| [grub](../../packages/grub.md) | Provider-agnostic CRUD | atom, capitan, sentinel |
| [edamame](../../packages/edamame.md) | Capability-driven queries | soy, astql, capitan |

## Package Details

### soy

Type-safe database access. Schema-validated queries with compile-time safety.

**Key capabilities:**
- Generic `Soy[T]` for type-safe results
- Builders: Select, Query, Create, Update, Delete, Aggregate, Compound
- Struct tags as schema definition
- Required WHERE for UPDATE/DELETE (prevents accidents)
- Named parameters (`:param` syntax)
- Metadata cached via sentinel (zero reflection on hot path)

**Framework dependencies:**
- astql: SQL validation (Layer 2)
- sentinel: Type metadata (Layer 0)
- dbml: Schema definition (Layer 0)
- capitan: Signal emission (Layer 0)

**Why Layer 3:** Ergonomic API over astql's validation layer.

### grub

Provider-agnostic storage. Write once, store anywhere.

**Key capabilities:**
- Nine providers: Redis, S3, MongoDB, DynamoDB, BoltDB, BadgerDB, Firestore, GCS, Azure Blob
- Generic `Service[T]` for type safety
- Global singleton cache with LRU eviction
- Field accessors (GetStringField, GetIntField, etc.)
- Lifecycle as optional interface (stateless backends don't need it)

**Framework dependencies:**
- atom: Type-safe struct decomposition (Layer 1)
- capitan: Signal emission (Layer 0)
- sentinel: Type metadata (Layer 0)

**Why Layer 3:** Storage abstraction for document/object/key-value stores. Complements soy (SQL) with non-relational storage.

### edamame

Capability-driven query factory. Named, introspectable database operations.

**Key capabilities:**
- Pure data specs (JSON-serialisable, LLM-safe)
- Generic `Factory[T]` wrapping soy
- Default CRUD auto-registered
- SQL caching by "type:name" key
- Recursive conditions (nested AND/OR)
- Thread-safe registry (RWMutex)

**Framework dependencies:**
- soy: Query building (Layer 3)
- astql: Validation (Layer 2)
- capitan: Events (Layer 0)

**Why Layer 3:** Semantic layer over soy. Named capabilities for discoverability. Safe for LLM integration.

## Data Access Patterns

Three complementary approaches:

```
┌─────────────────────────────────────────────────────────────┐
│                      Application Code                        │
├─────────────────┬─────────────────┬─────────────────────────┤
│    edamame      │      soy        │         grub            │
│  (capabilities) │ (SQL queries)   │   (document/KV)         │
├─────────────────┼─────────────────┼─────────────────────────┤
│   soy.Soy[T]    │     astql       │     atom.Atom           │
├─────────────────┴─────────────────┼─────────────────────────┤
│          dbml (schema)            │   sentinel (metadata)   │
└───────────────────────────────────┴─────────────────────────┘
```

| Use Case | Package | Why |
|----------|---------|-----|
| Relational data | soy | Schema validation, complex queries |
| LLM-driven queries | edamame | Safe specs, named capabilities |
| Document storage | grub | Provider flexibility, simple CRUD |
| Object storage | grub | S3/GCS/Azure abstraction |
| Key-value storage | grub | Redis/BadgerDB/BoltDB |

## Internal Dependencies

Within Layer 3, edamame depends on soy (no cycle):

```
edamame → soy → astql → dbml
              ↘ sentinel
```

grub is independent:

```
grub → atom → sentinel
     ↘ capitan
```
