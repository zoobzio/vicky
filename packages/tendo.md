# Tendo

Composable tensor computation for Go.

[GitHub](https://github.com/zoobzio/tendo)

## Vision

Tensors that compose. Tendo bridges the gap between heavy C++ libraries and limited pure-Go implementations. Write tensor operations that chain via pipz, emit signals via capitan, and run on CPU or CUDA without code changes. Numerical computation with the same vocabulary as the rest of the framework.

## Design Decisions

**Tensor = Storage + Shape**
Storage holds bytes on a device. Shape, stride, and offset describe the view. Views and slices share storage without copying. Non-contiguous layouts supported.

**Storage as interface**
CPU and CUDA backends implement the same contract. Write operations once, run anywhere. Device transfer is explicit.

**Operations return Chainable**
Every operation returns `pipz.Chainable[*Tensor]`. Build computation graphs with the same composition patterns used elsewhere. Type-safe, error-handling built in.

**Signal on every operation**
80+ capitan signals cover all operations. Foundation for autograd capture, logging, profiling. The data flows through the event system.

**Memory pooling**
Tensor allocation is expensive. Pool with LRU eviction reduces GC pressure. Per-device pools, statistics tracking.

**Context-driven behaviour**
Training vs inference mode. Pool selection. Trace IDs for autograd correlation. No global state, everything flows through context.

**Backend interface composition**
Separate interfaces: StorageOps, UnaryOps, BinaryOps, ReduceOps, MatrixOps. Backends implement what they support. Capability checking at runtime.

**Float16/BFloat16 support**
IEEE 754 conversions for precision trade-offs. Memory efficiency for large models.

## The Vocabulary

### Element-wise

| Operation | Purpose |
|-----------|---------|
| `Add`, `Sub`, `Mul`, `Div` | Arithmetic with broadcasting |
| `Neg`, `Abs`, `Square` | Unary transforms |
| `Exp`, `Log`, `Sqrt`, `Pow` | Mathematical functions |
| `Clamp`, `Where` | Conditional operations |

### Matrix

| Operation | Purpose |
|-----------|---------|
| `MatMul` | Matrix multiplication with batching |
| `Transpose`, `T` | Dimension reordering |

### Shape

| Operation | Purpose |
|-----------|---------|
| `Reshape`, `Squeeze`, `Unsqueeze` | Dimension manipulation |
| `Slice`, `Expand`, `Permute` | View operations |
| `Cat`, `Stack` | Tensor combination |

### Reductions

| Operation | Purpose |
|-----------|---------|
| `Sum`, `Mean`, `Prod` | Aggregations |
| `Max`, `Min`, `Var`, `Std` | Statistics |
| `ArgMax`, `ArgMin` | Index retrieval |

### Activations

| Operation | Purpose |
|-----------|---------|
| `ReLU`, `LeakyReLU` | Rectified linear |
| `Sigmoid`, `Tanh` | Bounded activations |
| `GELU`, `SiLU` | Modern activations |
| `Softmax`, `LogSoftmax` | Probability distributions |
| `Dropout` | Regularisation |

### Neural Network

| Operation | Purpose |
|-----------|---------|
| `Conv2d` | 2D convolution with groups |
| `MaxPool2d`, `AvgPool2d` | Spatial pooling |
| `BatchNorm2d`, `LayerNorm` | Normalisation |
| `MSELoss`, `L1Loss`, `CrossEntropyLoss`, `NLLLoss` | Loss functions |

### Device

| Operation | Purpose |
|-----------|---------|
| `ToCPU`, `ToCUDA` | Device transfer |
| `Contiguous` | Memory layout normalisation |

## Internal Architecture

```
Tensor
├── Storage (interface)
│   ├── CPUStorage ([]float32 on heap)
│   └── CUDAStorage (device pointer)
├── Shape []int
├── Stride []int
└── Offset int
        │
        ▼
Operations (pipz.Chainable[*Tensor])
        │
        ├── Dispatch to Backend
        │   ├── pkg/cpu (complete)
        │   └── pkg/cuda (in progress)
        │
        └── Emit Signal (capitan)
                │
                └── Autograd / Logging / Profiling
```

Memory Pool manages allocation. Context carries training mode, pool reference, trace ID.

## Code Organisation

| Category | Files |
|----------|-------|
| Core | `tensor.go`, `storage.go`, `shape.go`, `backend.go` |
| Constructors | `constructors.go`, `float16.go` |
| Operations | `ops.go`, `ops_elementwise.go`, `ops_matrix.go`, `ops_shape.go`, `ops_reduce.go` |
| Neural | `ops_activation.go`, `ops_norm.go`, `ops_conv.go`, `ops_pool.go`, `ops_loss.go` |
| Device | `ops_device.go`, `ops_dtype.go`, `ops_transfer.go` |
| Memory | `pool.go` |
| Signals | `signals.go` |
| Backends | `pkg/cpu/`, `pkg/cuda/` |

## Current State / Direction

CPU backend complete. CUDA infrastructure present, operations in progress.

Future considerations:
- Complete CUDA kernel implementations
- Autograd integration (signals prepared)
- ROCm backend
- Additional operations as patterns emerge

## Framework Context

**Dependencies**: pipz (composition), capitan (signals), clockz (timing).

**Role**: Numerical computation layer. Enables ML workloads, scientific computing, any domain requiring tensor operations. Same composition vocabulary as the rest of the framework.
