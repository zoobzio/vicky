# Layer 6: Intelligence

Reasoning and visualisation. Highest abstraction level for computation.

## Rule

Packages in this layer may depend on all lower layers (0-5).

## Packages

| Package | Role | Framework Deps |
|---------|------|----------------|
| [cogito](../../packages/cogito.md) | LLM reasoning chains | capitan, pipz, zyn, soy |
| [erd](../../packages/erd.md) | ERD generation | sentinel |

## Package Details

### cogito

LLM-powered reasoning chains. Thought-Note-Memory architecture.

**Key capabilities:**
- Thought-Note-Memory trilogy for reasoning state
- Two-phase reasoning: deterministic decision + creative introspection
- 15+ primitives across categories:
  - Decision: Decide, Analyze, Assess, Categorize, Prioritize
  - Control: Sift, Discern (semantic gates/routers)
  - Memory: Recall, Reflect, Checkpoint, Seek, Survey
  - Session: Compress, Truncate (token management)
  - Synthesis: Amplify, Converge
  - Utility: Forget, Reset, Restore
- Provider hierarchy (step → context → global)
- Vector embeddings with pgvector support
- SoyMemory for PostgreSQL persistence

**Framework dependencies:**
- capitan: Signals (Layer 0)
- pipz: Orchestration (Layer 1)
- zyn: LLM synapses (Layer 5)
- soy: Persistence (Layer 3)

**Why Layer 6:** Highest abstraction for LLM reasoning. Builds on zyn synapses to create reasoning chains with persistent semantic memory.

### erd

Entity Relationship Diagrams from Go types.

**Key capabilities:**
- Direct sentinel integration for type metadata
- Two workflows: schema-driven or manual builder
- Automatic relationship extraction from type graph
- Two output formats: Mermaid and GraphViz DOT
- Struct tag parsing (`erd:"pk,fk,uk,note:..."`)

**Framework dependencies:**
- sentinel: Type metadata (Layer 0)

**Why Layer 6:** Visualisation layer. Transforms type intelligence into diagrams.

## Architectural Position

Layer 6 represents the highest abstraction levels:

```
┌─────────────────────────────────────────────────────────────┐
│                    Layer 6: Intelligence                     │
│  ┌──────────────────────────┐  ┌──────────────────────────┐ │
│  │         cogito           │  │           erd            │ │
│  │   (reasoning chains)     │  │    (visualisation)       │ │
│  └────────────┬─────────────┘  └────────────┬─────────────┘ │
│               │                              │               │
│               ▼                              │               │
│  ┌──────────────────────────┐               │               │
│  │       zyn (LLM)          │               │               │
│  │   Layer 5: Services      │               │               │
│  └────────────┬─────────────┘               │               │
│               │                              │               │
│               ▼                              ▼               │
│  ┌──────────────────────────┐  ┌──────────────────────────┐ │
│  │    soy (persistence)     │  │   sentinel (metadata)    │ │
│  │   Layer 3: Data Access   │  │   Layer 0: Primitives    │ │
│  └──────────────────────────┘  └──────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### cogito's Layer Stack

```
cogito (reasoning chains)
    ↓
zyn (LLM calls) ─────────────────┐
    ↓                            │
pipz (composition) ←─────────────┤
    ↓                            │
soy (persistence) ←──────────────┤
    ↓                            │
capitan (events) ←───────────────┘
```

All pipz patterns available for reasoning:
- Retry failed LLM calls
- Timeout long reasoning chains
- Circuit break unreliable providers
- Rate limit API calls

### erd's Simplicity

```
erd (diagrams)
    ↓
sentinel (type metadata)
```

erd demonstrates that not all high-level packages need deep dependency trees. It simply transforms existing type intelligence into visual form.

## Use Cases

| Package | Use Case |
|---------|----------|
| cogito | AI agents, RAG systems, reasoning pipelines, semantic search |
| erd | Documentation generation, schema visualisation, architecture diagrams |

## Persistence and Memory

cogito introduces the concept of persistent semantic memory:

| Component | Storage | Purpose |
|-----------|---------|---------|
| Note | In-memory (Thought) | Ephemeral reasoning context |
| Memory | SoyMemory (PostgreSQL) | Persistent semantic store |
| Vector | pgvector | Similarity search |

This enables reasoning chains that accumulate knowledge across sessions.
