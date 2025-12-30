# Intelligence Story

*"How do I build LLM applications?"*

## The Flow

```
pipz → zyn → cogito
```

## The Packages

### pipz - Reliability Foundation

LLM calls are just pipelines. Same patterns apply.

```go
// Retry failed API calls
// Timeout long responses
// Circuit break unreliable providers
// Rate limit to avoid quotas
```

### zyn - Type-Safe LLM Calls

Eight synapses for common patterns. Provider-agnostic.

```go
// Extract structured data
type Person struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

synapse := zyn.NewExtraction[Person]()
result, err := synapse.Fire(ctx, session, "Extract person from: John Smith, john@example.com")
// result.Name = "John Smith", result.Email = "john@example.com"

// Classification
synapse := zyn.NewClassification[Category]()

// Binary decision
synapse := zyn.NewBinary()

// Sentiment analysis
synapse := zyn.NewSentiment()
```

**Synapses:**
| Synapse | Purpose | Default Temp |
|---------|---------|--------------|
| Binary | Yes/no decisions | 0.1 |
| Classification | Multi-class categorisation | 0.3 |
| Extraction | Structured data from text | 0.1 |
| Transform | Text transformation | 0.3 |
| Analyze | Analyze structured data | 0.2 |
| Convert | Schema conversion | 0.1 |
| Ranking | Order items by criteria | 0.2 |
| Sentiment | Emotional tone analysis | 0.2 |

**Providers:**
- OpenAI (Chat Completions)
- Anthropic (Messages API)
- Gemini (GenerateContent)

### cogito - Reasoning Chains

Thought-Note-Memory architecture. Semantic persistence.

```go
// Create thought with memory
memory := cogito.NewSoyMemory(db)
thought := cogito.NewThought(
    cogito.WithMemory(memory),
    cogito.WithEmbedder(embedder),
)

// Build reasoning chain
chain := pipz.NewSequence[*cogito.Thought](
    cogito.Analyze("initial_analysis", analyzeInput),
    cogito.Decide("route_decision", []string{"simple", "complex"}),
    cogito.NewDiscern("router",
        cogito.Route("simple", simpleHandler),
        cogito.Route("complex", complexHandler),
    ),
    cogito.Checkpoint("save_progress"),
)

result, err := chain.Process(ctx, thought)
```

**Primitives:**
| Category | Primitives |
|----------|------------|
| Decision | Decide, Analyze, Assess, Categorize, Prioritize |
| Control | Sift, Discern |
| Memory | Recall, Reflect, Checkpoint, Seek, Survey |
| Session | Compress, Truncate |
| Synthesis | Amplify, Converge |

## The Key Insight

**LLM calls are just pipz pipelines. Same reliability patterns apply.**

Retry failed API calls. Timeout long responses. Circuit break unreliable providers. Rate limit to avoid quotas. The vocabulary is the same.

```
┌─────────────────────────────────────────────────────────────┐
│                       cogito                                 │
│   ┌─────────────────────────────────────────────────────┐   │
│   │   Thought → Primitive → Primitive → ...             │   │
│   │              ↓            ↓                         │   │
│   │           Note         Note                         │   │
│   └─────────────────────────────────────────────────────┘   │
│                        ↓                                     │
│   ┌─────────────────────────────────────────────────────┐   │
│   │                      zyn                             │   │
│   │   Session → Synapse.Fire() → Structured Response    │   │
│   │                  ↓                                   │   │
│   │              Provider (OpenAI/Anthropic/Gemini)     │   │
│   └─────────────────────────────────────────────────────┘   │
│                        ↓                                     │
│   ┌─────────────────────────────────────────────────────┐   │
│   │                     pipz                             │   │
│   │   Retry → Timeout → CircuitBreaker → RateLimit      │   │
│   └─────────────────────────────────────────────────────┘   │
│                        ↓                                     │
│   ┌─────────────────────────────────────────────────────┐   │
│   │                     soy                              │   │
│   │   SoyMemory (PostgreSQL + pgvector)                 │   │
│   └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Two-Phase Reasoning

cogito separates deterministic and creative phases:

| Phase | Temperature | Purpose |
|-------|-------------|---------|
| Reasoning | Low (0.1) | Deterministic decisions, extraction |
| Introspection | Higher (0.3) | Creative semantic summaries |

Introspection disabled by default for cost control.

## Related Stories

- [Composition](composition.md) - pipz patterns enable LLM reliability
- [Data Access](data-access.md) - soy provides cogito persistence
