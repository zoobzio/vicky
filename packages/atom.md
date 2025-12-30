# Atom

Type-safe intermediate representation for Go.

[GitHub](https://github.com/zoobzio/atom)

## Vision

Atom is the framework's internal data currency. Instead of passing structs through the system, we pass Atoms - type-segregated maps that enable field-level operations without reflection on the hot path.

The flow: Database → Atom → Cache → Atom → Business Logic → Struct. Structs exist only at user-facing boundaries. Internally, the framework operates on Atoms.

## Why Atom Exists

**Field-level operations without reflection**
Given a Spec (type metadata), code can access `atom.Strings["email"]` to encrypt, mask, or transform fields without knowing the struct type at compile time. No reflection needed once atomized.

**Type-agnostic caching**
Cache layers store Atoms. Applications fetch Atoms directly if they have a Spec. The struct type doesn't need to be known until deatomization at the boundary.

**Safe type-agnostic data flow**
Packages like grub produce Atoms from database providers. Packages like soy operate on Atoms for database access. The struct is reconstructed only when users need it.

## Design Decisions

**Type segregation**
Each primitive type in its own map. Strings in `Strings`, ints in `Ints`. No `any`. Type-safe access without assertions.

**Registry with caching**
`Use[T]()` builds atomizer once, returns cached instance. Field plans computed at registration. Zero reflection on hot path.

**Width-safe conversion**
int8/16/32/64 all stored as int64. Overflow checking on restore. Same for uint and float variants.

**Nested support**
Structs and struct slices handled recursively. Stored in `Nested` and `NestedSlices` maps. Circular references handled via shell registration.

**Escape hatches**
`Atomizable` and `Deatomizable` interfaces bypass reflection. For code-generated types that want zero-reflection performance.

**Flatten/Unflatten**
Convert Atom to `map[string]any` for external interop. Unflatten reconstructs Atom from external data using Spec for type mapping.

**Maps unsupported**
Explicit error on map fields. Only primitives and their composites (pointers, slices, structs).

## Type-Segregated Maps

| Category | Maps |
|----------|------|
| Scalars | Strings, Ints, Uints, Floats, Bools, Times, Bytes |
| Pointers | StringPtrs, IntPtrs, UintPtrs, FloatPtrs, BoolPtrs, TimePtrs, BytePtrs |
| Slices | StringSlices, IntSlices, UintSlices, FloatSlices, BoolSlices, TimeSlices, ByteSlices |
| Nested | Nested (map[string]Atom), NestedSlices (map[string][]Atom) |

Each struct field maps to exactly one map based on its type.

## Internal Architecture

```
Database
    ↓
grub scans directly to Atom
    ↓
Atom (type-segregated maps + Spec)
    ↓
Cache layer (stores Atoms)
    ↓
Field-level operations (encrypt, mask, transform)
    ↓
Deatomize at boundary → *T
```

Atomizer registration:
```
Use[T]() → Atomizer[T] (cached in registry)
    ├── Atomize(*T) → *Atom
    └── Deatomize(*Atom) → *T
```

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Atom struct, Atomizable/Deatomizable interfaces |
| `atomizer.go` | Atomizer[T] generic type |
| `registry.go` | Use[T]() with caching |
| `planning.go` | Field plan building, type-to-map mapping |
| `encoding.go` | atomize/deatomize field operations |
| `flatten.go` | Flatten/Unflatten for external interop |
| `resolver.go` | reflectAtomizer implementation |

## Current State / Direction

Stable. Core atomization complete.

Future considerations:
- Code generation for Atomizable/Deatomizable
- scio integration (direct Atom consumption)

## Framework Context

**Dependencies**: sentinel (type metadata).

**Consumers**: grub (produces Atoms from databases), soy (database access), scio (direct Atom consumption - in development).

**Role**: Internal data representation. Enables field-level operations and type-agnostic data flow. Structs are for users; Atoms are for the framework.
