#!/usr/bin/env python3
"""Run data synthesis pipeline."""

import argparse
import sys
from pathlib import Path

# Add src to path for development
sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from trainer.synthesis import run_synthesis


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Synthesize training data from repositories and examples"
    )
    parser.add_argument(
        "--config",
        type=str,
        default="configs/data/synthesis.yaml",
        help="Path to synthesis config",
    )
    parser.add_argument(
        "--output-dir",
        type=str,
        help="Override output directory",
    )
    args = parser.parse_args()

    run_synthesis(config_path=args.config, output_dir=args.output_dir)


if __name__ == "__main__":
    main()
