# Erd

Entity Relationship Diagrams from Go types.

[GitHub](https://github.com/zoobzio/erd)

## Vision

Auto-generate ERDs from domain models. Sentinel extracts type metadata, erd renders diagrams. Documentation that stays synchronised with code. No manual diagram maintenance.

## Design Decisions

**Direct sentinel integration**
Uses sentinel for type reflection. Registers "erd" tag at init. Extracts relationships automatically from sentinel's type graph.

**Two workflows**
- Schema-driven: `FromSchema()` generates diagram from sentinel metadata
- Manual: Builder API for fine-grained control

**Fluent builder**
`NewDiagram()`, `NewEntity()`, `NewAttribute()`, `NewRelationship()`. Method chaining throughout. Returns pointers for composition.

**Multiple outputs**
- Mermaid: Lightweight, web-friendly `erDiagram` syntax
- GraphViz DOT: Publication-quality `digraph` with record nodes

**Filtered relationships**
Relationship fields excluded from entity attributes. Represented as diagram edges instead. Keeps visual model clean.

**Validation separated**
Diagrams can be incomplete during building. Must pass validation before output. Enables fail-fast error reporting.

## Sentinel → ERD Mapping

| Sentinel Relationship Kind | ERD Cardinality |
|----------------------------|-----------------|
| Reference (pointer) | OneToOne |
| Collection (slice) | OneToMany |
| Embedding | OneToOne |
| Map | ManyToMany |

Struct tag `erd:"pk,fk,uk,note:..."` parsed for key types and notes.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `types.go` | Domain types: Diagram, Entity, Attribute, Relationship |
| `builder.go` | Fluent constructors |
| `sentinel.go` | Sentinel → ERD conversion, tag parsing |
| `validate.go` | Structural validation with detailed errors |
| `mermaid.go` | Mermaid erDiagram output |
| `dot.go` | GraphViz DOT output |

## Current State / Direction

Stable. Core diagram generation complete.

Future considerations:
- Additional output formats
- Layout hints

## Framework Context

**Dependencies**: sentinel.

**Role**: Visualisation layer. Sentinel's type intelligence rendered as diagrams. ERDs that stay synchronised with code.
