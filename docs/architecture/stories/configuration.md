# Configuration Story

*"How do I manage configuration?"*

## The Packages

```
flux + flume
```

## flux - Reactive Configuration

Live configuration that reacts to changes. Auto-rollback on failure.

```go
type AppConfig struct {
    DatabaseURL string `json:"database_url" validate:"required,url"`
    MaxWorkers  int    `json:"max_workers" validate:"min=1,max=100"`
}

func (c *AppConfig) Validate() error {
    // Custom validation logic
    return nil
}

capacitor := flux.NewCapacitor[*AppConfig](
    fileWatcher,
    flux.WithDebounce(100*time.Millisecond),
    flux.WithRetry(3),
)

capacitor.Start(ctx, func(cfg *AppConfig) {
    // React to config changes
    updateWorkerPool(cfg.MaxWorkers)
})
```

**State Machine:**
| State | Meaning |
|-------|---------|
| Loading | Initial, no config yet |
| Healthy | Valid config applied |
| Degraded | Last change failed, previous config retained |
| Empty | Initial load failed, no valid config ever |

**Watcher Providers:**
| Provider | Source |
|----------|--------|
| file | File system (fsnotify) |
| redis | Redis keyspace notifications |
| consul | Consul blocking queries |
| etcd | etcd Watch API |
| nats | NATS JetStream KV |
| kubernetes | ConfigMap/Secret watch |
| zookeeper | ZooKeeper node watch |
| firestore | Firestore realtime listeners |

**Composite Sources:**
```go
// Layer config: defaults + environment + file
composite := flux.NewCompositeCapacitor[*AppConfig](
    []flux.Watcher{defaultsWatcher, envWatcher, fileWatcher},
    mergeConfigs,
)
```

## flume - Pipeline Configuration

Define pipelines in YAML. Hot-reload without restarts.

```yaml
# pipeline.yaml
root:
  type: sequence
  steps:
    - ref: validate
    - type: retry
      max_retries: 3
      processor:
        ref: process
    - ref: notify
```

```go
factory := flume.NewFactory[Request]()
factory.Add("validate", validateProcessor)
factory.Add("process", processProcessor)
factory.Add("notify", notifyProcessor)

binding := factory.Bind(identity, schema)

// Later: update pipeline without restart
binding.Update(newSchema)

// Rollback if needed
binding.Rollback()
```

**Binding Lifecycle:**
```
Bind(schema) → Processing
                  │
        ┌─────────┴─────────┐
        │                   │
   Update(new)         Rollback()
        │                   │
        ▼                   ▼
   New Pipeline     Previous Pipeline
```

## The Key Insight

**Configuration is reactive. State machine semantics.**

flux watches sources. Validation happens on every change. Invalid config? Retain previous. Empty on first failure? Report, don't crash. Degraded state visible to health checks.

```
┌─────────────────────────────────────────────────────────────┐
│                      Configuration                           │
│                                                              │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐               │
│   │   file   │   │  consul  │   │   etcd   │   ...         │
│   └────┬─────┘   └────┬─────┘   └────┬─────┘               │
│        │              │              │                      │
│        └──────────────┼──────────────┘                      │
│                       ▼                                      │
│               ┌──────────────┐                               │
│               │  Capacitor   │                               │
│               │  (watches)   │                               │
│               └──────┬───────┘                               │
│                      │                                       │
│                      ▼                                       │
│               ┌──────────────┐                               │
│               │  Validate()  │                               │
│               └──────┬───────┘                               │
│                      │                                       │
│         ┌────────────┼────────────┐                         │
│         │            │            │                         │
│         ▼            ▼            ▼                         │
│     Healthy      Degraded      Empty                        │
│   (new config)  (keep old)  (no config)                     │
└─────────────────────────────────────────────────────────────┘
```

## Related Stories

- [Composition](composition.md) - flume configures pipz pipelines
- [Observability](observability.md) - aperture hot-reloads config
