# Architecture

Documentation for the Vicky framework's design and structure.

## Sections

| Section | Purpose | Start Here |
|---------|---------|------------|
| [Layers](layers/index.md) | Abstraction hierarchy and dependency rules | "Where does each package live?" |
| [Stories](stories/index.md) | Coherent narratives for solving problems | "How do I solve X?" |
| [Patterns](patterns/index.md) | Recurring design patterns and rationale | "Why is it designed this way?" |

## Quick Reference

**29 packages** organised into **8 layers**:

```
Layer 7: periscope (endpoints)
Layer 6: cogito, erd (intelligence)
Layer 5: rocco, zyn, tendo, ago (services)
Layer 4: flux, flume, aperture, herald, aegis, sctx (infrastructure)
Layer 3: soy, grub, edamame (data)
Layer 2: astql, vectql, cereal (validation)
Layer 1: pipz, atom, streamz (core)
Layer 0: sentinel, capitan, clockz, dbml, ddml, vdml, openapi (primitives)
```

## Reading Paths

**New to the framework?** Start with [Layers](layers/index.md) to understand the structure.

**Solving a specific problem?** Browse [Stories](stories/index.md) to find relevant package combinations.

**Understanding design decisions?** Read [Patterns](patterns/index.md) for rationale.

