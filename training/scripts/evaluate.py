#!/usr/bin/env python3
"""Run evaluation pipeline."""

import argparse
import sys
from pathlib import Path

# Add src to path for development
sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from trainer.evaluation import run_evaluation


def main() -> None:
    parser = argparse.ArgumentParser(description="Evaluate trained model")
    parser.add_argument(
        "--checkpoint",
        type=str,
        required=True,
        help="Path to model checkpoint",
    )
    parser.add_argument(
        "--config",
        type=str,
        default="configs/training/default.yaml",
        help="Path to config (for data settings)",
    )
    parser.add_argument(
        "--output",
        type=str,
        help="Path to save detailed results JSON",
    )
    parser.add_argument(
        "--max-examples",
        type=int,
        help="Maximum number of examples to evaluate",
    )
    args = parser.parse_args()

    run_evaluation(
        checkpoint_path=args.checkpoint,
        config_path=args.config,
        output_path=args.output,
        max_examples=args.max_examples,
    )


if __name__ == "__main__":
    main()
