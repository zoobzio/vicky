#!/usr/bin/env python3
"""Run training pipeline."""

import argparse
import sys
from pathlib import Path

# Add src to path for development
sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from trainer.training import run_training


def main() -> None:
    parser = argparse.ArgumentParser(description="Train model with LoRA")
    parser.add_argument(
        "--config",
        type=str,
        default="configs/training/default.yaml",
        help="Path to training config",
    )
    parser.add_argument(
        "--resume",
        type=str,
        help="Path to checkpoint to resume from",
    )
    args = parser.parse_args()

    run_training(config_path=args.config, resume_from=args.resume)


if __name__ == "__main__":
    main()
