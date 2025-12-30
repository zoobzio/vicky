# Sentinel

Type intelligence for Go.

[GitHub](https://github.com/zoobzio/sentinel)

## Vision

Programs that *know* their data. Sentinel extracts comprehensive metadata from Go structs, caches it permanently, and discovers relationships between types. The foundation for type-aware tooling.

## Design Decisions

**Zero dependencies**
Foundation layer. Can't have external risk. Standard library only.

**Permanent caching**
Go types are immutable at runtime. Extract once, cache forever. Global singleton suffices.

**Module-aware scanning**
`Scan()` stays within module boundaries. Keeps focus on domain types, ignores external noise.

**Two modes**
- `Inspect[T]()` - single type, package-scoped relationships
- `Scan[T]()` - recursive, module-scoped, cycle-safe

Different tools need different granularity.

**Type normalisation**
`*User` and `User` resolve to the same metadata. Pointer indirection handled transparently.

**Automatic tag extraction**
Common tags pre-registered: `json`, `db`, `validate`, `scope`, `encrypt`, `redact`, `desc`, `example`. Extensible via `Tag()`.

## Internal Architecture

```
API Layer (api.go)
    Inspect / Scan / Tag / Browse / Lookup
              ↓
Extraction Layer (extraction.go)
    extractMetadata / extractFieldMetadata
              ↓
Relationship Layer (relationship.go)
    extractRelationships / module boundary checking
              ↓
Cache Layer (cache.go)
    PermanentCache with RWMutex protection
```

Global singleton holds cache, tag registry, and module path. Thread-safe via RWMutex.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Public surface, singleton management |
| `metadata.go` | Type definitions: Metadata, FieldMetadata, TypeRelationship, FieldKind |
| `extraction.go` | Field and metadata extraction logic |
| `relationship.go` | Type graph discovery, module boundary detection |
| `cache.go` | Cache interface and implementations |

## Current State / Direction

Stable. Core extraction and caching complete. Relationship discovery working.

Future considerations:
- Additional relationship kinds as patterns emerge
- Performance optimisation if profiling warrants

## Framework Context

**Dependencies**: None.

**Role**: Foundation for type-aware tooling. Enables packages to write code that understands data structures at a deep level - schema generation, documentation, validation, serialisation.
