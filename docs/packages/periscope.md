# Periscope

3D visualization for distributed systems.

[GitHub](https://github.com/zoobzio/periscope)

## Vision

See your architecture. Periscope renders distributed systems as interactive 3D voxel worlds. Services become geometry. State becomes colour and brightness. Interactions become animations. Connect an orchestrator, watch infrastructure come alive.

## Design Decisions

**Orchestrator controls what, client controls how**
JSON messages define elements, zones, and animations. Client handles camera and movement. No network latency in navigation.

**SDF-based geometry**
Elements defined by mathematical expressions (signed distance functions). Procedural shapes, resolution-independent. Falls back to explicit voxel lists when needed.

**Zone system**
Each element contains zones - independently colourable regions. A service's frontend and backend light up separately. Fine-grained state without geometry changes.

**GPU-resident state**
Geometry uploads once. Transforms and properties live in GPU buffers, updated per-frame. Smooth animations without re-voxelisation.

**Dirty flag pattern**
Geometry, properties, transforms tracked separately. Only rebuild what changed. No per-frame overhead for static scenes.

**Async socket, sync rendering**
Non-blocking message receive. Frame rendering continues during orchestrator latency. Queue buffers bursts.

**Client-side collision**
Sphere-box collision on occupied blocks. Responsive WASD movement without round-trips.

**Strict resource limits**
Max 100k voxels per element. Max 1k elements. Prevents runaway allocations from misconfigured orchestrators.

## Protocol

### World Messages

| Message | Purpose |
|---------|---------|
| `add_elements` | Define new geometry with zones |
| `remove_elements` | Remove by name |
| `camera` | Set zoom and view distance |

### Stream Messages

| Message | Purpose |
|---------|---------|
| `zone_properties` | Update colour, brightness by zone |
| `blocks` | Transient voxels (ephemeral state) |
| `animations` | Trigger play/pause/stop/resume |
| `camera` | Update view parameters |

## Internal Architecture

```
Unix Socket (JSON per line)
        │
        ▼
Message Parser (new format + legacy fallback)
        │
        ▼
World State
├── Elements (geometry, zones, animations)
├── Zone Properties (colour, brightness)
├── Element Transforms (position, scale, rotation)
└── Dirty Flags (geometry, properties, transforms)
        │
        ▼
Animation Manager (keyframe interpolation)
        │
        ▼
Render Pipeline (wgpu)
├── Static Mesh (culled cubes, zone IDs)
├── Property Buffer (per-zone colour/brightness)
├── Transform Buffer (per-element matrices)
└── Uniforms (camera, lighting)
        │
        ▼
Frame Output + egui Overlay
```

Camera and collision operate on client state. No GPU round-trip for movement.

## Code Organisation

| Category | Files |
|----------|-------|
| Entry | `main.rs` |
| Application | `app/mod.rs`, `app/update.rs`, `app/ui.rs` |
| Scene | `scene/world.rs`, `scene/element.rs`, `scene/animation.rs`, `scene/camera.rs`, `scene/collision.rs`, `scene/sdf.rs`, `scene/block.rs` |
| Rendering | `render/pipeline.rs`, `render/static_mesh.rs`, `render/uniforms.rs`, `render/overlay.rs`, `render/grid.rs`, `render/cube.rs` |
| Socket | `socket/listener.rs`, `socket/protocol.rs` |
| Input | `input/keyboard.rs` |

## Current State / Direction

Stable. Core rendering, animation, and socket communication complete.

Future considerations:
- Full easing function implementations
- Zone property animation interpolation
- Shader hot-reload
- Additional geometry primitives

## Framework Context

**Dependencies**: wgpu (GPU), winit (windowing), egui (UI), tokio (async), rayon (parallel voxelisation), fasteval (SDF expressions).

**Role**: Visualisation endpoint. Receives system state from orchestrators, renders interactive 3D representation. Companion to monitoring and observability tooling.
