# Trainer

LLM finetuning pipeline for Go framework documentation.

## Overview

This package provides a complete pipeline for finetuning Qwen 2.5 7B (or similar models) on your Go framework documentation and code patterns. It uses LoRA/QLoRA for efficient training and exports models in safetensors format.

## Requirements

- Python 3.10+
- CUDA-capable GPU (24GB+ VRAM recommended for QLoRA)
- ~50GB disk space for model weights and checkpoints

## Setup

```bash
# Create virtual environment
make setup
source .venv/bin/activate

# Install dependencies
make install

# Install dev dependencies (for testing/linting)
make dev
```

## Quick Start

### 1. Prepare Training Data

Create training examples in `data/raw/examples/` in JSONL format:

```jsonl
{"messages": [{"role": "user", "content": "How do I validate a schema?"}, {"role": "assistant", "content": "Use the Validate() method..."}]}
{"messages": [{"role": "user", "content": "What is the flux pattern?"}, {"role": "assistant", "content": "Flux is a data flow pattern..."}]}
```

For tool-use examples:

```jsonl
{"messages": [{"role": "user", "content": "Show me the aegis docs"}, {"role": "assistant", "tool_calls": [{"name": "lookup_docs", "arguments": {"package": "aegis"}}]}, {"role": "tool", "content": "aegis package documentation..."}, {"role": "assistant", "content": "Here's what I found..."}]}
```

### 2. Process Data

```bash
make synthesize
```

This will:
- Extract code/docs from configured GitHub repos (optional)
- Load your training examples
- Split into train/validation/test sets
- Save in HuggingFace datasets format

### 3. Train

```bash
# Full training
make train

# Debug run (quick test)
make train-debug
```

Training outputs:
- Checkpoints: `outputs/checkpoints/`
- Logs: `outputs/logs/`
- W&B dashboard (if configured)

### 4. Evaluate

```bash
python scripts/evaluate.py --checkpoint outputs/checkpoints/checkpoint-XXX
```

### 5. Export

```bash
python scripts/export.py --checkpoint outputs/checkpoints/checkpoint-XXX
```

Exports merged model to `outputs/models/` in safetensors format.

## Configuration

Configuration files in `configs/`:

- `model/qwen-7b-lora.yaml` - Model and LoRA settings
- `data/synthesis.yaml` - Data sources and processing
- `training/default.yaml` - Training hyperparameters
- `training/debug.yaml` - Quick local testing

### Key Settings

**LoRA rank (`r`)**: Higher = more capacity, more VRAM. Default 64 is a good balance.

**Quantization**: QLoRA (4-bit) enabled by default. Comment out `quantization` section in model config for full LoRA (needs ~48GB VRAM).

**Batch size**: Adjust `per_device_train_batch_size` and `gradient_accumulation_steps` based on available VRAM.

## Project Structure

```
training/
├── configs/           # YAML configuration files
├── src/trainer/       # Main package
│   ├── synthesis/     # Data processing
│   ├── training/      # Training pipeline
│   ├── evaluation/    # Evaluation metrics
│   └── export/        # Model export
├── scripts/           # CLI entry points
├── data/              # Training data (gitignored)
├── outputs/           # Checkpoints and models (gitignored)
└── tests/             # Unit tests
```

## Development

```bash
# Run tests
make test

# Run linting
make lint

# Format code
make format
```

## License

Apache 2.0
