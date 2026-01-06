# Layer 7: Endpoints

External consumers. The boundary between framework and outside world.

## Rule

This layer consumes framework output. Packages here may be written in any language.

## Packages

| Package | Role | Language | Framework Deps |
|---------|------|----------|----------------|
| [periscope](../../packages/periscope.md) | 3D visualisation | Rust | External (consumes JSON) |

## Package Details

### periscope

3D visualisation for distributed systems. Interactive voxel rendering.

**Key capabilities:**
- SDF-based geometry (procedural, resolution-independent)
- Zone system for fine-grained state colouring
- GPU-resident state (smooth animations)
- Keyframe animation system
- Client-side camera and collision
- Unix socket protocol (JSON per line)

**Technology stack:**
- wgpu: GPU abstraction
- winit: Windowing
- egui: Overlay UI
- tokio: Async runtime
- rayon: Parallel voxelisation
- fasteval: SDF expression parsing

**Why Layer 7:** External consumer. Receives state from orchestrators, provides visual representation.

## Protocol Architecture

periscope communicates via Unix socket with JSON messages:

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Framework                              │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐ │
│  │    capitan     │  │    herald      │  │   aperture     │ │
│  │   (events)     │  │  (messaging)   │  │ (observability)│ │
│  └───────┬────────┘  └───────┬────────┘  └───────┬────────┘ │
│          │                   │                   │          │
│          └───────────────────┼───────────────────┘          │
│                              │                              │
│                    ┌─────────▼─────────┐                    │
│                    │   Orchestrator    │                    │
│                    │   (Go adapter)    │                    │
│                    └─────────┬─────────┘                    │
└──────────────────────────────┼──────────────────────────────┘
                               │ JSON over Unix Socket
┌──────────────────────────────┼──────────────────────────────┐
│                              ▼                              │
│                    ┌───────────────────┐                    │
│                    │    periscope      │                    │
│                    │   (Rust client)   │                    │
│                    └───────────────────┘                    │
│                         Layer 7                              │
└─────────────────────────────────────────────────────────────┘
```

### Message Types

**World Messages** (geometry definition):
| Message | Purpose |
|---------|---------|
| `add_elements` | Define geometry with zones |
| `remove_elements` | Remove by name |
| `camera` | Set view parameters |

**Stream Messages** (state updates):
| Message | Purpose |
|---------|---------|
| `zone_properties` | Colour and brightness |
| `blocks` | Transient voxels |
| `animations` | Trigger animations |

## Integration Pattern

The orchestrator pattern bridges framework internals to periscope:

```go
// Example orchestrator (not part of periscope)
capitan.Hook(ctx, mySignal, func(e *capitan.Event) {
    msg := toPeriscopeMessage(e)
    socket.Write(msg)
})
```

This separation means:
- Framework packages remain Go-only
- Visualisation can evolve independently
- Protocol is language-agnostic
- Other visualisation tools can implement same protocol

## Architectural Significance

Layer 7 represents the boundary:

| Inside (Go) | Boundary | Outside |
|-------------|----------|---------|
| capitan events | JSON protocol | periscope rendering |
| herald messages | Socket transport | 3D visualisation |
| aperture telemetry | Orchestrator adapters | Interactive UI |

The framework doesn't depend on periscope. Periscope consumes framework output. This clean boundary enables:
- Independent deployment
- Language flexibility
- Protocol evolution
- Alternative implementations
