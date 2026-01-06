# Aegis

Mesh networking for Go.

[GitHub](https://github.com/zoobzio/aegis)

## Vision

Peer-to-peer mesh with gRPC. Nodes discover each other, form a mesh, communicate directly. TLS required. Remote function execution, rooms for pub/sub, consensus voting for mesh membership.

## Design Decisions

**TLS required**
No insecure connections. `LoadOrGenerateTLS` handles certificate generation per node. mTLS for peer authentication.

**Versioned topology**
Topology version increments on node add/remove. Peers sync by comparing versions. Higher version wins on merge.

**Consensus for membership**
Join requests require votes from existing nodes. Any node can veto. Timeout treated as rejection. Protects mesh integrity.

**Host-based rooms**
Room host manages membership list. Messages relay through host to all members. Simple pub/sub without coordination overhead.

**Function registry**
Named functions with string parameters and string result. Execute locally or remotely on any peer.

**Node types**
`generic`, `gateway`, `processor`, `storage`. Semantic hints for topology-aware routing.

## Node Composition

```
Node
├── PeerManager (gRPC client connections)
├── MeshServer (gRPC service implementation)
├── FunctionRegistry (named functions)
├── RoomManager (pub/sub groups)
├── Topology (versioned node map)
├── ConsensusManager (join voting)
└── TLSConfig (certificates)
```

## Capabilities

| Feature | Description |
|---------|-------------|
| Peer management | Add/remove peers, ping, health check |
| Remote functions | Register locally, execute on self or peers |
| Rooms | Create, invite, join, leave, broadcast messages |
| Topology sync | Version-based sync, merge with higher version |
| Consensus | Vote on join requests, veto with reason |

## gRPC MeshService

| Operation | Purpose |
|-----------|---------|
| `Ping` | Connectivity check with timestamp |
| `GetHealth` | Node health status |
| `GetNodeInfo` | Node metadata and health |
| `NotifyHealthChange` | Broadcast peer failures |
| `ExecuteFunction` | Remote function invocation |
| `SendMessage` | Peer-to-peer messaging |
| `CreateRoom`, `JoinRoom`, `LeaveRoom` | Room lifecycle |
| `InviteToRoom`, `SendRoomMessage` | Room communication |
| `SyncTopology` | Topology version exchange |
| `RequestJoinMesh` | Initiate consensus vote |

## Consensus Flow

```
JoinMesh(entryNode)
    ↓
InitiateJoinRequest → broadcast to existing nodes
    ↓
Each node: EvaluateJoinRequest → approve or veto
    ↓
All votes OR timeout
    ↓
Approved: add to topology, sync to peers
Rejected: return veto reason
```

**Veto reasons:** blacklisted, resource_limit, suspicious, version_mismatch, custom

## Code Organisation

| File | Responsibility |
|------|----------------|
| `node.go` | Node type, high-level API |
| `peer.go` | PeerManager, Peer, gRPC clients |
| `server.go` | MeshServer, gRPC service |
| `topology.go` | Topology, NodeInfo, JoinRequest, Vote |
| `consensus.go` | ConsensusManager, VoteTracker |
| `rooms.go` | Room, RoomManager |
| `functions.go` | FunctionRegistry, NodeFunction |
| `health.go` | HealthInfo, HealthChecker |
| `tls.go` | TLSConfig, LoadOrGenerateTLS |
| `mesh.proto` | gRPC service definition |

## Current State / Direction

Stable. Core mesh networking complete.

Future considerations:
- Raft-based leader election
- Automatic peer discovery

## Framework Context

**Dependencies**: google.golang.org/grpc, github.com/google/uuid.

**Role**: Peer-to-peer mesh networking. Distributed coordination without central server. Foundation for distributed processing across framework components.
