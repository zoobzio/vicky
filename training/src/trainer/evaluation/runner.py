"""Evaluation runner."""

import json
from pathlib import Path

import torch
import yaml
from datasets import load_from_disk
from peft import PeftModel
from transformers import AutoModelForCausalLM, AutoTokenizer

from trainer.evaluation.metrics import EvaluationResult, aggregate_metrics, compute_metrics
from trainer.training.config import load_config


def run_evaluation(
    checkpoint_path: str,
    config_path: str,
    output_path: str | None = None,
    max_examples: int | None = None,
) -> dict[str, float]:
    """
    Run evaluation on the test set.

    Args:
        checkpoint_path: Path to model checkpoint (LoRA adapter)
        config_path: Path to training config (for model and data settings)
        output_path: Optional path to save detailed results
        max_examples: Maximum number of examples to evaluate

    Returns:
        Aggregated metrics dictionary
    """
    print(f"Loading config from {config_path}...")
    config = load_config(config_path)

    # Load data config
    data_config = yaml.safe_load(Path(config.data_config_path).read_text())
    splits_dir = Path(data_config["output"]["splits_dir"])

    # Load test dataset
    print("Loading test dataset...")
    test_dataset = load_from_disk(str(splits_dir / "test"))
    if max_examples:
        test_dataset = test_dataset.select(range(min(max_examples, len(test_dataset))))
    print(f"  Evaluating on {len(test_dataset)} examples")

    # Load tokenizer
    print("Loading tokenizer...")
    tokenizer = AutoTokenizer.from_pretrained(
        config.model.name,
        trust_remote_code=config.model.trust_remote_code,
    )
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token

    # Load model with LoRA adapter
    print("Loading model...")
    base_model = AutoModelForCausalLM.from_pretrained(
        config.model.name,
        torch_dtype=torch.bfloat16,
        device_map="auto",
        trust_remote_code=config.model.trust_remote_code,
    )
    model = PeftModel.from_pretrained(base_model, checkpoint_path)
    model.eval()

    # Run evaluation
    print("Running evaluation...")
    results: list[EvaluationResult] = []

    for i, example in enumerate(test_dataset):
        messages = example["messages"]

        # Find the last user message and expected assistant response
        user_messages = []
        expected_response = ""

        for msg in messages:
            if msg["role"] in ("user", "system"):
                user_messages.append(msg)
            elif msg["role"] == "assistant":
                expected_response = msg["content"]

        if not user_messages or not expected_response:
            continue

        # Generate response
        input_text = tokenizer.apply_chat_template(
            user_messages,
            tokenize=False,
            add_generation_prompt=True,
        )
        inputs = tokenizer(input_text, return_tensors="pt").to(model.device)

        with torch.no_grad():
            outputs = model.generate(
                **inputs,
                max_new_tokens=512,
                do_sample=False,
                pad_token_id=tokenizer.pad_token_id,
            )

        generated_text = tokenizer.decode(
            outputs[0][inputs["input_ids"].shape[1] :],
            skip_special_tokens=True,
        )

        # Compute metrics
        scores = compute_metrics(expected_response, generated_text)

        results.append(
            EvaluationResult(
                input_text=input_text,
                expected=expected_response,
                generated=generated_text,
                scores=scores,
            )
        )

        if (i + 1) % 10 == 0:
            print(f"  Processed {i + 1}/{len(test_dataset)} examples")

    # Aggregate metrics
    aggregated = aggregate_metrics(results)

    print("\nResults:")
    for name, value in sorted(aggregated.items()):
        print(f"  {name}: {value:.4f}")

    # Save detailed results if requested
    if output_path:
        output_file = Path(output_path)
        output_file.parent.mkdir(parents=True, exist_ok=True)
        with open(output_file, "w") as f:
            json.dump(
                {
                    "aggregated": aggregated,
                    "examples": [
                        {
                            "input": r.input_text,
                            "expected": r.expected,
                            "generated": r.generated,
                            "scores": r.scores,
                        }
                        for r in results
                    ],
                },
                f,
                indent=2,
            )
        print(f"\nDetailed results saved to {output_path}")

    return aggregated
