"""Continued pre-training trainer."""

from pathlib import Path

import yaml
from datasets import load_from_disk
from trl import SFTConfig, SFTTrainer

from trainer.training.config import load_config
from trainer.training.model import setup_model, setup_tokenizer


def run_pretraining(
    config_path: str,
    data_dir: str | None = None,
    resume_from: str | None = None,
) -> None:
    """
    Run continued pre-training on raw text data.

    This uses standard causal language modeling (next token prediction)
    rather than instruction tuning.

    Args:
        config_path: Path to training config YAML
        data_dir: Override data directory (default: from pretrain config)
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
            tags=config.wandb.get("tags", []) + ["pretrain"],
        )

    # Load pretrain data config to get data path
    pretrain_config = yaml.safe_load(Path("configs/data/pretrain.yaml").read_text())
    data_path = Path(data_dir or pretrain_config["output"]["dir"]) / "dataset"

    # Load dataset
    print(f"Loading dataset from {data_path}...")
    dataset = load_from_disk(str(data_path))
    print(f"  Loaded {len(dataset)} chunks")

    # Split for eval (use small portion)
    split = dataset.train_test_split(test_size=0.02, seed=42)
    train_dataset = split["train"]
    eval_dataset = split["test"]
    print(f"  Train: {len(train_dataset)}, Eval: {len(eval_dataset)}")

    # Setup model and tokenizer
    print("Setting up model and tokenizer...")
    tokenizer = setup_tokenizer(config.model)
    model = setup_model(config.model, config.quantization, config.lora)

    # For pre-training, we don't need a formatting function
    # The data is already plain text in the "text" field

    # Setup SFT config for pre-training
    # (SFTTrainer with plain text is effectively continued pre-training)
    sft_config = SFTConfig(
        output_dir=config.training.get("output_dir", "outputs/checkpoints"),
        num_train_epochs=config.training.get("num_train_epochs", 1),
        per_device_train_batch_size=config.training.get("per_device_train_batch_size", 1),
        per_device_eval_batch_size=config.training.get("per_device_eval_batch_size", 1),
        gradient_accumulation_steps=config.training.get("gradient_accumulation_steps", 8),
        learning_rate=config.training.get("learning_rate", 1e-4),
        weight_decay=config.training.get("weight_decay", 0.01),
        warmup_ratio=config.training.get("warmup_ratio", 0.03),
        lr_scheduler_type=config.training.get("lr_scheduler_type", "cosine"),
        max_grad_norm=config.training.get("max_grad_norm", 1.0),
        eval_strategy=config.training.get("eval_strategy", "steps"),
        eval_steps=config.training.get("eval_steps", 100),
        save_strategy=config.training.get("save_strategy", "steps"),
        save_steps=config.training.get("save_steps", 200),
        save_total_limit=config.training.get("save_total_limit", 3),
        logging_dir=config.training.get("logging_dir", "outputs/logs"),
        logging_steps=config.training.get("logging_steps", 10),
        report_to=config.training.get("report_to", "none"),
        gradient_checkpointing=config.training.get("gradient_checkpointing", True),
        optim=config.training.get("optim", "paged_adamw_32bit"),
        bf16=config.training.get("bf16", True),
        tf32=config.training.get("tf32", True),
        # Pre-training specific
        max_length=config.sft.get("max_seq_length", 2048),
        packing=True,  # Pack sequences for efficiency
        dataset_text_field="text",
    )

    # Setup trainer
    print("Initializing trainer...")
    trainer = SFTTrainer(
        model=model,
        args=sft_config,
        train_dataset=train_dataset,
        eval_dataset=eval_dataset,
        processing_class=tokenizer,
    )

    # Train
    print("Starting pre-training...")
    trainer.train(resume_from_checkpoint=resume_from)

    # Save final checkpoint
    print("Saving final model...")
    trainer.save_model()

    if config.wandb:
        wandb.finish()

    print("Pre-training complete.")
