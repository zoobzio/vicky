# Layer 5: Domain Services

Application-facing capabilities.

## Rule

Packages in this layer may depend on Layers 0-4.

## Packages

| Package | Role | Framework Deps |
|---------|------|----------------|
| [rocco](../../packages/rocco.md) | HTTP framework | sentinel, capitan, openapi |
| [zyn](../../packages/zyn.md) | LLM orchestration | capitan, pipz, sentinel |
| [tendo](../../packages/tendo.md) | Tensor computation | pipz, capitan, clockz |
| [ago](../../packages/ago.md) | Event coordination | capitan, pipz, herald (optional) |

## Package Details

### rocco

Type-safe HTTP framework. Struct tags as single source of truth.

**Key capabilities:**
- Generic `Handler[In, Out]` with compile-time checking
- Automatic OpenAPI generation from struct tags
- Sentinel integration for type metadata
- Pre-defined typed errors (`Error[D]`)
- Chi router foundation
- SSE streaming support

**Framework dependencies:**
- sentinel: Type metadata (Layer 0)
- capitan: Signal emission (Layer 0)
- openapi: Specification types (Layer 0)

**Why Layer 5:** Application-facing HTTP layer. Bridges framework to external clients.

### zyn

LLM orchestration. Type-safe interactions with composable reliability.

**Key capabilities:**
- Eight synapses: Binary, Classification, Extraction, Transform, Analyze, Convert, Ranking, Sentiment
- Generic response types with compile-time checking
- JSON schema generation from Go types (via sentinel)
- Session-based context management
- pipz integration for reliability (retry, timeout, circuit breaker)
- Three providers: OpenAI, Anthropic, Gemini

**Framework dependencies:**
- capitan: Observability (Layer 0)
- pipz: Reliability patterns (Layer 1)
- sentinel: Schema generation (Layer 0)

**Why Layer 5:** LLM interaction layer. Same pipz vocabulary applied to AI workloads.

### tendo

Composable tensor computation. Numerical operations with framework vocabulary.

**Key capabilities:**
- 80+ operations (element-wise, matrix, shape, reductions, activations, neural network)
- `Storage` interface for CPU and CUDA backends
- Operations return `pipz.Chainable[*Tensor]`
- Signal emission for every operation (autograd preparation)
- Memory pooling with LRU eviction
- Float16/BFloat16 support

**Framework dependencies:**
- pipz: Composition patterns (Layer 1)
- capitan: Signal emission (Layer 0)
- clockz: Timing (Layer 0)

**Why Layer 5:** Numerical computation layer. ML and scientific computing with framework patterns.

### ago

Event-driven coordination. Sagas, request/response, distributed patterns.

**Key capabilities:**
- `Flow[T]` envelope with correlation context
- Saga primitives: SagaStep, Compensate (LIFO rollback)
- Request/response: Request, Await
- Emit and DeadLetter for event routing
- Idempotency tracking by correlation + step name
- Herald integration for distributed DLQ

**Framework dependencies:**
- capitan: Events (Layer 0)
- pipz: Pipeline composition (Layer 1)
- herald: Optional DLQ (Layer 4)

**Why Layer 5:** Coordination layer. Distributed transactions and async patterns.

## Application Domains

Layer 5 provides capabilities for different application types:

```
┌─────────────────────────────────────────────────────────────┐
│                      Applications                            │
├───────────────┬───────────────┬───────────────┬─────────────┤
│   Web APIs    │  AI/ML Apps   │  Compute Jobs │ Event-Driven│
│     ↓         │      ↓        │      ↓        │      ↓      │
│   rocco       │     zyn       │    tendo      │     ago     │
│   (HTTP)      │    (LLM)      │  (tensors)    │   (sagas)   │
├───────────────┴───────────────┴───────────────┴─────────────┤
│                    Layer 4: Infrastructure                   │
└─────────────────────────────────────────────────────────────┘
```

| Application Type | Primary Package | Supporting Packages |
|------------------|-----------------|---------------------|
| REST APIs | rocco | soy (data), cereal (serialisation) |
| LLM applications | zyn | cogito (reasoning), soy (persistence) |
| ML pipelines | tendo | streamz (data flow), grub (storage) |
| Event systems | ago | herald (messaging), flux (config) |

## Common Patterns

All Layer 5 packages share:

| Pattern | Implementation |
|---------|----------------|
| Type-safe generics | `Handler[In,Out]`, `Service[T]`, `Flow[T]`, `Tensor` |
| pipz composition | Reliability patterns, operation chaining |
| capitan signals | Observability without instrumentation |
| Provider abstraction | rocco (Chi), zyn (LLM APIs), tendo (backends) |
