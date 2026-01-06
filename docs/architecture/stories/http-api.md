# HTTP API Story

*"How do I build HTTP services?"*

## The Flow

```
sentinel → openapi → rocco
```

## The Packages

### sentinel - Type Metadata

Extract struct metadata once.

```go
// Tags drive everything
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=1,max=100"`
    Age      int    `json:"age" validate:"min=0,max=150"`
}

type CreateUserResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}
```

### openapi - Specification Types

OpenAPI 3.0 as Go structs.

```go
// rocco generates this automatically from struct tags
spec := openapi.OpenAPI{
    Paths: map[string]openapi.PathItem{
        "/users": {
            Post: &openapi.Operation{
                RequestBody: &openapi.RequestBody{...},
                Responses: map[string]*openapi.Response{...},
            },
        },
    },
}
```

### rocco - Type-Safe Handlers

Struct tags as single source of truth.

```go
handler := rocco.NewHandler(
    rocco.POST,
    "/users",
    func(ctx context.Context, req rocco.Request[CreateUserRequest]) (*CreateUserResponse, error) {
        // req.Body is typed, validated
        user := createUser(req.Body)
        return &CreateUserResponse{
            ID:        user.ID,
            Email:     user.Email,
            CreatedAt: user.CreatedAt.Format(time.RFC3339),
        }, nil
    },
    rocco.WithErrors(
        rocco.Error[ValidationError](400),
        rocco.Error[ConflictError](409),
    ),
)

engine := rocco.NewEngine()
engine.Register(handler)
engine.Start(":8080")
```

**OpenAPI Generation:**
```go
// GET /openapi returns spec
// Generated from struct tags - no drift
```

**Tag to OpenAPI mapping:**
| Tag | OpenAPI |
|-----|---------|
| `required` | required array |
| `min=N` / `max=N` | minLength/maxLength or minimum/maximum |
| `oneof=a b c` | enum |
| `email` | format: "email" |
| `uuid4` | format: "uuid" |

## The Key Insight

**Define once, get runtime validation AND API docs.**

Struct tags drive both go-playground/validator at runtime AND OpenAPI generation. No drift between code and documentation.

```
┌─────────────────────────────────────────────────────────────┐
│                       struct tags                            │
│   type Request struct {                                      │
│       Email string `json:"email" validate:"required,email"` │
│   }                                                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
┌─────────────────────────┐ ┌─────────────────────────┐
│   Runtime Validation    │ │   OpenAPI Generation    │
│  (go-playground/valid)  │ │  (sentinel + openapi)   │
└─────────────────────────┘ └─────────────────────────┘
              │                         │
              ▼                         ▼
┌─────────────────────────┐ ┌─────────────────────────┐
│   400 Bad Request       │ │   GET /openapi          │
│   {errors: [...]}       │ │   (JSON spec)           │
└─────────────────────────┘ └─────────────────────────┘
```

## Error Handling

Typed errors prevent leaking internals.

```go
// Declared errors return proper status
rocco.WithErrors(
    rocco.Error[NotFoundError](404),
    rocco.Error[ValidationError](400),
)

// Undeclared errors return 500 (no details leaked)
```

## SSE Streaming

```go
handler := rocco.NewStreamHandler(
    rocco.GET,
    "/events",
    func(ctx context.Context, req rocco.Request[Empty], stream rocco.Stream) error {
        for event := range events {
            stream.Send(event)
        }
        return nil
    },
)
```

## Related Stories

- [Type Intelligence](type-intelligence.md) - sentinel enables HTTP patterns
- [Observability](observability.md) - rocco emits to capitan
