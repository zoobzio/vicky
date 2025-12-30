# Flux

Reactive configuration synchronization for Go.

[GitHub](https://github.com/zoobzio/flux)

## Vision

Live configuration that reacts to external changes. Capacitor watches sources, validates data, and delivers updates with automatic rollback. Hot-reload without restarts. Composable pipelines for reliability patterns.

## Design Decisions

**Capacitor as core primitive**
Generic `Capacitor[T Validator]` watches a source. Type-safe throughout. Validation via interface gives full control to consumers.

**State machine semantics**
Four states: Loading, Healthy, Degraded, Empty. Degraded retains last good config. Empty indicates initial failure. Clear operational semantics.

**Watcher abstraction**
Interface abstracts change sources. Core provides ChannelWatcher for testing. pkg/ has 8 production implementations. Easy to add new sources.

**Pipz integration**
Processing pipeline built on pipz. Retry, backoff, timeout, circuit breaker, rate limiting. Same composable patterns as elsewhere in framework.

**Composite sources**
`CompositeCapacitor` watches multiple sources. Reducer merges parsed values. Enables layered config (defaults + overrides + environment).

**Debouncing**
Changes coalesced over configurable window. Prevents thundering herd from rapid updates. Configurable per capacitor.

**Sync mode for testing**
Disables debouncing and goroutines. Manual `Process()` calls. Deterministic test execution. clockz integration for time control.

**Capitan integration**
All state transitions and processing events emitted. Full observability of configuration lifecycle.

## Internal Architecture

```
Application
    ↓
Capacitor[T] / CompositeCapacitor[T]
    ↓
Watcher (source abstraction)
    ↓                          ↓
pkg/file, pkg/redis, ...   ChannelWatcher (test)
    ↓
Codec (JSON/YAML deserialization)
    ↓
T.Validate()
    ↓
Pipeline (pipz-based)
    ↓
Callback / Reducer
    ↓
State Machine (Loading → Healthy / Degraded / Empty)
```

## State Machine

| State | Meaning |
|-------|---------|
| Loading | Initial, no config yet |
| Healthy | Valid config applied |
| Degraded | Last change failed, previous config retained |
| Empty | Initial load failed, no valid config ever obtained |

## Watcher Providers

| Provider | Source |
|----------|--------|
| `pkg/file` | File system (fsnotify) |
| `pkg/redis` | Redis keyspace notifications |
| `pkg/consul` | Consul blocking queries |
| `pkg/etcd` | etcd Watch API |
| `pkg/nats` | NATS JetStream KV |
| `pkg/kubernetes` | ConfigMap/Secret watch |
| `pkg/zookeeper` | ZooKeeper node watch |
| `pkg/firestore` | Firestore realtime listeners |

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Capacitor definition, Start, Process, state access |
| `compose.go` | CompositeCapacitor for multiple sources |
| `state.go` | State enum and transitions |
| `options.go` | Pipeline options (With*) and middleware (Use*) |
| `codec.go` | JSON/YAML codecs |
| `signals.go` | Capitan signals |
| `metrics.go` | MetricsProvider interface |
| `request.go` | Request type for pipeline |
| `fields.go` | Capitan field keys |
| `error_ring.go` | Error history ring buffer |
| `channel_watcher.go` | ChannelWatcher for testing |

## Current State / Direction

Stable. Core reactive pattern complete.

Future considerations:
- Additional watcher providers
- Schema migration support

## Framework Context

**Dependencies**: pipz (processing pipeline), capitan (events), clockz (time abstraction).

**Role**: Reactive configuration layer. Hot-reload without restarts. Integrates with framework observability through capitan signals.
