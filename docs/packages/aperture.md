# Aperture

Config-driven observability for Go.

[GitHub](https://github.com/zoobzio/aperture)

## Vision

Full observability across the framework with zero instrumentation code. Every package emits through capitan. Aperture bridges that unified stream to OpenTelemetry. Configure what gets logged, metricsed, and traced - the data is already flowing.

## Design Decisions

**Name-based matching**
Config uses signal names, not Go types. Enables hot-reload without recompilation.

**Three independent handlers**
Logs, metrics, traces processed separately. Configure any subset. Each is optional.

**User-provided OTEL providers**
Aperture decides transformation. User decides exporters, batching, sampling, security policies. Sharp boundary between concerns.

**Automatic type conversion**
Primitives map directly to OTEL attributes. Custom types JSON-serialised. Zero-configuration for standard types.

**Span correlation via composite keys**
`signalName + correlationID + endSignalName`. Handles out-of-order delivery gracefully.

**Timeout-based cleanup**
Pending spans cleaned on 1-minute ticker. Configurable per trace. Prevents memory leaks from unmatched events.

**Hot-reload via Apply()**
Swap configuration atomically. Drain old observer, close it, create new one. No data loss.

## Internal Architecture

```
Capitan Event Stream
        │
        ▼
capitanObserver
        │
        ├── Log Handler ──────→ OTEL LoggerProvider
        │
        ├── Metrics Handler ──→ OTEL MeterProvider
        │
        └── Traces Handler ───→ OTEL TracerProvider
```

Events processed synchronously in observer callback. No additional buffering - batching delegated to OTEL providers.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Main API: New(), Apply(), Close() |
| `schema.go` | User-facing config types, YAML/JSON loaders |
| `config.go` | Internal config structures |
| `capitan.go` | Event router, capitanObserver |
| `metrics.go` | Counter/gauge/histogram/updowncounter handling |
| `traces.go` | Span correlation, pending span management |
| `transform.go` | Field → OTEL attribute conversion |
| `stdout.go` | Human-readable stdout logging |
| `internal.go` | Diagnostic signals for config issues |

## Current State / Direction

Stable. All three signal types working. Hot-reload working. Diagnostic signals for missing values and expired spans.

Future considerations:
- Additional metric types as OTEL evolves
- Context propagation enhancements

## Framework Context

**Dependencies**: capitan, OpenTelemetry (heavy, intentional - this is an observability bridge).

**Role**: The observability layer. Plug it in, configure it, observe everything.
