# Numerical Computation Story

*"How do I do ML/scientific computing?"*

## The Flow

```
pipz + capitan → tendo
```

## The Package

### tendo - Composable Tensor Computation

Tensors that compose. Same vocabulary as everything else.

```go
// Create tensors
x := tendo.Randn([]int{32, 784})  // batch of 32, 784 features
w := tendo.Randn([]int{784, 128}) // weight matrix

// Operations return Chainable
pipeline := pipz.NewSequence[*tendo.Tensor](
    tendo.MatMul(w),
    tendo.ReLU(),
    tendo.Dropout(0.5),
)

result, err := pipeline.Process(ctx, x)
```

**80+ Operations:**

| Category | Operations |
|----------|------------|
| Element-wise | Add, Sub, Mul, Div, Neg, Abs, Exp, Log, Sqrt, Pow, Clamp, Where |
| Matrix | MatMul, Transpose |
| Shape | Reshape, Squeeze, Unsqueeze, Slice, Expand, Permute, Cat, Stack |
| Reductions | Sum, Mean, Prod, Max, Min, Var, Std, ArgMax, ArgMin |
| Activations | ReLU, LeakyReLU, Sigmoid, Tanh, GELU, SiLU, Softmax, Dropout |
| Neural | Conv2d, MaxPool2d, AvgPool2d, BatchNorm2d, LayerNorm |
| Loss | MSELoss, L1Loss, CrossEntropyLoss, NLLLoss |
| Device | ToCPU, ToCUDA, Contiguous |

### Storage Architecture

```
Tensor
├── Storage (interface)
│   ├── CPUStorage ([]float32 on heap)
│   └── CUDAStorage (device pointer)
├── Shape []int
├── Stride []int
└── Offset int
```

- Views and slices share storage (no copy)
- Non-contiguous layouts supported
- Explicit device transfer (ToCPU, ToCUDA)

### Signal Emission

Every operation emits to capitan:

```go
// tendo internally emits
capitan.Emit(ctx, tendo.MatMulSignal,
    tendo.ShapeKey.Field(shape),
    tendo.DurationKey.Field(elapsed),
)
```

Foundation for:
- Autograd capture
- Profiling
- Logging
- Metrics

### Memory Pooling

```go
ctx := tendo.WithPool(ctx, pool)

// Allocations come from pool
tensor := tendo.Zeros(ctx, []int{1000, 1000})

// LRU eviction when full
// Per-device pools
// Statistics tracking
```

## The Key Insight

**Tensors compose like everything else in the framework.**

Operations return `pipz.Chainable[*Tensor]`. Same reliability patterns. Same observability via capitan. Same composition vocabulary.

```
┌─────────────────────────────────────────────────────────────┐
│                    Neural Network                            │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  pipz.Sequence[*Tensor]                               │  │
│  │    ↓                                                  │  │
│  │  tendo.MatMul(w1) → ReLU → Dropout → MatMul(w2) ...  │  │
│  └───────────────────────────────────────────────────────┘  │
│                         │                                    │
│                         ▼                                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Backend (pkg/cpu or pkg/cuda)                        │  │
│  └───────────────────────────────────────────────────────┘  │
│                         │                                    │
│                         ▼                                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  capitan signals (for autograd, profiling, logging)   │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Example: Simple MLP

```go
// Define layers
linear1 := tendo.MatMul(w1)
linear2 := tendo.MatMul(w2)

// Build network as pipeline
network := pipz.NewSequence[*tendo.Tensor](
    linear1,
    tendo.ReLU(),
    tendo.Dropout(0.5),
    linear2,
    tendo.Softmax(-1),
)

// Forward pass
output, err := network.Process(ctx, input)

// Loss
loss, err := tendo.CrossEntropyLoss()(ctx, output, targets)
```

## Backend Status

| Backend | Status |
|---------|--------|
| CPU | Complete |
| CUDA | In progress |
| ROCm | Future |

## Related Stories

- [Composition](composition.md) - pipz patterns for tensor operations
- [Observability](observability.md) - capitan signals for autograd
