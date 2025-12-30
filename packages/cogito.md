# Cogito

LLM-powered reasoning chains for Go.

[GitHub](https://github.com/zoobzio/cogito)

## Vision

Composable primitives for reasoning pipelines with semantic memory. Thought-Note-Memory architecture. Reasoning steps accumulate context. Memory persists across sessions. Built on the pipz vocabulary and zyn synapses.

## Design Decisions

**Thought-Note-Memory trilogy**
- Thought: Rolling context with append-only Notes and LLM Session
- Note: Atomic unit of information with optional vector embeddings
- Memory: Persistent storage with semantic search

Separates ephemeral runtime state from persistent semantic state.

**Two-phase reasoning**
1. Reasoning phase: Deterministic (low temperature) decision/extraction
2. Introspection phase: Creative (higher temperature) semantic summary

Introspection disabled by default. Enables cost control while enriching context when needed.

**Provider hierarchy**
1. Explicit parameter (step-level)
2. Context value (`cogito.WithProvider(ctx, p)`)
3. Global fallback (`cogito.SetProvider(p)`)

Flexible configuration from simple defaults to fine-grained control.

**Pipz integration**
All primitives implement `pipz.Chainable[*Thought]`. Compose with Sequence, Filter, Switch, Concurrent. Full pipz vocabulary available.

**Zyn synapses**
Uses Binary, Classification, Extract, Transform, Sentiment for LLM operations. Provider-agnostic.

**Vector embeddings**
Custom Vector type with pgvector support. Semantic search across historical notes. SoyMemory provides PostgreSQL backend.

## Primitives

| Category | Primitives | Purpose |
|----------|------------|---------|
| Decision | Decide, Analyze, Assess, Categorize, Prioritize | Structured reasoning outputs |
| Control | Sift, Discern | Semantic gates and routers |
| Memory | Recall, Reflect, Checkpoint, Seek, Survey | Persistence and retrieval |
| Session | Compress, Truncate | Token management |
| Synthesis | Amplify, Converge | Iterative refinement, parallel synthesis |
| Utility | Forget, Reset, Restore | State management |

## Internal Architecture

```
Thought (mutable state)
├── Notes (append-only log)
├── Session (zyn.Session for LLM history)
├── Memory (persistence interface)
└── Embedder (vector generation)
          ↓
Primitive.Process(ctx, thought)
          ↓
zyn Synapse (LLM call)
          ↓
Note appended → Optional introspection → Updated Thought
```

## Code Organisation

| Category | Files |
|----------|-------|
| Core | `thought.go`, `memory.go`, `provider.go`, `embedder.go`, `vector.go`, `signals.go` |
| Storage | `soy.go` (PostgreSQL), `config.go` |
| Decision | `decide.go`, `analyze.go`, `assess.go`, `categorize.go`, `prioritize.go` |
| Control | `sift.go`, `discern.go` |
| Memory | `recall.go`, `reflect.go`, `checkpoint.go`, `seek.go`, `survey.go` |
| Session | `compress.go`, `truncate.go` |
| Synthesis | `amplify.go`, `converge.go` |
| Utility | `forget.go`, `reset.go`, `restore.go`, `helpers.go` |

## Current State / Direction

Stable. Core primitives complete. SoyMemory provides production-ready persistence.

Future considerations:
- Additional memory backends
- Distributed thought coordination

## Framework Context

**Dependencies**: capitan (signals), pipz (orchestration), zyn (LLM), soy (persistence).

**Role**: Reasoning layer. The pipz vocabulary applied to LLM reasoning chains. Zyn synapses for LLM calls. Semantic memory for persistent, searchable context across sessions.
