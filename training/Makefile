.PHONY: setup install dev lint test train evaluate export clean pretrain-data pretrain

# Environment setup
setup:
	python -m venv .venv
	@echo "Run 'source .venv/bin/activate' then 'make install'"

install:
	pip install -e .

dev:
	pip install -e ".[dev]"

# Code quality
lint:
	ruff check src tests
	ruff format --check src tests
	mypy src

format:
	ruff check --fix src tests
	ruff format src tests

test:
	pytest tests -v

test-cov:
	pytest tests -v --cov=trainer --cov-report=term-missing

# Data synthesis
synthesize:
	python scripts/synthesize.py $(ARGS)

# Pre-training (continued pre-training on raw code/docs)
pretrain-data:
	python scripts/pretrain_data.py $(ARGS)

pretrain:
	python scripts/pretrain.py $(ARGS)

pretrain-local:
	python scripts/pretrain.py --config configs/training/pretrain-local.yaml $(ARGS)

# Instruction fine-tuning
train:
	python scripts/train.py $(ARGS)

train-debug:
	python scripts/train.py --config configs/training/debug.yaml $(ARGS)

train-local:
	python scripts/train.py --config configs/training/local.yaml $(ARGS)

# Evaluation
evaluate:
	python scripts/evaluate.py $(ARGS)

# Export
export:
	python scripts/export.py $(ARGS)

# Cleanup
clean:
	rm -rf outputs/checkpoints/*
	rm -rf outputs/logs/*
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete 2>/dev/null || true

clean-all: clean
	rm -rf outputs/models/*
	rm -rf data/processed/*
	rm -rf data/splits/*
	rm -rf .venv
