"""Extract all files from repositories for pre-training."""

import fnmatch
import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path
from typing import Iterator


@dataclass
class ExtractedFile:
    """A file extracted from a repository."""

    repo_url: str
    path: str
    content: str
    size_bytes: int


class RepoExtractor:
    """Extract files from GitHub repositories for pre-training."""

    def __init__(
        self,
        cache_dir: Path | None = None,
        include_patterns: list[str] | None = None,
        exclude_patterns: list[str] | None = None,
        max_size_bytes: int = 102400,
    ) -> None:
        self.cache_dir = cache_dir or Path(tempfile.gettempdir()) / "trainer-cache"
        self.cache_dir.mkdir(parents=True, exist_ok=True)
        self.include_patterns = include_patterns or ["**/*"]
        self.exclude_patterns = exclude_patterns or []
        self.max_size_bytes = max_size_bytes

    def extract(
        self,
        repo_url: str,
        branch: str = "main",
    ) -> Iterator[ExtractedFile]:
        """
        Extract all matching files from a repository.

        Args:
            repo_url: GitHub repository URL
            branch: Branch to extract from

        Yields:
            ExtractedFile objects for each matching file
        """
        repo_dir = self._clone_or_update(repo_url, branch)

        # Walk all files
        for file_path in repo_dir.rglob("*"):
            if not file_path.is_file():
                continue

            rel_path = file_path.relative_to(repo_dir)

            # Check exclusions first
            if self._is_excluded(rel_path):
                continue

            # Check inclusions
            if not self._is_included(rel_path):
                continue

            # Check file size
            size = file_path.stat().st_size
            if size > self.max_size_bytes:
                continue

            # Try to read as text
            try:
                content = file_path.read_text(encoding="utf-8")
                # Skip files that look binary (high ratio of non-printable chars)
                if self._looks_binary(content):
                    continue

                yield ExtractedFile(
                    repo_url=repo_url,
                    path=str(rel_path),
                    content=content,
                    size_bytes=size,
                )
            except (UnicodeDecodeError, PermissionError):
                # Skip binary or unreadable files
                continue

    def _clone_or_update(self, repo_url: str, branch: str) -> Path:
        """Clone a repository or update if already cached."""
        import os

        repo_name = repo_url.rstrip("/").split("/")[-1].replace(".git", "")
        org_name = repo_url.rstrip("/").split("/")[-2]
        repo_dir = self.cache_dir / f"{org_name}_{repo_name}"

        # Bypass user git config (e.g., URL rewrites) to ensure HTTPS works
        env = {**os.environ, "GIT_CONFIG_GLOBAL": "/dev/null", "GIT_CONFIG_SYSTEM": "/dev/null"}

        if repo_dir.exists():
            subprocess.run(
                ["git", "fetch", "origin", branch],
                cwd=repo_dir,
                check=True,
                capture_output=True,
                env=env,
            )
            subprocess.run(
                ["git", "checkout", f"origin/{branch}"],
                cwd=repo_dir,
                check=True,
                capture_output=True,
                env=env,
            )
        else:
            subprocess.run(
                [
                    "git",
                    "clone",
                    "--depth",
                    "1",
                    "--branch",
                    branch,
                    repo_url,
                    str(repo_dir),
                ],
                check=True,
                capture_output=True,
                env=env,
            )

        return repo_dir

    def _is_included(self, rel_path: Path) -> bool:
        """Check if file matches any include pattern."""
        path_str = str(rel_path)
        for pattern in self.include_patterns:
            if fnmatch.fnmatch(path_str, pattern):
                return True
            # Also check with ** prefix for nested matches
            if fnmatch.fnmatch(path_str, pattern.lstrip("*/")):
                return True
        return False

    def _is_excluded(self, rel_path: Path) -> bool:
        """Check if file matches any exclude pattern."""
        path_str = str(rel_path)
        for pattern in self.exclude_patterns:
            if fnmatch.fnmatch(path_str, pattern):
                return True
            # Check each path component for directory exclusions
            for part in rel_path.parts:
                if fnmatch.fnmatch(part, pattern.strip("*/")):
                    return True
        return False

    def _looks_binary(self, content: str, threshold: float = 0.1) -> bool:
        """Check if content looks like binary data."""
        if not content:
            return False
        # Count non-printable characters (excluding common whitespace)
        non_printable = sum(
            1 for c in content[:1000]  # Check first 1000 chars
            if not c.isprintable() and c not in "\n\r\t"
        )
        return (non_printable / min(len(content), 1000)) > threshold
