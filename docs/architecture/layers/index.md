# Layers

Abstraction hierarchy for the framework. Eight layers, strict dependency rules.

## Overview

| Layer | Name | Purpose | Package Count |
|-------|------|---------|---------------|
| 0 | [Primitives](0-primitives.md) | Zero framework dependencies | 7 |
| 1 | [Core](1-core.md) | Foundational patterns | 3 |
| 2 | [Validation](2-validation.md) | Transform and validate | 3 |
| 3 | [Data](3-data.md) | Storage abstraction | 3 |
| 4 | [Infrastructure](4-infrastructure.md) | Distributed concerns | 6 |
| 5 | [Services](5-services.md) | Application capabilities | 4 |
| 6 | [Intelligence](6-intelligence.md) | Reasoning and visualisation | 2 |
| 7 | [Endpoints](7-endpoints.md) | External consumers | 1 |

## Dependency Rules

```
Layer N may depend on Layers 0..(N-1)
Layer N must NOT depend on Layers (N+1)..7
```

Strict downward dependencies. No cycles. Each layer builds on those below.

## Visual Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│ Layer 7: Endpoints                                          │
│   periscope                                                 │
├─────────────────────────────────────────────────────────────┤
│ Layer 6: Intelligence                                       │
│   cogito, erd                                               │
├─────────────────────────────────────────────────────────────┤
│ Layer 5: Services                                           │
│   rocco, zyn, tendo, ago                                    │
├─────────────────────────────────────────────────────────────┤
│ Layer 4: Infrastructure                                     │
│   flux, flume, aperture, herald, aegis, sctx                │
├─────────────────────────────────────────────────────────────┤
│ Layer 3: Data                                               │
│   soy, grub, edamame                                        │
├─────────────────────────────────────────────────────────────┤
│ Layer 2: Validation                                         │
│   astql, vectql, cereal                                     │
├─────────────────────────────────────────────────────────────┤
│ Layer 1: Core                                               │
│   pipz, atom, streamz                                       │
├─────────────────────────────────────────────────────────────┤
│ Layer 0: Primitives                                         │
│   sentinel, capitan, clockz, dbml, ddml, vdml, openapi      │
└─────────────────────────────────────────────────────────────┘
```

## Reading Guide

Start at Layer 0. Each successive layer assumes familiarity with those below.
