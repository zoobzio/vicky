package contracts

import (
	"context"

	"github.com/zoobzio/vicky/external/github"
)

// GitHub defines the contract for GitHub API operations.
type GitHub interface {
	// GetTree retrieves the full repository tree at a given ref.
	// Returns all files and directories recursively.
	GetTree(ctx context.Context, token, owner, repo, ref string) ([]github.TreeEntry, error)

	// GetFileContent retrieves the content of a single file.
	// Returns an error if the file exceeds GitHub's size limit (100MB).
	GetFileContent(ctx context.Context, token, owner, repo, path, ref string) (*github.FileContent, error)

	// GetFileContentBatch retrieves multiple files concurrently.
	// Skips files that fail and returns successful results.
	GetFileContentBatch(ctx context.Context, token, owner, repo, ref string, paths []string) ([]*github.FileContent, error)
}
