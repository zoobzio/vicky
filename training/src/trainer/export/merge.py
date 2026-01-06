"""Merge LoRA adapters into base model."""

from pathlib import Path

import torch
from peft import PeftModel
from transformers import AutoModelForCausalLM, AutoTokenizer, PreTrainedModel


def merge_lora_weights(
    base_model_name: str,
    adapter_path: str,
    trust_remote_code: bool = True,
) -> tuple[PreTrainedModel, AutoTokenizer]:
    """
    Merge LoRA adapter weights into base model.

    Args:
        base_model_name: HuggingFace model name/path for base model
        adapter_path: Path to LoRA adapter checkpoint
        trust_remote_code: Whether to trust remote code for model loading

    Returns:
        Tuple of (merged_model, tokenizer)
    """
    print(f"Loading base model: {base_model_name}")
    base_model = AutoModelForCausalLM.from_pretrained(
        base_model_name,
        torch_dtype=torch.bfloat16,
        device_map="cpu",  # Load on CPU for merging
        trust_remote_code=trust_remote_code,
    )

    print(f"Loading LoRA adapter: {adapter_path}")
    model = PeftModel.from_pretrained(base_model, adapter_path)

    print("Merging weights...")
    merged_model = model.merge_and_unload()

    print("Loading tokenizer...")
    tokenizer = AutoTokenizer.from_pretrained(
        base_model_name,
        trust_remote_code=trust_remote_code,
    )

    return merged_model, tokenizer


def get_base_model_name(adapter_path: str) -> str:
    """Extract base model name from adapter config."""
    import json

    config_path = Path(adapter_path) / "adapter_config.json"
    if config_path.exists():
        config = json.loads(config_path.read_text())
        return config.get("base_model_name_or_path", "")
    raise ValueError(f"No adapter_config.json found at {adapter_path}")
