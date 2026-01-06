"""Model and tokenizer setup."""

import torch
from peft import LoraConfig as PeftLoraConfig
from peft import get_peft_model, prepare_model_for_kbit_training
from transformers import (
    AutoModelForCausalLM,
    AutoTokenizer,
    BitsAndBytesConfig,
    PreTrainedModel,
    PreTrainedTokenizer,
)

from trainer.training.config import LoraConfig, ModelConfig, QuantizationConfig


def get_torch_dtype(dtype_str: str) -> torch.dtype:
    """Convert string dtype to torch dtype."""
    dtype_map = {
        "float32": torch.float32,
        "float16": torch.float16,
        "bfloat16": torch.bfloat16,
    }
    return dtype_map.get(dtype_str, torch.bfloat16)


def setup_tokenizer(model_config: ModelConfig) -> PreTrainedTokenizer:
    """Setup tokenizer for training."""
    tokenizer = AutoTokenizer.from_pretrained(
        model_config.name,
        trust_remote_code=model_config.trust_remote_code,
    )

    # Ensure padding token is set
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token

    tokenizer.padding_side = "right"

    return tokenizer


def setup_model(
    model_config: ModelConfig,
    quantization_config: QuantizationConfig | None,
    lora_config: LoraConfig,
) -> PreTrainedModel:
    """Setup model with quantization and LoRA."""
    torch_dtype = get_torch_dtype(model_config.torch_dtype)

    # Setup quantization if configured
    bnb_config = None
    if quantization_config:
        bnb_config = BitsAndBytesConfig(
            load_in_4bit=quantization_config.load_in_4bit,
            bnb_4bit_compute_dtype=get_torch_dtype(quantization_config.bnb_4bit_compute_dtype),
            bnb_4bit_quant_type=quantization_config.bnb_4bit_quant_type,
            bnb_4bit_use_double_quant=quantization_config.bnb_4bit_use_double_quant,
        )

    # Load base model
    load_kwargs = {
        "quantization_config": bnb_config,
        "torch_dtype": torch_dtype,
        "device_map": model_config.device_map,
        "trust_remote_code": model_config.trust_remote_code,
    }
    if model_config.attn_implementation:
        load_kwargs["attn_implementation"] = model_config.attn_implementation

    model = AutoModelForCausalLM.from_pretrained(model_config.name, **load_kwargs)

    # Prepare for k-bit training if quantized
    if quantization_config:
        model = prepare_model_for_kbit_training(model)

    # Setup LoRA
    peft_config = PeftLoraConfig(
        r=lora_config.r,
        lora_alpha=lora_config.lora_alpha,
        lora_dropout=lora_config.lora_dropout,
        bias=lora_config.bias,
        task_type=lora_config.task_type,
        target_modules=lora_config.target_modules,
    )

    model = get_peft_model(model, peft_config)

    # Print trainable parameters
    model.print_trainable_parameters()

    return model
