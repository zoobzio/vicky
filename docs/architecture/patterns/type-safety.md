# Type Safety via Generics

*`Service[T]` not `interface{}`.*

## The Problem

Go's `interface{}` (now `any`) loses type information. Runtime type assertions are error-prone. Compile-time checking is preferable.

## The Pattern

Use generics to preserve types through the entire chain. No type assertions on hot paths.

```go
// Generic service
type Service[T any] struct {
    provider Provider
}

// Type flows through
func (s *Service[T]) Get(ctx context.Context, key string) (*T, error) {
    data, err := s.provider.Get(ctx, key)
    if err != nil {
        return nil, err
    }
    return deserialize[T](data)
}

// Usage - type safety preserved
users := grub.NewService[User](provider)
user, err := users.Get(ctx, "user-123")
// user is *User, not interface{}
```

## Where It's Used

| Package | Generic Type | Purpose |
|---------|--------------|---------|
| soy | `Soy[T]` | Database records |
| grub | `Service[T]` | Storage records |
| herald | `Publisher[T]`, `Subscriber[T]` | Messages |
| zyn | Synapses | LLM responses |
| streamz | `Result[T]` | Stream items |
| pipz | `Chainable[T]` | Pipeline values |
| cogito | `Flow[T]` | Reasoning state |
| edamame | `Factory[T]` | Query results |

## Implementation Details

### Compile-Time Checking

```go
// This won't compile
users := soy.New[User](db)
orders, err := users.Query().Exec(ctx) // Type mismatch: expects User, got Order

// Caught at compile time, not runtime
```

### Constraint Interfaces

```go
// Require specific capabilities
type Validator interface {
    Validate() error
}

type Service[T Validator] struct { ... }

// T must implement Validate()
```

### Generic Chains

```go
// Pipelines preserve type through all operations
pipeline := pipz.NewSequence[Order](
    pipz.Apply(validateOrder),    // Order → Order
    pipz.Apply(enrichOrder),      // Order → Order
    pipz.Apply(persistOrder),     // Order → Order
)

result, err := pipeline.Process(ctx, order)
// result is Order, not interface{}
```

## Why This Pattern

**Compile-time safety.** Type errors caught by compiler, not at runtime.

**IDE support.** Autocomplete works. Go to definition works. Refactoring works.

**No runtime assertions.** No `value.(Type)` that might panic.

**Self-documenting.** `Service[User]` is clearer than `Service` with user passed as `interface{}`.

## Pattern Variations

### Result Types (streamz)

```go
type Result[T any] struct {
    Value T
    Err   error
}

// Single channel carries success and error
func Filter[T any](in <-chan Result[T], pred func(T) bool) <-chan Result[T]
```

### Constraint Combinations

```go
type Cloner[T any] interface {
    Clone() T
}

// T must be cloneable
type Serializer[T Cloner[T]] struct { ... }
```

## Trade-offs

**Instantiation cost.** Each `Service[User]` is a different type. Code bloat possible (though Go handles this well).

**Learning curve.** Generic patterns take time to master.

**Constraint limitations.** Some patterns hard to express with Go's type system.

## Before/After

**Before (interface{}):**
```go
func (s *Service) Get(key string) (interface{}, error) {
    data, _ := s.provider.Get(key)
    return data, nil
}

// Usage
result, _ := service.Get("key")
user := result.(*User) // Runtime panic if wrong type
```

**After (generics):**
```go
func (s *Service[T]) Get(key string) (*T, error) {
    data, _ := s.provider.Get(key)
    return deserialize[T](data)
}

// Usage
user, _ := service.Get("key")
// user is *User, compile-time guaranteed
```

## Related Patterns

- [Chainable Composition](chainable-composition.md) - Generics preserve type through chains
- [Registry Caching](registry-caching.md) - Generic registries
