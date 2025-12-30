# Layered Architecture

*Each layer has single responsibility.*

## The Problem

Large frameworks become tangled. Packages depend on each other in unpredictable ways. Changes ripple unexpectedly. Testing becomes difficult.

## The Pattern

Organise packages into layers. Each layer has clear responsibility. Dependencies flow downward only.

```
Layer N may depend on Layers 0..(N-1)
Layer N must NOT depend on Layers (N+1)..7
```

## The Framework Layers

| Layer | Name | Responsibility |
|-------|------|----------------|
| 0 | Primitives | Zero framework deps |
| 1 | Core | Foundational patterns |
| 2 | Validation | Transform and validate |
| 3 | Data | Storage abstraction |
| 4 | Infrastructure | Distributed concerns |
| 5 | Services | Application capabilities |
| 6 | Intelligence | Reasoning and visualisation |
| 7 | Endpoints | External consumers |

## Example: SQL Path

```
soy (Layer 3 - ergonomic API)
    ↓
astql (Layer 2 - validation)
    ↓
dbml (Layer 0 - schema)
    ↓
sqlx (external - execution)
```

Each layer has single responsibility:
- dbml: Schema definition
- astql: Query validation
- soy: User-facing API

## Example: Storage Path

```
grub (Layer 3 - CRUD abstraction)
    ↓
atom (Layer 1 - type decomposition)
    ↓
sentinel (Layer 0 - type metadata)
```

## Why This Pattern

**Predictable dependencies.** Know what a package can depend on by its layer.

**Testable layers.** Lower layers tested in isolation. Higher layers mock lower.

**Clear boundaries.** Responsibility contained. Changes localised.

**Evolution.** Add new packages at appropriate layer without breaking existing code.

## Dependency Rules

**Allowed:**
```go
// soy (Layer 3) imports astql (Layer 2)
import "github.com/zoobzio/astql"

// astql (Layer 2) imports dbml (Layer 0)
import "github.com/zoobzio/dbml"
```

**Forbidden:**
```go
// dbml (Layer 0) importing soy (Layer 3)
import "github.com/zoobzio/soy"  // VIOLATION
```

## Internal Layer Dependencies

Packages within same layer can depend on each other (no cycles):

```
edamame (Layer 3) → soy (Layer 3)  ✓
soy (Layer 3) → edamame (Layer 3)  ✗ (would create cycle)
```

## Layer Characteristics

### Layer 0: Primitives
- Zero framework dependencies
- Only stdlib or minimal external (yaml.v3)
- Can be used independently

### Layers 1-2: Core Infrastructure
- Build on primitives
- Define patterns used everywhere
- Still relatively independent

### Layers 3-5: Application Services
- Build on everything below
- Provide user-facing capabilities
- Most feature-rich

### Layers 6-7: High-Level
- Maximum abstraction
- Consume from all layers
- Application-specific

## Trade-offs

**Rigid structure.** Sometimes a Layer 2 package wants a Layer 4 feature. Must restructure or accept limitation.

**Layer inflation.** Temptation to add layers. Keep minimal.

**Cross-cutting concerns.** Some features (logging, config) span all layers. Handled via patterns (signal emission, dependency injection).

## Visual Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│                    Layer 7: Endpoints                        │
├─────────────────────────────────────────────────────────────┤
│                   Layer 6: Intelligence                      │
├─────────────────────────────────────────────────────────────┤
│                    Layer 5: Services                         │
├─────────────────────────────────────────────────────────────┤
│                  Layer 4: Infrastructure                     │
├─────────────────────────────────────────────────────────────┤
│                      Layer 3: Data                           │
├─────────────────────────────────────────────────────────────┤
│                   Layer 2: Validation                        │
├─────────────────────────────────────────────────────────────┤
│                      Layer 1: Core                           │
├─────────────────────────────────────────────────────────────┤
│                   Layer 0: Primitives                        │
└─────────────────────────────────────────────────────────────┘
         ↑ Dependencies flow upward (allowed)
         ↓ Dependencies flow downward (forbidden)
```

## Related Patterns

- [Provider Abstraction](provider-abstraction.md) - Providers at specific layers
- [Schema-Driven Validation](schema-driven.md) - Schema languages at Layer 0
