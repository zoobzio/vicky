---
name: add-contract
description: Create a contract interface in the contracts/ package
---

# Add Contract

You are creating a contract - an interface that defines the API a handler needs from an implementation. Contracts are the bridge between handlers and their dependencies (stores, API clients, service wrappers, etc.).

## Prerequisites

**The implementation must exist first.** A contract is extracted from what an implementation can do, not invented in isolation.

- **Implementation doesn't exist?** Trigger the appropriate skill (`/add-store`, `/add-client`, etc.) to create it first, then return to complete the contract.
- **Implementation exists but missing methods?** You can add methods to the implementation as part of this skill (with user approval in the spec).

## Technical Foundation

Contracts live in `contracts/` as interfaces:

```go
package contracts

import (
    "context"

    "github.com/yourorg/yourapp/models"
)

// Users defines the contract for user operations.
type Users interface {
    // Get retrieves a user by primary key.
    Get(ctx context.Context, key string) (*models.User, error)
    // Set creates or updates a user.
    Set(ctx context.Context, key string, user *models.User) error
    // GetByLogin retrieves a user by GitHub login.
    GetByLogin(ctx context.Context, login string) (*models.User, error)
}
```

**Key principles:**

- Every method takes `context.Context` as first parameter
- Every method returns `error` as last return value
- Method names describe what they do: `Get`, `Set`, `GetByLogin`, `ListByVersion`, `FindSimilar`
- Each method has a doc comment explaining its purpose

### What Satisfies Contracts

Contracts can be satisfied by any struct that implements the interface:

- **Stores** - database-backed implementations (`stores/users.go` → `contracts.Users`)
- **API clients** - external service wrappers (`external/github/client.go` → `contracts.GitHub`)
- **gRPC clients** - service clients (`external/indexer/client.go` → `contracts.Indexer`)
- **Mocks** - test doubles (`testing/mocks.go` → any contract)

### Registration and Usage

```go
// In main.go - register implementation against contract
sum.Register[contracts.Users](k, allStores.Users)
sum.Register[contracts.GitHub](k, github.NewClient())

// In handlers - retrieve by contract type
users := sum.MustUse[contracts.Users](req.Context)
user, err := users.Get(req.Context, userID)
```

### Dependency Direction

```
implementation (store, client)  ←  contract  ←  handler
       no deps on contracts        interface     uses contract
```

- Implementations have NO dependency on contracts or handlers
- Handlers import contracts (to declare what they need)

## Your Task

Understand what the user is building:

1. **What implementation exists?** (store, client, service wrapper)
2. **What handler will consume this?** (or what operations are needed)
3. **What specific methods are needed?** Not generic CRUD - the exact operations the handler requires

The contract should expose exactly what's needed, nothing more.

## Before Writing Code

Produce a spec for approval:

```
## Contract: [Name]

**Implementation:** [What struct satisfies this]

**Consumer:** [What handler/service will use this]

**Methods (existing):**

[MethodName](ctx context.Context, [params]) ([returns], error)
  Purpose: [What this method does]

**Methods (to be added to implementation):**

[NewMethodName](ctx context.Context, [params]) ([returns], error)
  Purpose: [What this method does]
  Implementation: [Brief description of how it will work]
```

If all methods already exist on the implementation, omit the "to be added" section.

Wait for the user to approve before writing any files.

## After Approval

1. If methods need to be added to the implementation, add them first
2. Create `contracts/[name].go` with the interface
3. Verify the implementation satisfies the interface (compile check)
4. Update `cmd/app/main.go` to register the implementation:
   ```go
   sum.Register[contracts.Name](k, implementation)
   ```
