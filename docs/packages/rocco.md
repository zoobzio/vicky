# Rocco

Type-safe HTTP framework for Go.

[GitHub](https://github.com/zoobzio/rocco)

## Vision

Struct tags as single source of truth. Define your types, get validation AND OpenAPI documentation. No drift between code and docs. Leverage sentinel's type intelligence to generate accurate API specifications automatically.

## Design Decisions

**Generic handlers**
`Handler[In, Out]`. Compile-time type checking. Request.Body is concrete, not `interface{}`.

**Sentinel integration**
Type metadata extracted at handler creation via `sentinel.Scan[T]()`. Struct tags drive both runtime validation AND OpenAPI schema generation.

**Sentinel errors**
Pre-defined typed errors with `Error[D]`. Must declare expected errors via `WithErrors()`. Undeclared errors return 500 (prevents leaking internal details).

**Chi foundation**
Wrap battle-tested router rather than reimplementing. Users can access underlying Chi for advanced cases.

**Capitan events**
Observability hooks throughout lifecycle. Engine, handler, auth, streaming events. No hard dependency on specific backends.

**Declarative everything**
All behaviour declared in code. No hidden config files. Explicit over magic.

## Internal Architecture

```
HTTP Request → Chi Router → Middleware → Handler.Process()
    ├── Extract path/query params
    ├── Read and parse body
    ├── Validate input (struct tags → go-playground/validator)
    ├── Execute handler function
    ├── Handle errors (check if declared, map to response)
    └── Marshal and write response
```

Auth, authorization, and rate limiting added as middleware when configured.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `engine.go` | Server lifecycle, routing, middleware orchestration |
| `handler.go` | Handler[In,Out], request processing pipeline |
| `errors.go` | Error[D], built-in errors, details types |
| `docs.go` | OpenAPI generation (uses sentinel for type metadata) |
| `events.go` | Capitan signals and field keys |
| `stream.go` | SSE streaming implementation |
| `identity.go` | Identity interface for auth |
| `request.go` | Request[In], parameter extraction |

## OpenAPI Generation

Validation tags map to OpenAPI constraints:

| Tag | OpenAPI |
|-----|---------|
| `required` | required array |
| `min=N` / `max=N` | minLength/maxLength or minimum/maximum |
| `oneof=a b c` | enum |
| `email` | format: "email" |
| `uuid4` | format: "uuid" |

Generated once on first `/openapi` request, cached thereafter.

## Current State / Direction

Stable. Core HTTP handling and OpenAPI generation complete.

Future considerations:
- Additional validation tag mappings
- Enhanced streaming capabilities

## Framework Context

**Dependencies**: chi, go-playground/validator, capitan, sentinel, openapi.

**Role**: HTTP layer. Type-safe handlers with automatic documentation. Sentinel powers the "single source of truth" - struct definitions drive both runtime behaviour and API specifications.
