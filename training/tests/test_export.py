"""Tests for export pipeline."""

import json
import tempfile
from pathlib import Path

import pytest

from trainer.export.merge import get_base_model_name


class TestExport:
    """Tests for export functionality."""

    def test_get_base_model_name(self, tmp_path: Path) -> None:
        """Test extracting base model name from adapter config."""
        adapter_config = {
            "base_model_name_or_path": "Qwen/Qwen2.5-7B-Instruct",
            "r": 64,
        }
        config_path = tmp_path / "adapter_config.json"
        config_path.write_text(json.dumps(adapter_config))

        result = get_base_model_name(str(tmp_path))
        assert result == "Qwen/Qwen2.5-7B-Instruct"

    def test_get_base_model_name_missing_config(self, tmp_path: Path) -> None:
        """Test error when adapter config is missing."""
        with pytest.raises(ValueError, match="No adapter_config.json"):
            get_base_model_name(str(tmp_path))
