# Patterns

Recurring design patterns and their rationale. Each pattern answers "Why is it designed this way?"

## Overview

| Pattern | Purpose | Primary Packages |
|---------|---------|------------------|
| [Provider Abstraction](provider-abstraction.md) | Single interface, multiple backends | herald, grub, flux, astql, vectql, zyn, tendo |
| [Schema-Driven Validation](schema-driven.md) | Schema as allowlist | dbml, vdml, astql, vectql, flume, rocco |
| [Chainable Composition](chainable-composition.md) | Uniform composition interface | pipz, tendo, cogito, zyn, flux, flume, herald, ago |
| [Signal Emission](signal-emission.md) | Observability without instrumentation | capitan, all packages |
| [Type Safety](type-safety.md) | Generics over interface{} | soy, grub, herald, zyn, streamz, pipz |
| [Clock Abstraction](clock-abstraction.md) | Deterministic time testing | clockz, streamz, pipz, flux |
| [Registry Caching](registry-caching.md) | Extract once, cache forever | sentinel, atom, soy, flume |
| [Layered Architecture](layered-architecture.md) | Single responsibility per layer | dbml→astql→soy, sentinel→atom→grub |

## Pattern Selection Guide

| Problem | Pattern | Example |
|---------|---------|---------|
| Need multiple backends | Provider Abstraction | grub with 9 storage providers |
| Need to prevent bad data | Schema-Driven Validation | astql with DBML allowlist |
| Need composable operations | Chainable Composition | pipz with Retry, Timeout, etc. |
| Need observability | Signal Emission | capitan throughout framework |
| Need compile-time safety | Type Safety via Generics | `Service[T]` not `interface{}` |
| Need deterministic tests | Clock Abstraction | clockz FakeClock |
| Need zero-reflection hot path | Registry Caching | sentinel metadata caching |
| Need separation of concerns | Layered Architecture | soy → astql → dbml |

## Pattern Relationships

```
┌─────────────────────────────────────────────────────────────┐
│                  Signal Emission                             │
│  (observability foundation - used by all other patterns)    │
└──────────────────────────┬──────────────────────────────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│    Provider     │ │    Chainable    │ │   Type Safety   │
│   Abstraction   │ │   Composition   │ │  via Generics   │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         │                   ▼                   │
         │         ┌─────────────────┐           │
         │         │     Schema      │           │
         │         │    Validation   │           │
         │         └────────┬────────┘           │
         │                  │                    │
         └──────────────────┼────────────────────┘
                            │
              ┌─────────────┴─────────────┐
              │                           │
              ▼                           ▼
     ┌─────────────────┐       ┌─────────────────┐
     │    Registry     │       │    Layered      │
     │    Caching      │       │  Architecture   │
     └─────────────────┘       └─────────────────┘
              │                           │
              └─────────────┬─────────────┘
                            │
                            ▼
                   ┌─────────────────┐
                   │     Clock       │
                   │   Abstraction   │
                   └─────────────────┘
```

## Reading Guide

Patterns are design-oriented. Read to understand why the framework is structured as it is.

Each pattern explains:
- The problem it solves
- How it works
- Where it's used in the framework
- Trade-offs and alternatives
