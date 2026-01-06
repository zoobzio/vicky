"""Convert and export models to safetensors format."""

from datetime import datetime
from pathlib import Path

from safetensors.torch import save_model

from trainer.export.merge import get_base_model_name, merge_lora_weights


def run_export(
    checkpoint_path: str,
    output_dir: str = "outputs/models",
    model_name: str | None = None,
) -> str:
    """
    Export a trained model to safetensors format.

    Args:
        checkpoint_path: Path to LoRA checkpoint
        output_dir: Directory to save exported model
        model_name: Name for the exported model (auto-generated if None)

    Returns:
        Path to exported model directory
    """
    # Get base model name from adapter config
    base_model_name = get_base_model_name(checkpoint_path)
    if not base_model_name:
        raise ValueError(
            "Could not determine base model. "
            "Ensure adapter_config.json contains base_model_name_or_path"
        )

    # Merge weights
    merged_model, tokenizer = merge_lora_weights(
        base_model_name=base_model_name,
        adapter_path=checkpoint_path,
    )

    # Generate model name if not provided
    if model_name is None:
        timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
        base_short = base_model_name.split("/")[-1]
        model_name = f"{base_short}-finetuned-{timestamp}"

    # Create output directory
    output_path = Path(output_dir) / model_name
    output_path.mkdir(parents=True, exist_ok=True)

    # Save model in safetensors format
    print(f"Saving model to {output_path}...")

    # Save using HuggingFace's save_pretrained which uses safetensors by default
    merged_model.save_pretrained(
        output_path,
        safe_serialization=True,
    )

    # Save tokenizer
    tokenizer.save_pretrained(output_path)

    # Save model card
    model_card = f"""---
base_model: {base_model_name}
library_name: transformers
license: apache-2.0
tags:
- finetuned
- go
- documentation
---

# {model_name}

This model was finetuned from [{base_model_name}](https://huggingface.co/{base_model_name})
for Go framework documentation assistance.

## Usage

```python
from transformers import AutoModelForCausalLM, AutoTokenizer

model = AutoModelForCausalLM.from_pretrained("{model_name}")
tokenizer = AutoTokenizer.from_pretrained("{model_name}")
```

## Training

Finetuned using LoRA with the trainer pipeline.
"""

    (output_path / "README.md").write_text(model_card)

    print(f"Export complete: {output_path}")
    print(f"  - Model weights: {output_path / 'model.safetensors'}")
    print(f"  - Tokenizer: {output_path / 'tokenizer.json'}")
    print(f"  - Config: {output_path / 'config.json'}")

    return str(output_path)
