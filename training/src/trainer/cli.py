"""CLI entry points for trainer commands."""

import argparse
import sys


def synthesize() -> None:
    """Run data synthesis pipeline."""
    parser = argparse.ArgumentParser(description="Synthesize training data")
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

    from trainer.synthesis import run_synthesis

    run_synthesis(config_path=args.config, output_dir=args.output_dir)


def train() -> None:
    """Run training pipeline."""
    parser = argparse.ArgumentParser(description="Train model")
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

    from trainer.training import run_training

    run_training(config_path=args.config, resume_from=args.resume)


def evaluate() -> None:
    """Run evaluation pipeline."""
    parser = argparse.ArgumentParser(description="Evaluate model")
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
    args = parser.parse_args()

    from trainer.evaluation import run_evaluation

    run_evaluation(checkpoint_path=args.checkpoint, config_path=args.config)


def export() -> None:
    """Run export pipeline."""
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

    from trainer.export import run_export

    run_export(
        checkpoint_path=args.checkpoint,
        output_dir=args.output_dir,
        model_name=args.model_name,
    )


if __name__ == "__main__":
    print("Use trainer-synthesize, trainer-train, trainer-evaluate, or trainer-export")
    sys.exit(1)
