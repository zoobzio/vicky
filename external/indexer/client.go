// Package indexer provides a gRPC client implementing contracts.Indexer.
package indexer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zoobzio/pipz"
	pb "github.com/zoobzio/vicky/proto/indexer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zoobzio/vicky/models"
)

// Resilience configuration.
const (
	indexTimeout          = 5 * time.Minute // indexing can take a while
	indexMaxAttempts      = 2               // fewer retries since indexing is expensive
	indexBackoffDelay     = 500 * time.Millisecond
	indexFailureThreshold = 3
	indexResetTimeout     = 60 * time.Second
)

// Request contains the information needed to run a SCIP indexer.
type Request struct {
	JobID        int64
	RepositoryID int64
	VersionID    int64
	UserID       int64
	Owner        string
	RepoName     string
	Tag          string
	CommitSHA    string
	Language     models.Language
}

// Result contains the output from a SCIP indexer.
type Result struct {
	JobID     int64
	VersionID int64
	IndexData []byte // raw SCIP protobuf
	Error     string
}

// indexCall carries request and response through the pipeline.
type indexCall struct {
	request Request
	result  *Result
}

func (c *indexCall) Clone() *indexCall {
	clone := *c
	if c.result != nil {
		r := *c.result
		if c.result.IndexData != nil {
			r.IndexData = make([]byte, len(c.result.IndexData))
			copy(r.IndexData, c.result.IndexData)
		}
		clone.result = &r
	}
	return &clone
}

// langPipeline holds per-language connection and circuit breaker.
type langPipeline struct {
	addr     string
	conn     *grpc.ClientConn
	pipeline pipz.Chainable[*indexCall]
}

// Client implements contracts.Indexer by dispatching to language-specific
// gRPC indexer sidecars.
type Client struct {
	mu        sync.Mutex
	pipelines map[models.Language]*langPipeline
}

// NewClient creates a new indexer client with the given language-to-address map.
func NewClient(addrs map[models.Language]string) *Client {
	c := &Client{
		pipelines: make(map[models.Language]*langPipeline, len(addrs)),
	}

	for lang, addr := range addrs {
		c.pipelines[lang] = &langPipeline{addr: addr}
	}

	return c
}

// getPipeline returns or creates the pipeline for a language.
func (c *Client) getPipeline(lang models.Language) (*langPipeline, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	lp, ok := c.pipelines[lang]
	if !ok {
		return nil, fmt.Errorf("no indexer configured for language %s", lang)
	}

	// Lazily initialize connection and pipeline
	if lp.pipeline == nil {
		conn, err := grpc.NewClient(lp.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("dial indexer %s at %s: %w", lang, lp.addr, err)
		}
		lp.conn = conn
		lp.pipeline = c.buildPipeline(lang, conn)
	}

	return lp, nil
}

// buildPipeline constructs the resilient processing pipeline for a language.
func (c *Client) buildPipeline(lang models.Language, conn *grpc.ClientConn) pipz.Chainable[*indexCall] {
	client := pb.NewIndexerServiceClient(conn)

	processorID := pipz.NewIdentity(fmt.Sprintf("indexer.%s.call", lang), "gRPC call to indexer sidecar")
	timeoutID := pipz.NewIdentity(fmt.Sprintf("indexer.%s.timeout", lang), "Timeout for indexer calls")
	backoffID := pipz.NewIdentity(fmt.Sprintf("indexer.%s.backoff", lang), "Backoff retry for indexer calls")
	breakerID := pipz.NewIdentity(fmt.Sprintf("indexer.%s.breaker", lang), "Circuit breaker for indexer sidecar")

	processor := pipz.Apply(processorID, func(ctx context.Context, call *indexCall) (*indexCall, error) {
		res, err := client.Index(ctx, &pb.IndexRequest{
			JobId:        call.request.JobID,
			RepositoryId: call.request.RepositoryID,
			VersionId:    call.request.VersionID,
			UserId:       call.request.UserID,
			Owner:        call.request.Owner,
			RepoName:     call.request.RepoName,
			Tag:          call.request.Tag,
			CommitSha:    call.request.CommitSHA,
			Language:     string(call.request.Language),
		})
		if err != nil {
			return call, fmt.Errorf("indexer %s: %w", lang, err)
		}

		call.result = &Result{
			JobID:     res.JobId,
			VersionID: res.VersionId,
			IndexData: res.IndexData,
			Error:     res.Error,
		}
		return call, nil
	})

	return pipz.NewCircuitBreaker(breakerID,
		pipz.NewBackoff(backoffID,
			pipz.NewTimeout(timeoutID, processor, indexTimeout),
			indexMaxAttempts, indexBackoffDelay,
		),
		indexFailureThreshold, indexResetTimeout,
	)
}

// Index calls the appropriate language-specific indexer sidecar.
func (c *Client) Index(ctx context.Context, req Request) (*Result, error) {
	lp, err := c.getPipeline(req.Language)
	if err != nil {
		return nil, err
	}

	call := &indexCall{request: req}

	result, err := lp.pipeline.Process(ctx, call)
	if err != nil {
		return nil, err
	}

	return result.result, nil
}

// Supports returns true if an indexer address is configured for the given language.
func (c *Client) Supports(language models.Language) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.pipelines[language]
	return ok
}

// Close shuts down all open gRPC connections and pipelines.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var firstErr error
	for _, lp := range c.pipelines {
		if lp.pipeline != nil {
			if err := lp.pipeline.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		if lp.conn != nil {
			if err := lp.conn.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
			lp.conn = nil
		}
	}
	return firstErr
}
