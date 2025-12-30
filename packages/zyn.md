# Zyn

LLM orchestration for Go.

[GitHub](https://github.com/zoobzio/zyn)

## Vision

Type-safe LLM interactions with composable reliability. Define your response types, get JSON schema generation, validation, and structured outputs. Reliability patterns via pipz. Provider-agnostic.

## Design Decisions

**Generic synapses**
`Extract[T]`, `Analyze[T]`, `Convert[T,U]`. Response types validated at compile time. LLM outputs validated via `Validator` interface before returning.

**Structured prompts**
No string interpolation. Consistent ordering: Task → Input → Context → Examples → Schema → Constraints. Debuggable, auditable. Enables provider-side prompt caching.

**Service[T] abstraction**
Single point for parsing, validation, session updates, hook emission. Transactional semantics: session only updated after successful response.

**Pipz integration**
All reliability through pipz pipeline. Options compose naturally:
```go
WithRetry(3), WithTimeout(10s), WithCircuitBreaker(5, 30s), WithRateLimit(10)
```

**Session-based context**
Full message history sent with every call. Automatic prompt caching (providers cache 5min-1hr). Thread-safe via RWMutex.

**Provider interface**
Two methods: `Call()` and `Name()`. Minimal surface. OpenAI, Anthropic, Gemini included. Custom providers easy to add.

**Temperature defaults**
Each synapse has a default suited to its purpose. Deterministic tasks (0.1), analytical (0.2), creative (0.3). Overridable per-request.

## Synapse Types

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

## Internal Architecture

```
Application → Synapse.Fire(ctx, session, input)
                  ↓
              Service[T]
                  ↓
              pipz Pipeline (retry, timeout, circuit breaker, etc.)
                  ↓
              Provider.Call() → LLM API
                  ↓
              Parse JSON → Validate → Update Session
                  ↓
              Return typed response
```

## Providers

| Provider | Package | API |
|----------|---------|-----|
| OpenAI | `pkg/openai` | Chat Completions with JSON mode |
| Anthropic | `pkg/anthropic` | Messages API |
| Gemini | `pkg/gemini` | GenerateContent |

Each provider implements the two-method interface. HTTP done with standard library.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Provider, Validator, Message interfaces |
| `service.go` | Generic Service[T] execution |
| `schema.go` | JSON schema from Go types (uses sentinel) |
| `session.go` | Message history, token tracking |
| `prompt.go` | Structured prompt building |
| `options.go` | Reliability options wrapping pipz |
| `hooks.go` | Capitan signals and field keys |
| `[synapse].go` | Eight synapse implementations |
| `pkg/` | Provider implementations |

## Current State / Direction

Stable. Eight synapses covering common LLM interaction patterns.

Future considerations:
- Additional synapse types as patterns emerge
- Streaming support

## Framework Context

**Dependencies**: capitan (observability), pipz (reliability), sentinel (schema generation).

**Role**: LLM layer. Type-safe interactions with composable reliability. The pipz vocabulary applied to LLM calls. Extensible via providers and custom synapses.
