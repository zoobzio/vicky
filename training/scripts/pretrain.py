#!/usr/bin/env python3
"""Run continued pre-training."""

import argparse
import sys
from pathlib import Path

# Add src to path for development
sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from trainer.pretrain import run_pretraining


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Run continued pre-training on repository data"
    )
    parser.add_argument(
        "--config",
        type=str,
        default="configs/training/pretrain.yaml",
        help="Path to training config",
    )
    parser.add_argument(
        "--data-dir",
        type=str,
        help="Override data directory",
    )
    parser.add_argument(
        "--resume",
        type=str,
        help="Path to checkpoint to resume from",
    )
    args = parser.parse_args()

    run_pretraining(
        config_path=args.config,
        data_dir=args.data_dir,
        resume_from=args.resume,
    )


if __name__ == "__main__":
    main()
