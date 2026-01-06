# Layer 0: Primitives

Zero framework dependencies. The bedrock.

## Rule

Packages in this layer MUST NOT depend on other framework packages. Only standard library or minimal external dependencies (yaml.v3).

## Packages

| Package | Role | External Deps |
|---------|------|---------------|
| [sentinel](../../packages/sentinel.md) | Type intelligence extraction | None (stdlib) |
| [capitan](../../packages/capitan.md) | Event coordination backbone | None (stdlib) |
| [clockz](../../packages/clockz.md) | Deterministic time abstraction | None (stdlib) |
| [dbml](../../packages/dbml.md) | Relational schema language | yaml.v3 |
| [ddml](../../packages/ddml.md) | Document schema language | yaml.v3 |
| [vdml](../../packages/vdml.md) | Vector schema language | yaml.v3 |
| [openapi](../../packages/openapi.md) | API specification types | yaml.v3 |

## Package Details

### sentinel

Type intelligence for Go. Extracts comprehensive metadata from structs, caches permanently, discovers relationships between types.

**Key capabilities:**
- `Inspect[T]()` - Single type, package-scoped relationships
- `Scan[T]()` - Recursive, module-scoped, cycle-safe
- Automatic tag extraction (json, db, validate, etc.)
- Permanent caching (types immutable at runtime)

**Why Layer 0:** Foundation for type-aware tooling. Used by atom, cereal, soy, erd.

### capitan

Type-safe event coordination. Unified event stream across the framework.

**Key capabilities:**
- Type-safe fields via generic `Key[T]`
- Per-signal worker goroutines with isolation
- Buffered channels with backpressure policies
- Panic recovery per listener
- Sync mode for deterministic testing

**Why Layer 0:** Event backbone. All packages emit through capitan.

### clockz

Deterministic time control. One interface, two implementations.

**Key capabilities:**
- `Clock` interface mirrors `time` package
- `RealClock` for production (zero overhead)
- `FakeClock` for tests (explicit time control)
- `Advance()` fires waiters deterministically
- Context integration (`WithTimeout`, `WithDeadline`)

**Why Layer 0:** Enables deterministic testing of time-dependent code. Used by pipz, streamz, flux.

### dbml

Database Markup Language. Programmatic schema definition for relational databases.

**Key capabilities:**
- Fluent builder API
- Tables, columns, indexes, references, enums
- Multi-schema support (PostgreSQL conventions)
- Validation before generation
- JSON/YAML serialisation

**Why Layer 0:** Schema language for SQL databases. Used by astql for validation.

### ddml

Document Database Markup Language. Schema definition for document databases.

**Key capabilities:**
- Fluent builder API
- Collections, nested documents, arrays
- TTL indexes, sparse indexes, geo indexes
- References between collections
- JSON/YAML serialisation

**Why Layer 0:** Schema language for document databases. Parallel to dbml.

### vdml

Vector Database Markup Language. Schema definition for vector databases.

**Key capabilities:**
- Collections with embeddings and metadata
- Distance metrics (Cosine, Euclidean, DotProduct)
- Index types (HNSW, IVFFlat, IVFPQ, Flat)
- Metadata field types and indexing
- JSON/YAML serialisation

**Why Layer 0:** Schema language for vector databases. Will integrate with vectql.

### openapi

OpenAPI 3.0 types. Complete specification as Go structs.

**Key capabilities:**
- All OpenAPI 3.0 constructs as structs
- Paths, operations, schemas, components
- JSON Schema support (composition, validation)
- JSON/YAML serialisation
- No $ref resolution (caller responsibility)

**Why Layer 0:** Type definitions for API specs. Used by rocco for documentation generation.

## Architectural Significance

Layer 0 packages are the foundation. They define:
- How types are understood (sentinel)
- How events flow (capitan)
- How time is controlled (clockz)
- How schemas are expressed (dbml, ddml, vdml, openapi)

Every layer above builds on these primitives.
