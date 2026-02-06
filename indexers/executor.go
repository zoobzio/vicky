package indexers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Executor runs a SCIP indexer CLI and returns the raw index data.
type Executor interface {
	Execute(ctx context.Context, workDir string) ([]byte, error)
}

// GoExecutor runs scip-go against a Go project.
type GoExecutor struct{}

// Execute runs scip-go and returns the SCIP index bytes.
func (e *GoExecutor) Execute(ctx context.Context, workDir string) ([]byte, error) {
	outputPath := filepath.Join(workDir, "index.scip")

	cmd := exec.CommandContext(ctx, "scip-go",
		"--project-root", workDir,
		"--output", outputPath,
	)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("scip-go: %w", err)
	}

	return os.ReadFile(outputPath)
}

// TypeScriptExecutor runs scip-typescript against a TypeScript/JavaScript project.
type TypeScriptExecutor struct{}

// Execute runs scip-typescript and returns the SCIP index bytes.
func (e *TypeScriptExecutor) Execute(ctx context.Context, workDir string) ([]byte, error) {
	outputPath := filepath.Join(workDir, "index.scip")

	// Attempt npm install for dependency resolution
	npmInstall := exec.CommandContext(ctx, "npm", "install", "--ignore-scripts")
	npmInstall.Dir = workDir
	_ = npmInstall.Run() // best-effort; some repos may not need it

	cmd := exec.CommandContext(ctx, "scip-typescript", "index",
		"--output", outputPath,
		"--infer-tsconfig",
	)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("scip-typescript: %w", err)
	}

	return os.ReadFile(outputPath)
}
