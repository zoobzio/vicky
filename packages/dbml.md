# Dbml

Database Markup Language for Go.

[GitHub](https://github.com/zoobzio/dbml)

## Vision

Programmatic construction of DBML schemas. Fluent API for building database schemas in Go. Serialize to DBML syntax, JSON, or YAML. Schema-as-code.

## Design Decisions

**Fluent builder**
Method chaining throughout. `NewProject().WithDatabaseType().AddTable()`. Readable, progressive construction.

**Pointer-based optionals**
Distinguishes "not set" from "set to empty". Clear intent. Only non-nil values included in output.

**Schema namespacing**
Tables and enums keyed by "schema.name" internally. Supports multi-schema databases. Default schema is "public" (PostgreSQL convention).

**Inline vs standalone refs**
- Inline: `column.WithRef()` for simple foreign keys
- Standalone: `project.AddRef()` for complex cases (composite keys, referential actions)

**Validation separated**
Validate before generate. Multi-layer validation (project, table, column, index, ref, enum). Early error detection.

**Minimal dependencies**
Only yaml.v3 for YAML serialization. Security through simplicity. No network, no filesystem.

## Data Model

```
Project
├── Tables (map[string]*Table)
│   ├── Columns []*Column
│   │   └── InlineRef *InlineRef
│   └── Indexes []*Index
├── Enums (map[string]*Enum)
├── Refs []*Ref
│   ├── Left *RefEndpoint
│   └── Right *RefEndpoint
└── TableGroups []*TableGroup
```

## Code Organisation

| File | Responsibility |
|------|----------------|
| `types.go` | Core type definitions |
| `builder.go` | Fluent constructors (43 builder methods) |
| `generator.go` | DBML syntax generation |
| `validate.go` | Schema validation |
| `serialize.go` | JSON/YAML serialization |

## Generation Pipeline

```
Validate() → error (stop if invalid)
    ↓
Generate() → string (DBML syntax)
    ↓
ToJSON() / ToYAML() → []byte (alternative formats)
```

## Current State / Direction

Stable. Generation complete.

Future considerations:
- DBML parsing (inverse operation)
- Schema diff/migration support

## Framework Context

**Dependencies**: yaml.v3 only.

**Role**: Schema language. Used by astql for SQL validation. Soy generates DBML from struct tags. Foundation for schema-driven tooling across the framework.
