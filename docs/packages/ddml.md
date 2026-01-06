# Ddml

Document Database Markup Language for Go.

[GitHub](https://github.com/zoobzio/ddml)

## Vision

Programmatic schema definition for document databases. Fluent builder API for defining collections, nested documents, arrays, and indexes. Generate DDML syntax. Schema-as-code for document databases.

## Design Decisions

**Fluent builder**
`NewSchema().AddCollection().AddField()`. Chainable methods returning pointers. Readable, progressive construction.

**Recursive field structure**
Field contains `[]*Field` for nested documents. `ArrayOf *Field` for typed arrays. Natural mapping to document database hierarchies.

**Type-safe enums**
`FieldType`, `IndexType`, `IndexOrder` as typed strings with constants. Compile-time safety for valid values.

**Document-specific features**
TTL indexes, sparse indexes, geo indexes (2d/2dsphere), references between collections. Models MongoDB-style document databases.

**Pointer optionals**
`*string`, `*int` for optional fields (Note, TTL, Default). Distinguishes unset from empty/zero.

**Separation of concerns**
types, builder, generator, serialize, validate in separate files. Clear responsibilities.

**Zero-dependency core**
Only yaml.v3 for YAML serialization. Pure Go, minimal attack surface.

## Data Model

```
Schema
├── Collections (map[string]*Collection)
│   ├── Fields []*Field (recursive)
│   │   ├── Fields []*Field (nested documents)
│   │   ├── ArrayOf *Field (array element type)
│   │   ├── Ref *FieldRef (collection references)
│   │   └── Name, Type, Required, Unique, PrimaryKey, Default, Note
│   └── Indexes []*Index
│       └── Fields, Type, Unique, Sparse, TTL, Name
└── Enums (map[string]*Enum)
    └── Name, Values, Note
```

## Constants

**Field Types**: TypeString, TypeInt, TypeFloat, TypeBool, TypeDate, TypeObjectID, TypeObject, TypeArray, TypeEnum, TypeBinary, TypeGeoPoint

**Index Types**: IndexText, IndexGeo2D, IndexGeo2DSphere, IndexHashed

**Index Orders**: Ascending, Descending

## Code Organisation

| File | Responsibility |
|------|----------------|
| `types.go` | Type definitions and constants |
| `builder.go` | Constructors and fluent builder methods |
| `generator.go` | DDML syntax generation |
| `serialize.go` | JSON/YAML serialization |
| `validate.go` | Hierarchical validation |

## Current State / Direction

Stable. Core schema definition complete.

Future considerations:
- DDML parser (reverse generation)
- Schema composition

## Framework Context

**Dependencies**: yaml.v3 only.

**Role**: Schema language for document databases. Parallel to dbml for relational databases and vdml for vector databases. Completes the schema language trilogy.
