# Registry with Caching

*Extract once, cache forever.*

## The Problem

Reflection is expensive. Type metadata extraction via `reflect` has significant overhead. Repeating it on every operation kills performance.

## The Pattern

Extract metadata once at registration time. Cache forever. Zero reflection on hot path.

```go
// Registration (happens once)
atomizer := atom.Use[User]()  // Reflection here

// Usage (happens many times)
atom := atomizer.Atomize(&user)     // No reflection
user := atomizer.Deatomize(atom)    // No reflection
```

## Where It's Used

| Package | What's Cached | Registration |
|---------|---------------|--------------|
| sentinel | Type metadata, relationships | `Inspect[T]()`, `Scan[T]()` |
| atom | Field plans, type mappings | `Use[T]()` |
| soy | Table metadata, column info | `New[T]()` |
| flume | Processor identities | `Factory.Add()` |
| cereal | Field plans, sanitisation rules | `NewSerializer[T]()` |

## Implementation Details

### Global Registry

```go
var registry = struct {
    sync.RWMutex
    atomizers map[reflect.Type]*atomizer
}{
    atomizers: make(map[reflect.Type]*atomizer),
}

func Use[T any]() *Atomizer[T] {
    typ := reflect.TypeOf((*T)(nil)).Elem()

    registry.RLock()
    if a, ok := registry.atomizers[typ]; ok {
        registry.RUnlock()
        return (*Atomizer[T])(a)
    }
    registry.RUnlock()

    // Build atomizer with reflection
    a := buildAtomizer[T]()

    registry.Lock()
    registry.atomizers[typ] = a
    registry.Unlock()

    return a
}
```

### Permanent Caching

Go types are immutable at runtime. Once compiled, struct fields don't change. Cache forever is safe.

```go
// sentinel's permanent cache
type PermanentCache struct {
    mu   sync.RWMutex
    data map[reflect.Type]*Metadata
}

// Never evicts - types don't change at runtime
```

### Registration Patterns

**Explicit registration:**
```go
flume.Add("validate", validateProcessor)
flume.Add("transform", transformProcessor)
```

**Implicit registration (on first use):**
```go
meta := sentinel.Inspect[User]()  // Registers if not cached
```

**Init-time registration:**
```go
func init() {
    sentinel.Tag("custom")  // Register custom tag
}
```

## Why This Pattern

**Performance.** Reflection once, not per-operation. Critical for hot paths.

**Simplicity.** Users don't manage caches. Registration is automatic or simple.

**Consistency.** Same metadata everywhere. No drift between uses.

**Thread-safe.** RWMutex allows concurrent reads (hot path) with serialised writes (registration).

## Cache Invalidation

**Answer: Don't.**

Types don't change at runtime. Cache forever. If types change, restart the application (new compilation).

## Pattern Variations

### Lazy vs Eager

**Lazy (sentinel, atom):**
```go
// Cache on first use
meta := sentinel.Inspect[User]()
```

**Eager (flume):**
```go
// Cache at registration
factory.Add("processor", p)
```

### Identity Management (flume)

```go
// Same name = same identity (cached)
id1 := factory.Identity("validate", "Validates input")
id2 := factory.Identity("validate", "Validates input")
// id1 == id2 (same cached instance)

// Different description = panic
id3 := factory.Identity("validate", "Different description")
// PANIC: inconsistent identity
```

## Trade-offs

**Memory.** Caches grow. Not a problem for finite type sets.

**Startup cost.** Eager registration adds to startup time.

**Global state.** Singletons can complicate testing. Most packages provide test hooks.

## Example: sentinel Flow

```go
// First call: extract with reflection
meta := sentinel.Inspect[User]()

// Internal:
// 1. Check cache (RLock) - miss
// 2. Extract metadata (reflection)
// 3. Store in cache (Lock)
// 4. Return metadata

// Second call: cache hit
meta := sentinel.Inspect[User]()

// Internal:
// 1. Check cache (RLock) - hit
// 2. Return cached metadata
// No reflection
```

## Related Patterns

- [Schema-Driven Validation](schema-driven.md) - Schemas cached in registry
- [Type Safety](type-safety.md) - Generic registries
