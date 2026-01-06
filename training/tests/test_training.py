"""Tests for training pipeline."""

import tempfile
from pathlib import Path

import pytest
import yaml

from trainer.training.config import (
    LoraConfig,
    ModelConfig,
    QuantizationConfig,
    TrainingConfig,
    load_config,
)


class TestConfig:
    """Tests for configuration loading."""

    def test_load_model_config(self) -> None:
        """Test ModelConfig creation."""
        config = ModelConfig(
            name="test/model",
            torch_dtype="bfloat16",
            device_map="auto",
        )
        assert config.name == "test/model"
        assert config.torch_dtype == "bfloat16"

    def test_load_lora_config_defaults(self) -> None:
        """Test LoraConfig defaults."""
        config = LoraConfig()
        assert config.r == 64
        assert config.lora_alpha == 16
        assert "q_proj" in config.target_modules

    def test_load_quantization_config(self) -> None:
        """Test QuantizationConfig creation."""
        config = QuantizationConfig(
            load_in_4bit=True,
            bnb_4bit_quant_type="nf4",
        )
        assert config.load_in_4bit is True
        assert config.bnb_4bit_quant_type == "nf4"

    def test_load_config_from_file(self, tmp_path: Path) -> None:
        """Test loading complete config from YAML files."""
        # Create model config
        model_config = {
            "model": {
                "name": "test/model",
                "torch_dtype": "bfloat16",
            },
            "lora": {
                "r": 32,
                "lora_alpha": 16,
            },
        }
        model_config_path = tmp_path / "model.yaml"
        model_config_path.write_text(yaml.dump(model_config))

        # Create training config
        training_config = {
            "model_config": str(model_config_path),
            "training": {
                "num_train_epochs": 3,
                "learning_rate": 2e-4,
            },
            "sft": {
                "max_seq_length": 4096,
            },
        }
        training_config_path = tmp_path / "training.yaml"
        training_config_path.write_text(yaml.dump(training_config))

        # Load and verify
        config = load_config(str(training_config_path))
        assert config.model.name == "test/model"
        assert config.lora.r == 32
        assert config.training["num_train_epochs"] == 3
