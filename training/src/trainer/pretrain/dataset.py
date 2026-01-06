"""Create pre-training dataset from extracted files."""

from pathlib import Path
from typing import Iterator

import yaml
from datasets import Dataset
from transformers import PreTrainedTokenizer

from trainer.pretrain.extractor import ExtractedFile, RepoExtractor


def create_pretrain_dataset(
    config_path: str,
    tokenizer: PreTrainedTokenizer,
    output_dir: str | None = None,
) -> Dataset:
    """
    Create a pre-training dataset from repository files.

    Args:
        config_path: Path to pretrain config YAML
        tokenizer: Tokenizer for chunking
        output_dir: Override output directory

    Returns:
        HuggingFace Dataset ready for pre-training
    """
    config = yaml.safe_load(Path(config_path).read_text())

    # Setup extractor
    extractor = RepoExtractor(
        include_patterns=config["files"]["include"],
        exclude_patterns=config["files"]["exclude"],
        max_size_bytes=config["files"]["max_size_bytes"],
    )

    # Extract all files
    all_files: list[ExtractedFile] = []
    for repo_config in config["sources"]["repositories"]:
        print(f"Extracting from {repo_config['url']}...")
        files = list(extractor.extract(
            repo_url=repo_config["url"],
            branch=repo_config.get("branch", "main"),
        ))
        print(f"  Extracted {len(files)} files")
        all_files.extend(files)

    print(f"Total files extracted: {len(all_files)}")

    if not all_files:
        raise ValueError("No files extracted. Check your repository configuration.")

    # Format files into documents
    file_separator = config["format"]["file_separator"]
    include_path = config["format"]["include_path"]

    documents = []
    for f in all_files:
        if include_path:
            doc = file_separator.format(path=f.path) + f.content
        else:
            doc = f.content
        documents.append(doc)

    # Concatenate all documents
    full_text = "\n".join(documents)
    print(f"Total text length: {len(full_text):,} characters")

    # Chunk the text
    chunk_size = config["output"]["chunk_size"]
    chunk_overlap = config["output"]["chunk_overlap"]

    chunks = list(chunk_text(
        full_text,
        tokenizer,
        chunk_size=chunk_size,
        overlap=chunk_overlap,
    ))
    print(f"Created {len(chunks)} chunks")

    # Create dataset
    dataset = Dataset.from_dict({"text": chunks})

    # Save if output_dir specified
    out_dir = Path(output_dir or config["output"]["dir"])
    out_dir.mkdir(parents=True, exist_ok=True)
    dataset.save_to_disk(str(out_dir / "dataset"))
    print(f"Saved dataset to {out_dir / 'dataset'}")

    # Save metadata
    metadata = {
        "num_files": len(all_files),
        "num_chunks": len(chunks),
        "total_chars": len(full_text),
        "chunk_size": chunk_size,
        "chunk_overlap": chunk_overlap,
        "repositories": [r["url"] for r in config["sources"]["repositories"]],
    }
    (out_dir / "metadata.yaml").write_text(yaml.dump(metadata))

    return dataset


def chunk_text(
    text: str,
    tokenizer: PreTrainedTokenizer,
    chunk_size: int = 2048,
    overlap: int = 256,
) -> Iterator[str]:
    """
    Chunk text into pieces of approximately chunk_size tokens.

    Uses the tokenizer to ensure accurate token counts.
    Chunks overlap to maintain context continuity.
    """
    # Tokenize full text
    tokens = tokenizer.encode(text, add_special_tokens=False)
    total_tokens = len(tokens)

    if total_tokens == 0:
        return

    print(f"Total tokens: {total_tokens:,}")

    # Generate chunks with overlap
    start = 0
    while start < total_tokens:
        end = min(start + chunk_size, total_tokens)
        chunk_tokens = tokens[start:end]

        # Decode back to text
        chunk_text = tokenizer.decode(chunk_tokens, skip_special_tokens=True)
        yield chunk_text

        # Move start, accounting for overlap
        start = end - overlap
        if start >= total_tokens - overlap:
            break  # Avoid tiny final chunks
