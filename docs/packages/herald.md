# Herald

Distributed messaging for Go.

[GitHub](https://github.com/zoobzio/herald)

## Vision

Bridge the unified event stream to external systems. Capitan handles in-process coordination. Herald extends that across process boundaries. Emit locally, publish to Kafka/NATS/SQS/etc. Subscribe externally, receive as local events. Same event contracts, distributed.

## Design Decisions

**Bidirectional**
Publisher observes capitan signals → publishes to broker. Subscriber consumes broker → emits to capitan. Symmetric patterns across services.

**Provider interface minimalism**
Four methods: Publish, Subscribe, Ping, Close. Enough for 11 different brokers without abstraction bloat.

**Type-safe generics**
`Publisher[T]` and `Subscriber[T]`. Compile-time checking, no runtime type assertions.

**Pipz integration**
Retry, backoff, circuit breaker, rate limiting composed from pipz primitives. Herald doesn't reimplement resilience patterns.

**Envelope abstraction**
Value + metadata travel together through middleware. No type assertion gymnastics.

**Errors through capitan**
All failures emit to `herald.ErrorSignal`. Unified error handling across the application.

**JSON default, codecs pluggable**
Language-neutral default. Protobuf, MessagePack available via Codec interface.

**Context-native cancellation**
Publisher and subscriber respect context for graceful shutdown and timeouts.

## Providers

| Provider | Package | Notes |
|----------|---------|-------|
| Kafka | `pkg/kafka` | High throughput, ordered within partition |
| NATS | `pkg/nats` | Lightweight, low latency |
| JetStream | `pkg/jetstream` | Persistent NATS with ordering |
| Google Pub/Sub | `pkg/pubsub` | GCP managed |
| Redis Streams | `pkg/redis` | In-memory with persistence |
| AWS SQS | `pkg/sqs` | AWS managed queues |
| RabbitMQ/AMQP | `pkg/amqp` | Traditional message broker |
| AWS SNS | `pkg/sns` | Pub/sub fanout (publish only) |
| BoltDB | `pkg/bolt` | Embedded, local queues |
| Firestore | `pkg/firestore` | GCP document store |
| io | `pkg/io` | Testing, CLI piping |

Each provider implements the four-method interface. Broker-specific dependencies isolated to their package.

## Internal Architecture

```
Publish Flow:
Capitan Event → Publisher → Pipeline (middleware) → Provider → Broker

Subscribe Flow:
Broker → Provider → Subscriber → Pipeline (middleware) → Capitan Event
```

Pipeline layers compose: middleware → retry → backoff → terminal operation.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Provider interface, Metadata, Message, Result types |
| `publisher.go` | Publisher[T] implementation |
| `subscriber.go` | Subscriber[T] implementation |
| `codec.go` | Codec interface, JSONCodec |
| `options.go` | Reliability options wrapping pipz |
| `processing.go` | Middleware processors |
| `pkg/` | 11 provider implementations |

## Current State / Direction

Stable. 11 providers working. Core abstractions settled.

Future considerations:
- Additional providers as needs arise
- Enhanced broker-specific features where abstraction permits

## Framework Context

**Dependencies**: capitan, pipz.

**Role**: Distributed extension of the event stream. Local events become distributed messages. External messages become local events. Same contracts, across services.
