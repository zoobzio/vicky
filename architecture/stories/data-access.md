# Data Access Story

*"How do I access databases?"*

## Two Flows

**SQL Path:**
```
dbml → astql → soy → edamame
```

**Non-SQL Path:**
```
sentinel → atom → grub
```

## SQL Flow

### dbml - Schema Definition

Programmatic schema as Go code.

```go
project := dbml.NewProject("myapp").
    WithDatabaseType(dbml.PostgreSQL).
    AddTable("users").
        AddColumn("id").WithType("uuid").PrimaryKey().Done().
        AddColumn("email").WithType("varchar(255)").Unique().Done().
    Done()
```

### astql - Query Validation

Schema as allowlist. Validate at construction.

```go
instance := astql.New(schema, astql.WithPostgres())

// Safe - validated against schema
query := instance.Query().
    Select(instance.F("email"), instance.F("created_at")).
    From(instance.T("users")).
    Where(instance.F("id").Eq(instance.P("user_id")))

sql, params := instance.Render(query)
// SELECT "email", "created_at" FROM "users" WHERE "id" = :user_id
```

### soy - Ergonomic API

Type-safe results. Zero reflection on hot path.

```go
users := soy.New[User](db, renderer)

// Single record
user, err := users.Select().
    Where(soy.Eq("id", userID)).
    Exec(ctx, map[string]any{"id": userID})

// Multiple records
all, err := users.Query().
    Where(soy.Gt("created_at", since)).
    OrderBy("created_at", soy.DESC).
    Exec(ctx, params)
```

### edamame - Named Capabilities

Pure data specs for LLM integration.

```go
factory := edamame.NewFactory[User](soy)
factory.RegisterQuery("active_users", edamame.QuerySpec{
    Where: []edamame.ConditionSpec{{Field: "active", Op: "=", Value: true}},
    OrderBy: []edamame.OrderSpec{{Field: "created_at", Direction: "DESC"}},
})

// LLM-safe: spec as JSON
json, _ := factory.SpecJSON()

// Execute by name
users, err := factory.ExecQuery(ctx, "active_users", params)
```

## Non-SQL Flow

### sentinel - Type Metadata

Extract once, cache forever.

```go
meta := sentinel.Inspect[User]()
// Fields, tags, relationships available
```

### atom - Type Decomposition

Structs to typed maps.

```go
atomizer := atom.Use[User]()
atom := atomizer.Atomize(&user)
// atom.Strings["email"], atom.Times["created_at"], etc.
```

### grub - Provider-Agnostic Storage

Same code, any backend.

```go
// Redis
service := grub.NewService[User](redisProvider)

// Or S3, MongoDB, DynamoDB, etc.
service := grub.NewService[User](s3Provider)

// Same API
user, err := service.Get(ctx, userID)
err = service.Set(ctx, userID, &user)
```

## The Key Insight

**Schema-first validation. Types drive behaviour.**

SQL path: DBML schema validates all identifiers. Injection-resistant by design.

Non-SQL path: Sentinel extracts type structure. Atom normalises to typed maps. Grub abstracts storage.

```
┌─────────────────────────────────────────────────────────────┐
│                        SQL Path                              │
│                                                              │
│   struct tags → dbml → astql → soy → edamame                │
│       │          │       │       │        │                 │
│   schema def   schema  query   typed    named               │
│                        valid  results  capabilities         │
├─────────────────────────────────────────────────────────────┤
│                      Non-SQL Path                            │
│                                                              │
│   struct → sentinel → atom → grub                           │
│     │         │         │       │                           │
│   type def  metadata  typed   9 storage                     │
│                       maps    backends                       │
└─────────────────────────────────────────────────────────────┘
```

## Related Stories

- [Intelligence](intelligence.md) - cogito uses soy for persistence
- [Type Intelligence](type-intelligence.md) - sentinel enables data access
