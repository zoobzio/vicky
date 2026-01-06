"""Training configuration loading and validation."""

from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

import yaml


@dataclass
class ModelConfig:
    """Model configuration."""

    name: str
    torch_dtype: str = "bfloat16"
    device_map: str = "auto"
    trust_remote_code: bool = True
    attn_implementation: str | None = None  # "sdpa", "flash_attention_2", or None for default


@dataclass
class QuantizationConfig:
    """Quantization configuration for QLoRA."""

    load_in_4bit: bool = True
    bnb_4bit_compute_dtype: str = "bfloat16"
    bnb_4bit_quant_type: str = "nf4"
    bnb_4bit_use_double_quant: bool = True


@dataclass
class LoraConfig:
    """LoRA configuration."""

    r: int = 64
    lora_alpha: int = 16
    lora_dropout: float = 0.1
    bias: str = "none"
    task_type: str = "CAUSAL_LM"
    target_modules: list[str] = field(
        default_factory=lambda: [
            "q_proj",
            "k_proj",
            "v_proj",
            "o_proj",
            "gate_proj",
            "up_proj",
            "down_proj",
        ]
    )


@dataclass
class TrainingConfig:
    """Complete training configuration."""

    model: ModelConfig
    quantization: QuantizationConfig | None
    lora: LoraConfig
    training: dict[str, Any]
    sft: dict[str, Any]
    wandb: dict[str, Any] | None
    data_config_path: str


def load_config(config_path: str) -> TrainingConfig:
    """Load and validate training configuration."""
    config_file = Path(config_path)
    config = yaml.safe_load(config_file.read_text())

    # Load referenced model config
    model_config_path = Path(config.get("model_config", "configs/model/qwen-7b-lora.yaml"))
    model_config = yaml.safe_load(model_config_path.read_text())

    # Build configuration objects
    model = ModelConfig(**model_config["model"])

    quantization = None
    if "quantization" in model_config:
        quantization = QuantizationConfig(**model_config["quantization"])

    lora = LoraConfig(**model_config["lora"])

    return TrainingConfig(
        model=model,
        quantization=quantization,
        lora=lora,
        training=config.get("training", {}),
        sft=config.get("sft", {}),
        wandb=config.get("wandb"),
        data_config_path=config.get("data_config", "configs/data/synthesis.yaml"),
    )
