# Layer 4: Infrastructure

Distributed systems concerns.

## Rule

Packages in this layer may depend on Layers 0-3.

## Packages

| Package | Role | Framework Deps |
|---------|------|----------------|
| [flux](../../packages/flux.md) | Reactive configuration | pipz, capitan, clockz |
| [flume](../../packages/flume.md) | Dynamic pipeline factory | pipz, capitan |
| [aperture](../../packages/aperture.md) | OpenTelemetry bridge | capitan |
| [herald](../../packages/herald.md) | Distributed messaging | capitan, pipz |
| [aegis](../../packages/aegis.md) | Mesh networking | (external: gRPC) |
| [sctx](../../packages/sctx.md) | Certificate security | capitan |

## Package Details

### flux

Reactive configuration synchronisation. Hot-reload without restarts.

**Key capabilities:**
- `Capacitor[T]` watches sources, validates, delivers updates
- State machine: Loading → Healthy / Degraded / Empty
- Eight watcher providers (file, Redis, Consul, etcd, NATS, Kubernetes, ZooKeeper, Firestore)
- Automatic rollback on validation failure
- pipz integration for resilience patterns
- Debouncing for rapid updates

**Framework dependencies:**
- pipz: Processing pipeline (Layer 1)
- capitan: Signal emission (Layer 0)
- clockz: Time abstraction (Layer 0)

**Why Layer 4:** Operational infrastructure. Configuration management across distributed systems.

### flume

Schema-driven pipeline construction. Define pipelines in YAML/JSON.

**Key capabilities:**
- YAML/JSON schemas reference registered processors
- Hot-reloading via `Binding.Update()`
- Rollback to previous versions
- Two-phase validation (syntax + references)
- Cycle detection at build time
- 14 connector types mapping to pipz primitives

**Framework dependencies:**
- pipz: Pipeline primitives (Layer 1)
- capitan: Events (Layer 0)

**Why Layer 4:** Configuration layer over pipz. Bridge between pipeline design and runtime.

### aperture

Config-driven observability. Bridge capitan to OpenTelemetry.

**Key capabilities:**
- Three handlers: Logs, metrics, traces
- Name-based matching (hot-reload without recompilation)
- User-provided OTEL providers (separation of concerns)
- Automatic type conversion to OTEL attributes
- Span correlation via composite keys
- Timeout-based cleanup for pending spans

**Framework dependencies:**
- capitan: Event stream (Layer 0)

**Why Layer 4:** Observability bridge. The "O" in distributed observability.

### herald

Distributed messaging. Extend capitan across process boundaries.

**Key capabilities:**
- 11 providers (Kafka, NATS, JetStream, Pub/Sub, Redis Streams, SQS, AMQP, SNS, BoltDB, Firestore, io)
- Bidirectional: publish local events, subscribe to external
- Generic `Publisher[T]` / `Subscriber[T]`
- pipz-based resilience (retry, backoff, circuit breaker)
- Envelope abstraction for metadata
- JSON default, pluggable codecs

**Framework dependencies:**
- capitan: Event coordination (Layer 0)
- pipz: Resilience patterns (Layer 1)

**Why Layer 4:** Distributed event streaming. Local events become distributed messages.

### aegis

Peer-to-peer mesh with gRPC.

**Key capabilities:**
- TLS required (mTLS for peer authentication)
- Versioned topology (sync by comparing versions)
- Consensus voting for mesh membership
- Host-based rooms for pub/sub
- Function registry (local and remote execution)
- Node types: generic, gateway, processor, storage

**Framework dependencies:**
- gRPC (external)

**Why Layer 4:** Mesh networking infrastructure. Distributed coordination without central server.

### sctx

Certificate-based security contexts.

**Key capabilities:**
- Zero-trust authentication from PKI
- Assertion-based proof (prevents replay)
- Generic `Context[M]` for type-safe metadata
- Delegatable guards (permission validators)
- Dual crypto: Ed25519 (fast) or ECDSA P-256 (FIPS)
- Certificate fingerprint binding

**Framework dependencies:**
- capitan: Events (Layer 0)

**Why Layer 4:** Security infrastructure. Certificate-driven authentication.

## Infrastructure Patterns

Layer 4 addresses distributed systems concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                      Application                             │
├─────────────┬─────────────┬─────────────┬──────────────────┤
│    flux     │   flume     │   herald    │      aegis       │
│  (config)   │ (pipelines) │ (messaging) │     (mesh)       │
├─────────────┴──────┬──────┴─────────────┴──────────────────┤
│      aperture      │              sctx                      │
│   (observability)  │           (security)                   │
├────────────────────┴───────────────────────────────────────┤
│                    capitan (events)                         │
└────────────────────────────────────────────────────────────┘
```

| Concern | Package | Mechanism |
|---------|---------|-----------|
| Configuration | flux | Reactive watchers, validation, rollback |
| Pipeline configuration | flume | Schema-driven, hot-reload |
| Distributed events | herald | Broker abstraction, 11 providers |
| Observability | aperture | OTEL bridge, config-driven |
| Mesh networking | aegis | gRPC, consensus, rooms |
| Authentication | sctx | Certificate-derived tokens |
