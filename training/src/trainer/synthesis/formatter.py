"""Format training data into datasets format."""

import json
from dataclasses import dataclass
from pathlib import Path
from typing import Any

from datasets import Dataset


@dataclass
class ConversationMessage:
    """A single message in a conversation."""

    role: str
    content: str
    tool_calls: list[dict[str, Any]] | None = None


@dataclass
class Conversation:
    """A training conversation."""

    messages: list[ConversationMessage]
    metadata: dict[str, Any] | None = None


class DatasetFormatter:
    """Format conversations into training datasets."""

    def __init__(self, system_message: str | None = None) -> None:
        self.system_message = system_message

    def load_examples(self, input_path: Path) -> list[Conversation]:
        """
        Load training examples from a file or directory.

        Supports:
        - .jsonl files with one conversation per line
        - .json files with a list of conversations
        - Directories containing multiple .jsonl/.json files

        Expected format:
        {
            "messages": [
                {"role": "user", "content": "..."},
                {"role": "assistant", "content": "..."},
                ...
            ]
        }
        """
        conversations: list[Conversation] = []

        if input_path.is_dir():
            for file_path in input_path.glob("**/*.json*"):
                conversations.extend(self._load_file(file_path))
        else:
            conversations.extend(self._load_file(input_path))

        return conversations

    def _load_file(self, file_path: Path) -> list[Conversation]:
        """Load conversations from a single file."""
        conversations: list[Conversation] = []
        content = file_path.read_text(encoding="utf-8")

        if file_path.suffix == ".jsonl":
            for line in content.strip().split("\n"):
                if line:
                    conversations.append(self._parse_conversation(json.loads(line)))
        else:
            data = json.loads(content)
            if isinstance(data, list):
                for item in data:
                    conversations.append(self._parse_conversation(item))
            else:
                conversations.append(self._parse_conversation(data))

        return conversations

    def _parse_conversation(self, data: dict[str, Any]) -> Conversation:
        """Parse a raw conversation dict into a Conversation object."""
        messages = []
        for msg in data.get("messages", []):
            messages.append(
                ConversationMessage(
                    role=msg["role"],
                    content=msg.get("content", ""),
                    tool_calls=msg.get("tool_calls"),
                )
            )
        return Conversation(
            messages=messages,
            metadata=data.get("metadata"),
        )

    def format_for_training(
        self, conversations: list[Conversation]
    ) -> Dataset:
        """
        Format conversations into a HuggingFace Dataset.

        Converts conversations to the chat template format expected by
        the tokenizer's apply_chat_template method.
        """
        formatted_data = []

        for conv in conversations:
            messages = []

            # Add system message if configured
            if self.system_message:
                messages.append({"role": "system", "content": self.system_message})

            # Add conversation messages
            for msg in conv.messages:
                message_dict: dict[str, Any] = {
                    "role": msg.role,
                    "content": msg.content,
                }
                if msg.tool_calls:
                    message_dict["tool_calls"] = msg.tool_calls
                messages.append(message_dict)

            formatted_data.append({"messages": messages})

        return Dataset.from_list(formatted_data)

    def create_splits(
        self,
        dataset: Dataset,
        train_ratio: float = 0.9,
        val_ratio: float = 0.05,
        test_ratio: float = 0.05,
        seed: int = 42,
    ) -> dict[str, Dataset]:
        """Split dataset into train/val/test sets."""
        assert abs(train_ratio + val_ratio + test_ratio - 1.0) < 1e-6

        n = len(dataset)

        # For small datasets, ensure minimum 1 example per split
        if n < 10:
            # Very small: use all for training, duplicate for val/test
            print(f"  Warning: Small dataset ({n} examples), using all for train, duplicating for val/test")
            return {
                "train": dataset,
                "validation": dataset.select(range(min(1, n))),
                "test": dataset.select(range(min(1, n))),
            }

        # Calculate split sizes ensuring minimums
        test_size = max(1, int(n * test_ratio))
        val_size = max(1, int(n * val_ratio))
        train_size = n - test_size - val_size

        # First split: train vs (val + test)
        train_test = dataset.train_test_split(
            test_size=(val_size + test_size),
            seed=seed,
        )

        # Second split: val vs test
        if len(train_test["test"]) >= 2:
            val_test = train_test["test"].train_test_split(
                test_size=test_size,
                seed=seed,
            )
            return {
                "train": train_test["train"],
                "validation": val_test["train"],
                "test": val_test["test"],
            }
        else:
            # Edge case: not enough for separate val/test
            return {
                "train": train_test["train"],
                "validation": train_test["test"],
                "test": train_test["test"],
            }
