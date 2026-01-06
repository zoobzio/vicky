# Stories

Coherent narratives for solving problems. Each story answers "How do I...?"

## Overview

| Story | Question | Primary Flow |
|-------|----------|--------------|
| [Observability](observability.md) | How do I observe my application? | capitan → herald → aperture |
| [Data Access](data-access.md) | How do I access databases? | dbml → astql → soy → edamame |
| [Composition](composition.md) | How do I build reliable pipelines? | pipz → flume / streamz |
| [Intelligence](intelligence.md) | How do I build LLM applications? | pipz → zyn → cogito |
| [Configuration](configuration.md) | How do I manage configuration? | flux + flume |
| [Distributed](distributed.md) | How do I coordinate across services? | capitan → herald → ago / aegis |
| [HTTP API](http-api.md) | How do I build HTTP services? | sentinel → openapi → rocco |
| [Type Intelligence](type-intelligence.md) | How do I work with types? | sentinel → atom → cereal / erd |
| [Numerical](numerical.md) | How do I do ML/scientific computing? | pipz + capitan → tendo |

## Reading Guide

Stories are solution-oriented. Start with the problem you're solving.

Each story shows:
- The package flow (which packages, in what order)
- The key insight (what makes this approach work)
- How packages connect together

## Cross-Cutting Themes

Several themes appear across multiple stories:

| Theme | Stories | Enabling Package |
|-------|---------|------------------|
| Signals | Observability, Distributed, all others | capitan |
| Composition | Composition, Intelligence, Numerical | pipz |
| Type safety | Data Access, HTTP API, Type Intelligence | sentinel |
| Provider abstraction | Data Access, Observability, Intelligence | (pattern) |

## Story Dependencies

```
                    ┌─────────────────┐
                    │  Observability  │ (foundation - used by all)
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│   Data Access   │ │   Composition   │ │ Type Intelligence│
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         │         ┌─────────┴─────────┐         │
         │         │                   │         │
         ▼         ▼                   ▼         ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│   Intelligence  │ │  Configuration  │ │    HTTP API     │
└─────────────────┘ └────────┬────────┘ └─────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   Distributed   │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │    Numerical    │ (uses all patterns)
                    └─────────────────┘
```
