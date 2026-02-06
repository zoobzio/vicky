package github

import (
	"context"
	"sync"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/zoobzio/pipz"
	"golang.org/x/oauth2"
)

// Resilience configuration.
const (
	ghTimeout          = 30 * time.Second
	ghMaxAttempts      = 3
	ghBackoffDelay     = 200 * time.Millisecond
	ghFailureThreshold = 5
	ghResetTimeout     = 30 * time.Second

	// GitHub rate limit: 5000 requests/hour for authenticated users = ~1.4 req/sec
	// We use a slightly lower rate to leave headroom.
	ghRatePerSecond = 1.0
	ghRateBurst     = 10
)

// Pipeline identities.
var (
	treeProcessorID   = pipz.NewIdentity("github.tree.call", "GitHub API call for tree")
	treeTimeoutID     = pipz.NewIdentity("github.tree.timeout", "Timeout for tree calls")
	treeBackoffID     = pipz.NewIdentity("github.tree.backoff", "Backoff retry for tree calls")
	treeBreakerID     = pipz.NewIdentity("github.tree.breaker", "Circuit breaker for tree calls")
	treeRateLimiterID = pipz.NewIdentity("github.tree.ratelimit", "Rate limiter for tree calls")

	fileProcessorID   = pipz.NewIdentity("github.file.call", "GitHub API call for file content")
	fileTimeoutID     = pipz.NewIdentity("github.file.timeout", "Timeout for file calls")
	fileBackoffID     = pipz.NewIdentity("github.file.backoff", "Backoff retry for file calls")
	fileBreakerID     = pipz.NewIdentity("github.file.breaker", "Circuit breaker for file calls")
	fileRateLimiterID = pipz.NewIdentity("github.file.ratelimit", "Rate limiter for file calls")
)

const defaultWorkers = 10

// TreeEntry represents a file or directory in a GitHub repository tree.
type TreeEntry struct {
	Path string
	Mode string
	Type string // "blob" or "tree"
	SHA  string
	Size int64 // only for blobs
}

// FileContent represents the content of a file from GitHub.
type FileContent struct {
	Path    string
	SHA     string
	Content []byte
	Size    int64
}

// treeCall carries request and response through the pipeline.
type treeCall struct {
	token   string
	owner   string
	repo    string
	ref     string
	entries []TreeEntry
}

func (c *treeCall) Clone() *treeCall {
	clone := *c
	if c.entries != nil {
		clone.entries = make([]TreeEntry, len(c.entries))
		copy(clone.entries, c.entries)
	}
	return &clone
}

// fileCall carries request and response through the pipeline.
type fileCall struct {
	token   string
	owner   string
	repo    string
	ref     string
	path    string
	content *FileContent
}

func (c *fileCall) Clone() *fileCall {
	clone := *c
	if c.content != nil {
		fc := *c.content
		if c.content.Content != nil {
			fc.Content = make([]byte, len(c.content.Content))
			copy(fc.Content, c.content.Content)
		}
		clone.content = &fc
	}
	return &clone
}

// Client implements contracts.GitHub using google/go-github.
type Client struct {
	workers      int
	treePipeline pipz.Chainable[*treeCall]
	filePipeline pipz.Chainable[*fileCall]
}

// NewClient creates a new GitHub API client.
func NewClient() *Client {
	c := &Client{
		workers: defaultWorkers,
	}
	c.treePipeline = c.buildTreePipeline()
	c.filePipeline = c.buildFilePipeline()
	return c
}

// buildTreePipeline constructs the resilient processing pipeline for GetTree.
func (c *Client) buildTreePipeline() pipz.Chainable[*treeCall] {
	processor := pipz.Apply(treeProcessorID, func(ctx context.Context, call *treeCall) (*treeCall, error) {
		gh := newGitHubClient(ctx, call.token)

		tree, _, err := gh.Git.GetTree(ctx, call.owner, call.repo, call.ref, true)
		if err != nil {
			return call, err
		}

		call.entries = make([]TreeEntry, 0, len(tree.Entries))
		for _, e := range tree.Entries {
			call.entries = append(call.entries, TreeEntry{
				Path: e.GetPath(),
				Mode: e.GetMode(),
				Type: e.GetType(),
				SHA:  e.GetSHA(),
				Size: int64(e.GetSize()),
			})
		}
		return call, nil
	})

	return pipz.NewRateLimiter(treeRateLimiterID, ghRatePerSecond, ghRateBurst,
		pipz.NewCircuitBreaker(treeBreakerID,
			pipz.NewBackoff(treeBackoffID,
				pipz.NewTimeout(treeTimeoutID, processor, ghTimeout),
				ghMaxAttempts, ghBackoffDelay,
			),
			ghFailureThreshold, ghResetTimeout,
		),
	)
}

// buildFilePipeline constructs the resilient processing pipeline for GetFileContent.
func (c *Client) buildFilePipeline() pipz.Chainable[*fileCall] {
	processor := pipz.Apply(fileProcessorID, func(ctx context.Context, call *fileCall) (*fileCall, error) {
		gh := newGitHubClient(ctx, call.token)

		opts := &github.RepositoryContentGetOptions{Ref: call.ref}
		content, _, _, err := gh.Repositories.GetContents(ctx, call.owner, call.repo, call.path, opts)
		if err != nil {
			return call, err
		}

		decoded, err := content.GetContent()
		if err != nil {
			return call, err
		}

		call.content = &FileContent{
			Path:    content.GetPath(),
			SHA:     content.GetSHA(),
			Content: []byte(decoded),
			Size:    int64(content.GetSize()),
		}
		return call, nil
	})

	return pipz.NewRateLimiter(fileRateLimiterID, ghRatePerSecond, ghRateBurst,
		pipz.NewCircuitBreaker(fileBreakerID,
			pipz.NewBackoff(fileBackoffID,
				pipz.NewTimeout(fileTimeoutID, processor, ghTimeout),
				ghMaxAttempts, ghBackoffDelay,
			),
			ghFailureThreshold, ghResetTimeout,
		),
	)
}

// newGitHubClient creates an authenticated github.Client for the given token.
func newGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// GetTree retrieves the full repository tree at a given ref.
func (c *Client) GetTree(ctx context.Context, token, owner, repo, ref string) ([]TreeEntry, error) {
	call := &treeCall{
		token: token,
		owner: owner,
		repo:  repo,
		ref:   ref,
	}

	result, err := c.treePipeline.Process(ctx, call)
	if err != nil {
		return nil, err
	}

	return result.entries, nil
}

// GetFileContent retrieves the content of a single file.
func (c *Client) GetFileContent(ctx context.Context, token, owner, repo, path, ref string) (*FileContent, error) {
	call := &fileCall{
		token: token,
		owner: owner,
		repo:  repo,
		ref:   ref,
		path:  path,
	}

	result, err := c.filePipeline.Process(ctx, call)
	if err != nil {
		return nil, err
	}

	return result.content, nil
}

// GetFileContentBatch retrieves multiple files concurrently.
// Uses the resilient GetFileContent pipeline internally.
func (c *Client) GetFileContentBatch(ctx context.Context, token, owner, repo, ref string, paths []string) ([]*FileContent, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	type result struct {
		content *FileContent
		err     error
	}

	results := make(chan result, len(paths))
	pathChan := make(chan string, len(paths))

	for _, p := range paths {
		pathChan <- p
	}
	close(pathChan)

	workers := c.workers
	if len(paths) < workers {
		workers = len(paths)
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				content, err := c.GetFileContent(ctx, token, owner, repo, path, ref)
				results <- result{content: content, err: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var contents []*FileContent
	for r := range results {
		if r.err == nil && r.content != nil {
			contents = append(contents, r.content)
		}
	}

	return contents, nil
}

// Close shuts down the pipelines.
func (c *Client) Close() error {
	var errs []error

	if c.treePipeline != nil {
		if err := c.treePipeline.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.filePipeline != nil {
		if err := c.filePipeline.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
