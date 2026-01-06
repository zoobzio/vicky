# Vdml

Vector Database Markup Language for Go.

[GitHub](https://github.com/zoobzio/vdml)

## Vision

Programmatic schema definition for vector databases. Fluent builder API for defining collections, embeddings, metadata, and indexes. Generate VDML syntax. Schema-as-code for vector databases.

## Design Decisions

**Fluent builder**
`NewSchema().AddCollection().AddEmbedding()`. Chainable methods returning pointers. Readable, progressive construction.

**Type-safe enums**
`DistanceMetric`, `MetadataType`, `IndexType` as typed strings with constants. Validation ensures only valid values used.

**Pointer optionals**
`*string` for optional fields (Note, Name on Index). Distinguishes unset from empty.

**Separation of concerns**
types, builder, generator, serialize, validate in separate files. Clear responsibilities.

**Zero-dependency core**
Only yaml.v3 for YAML serialization. Pure Go, minimal attack surface.

## Data Model

```
Schema
└── Collections (map[string]*Collection)
    ├── Embeddings []*Embedding
    │   └── Name, Dimensions, Metric, Note
    ├── Metadata []*MetadataField
    │   └── Name, Type, Indexed, Required, Note
    ├── Indexes []*Index
    │   └── Type, Params (map), Name, Note
    └── Settings, Name, Note
```

## Constants

**Distance Metrics**: Cosine, Euclidean, DotProduct

**Metadata Types**: TypeString, TypeInt, TypeFloat, TypeBool, TypeStringArray, TypeIntArray, TypeFloatArray

**Index Types**: HNSW, IVFFlat, IVFPQ, Flat

## Code Organisation

| File | Responsibility |
|------|----------------|
| `types.go` | Type definitions and constants |
| `builder.go` | Constructors and fluent builder methods |
| `generator.go` | VDML syntax generation |
| `serialize.go` | JSON/YAML serialization |
| `validate.go` | Hierarchical validation |

## Current State / Direction

Stable. Core schema definition complete.

Future considerations:
- VDML parser (reverse generation)
- Schema composition

## Framework Context

**Dependencies**: yaml.v3 only.

**Role**: Schema language for vector databases. Parallel to dbml for relational databases. Will integrate with vectql for vector database operations.
