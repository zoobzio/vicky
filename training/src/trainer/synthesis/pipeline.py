"""Main synthesis pipeline orchestration."""

from pathlib import Path

import yaml

from trainer.synthesis.extractor import GitHubExtractor
from trainer.synthesis.formatter import DatasetFormatter


def run_synthesis(
    config_path: str,
    output_dir: str | None = None,
) -> None:
    """
    Run the data synthesis pipeline.

    Args:
        config_path: Path to synthesis config YAML
        output_dir: Override output directory from config
    """
    # Load config
    config = yaml.safe_load(Path(config_path).read_text())

    # Setup paths
    processed_dir = Path(output_dir or config["output"]["processed_dir"])
    splits_dir = Path(config["output"]["splits_dir"])
    processed_dir.mkdir(parents=True, exist_ok=True)
    splits_dir.mkdir(parents=True, exist_ok=True)

    # Extract from repositories if configured
    extractor = GitHubExtractor()
    for repo_config in config["sources"].get("repositories", []):
        print(f"Extracting from {repo_config['url']}...")
        for extracted_file in extractor.extract(
            repo_url=repo_config["url"],
            branch=repo_config.get("branch", "main"),
            include_patterns=repo_config.get("include"),
            exclude_patterns=repo_config.get("exclude"),
        ):
            # Save extracted files for reference
            out_path = processed_dir / "extracted" / extracted_file.path
            out_path.parent.mkdir(parents=True, exist_ok=True)
            out_path.write_text(extracted_file.content)
        print(f"  Extracted to {processed_dir / 'extracted'}")

    # Format training examples
    formatter = DatasetFormatter(
        system_message=config["format"].get("system_message"),
    )

    examples_dir = Path(config["input"]["examples_dir"])
    if examples_dir.exists():
        print(f"Loading examples from {examples_dir}...")
        conversations = formatter.load_examples(examples_dir)
        print(f"  Loaded {len(conversations)} conversations")

        if conversations:
            # Format and split
            dataset = formatter.format_for_training(conversations)
            splits = formatter.create_splits(
                dataset,
                train_ratio=config["output"]["train_ratio"],
                val_ratio=config["output"]["val_ratio"],
                test_ratio=config["output"]["test_ratio"],
                seed=config["output"]["seed"],
            )

            # Save splits
            for split_name, split_data in splits.items():
                split_path = splits_dir / split_name
                split_data.save_to_disk(str(split_path))
                print(f"  Saved {split_name}: {len(split_data)} examples to {split_path}")
    else:
        print(f"No examples directory found at {examples_dir}")
        print("Create training examples in JSONL format and place them there.")

    print("Synthesis complete.")
