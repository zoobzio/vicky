# Grub

Atom storage providers for Go.

[GitHub](https://github.com/zoobzio/grub)

## Vision

Thin bridge between storage backends and Atoms. Providers implement a minimal interface for storing and retrieving Atoms. No generic wrappers, no caching layer - just the Provider contract and backend implementations.

## Design Decisions

**Provider interface**
Five operations: Get, Set, Exists, List, Delete. All operate on Atoms directly. Table registration with Spec enables type-aware serialisation.

**Table registration**
Each type registered as a Table with name and Spec (type metadata). Providers use Spec for atom.Unflatten on retrieval.

**Flatten/Unflatten serialisation**
Atoms flattened to `map[string]any` for JSON storage. Unflatten reconstructs Atom using registered Spec. Works across all backends.

**Lifecycle as optional interface**
Connect, Close, Health separate from Provider. Stateless backends (S3, GCS) don't need lifecycle. Type assertion signals support.

## Provider Interface

```go
type Provider interface {
    Register(table Table) error
    Tables() []Table
    Get(ctx, table, key string) (*atom.Atom, error)
    Set(ctx, table, key string, data *atom.Atom) error
    Exists(ctx, table, key string) (bool, error)
    List(ctx, table, cursor string, limit int) (keys []string, nextCursor string, err error)
    Delete(ctx, table, key string) error
}
```

## Providers

| Provider | Backend | Storage Format |
|----------|---------|----------------|
| `redis` | Redis | Key prefix, JSON |
| `s3` | AWS S3 | Path prefix, JSON |
| `mongo` | MongoDB | Collections, BSON |
| `dynamo` | DynamoDB | Tables, AttributeValue |
| `bolt` | BoltDB | Buckets, JSON |
| `badger` | BadgerDB | Embedded KV, JSON |
| `firestore` | Firestore | Collections |
| `gcs` | Google Cloud Storage | Bucket objects, JSON |
| `azure` | Azure Blob | Container objects, JSON |

## Internal Architecture

```
Application
    ↓
Provider.Register(Table{Name, Spec})
    ↓
Provider.Set(ctx, table, key, atom)
    ↓
atom.Flatten() → map[string]any → JSON → Backend

Backend → JSON → map[string]any → atom.Unflatten(data, Spec) → *Atom
    ↓
Provider.Get(ctx, table, key) → *Atom
```

## Code Organisation

| Category | Files |
|----------|-------|
| Core | `api.go` (Provider, Lifecycle, Table, errors) |
| Internal | `internal/registry/` (table tracking) |
| Providers | `pkg/redis/`, `pkg/s3/`, `pkg/mongo/`, `pkg/dynamo/`, `pkg/bolt/`, `pkg/badger/`, `pkg/firestore/`, `pkg/gcs/`, `pkg/azure/` |
| Testing | `testing/helpers.go` |

## Current State / Direction

Stable. Nine providers complete.

Future considerations:
- Additional providers as demand emerges
- Batch operations

## Framework Context

**Dependencies**: atom.

**Role**: Storage layer for Atoms. Providers bridge storage backends to the Atom representation. Part of the Database → Atom → Cache flow described in atom docs.
