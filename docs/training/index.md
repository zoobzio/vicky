# Training

Documentation for fine-tuning Vicky, the documentation assistant.

## Overview

Vicky is trained in two stages:

1. **Continued pre-training** - Raw code and documentation teaches the model what the framework looks like
2. **Instruction fine-tuning** - Q&A examples teach the model how to help developers

Each package contributes training examples that are aggregated during the training pipeline.

## Guides

- [Character Profile](character.md) - Who Vicky is
- [Voice Guide](voice.md) - How Vicky's character manifests in language
- [Data Guidelines](data.md) - Format, structure, and quantity requirements for training examples
- Tools Guide (TODO) - What tools Vicky has access to and when to use them

## Quick Reference

**Per-package target:** 100-200 examples (multi-turn preferred)

**Location in each repo:**
```
<package>/
└── training/
    └── examples.jsonl
```

**Format:**
```json
{"messages": [{"role": "user", "content": "..."}, {"role": "assistant", "content": "..."}, ...]}
```

**Voice:** Direct, engaged on-topic, disengaged off-topic. No filler, no sycophancy.
