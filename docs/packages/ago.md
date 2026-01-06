# Ago

Event-driven coordination patterns for Go.

[GitHub](https://github.com/zoobzio/ago)

## Vision

Bridge capitan events with pipz pipelines. Distributed sagas, request/response, and stateful coordination. Flow wraps payload with correlation context. Primitives implement pipz.Chainable. Composable event-driven patterns.

## Design Decisions

**Flow as envelope**
`Flow[T]` wraps typed payload with correlation ID, causation ID, metadata, accumulated state, and error list. All primitives operate on Flow. Implements Cloner for parallel processing.

**Correlation-based routing**
CorrelationID links related events across steps. Request/Await filter by correlation. Sagas track steps by correlation.

**Compensation stack**
Each SagaStep serializes payload for potential rollback. Compensate pops in LIFO order. Stored in SagaState for restart recovery.

**Idempotency built in**
Keys generated as `correlationID:stepName`. Store tracks compensated steps. Safe retry on crashes.

**Exclusive access via WithSaga**
Store callback ensures only one caller modifies saga state. Signal emission happens after lock release.

**Timeout support**
Sagas can have timeout. Request/Await have configurable timeout. `RecoverSagas` compensates expired sagas at startup.

**Herald integration**
DeadLetter can publish to external broker via herald.Provider. Local capitan signal always emitted.

## Primitives

| Primitive | Purpose |
|-----------|---------|
| `SagaStep` | Execute step, register compensation, emit execute signal |
| `Compensate` | Run compensation stack in reverse order |
| `Emit` | Emit capitan signal with flow context |
| `Request` | Send request, wait for correlated response |
| `Await` | Wait for correlated event on signal |
| `DeadLetter` | Route failures to DLQ |

All primitives implement `pipz.Chainable[*Flow[T]]`. Composable via pipz topology or flume schema.

## Saga Lifecycle

```
Pending → Running (steps executing)
              ↓
    Each step pushes CompensationRecord
              ↓
    On failure → Compensating (LIFO)
              ↓
    Completed or Failed
```

## Internal Architecture

```
Capitan Event
    ↓
NewFromEvent[T](event, key) → *Flow[T]
    ↓
Pipeline (SagaStep → SagaStep → ...)
    ↓
On error: Compensate (reverse order)
    ↓
Emit / DeadLetter
```

Recovery at startup:
```
RecoverSagas[T]() → ListIncompleteSagas → Compensate each
```

## Code Organisation

| File | Responsibility |
|------|----------------|
| `flow.go` | Flow[T] type, correlation, accumulated state |
| `saga_step.go` | SagaStep primitive |
| `compensate.go` | Compensate primitive |
| `emit.go` | Emit primitive |
| `request.go` | Request/response primitive |
| `await.go` | Await primitive |
| `dead_letter.go` | DeadLetter primitive with herald integration |
| `store.go` | Store interface, SagaState, PendingState |
| `recovery.go` | RecoverSagas startup recovery |
| `signals.go` | Capitan signals and keys |

## Current State / Direction

Stable. Core saga and request/response patterns complete.

Future considerations:
- Additional store implementations
- Saga versioning

## Framework Context

**Dependencies**: capitan (events), pipz (pipeline composition), herald (optional DLQ).

**Role**: Event-driven coordination layer. Sagas for distributed transactions. Request/response for synchronous-over-async. Bridges local events to distributed patterns.
