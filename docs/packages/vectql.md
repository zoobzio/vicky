# Vectql

Query builder for vector databases.

[GitHub](https://github.com/zoobzio/vectql)

## Vision

Type-safe query construction for vector databases. AST-based architecture like astql. VDML schema validation optional. Provider-agnostic queries rendered to Pinecone, Qdrant, Milvus, Weaviate.

## Design Decisions

**AST-based architecture**
VectorAST represents queries as data. Renderer interface converts to provider-specific format. Same pattern as astql for SQL.

**Two modes**
Schemaless (function constructors like `C()`, `M()`, `P()`) or validated (VECTQL instance with VDML schema). Schema mode validates collections, embeddings, metadata at construction time.

**Fluent builder**
`Search()`, `Upsert()`, `Delete()`, `Fetch()`, `Update()`. Method chaining throughout. Returns `*Builder` for composition.

**Parameterized queries**
All values via `Param` type. No string interpolation. Injection patterns detected and rejected. Safe query construction.

**Rich filter system**
`FilterCondition` for simple comparisons. `FilterGroup` for AND/OR/NOT composition. `RangeFilter` for numeric ranges. `GeoFilter` for spatial queries.

**Complexity limits**
`MaxFilterDepth` (5), `MaxBatchSize` (100), `MaxTopK` (10000), `MaxIDsPerFetch` (1000). Prevents runaway queries. Validated at build time.

**Provider capabilities**
Renderer reports `SupportsOperation()`, `SupportsFilter()`, `SupportsMetric()`. Fail-fast on unsupported features rather than runtime errors.

## Internal Architecture

```
Application
    ↓
Builder (fluent API)
    ↓
VectorAST (validated)
    ↓
Renderer (provider-specific)
    ↓
QueryResult (provider format + params)
```

With VDML validation:
```
VDML Schema → VECTQL instance → validated C(), E(), M(), P()
```

## Operations

| Operation | Purpose |
|-----------|---------|
| SEARCH | Similarity search with vector |
| UPSERT | Insert or update vectors |
| DELETE | Remove vectors by ID or filter |
| FETCH | Retrieve vectors by ID |
| UPDATE | Modify metadata on existing vectors |

## Provider Renderers

| Provider | Package |
|----------|---------|
| Pinecone | `pkg/pinecone` |
| Qdrant | `pkg/qdrant` |
| Milvus | `pkg/milvus` |
| Weaviate | `pkg/weaviate` |

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Type re-exports and constants |
| `builder.go` | Fluent query builder |
| `expressions.go` | Filter expression constructors |
| `instance.go` | VECTQL with VDML schema validation |
| `renderer.go` | Renderer interface |
| `internal/types/` | AST, filter, param, result types |

## Current State / Direction

Stable. Core query building complete.

Future considerations:
- Additional provider renderers
- Hybrid search support

## Framework Context

**Dependencies**: vdml (optional schema validation).

**Role**: Vector database equivalent of astql. Type-safe query building with provider abstraction. Schema validation via vdml integration.
