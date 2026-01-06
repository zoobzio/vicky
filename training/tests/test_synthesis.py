"""Tests for data synthesis pipeline."""

import json
import tempfile
from pathlib import Path

import pytest

from trainer.synthesis.formatter import Conversation, ConversationMessage, DatasetFormatter


class TestDatasetFormatter:
    """Tests for DatasetFormatter."""

    def test_load_jsonl_examples(self, tmp_path: Path) -> None:
        """Test loading examples from JSONL file."""
        # Create test data
        examples = [
            {
                "messages": [
                    {"role": "user", "content": "What is Go?"},
                    {"role": "assistant", "content": "Go is a programming language."},
                ]
            },
            {
                "messages": [
                    {"role": "user", "content": "How do I declare a variable?"},
                    {"role": "assistant", "content": "Use var or :="},
                ]
            },
        ]

        # Write to file
        jsonl_file = tmp_path / "examples.jsonl"
        with open(jsonl_file, "w") as f:
            for example in examples:
                f.write(json.dumps(example) + "\n")

        # Load and verify
        formatter = DatasetFormatter()
        conversations = formatter.load_examples(jsonl_file)

        assert len(conversations) == 2
        assert conversations[0].messages[0].role == "user"
        assert conversations[0].messages[0].content == "What is Go?"

    def test_load_json_examples(self, tmp_path: Path) -> None:
        """Test loading examples from JSON file."""
        examples = [
            {
                "messages": [
                    {"role": "user", "content": "Question?"},
                    {"role": "assistant", "content": "Answer."},
                ]
            }
        ]

        json_file = tmp_path / "examples.json"
        json_file.write_text(json.dumps(examples))

        formatter = DatasetFormatter()
        conversations = formatter.load_examples(json_file)

        assert len(conversations) == 1

    def test_format_for_training(self) -> None:
        """Test formatting conversations for training."""
        conversations = [
            Conversation(
                messages=[
                    ConversationMessage(role="user", content="Hello"),
                    ConversationMessage(role="assistant", content="Hi there"),
                ]
            )
        ]

        formatter = DatasetFormatter(system_message="You are helpful.")
        dataset = formatter.format_for_training(conversations)

        assert len(dataset) == 1
        assert dataset[0]["messages"][0]["role"] == "system"
        assert dataset[0]["messages"][1]["role"] == "user"

    def test_create_splits(self) -> None:
        """Test dataset splitting."""
        conversations = [
            Conversation(
                messages=[
                    ConversationMessage(role="user", content=f"Q{i}"),
                    ConversationMessage(role="assistant", content=f"A{i}"),
                ]
            )
            for i in range(100)
        ]

        formatter = DatasetFormatter()
        dataset = formatter.format_for_training(conversations)
        splits = formatter.create_splits(dataset, train_ratio=0.8, val_ratio=0.1, test_ratio=0.1)

        assert "train" in splits
        assert "validation" in splits
        assert "test" in splits
        assert len(splits["train"]) == 80
        assert len(splits["validation"]) == 10
        assert len(splits["test"]) == 10
