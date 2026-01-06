"""Evaluation pipeline."""

from trainer.evaluation.metrics import compute_metrics
from trainer.evaluation.runner import run_evaluation

__all__ = ["compute_metrics", "run_evaluation"]
