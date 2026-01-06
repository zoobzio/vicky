# OpenAPI

OpenAPI 3.0 types for Go.

[GitHub](https://github.com/zoobzio/openapi)

## Vision

Complete Go implementation of OpenAPI 3.0 specification. Pure data structures with JSON/YAML serialization. Foundation for API documentation tools.

## Design Decisions

**Struct-based**
All OpenAPI constructs modelled as Go structs with JSON/YAML struct tags. Idiomatic Go. Leverages standard marshaling.

**Pointer optionals**
`*bool`, `*int`, `*string` distinguish "not set" from zero value. Follows OpenAPI conventions for optional fields.

**Map collections**
`map[string]Type` for named collections (Paths, Schemas, Components). Enables direct lookup. Aligns with OpenAPI structure.

**No validation**
Pure data structures. Validation delegated to callers. Keeps library focused and flexible.

**Minimal dependencies**
Only yaml.v3 for YAML support. No network, no filesystem, no code execution.

**No $ref resolution**
JSON Schema `$ref` references stored as strings. Resolution is caller responsibility. Prevents infinite loops, allows flexible strategies.

## Type Hierarchy

```
OpenAPI
├── Info, Contact, License
├── Servers []Server
│   └── ServerVariable
├── Paths map[string]PathItem
│   └── Operations (Get, Post, Put, Delete, Patch, Options, Head)
├── Components
│   ├── Schemas map[string]*Schema
│   ├── Responses map[string]*Response
│   ├── Parameters map[string]*Parameter
│   ├── SecuritySchemes map[string]*SecurityScheme
│   ├── RequestBodies map[string]*RequestBody
│   ├── Headers, Links, Examples, Callbacks
├── Security []SecurityRequirement
└── Tags []Tag
```

## Schema Type

Comprehensive JSON Schema support:
- Basic types: string, integer, number, boolean, array, object
- Composition: allOf, oneOf, anyOf, not
- Validation: min/max, length, patterns, enums
- Advanced: discriminators, XML metadata, read/write only

## Code Organisation

| File | Responsibility |
|------|----------------|
| `openapi.go` | All type definitions (~286 lines, 35+ structs) |
| `json.go` | ToJSON/FromJSON helpers |
| `yaml.go` | ToYAML/FromYAML helpers |

Single package. No subpackages.

## Current State / Direction

Stable. Complete OpenAPI 3.0 coverage.

Future considerations:
- OpenAPI 3.1 support

## Framework Context

**Dependencies**: yaml.v3 only.

**Role**: Type definitions for OpenAPI specifications. Used by rocco for API documentation generation from handler metadata.
