# Type Intelligence Story

*"How do I work with types programmatically?"*

## The Flows

```
sentinel → atom → cereal
sentinel → erd
```

## The Packages

### sentinel - Type Metadata Extraction

Programs that *know* their data.

```go
// Extract metadata once
meta := sentinel.Inspect[User]()

// Access field information
for _, field := range meta.Fields {
    fmt.Printf("%s: %s (%v)\n",
        field.Name,
        field.Type,
        field.Tags,
    )
}

// Discover relationships
meta := sentinel.Scan[User]()
for _, rel := range meta.Relationships {
    fmt.Printf("%s → %s (%s)\n",
        rel.Source.Name,
        rel.Target.Name,
        rel.Kind,
    )
}
```

**Two Modes:**
| Mode | Scope | Use Case |
|------|-------|----------|
| `Inspect[T]()` | Single type, package-scoped | Quick metadata lookup |
| `Scan[T]()` | Recursive, module-scoped | Relationship discovery |

**Automatic Tags:**
- json, db, validate, scope, encrypt, redact, desc, example
- Extensible via `sentinel.Tag("custom")`

### atom - Type Decomposition

Structs to typed maps. No `map[string]any`.

```go
atomizer := atom.Use[User]()

// Decompose
atom := atomizer.Atomize(&user)
// atom.Strings["email"]
// atom.Ints["age"]
// atom.Times["created_at"]

// Reconstruct
user := atomizer.Deatomize(atom)
```

**21 Typed Tables:**
| Category | Tables |
|----------|--------|
| Scalars | strings, ints, uints, floats, bools, times, bytes |
| Pointers | string_ptrs, int_ptrs, uint_ptrs, float_ptrs, bool_ptrs, time_ptrs, byte_ptrs |
| Slices | string_slices, int_slices, uint_slices, float_slices, bool_slices, time_slices, byte_slices |

### cereal - Sanitised Marshaling

Serialise with protection.

```go
type User struct {
    ID       string `json:"id"`
    Email    string `json:"email" mask:"email"`
    Password string `json:"password" redact:"***"`
    SSN      string `json:"ssn" encrypt:"pii"`
}

serializer := cereal.NewSerializer[User](jsonCereal,
    cereal.WithEncryptor("pii", aesEncryptor),
)

// Marshal with sanitisation
data, _ := serializer.Marshal(&user)
// {"id":"123","email":"j***@example.com","password":"***","ssn":"<encrypted>"}

// Unmarshal with restoration (decrypt only)
serializer.Unmarshal(data, &user)
```

**Sanitisation Tags:**
| Tag | Effect | Reversible |
|-----|--------|------------|
| `encrypt:"context"` | AES-GCM/RSA encryption | Yes |
| `hash:"sha256"` | Hex-encoded hash | No |
| `redact:"***"` | Replace with literal | No |
| `mask:"email"` | Content-aware masking | No |

### erd - Diagram Generation

Visualise type relationships.

```go
// From sentinel metadata
diagram := erd.FromSchema[User]()

// To Mermaid
mermaid := diagram.ToMermaid()

// To GraphViz DOT
dot := diagram.ToDOT()
```

**Sentinel → ERD Mapping:**
| Sentinel Relationship | ERD Cardinality |
|-----------------------|-----------------|
| Reference (pointer) | OneToOne |
| Collection (slice) | OneToMany |
| Embedding | OneToOne |
| Map | ManyToMany |

## The Key Insight

**Types are data. Inspect, transform, visualise.**

sentinel extracts. atom transforms. cereal protects. erd visualises. All from the same type definitions.

```
┌─────────────────────────────────────────────────────────────┐
│                    type User struct {                        │
│                        ID    string                          │
│                        Email string `mask:"email"`           │
│                        ...                                   │
│                    }                                         │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   sentinel   │
                    │  (extract)   │
                    └──────┬───────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
  │     atom     │  │    cereal    │  │     erd      │
  │ (typed maps) │  │  (sanitise)  │  │ (visualise)  │
  └──────────────┘  └──────────────┘  └──────────────┘
         │                 │                 │
         ▼                 ▼                 ▼
     Storage           JSON/XML          Diagrams
    Backends          (protected)       (Mermaid/DOT)
```

## Use Cases

| Package | Use Case |
|---------|----------|
| sentinel | Schema generation, validation, documentation |
| atom | Document databases, key-value stores |
| cereal | Secure transport, logging, APIs |
| erd | Architecture docs, schema visualisation |

## Related Stories

- [Data Access](data-access.md) - atom enables grub storage
- [HTTP API](http-api.md) - sentinel enables OpenAPI generation
