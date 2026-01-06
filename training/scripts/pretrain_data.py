#!/usr/bin/env python3
"""Create pre-training dataset from repositories."""

import argparse
import sys
from pathlib import Path

# Add src to path for development
sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from transformers import AutoTokenizer

from trainer.pretrain import create_pretrain_dataset


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Create pre-training dataset from repository files"
    )
    parser.add_argument(
        "--config",
        type=str,
        default="configs/data/pretrain.yaml",
        help="Path to pretrain data config",
    )
    parser.add_argument(
        "--model",
        type=str,
        default="Qwen/Qwen2.5-7B-Instruct",
        help="Model name for tokenizer (used for chunking)",
    )
    parser.add_argument(
        "--output-dir",
        type=str,
        help="Override output directory",
    )
    args = parser.parse_args()

    # Load tokenizer for chunking
    print(f"Loading tokenizer from {args.model}...")
    tokenizer = AutoTokenizer.from_pretrained(args.model, trust_remote_code=True)

    # Create dataset
    create_pretrain_dataset(
        config_path=args.config,
        tokenizer=tokenizer,
        output_dir=args.output_dir,
    )


if __name__ == "__main__":
    main()
