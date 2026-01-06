#!/usr/bin/env python3
"""Run export pipeline."""

import argparse
import sys
from pathlib import Path

# Add src to path for development
sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from trainer.export import run_export


def main() -> None:
    parser = argparse.ArgumentParser(description="Export model to safetensors")
    parser.add_argument(
        "--checkpoint",
        type=str,
        required=True,
        help="Path to LoRA checkpoint",
    )
    parser.add_argument(
        "--output-dir",
        type=str,
        default="outputs/models",
        help="Output directory for merged model",
    )
    parser.add_argument(
        "--model-name",
        type=str,
        help="Name for exported model (default: auto-generated)",
    )
    args = parser.parse_args()

    run_export(
        checkpoint_path=args.checkpoint,
        output_dir=args.output_dir,
        model_name=args.model_name,
    )


if __name__ == "__main__":
    main()
