"""Custom evaluation metrics."""

from dataclasses import dataclass


@dataclass
class EvaluationResult:
    """Result from evaluating a single example."""

    input_text: str
    expected: str
    generated: str
    scores: dict[str, float]


def compute_metrics(expected: str, generated: str) -> dict[str, float]:
    """
    Compute evaluation metrics for a single generation.

    Currently implements basic metrics. Extend with:
    - BLEU/ROUGE for text similarity
    - Code-specific metrics (syntax validity, etc.)
    - Custom domain metrics
    """
    metrics: dict[str, float] = {}

    # Exact match
    metrics["exact_match"] = 1.0 if expected.strip() == generated.strip() else 0.0

    # Length ratio (generated / expected)
    if len(expected) > 0:
        metrics["length_ratio"] = len(generated) / len(expected)
    else:
        metrics["length_ratio"] = 0.0

    # Token overlap (simple word overlap)
    expected_tokens = set(expected.lower().split())
    generated_tokens = set(generated.lower().split())

    if expected_tokens:
        intersection = expected_tokens & generated_tokens
        metrics["token_precision"] = len(intersection) / len(generated_tokens) if generated_tokens else 0.0
        metrics["token_recall"] = len(intersection) / len(expected_tokens)
        if metrics["token_precision"] + metrics["token_recall"] > 0:
            metrics["token_f1"] = (
                2 * metrics["token_precision"] * metrics["token_recall"]
                / (metrics["token_precision"] + metrics["token_recall"])
            )
        else:
            metrics["token_f1"] = 0.0
    else:
        metrics["token_precision"] = 0.0
        metrics["token_recall"] = 0.0
        metrics["token_f1"] = 0.0

    return metrics


def aggregate_metrics(results: list[EvaluationResult]) -> dict[str, float]:
    """Aggregate metrics across multiple results."""
    if not results:
        return {}

    # Collect all metric names
    metric_names = set()
    for result in results:
        metric_names.update(result.scores.keys())

    # Compute averages
    aggregated: dict[str, float] = {}
    for name in metric_names:
        values = [r.scores.get(name, 0.0) for r in results]
        aggregated[f"{name}_mean"] = sum(values) / len(values)

    aggregated["num_examples"] = float(len(results))
    return aggregated
