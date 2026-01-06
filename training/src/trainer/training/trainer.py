"""Training loop and SFTTrainer wrapper."""

from pathlib import Path

import yaml
from datasets import load_from_disk
from trl import SFTConfig, SFTTrainer

from trainer.training.config import load_config
from trainer.training.model import setup_model, setup_tokenizer


def run_training(
    config_path: str,
    resume_from: str | None = None,
) -> None:
    """
    Run the training pipeline.

    Args:
        config_path: Path to training config YAML
        resume_from: Optional checkpoint path to resume from
    """
    print(f"Loading config from {config_path}...")
    config = load_config(config_path)

    # Setup wandb if configured
    if config.wandb:
        import wandb

        wandb.init(
            project=config.wandb.get("project", "trainer"),
            name=config.wandb.get("name"),
            tags=config.wandb.get("tags", []),
        )

    # Load data config for paths
    data_config = yaml.safe_load(Path(config.data_config_path).read_text())
    splits_dir = Path(data_config["output"]["splits_dir"])

    # Load datasets
    print("Loading datasets...")
    train_dataset = load_from_disk(str(splits_dir / "train"))
    eval_dataset = load_from_disk(str(splits_dir / "validation"))
    print(f"  Train: {len(train_dataset)} examples")
    print(f"  Eval: {len(eval_dataset)} examples")

    # Setup model and tokenizer
    print("Setting up model and tokenizer...")
    tokenizer = setup_tokenizer(config.model)
    model = setup_model(config.model, config.quantization, config.lora)

    # Formatting function for chat template
    def formatting_func(example: dict) -> str:
        """Format a single example using chat template."""
        messages = example["messages"]
        text = tokenizer.apply_chat_template(
            messages,
            tokenize=False,
            add_generation_prompt=False,
        )
        return text

    # Setup SFT config (combines training args + SFT-specific args)
    sft_config = SFTConfig(
        output_dir=config.training.get("output_dir", "outputs/checkpoints"),
        num_train_epochs=config.training.get("num_train_epochs", 3),
        per_device_train_batch_size=config.training.get("per_device_train_batch_size", 4),
        per_device_eval_batch_size=config.training.get("per_device_eval_batch_size", 4),
        gradient_accumulation_steps=config.training.get("gradient_accumulation_steps", 4),
        learning_rate=config.training.get("learning_rate", 2e-4),
        weight_decay=config.training.get("weight_decay", 0.01),
        warmup_ratio=config.training.get("warmup_ratio", 0.03),
        lr_scheduler_type=config.training.get("lr_scheduler_type", "cosine"),
        max_grad_norm=config.training.get("max_grad_norm", 0.3),
        eval_strategy=config.training.get("eval_strategy", "steps"),
        eval_steps=config.training.get("eval_steps", 100),
        save_strategy=config.training.get("save_strategy", "steps"),
        save_steps=config.training.get("save_steps", 100),
        save_total_limit=config.training.get("save_total_limit", 3),
        logging_dir=config.training.get("logging_dir", "outputs/logs"),
        logging_steps=config.training.get("logging_steps", 10),
        report_to=config.training.get("report_to", "wandb"),
        gradient_checkpointing=config.training.get("gradient_checkpointing", True),
        optim=config.training.get("optim", "paged_adamw_32bit"),
        bf16=config.training.get("bf16", True),
        tf32=config.training.get("tf32", True),
        remove_unused_columns=False,
        # SFT-specific settings
        max_length=config.sft.get("max_seq_length", 4096),
        packing=config.sft.get("packing", False),
        dataset_text_field=config.sft.get("dataset_text_field", "text"),
    )

    # Setup trainer
    print("Initializing trainer...")
    trainer = SFTTrainer(
        model=model,
        args=sft_config,
        train_dataset=train_dataset,
        eval_dataset=eval_dataset,
        processing_class=tokenizer,
        formatting_func=formatting_func,
    )

    # Train
    print("Starting training...")
    trainer.train(resume_from_checkpoint=resume_from)

    # Save final checkpoint
    print("Saving final model...")
    trainer.save_model()

    if config.wandb:
        wandb.finish()

    print("Training complete.")
