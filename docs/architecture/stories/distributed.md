# Distributed Coordination Story

*"How do I coordinate across services?"*

## The Flow

```
capitan → herald → ago / aegis
```

## The Packages

### capitan - Local Events

Foundation. Every package emits here.

```go
// Emit locally
capitan.Emit(ctx, OrderPlaced,
    OrderIDKey.Field(order.ID),
    AmountKey.Field(order.Amount),
)

// Hook locally
capitan.Hook(ctx, OrderPlaced, func(e *capitan.Event) {
    // Handle locally
})
```

### herald - Distributed Messaging

Extend events across process boundaries.

```go
// Publish local events to Kafka
publisher := herald.NewPublisher[OrderEvent](kafkaProvider,
    herald.WithRetry(3),
    herald.WithCircuitBreaker(5, 30*time.Second),
)

capitan.Hook(ctx, OrderPlaced, func(e *capitan.Event) {
    publisher.Publish(ctx, toOrderEvent(e))
})

// Subscribe external events
subscriber := herald.NewSubscriber[PaymentEvent](kafkaProvider)
subscriber.Subscribe(ctx, func(msg herald.Message[PaymentEvent]) {
    capitan.Emit(ctx, PaymentReceived, ...)
})
```

**11 Providers:**
| Provider | Use Case |
|----------|----------|
| Kafka | High throughput, ordering |
| NATS | Low latency |
| JetStream | Persistent NATS |
| Pub/Sub | GCP managed |
| Redis Streams | In-memory with persistence |
| SQS | AWS managed |
| AMQP | Traditional broker |
| SNS | Fanout |
| BoltDB | Embedded |
| Firestore | GCP document |
| io | Testing |

### ago - Distributed Sagas

Compensating transactions. Request/response over events.

```go
// Define saga steps
saga := pipz.NewSequence[*ago.Flow[Order]](
    ago.SagaStep("reserve_inventory",
        reserveInventory,
        cancelReservation, // compensation
    ),
    ago.SagaStep("charge_payment",
        chargePayment,
        refundPayment, // compensation
    ),
    ago.SagaStep("ship_order",
        shipOrder,
        cancelShipment, // compensation
    ),
)

// On failure, compensate in reverse order
flow := ago.NewFlow(order)
result, err := saga.Process(ctx, flow)
if err != nil {
    // Compensations already executed (LIFO)
}
```

**Saga Lifecycle:**
```
Pending → Running → Completed
                 ↘ Compensating → Failed
```

**Request/Response:**
```go
// Send request, await correlated response
response := ago.Request[Resp](ctx, requestSignal, responseSignal,
    ago.WithTimeout(5*time.Second),
)
```

### aegis - P2P Mesh

Peer-to-peer coordination. No central server.

```go
node := aegis.NewNode(config)
node.JoinMesh(ctx, entryPoint)

// Remote function execution
result, err := node.ExecuteOn(ctx, peerID, "process_order", args)

// Pub/sub rooms
room := node.CreateRoom("orders")
room.Broadcast(ctx, message)

// Consensus voting
node.RequestJoin(ctx, newPeer) // other nodes vote
```

## The Key Insight

**Local events become distributed messages. Same contracts.**

capitan defines the event contract. herald transports it. ago adds saga semantics. Same code structure, distributed execution.

```
┌─────────────────────────────────────────────────────────────┐
│                      Service A                               │
│  ┌──────────────┐                                           │
│  │   capitan    │ ← All local events                        │
│  └──────┬───────┘                                           │
│         │                                                    │
│         ▼                                                    │
│  ┌──────────────┐       ┌──────────────┐                    │
│  │   herald     │──────▶│    Kafka     │                    │
│  │  (publish)   │       │   (broker)   │                    │
│  └──────────────┘       └──────┬───────┘                    │
└─────────────────────────────────┼────────────────────────────┘
                                  │
┌─────────────────────────────────┼────────────────────────────┐
│                      Service B  │                            │
│  ┌──────────────┐       ┌──────▼───────┐                    │
│  │   herald     │◀──────│    Kafka     │                    │
│  │ (subscribe)  │       │   (broker)   │                    │
│  └──────┬───────┘       └──────────────┘                    │
│         │                                                    │
│         ▼                                                    │
│  ┌──────────────┐                                           │
│  │   capitan    │ ← Becomes local event                     │
│  └──────────────┘                                           │
└─────────────────────────────────────────────────────────────┘
```

## Related Stories

- [Observability](observability.md) - herald extends event observability
- [Composition](composition.md) - ago uses pipz for saga orchestration
