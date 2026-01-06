"""Extract code and documentation from GitHub repositories."""

from dataclasses import dataclass
from pathlib import Path
import subprocess
import tempfile
from typing import Iterator


@dataclass
class ExtractedFile:
    """A file extracted from a repository."""

    repo_url: str
    path: str
    content: str


class GitHubExtractor:
    """Extract files from GitHub repositories."""

    def __init__(self, cache_dir: Path | None = None) -> None:
        self.cache_dir = cache_dir or Path(tempfile.gettempdir()) / "trainer-cache"
        self.cache_dir.mkdir(parents=True, exist_ok=True)

    def extract(
        self,
        repo_url: str,
        branch: str = "main",
        include_patterns: list[str] | None = None,
        exclude_patterns: list[str] | None = None,
    ) -> Iterator[ExtractedFile]:
        """
        Extract files from a GitHub repository.

        Args:
            repo_url: GitHub repository URL
            branch: Branch to extract from
            include_patterns: Glob patterns to include (e.g., ["**/*.go", "**/*.md"])
            exclude_patterns: Glob patterns to exclude (e.g., ["**/vendor/**"])

        Yields:
            ExtractedFile objects for each matching file
        """
        include_patterns = include_patterns or ["**/*.go", "**/*.md"]
        exclude_patterns = exclude_patterns or []

        # Clone or update repository
        repo_dir = self._clone_or_update(repo_url, branch)

        # Find matching files
        for pattern in include_patterns:
            for file_path in repo_dir.glob(pattern):
                if file_path.is_file() and not self._is_excluded(
                    file_path, repo_dir, exclude_patterns
                ):
                    try:
                        content = file_path.read_text(encoding="utf-8")
                        yield ExtractedFile(
                            repo_url=repo_url,
                            path=str(file_path.relative_to(repo_dir)),
                            content=content,
                        )
                    except UnicodeDecodeError:
                        # Skip binary files
                        continue

    def _clone_or_update(self, repo_url: str, branch: str) -> Path:
        """Clone a repository or update if already cached."""
        # Create a safe directory name from URL
        repo_name = repo_url.rstrip("/").split("/")[-1].replace(".git", "")
        org_name = repo_url.rstrip("/").split("/")[-2]
        repo_dir = self.cache_dir / f"{org_name}_{repo_name}"

        if repo_dir.exists():
            # Update existing clone
            subprocess.run(
                ["git", "fetch", "origin", branch],
                cwd=repo_dir,
                check=True,
                capture_output=True,
            )
            subprocess.run(
                ["git", "checkout", f"origin/{branch}"],
                cwd=repo_dir,
                check=True,
                capture_output=True,
            )
        else:
            # Fresh clone
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
            )

        return repo_dir

    def _is_excluded(
        self, file_path: Path, repo_dir: Path, exclude_patterns: list[str]
    ) -> bool:
        """Check if a file matches any exclusion pattern."""
        rel_path = file_path.relative_to(repo_dir)
        for pattern in exclude_patterns:
            if rel_path.match(pattern):
                return True
        return False
