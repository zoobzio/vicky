// Package chunker provides a gRPC client implementing contracts.Chunker.
package chunker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zoobzio/pipz"
	pb "github.com/zoobzio/vicky/proto/chunker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zoobzio/vicky/models"
)

// Resilience configuration.
const (
	chunkTimeout         = 30 * time.Second
	chunkMaxAttempts     = 3
	chunkBackoffDelay    = 100 * time.Millisecond
	chunkFailureThreshold = 5
	chunkResetTimeout    = 30 * time.Second
)

// Pipeline identities.
var (
	chunkProcessorID     = pipz.NewIdentity("chunker.call", "gRPC call to chunker sidecar")
	chunkTimeoutID       = pipz.NewIdentity("chunker.timeout", "Timeout for chunker calls")
	chunkBackoffID       = pipz.NewIdentity("chunker.backoff", "Backoff retry for chunker calls")
	chunkCircuitBreakerID = pipz.NewIdentity("chunker.breaker", "Circuit breaker for chunker sidecar")
)

// Result represents a single chunk produced by a language-aware chunker.
type Result struct {
	Content   string
	Symbol    string
	Kind      models.ChunkKind
	StartLine int
	EndLine   int
	Context   []string
}

// chunkCall carries request and response through the pipeline.
type chunkCall struct {
	language string
	filename string
	content  []byte
	results  []Result
}

func (c *chunkCall) Clone() *chunkCall {
	clone := *c
	if c.content != nil {
		clone.content = make([]byte, len(c.content))
		copy(clone.content, c.content)
	}
	if c.results != nil {
		clone.results = make([]Result, len(c.results))
		copy(clone.results, c.results)
	}
	return &clone
}

// Client implements contracts.Chunker by dispatching to a gRPC chunker sidecar.
type Client struct {
	addr     string
	pipeline pipz.Chainable[*chunkCall]

	mu   sync.Mutex
	conn *grpc.ClientConn
}

// NewClient creates a new chunker client with the given sidecar address.
func NewClient(addr string) *Client {
	c := &Client{addr: addr}
	c.pipeline = c.buildPipeline()
	return c
}

// buildPipeline constructs the resilient processing pipeline.
func (c *Client) buildPipeline() pipz.Chainable[*chunkCall] {
	processor := pipz.Apply(chunkProcessorID, c.doChunk)

	return pipz.NewCircuitBreaker(chunkCircuitBreakerID,
		pipz.NewBackoff(chunkBackoffID,
			pipz.NewTimeout(chunkTimeoutID, processor, chunkTimeout),
			chunkMaxAttempts, chunkBackoffDelay,
		),
		chunkFailureThreshold, chunkResetTimeout,
	)
}

// doChunk performs the actual gRPC call.
func (c *Client) doChunk(ctx context.Context, call *chunkCall) (*chunkCall, error) {
	client, err := c.dial()
	if err != nil {
		return call, err
	}

	res, err := client.Chunk(ctx, &pb.ChunkRequest{
		Language: call.language,
		Filename: call.filename,
		Content:  call.content,
	})
	if err != nil {
		return call, fmt.Errorf("chunker: %w", err)
	}

	call.results = make([]Result, len(res.Chunks))
	for i, ch := range res.Chunks {
		call.results[i] = Result{
			Content:   ch.Content,
			Symbol:    ch.Symbol,
			Kind:      models.ChunkKind(ch.Kind),
			StartLine: int(ch.StartLine),
			EndLine:   int(ch.EndLine),
			Context:   ch.Context,
		}
	}

	return call, nil
}

// dial lazily establishes a gRPC connection to the chunker sidecar.
func (c *Client) dial() (pb.ChunkerServiceClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return pb.NewChunkerServiceClient(c.conn), nil
	}

	conn, err := grpc.NewClient(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial chunker at %s: %w", c.addr, err)
	}

	c.conn = conn
	return pb.NewChunkerServiceClient(conn), nil
}

// Chunk splits file content into semantic chunks via the chunker sidecar.
func (c *Client) Chunk(ctx context.Context, language string, filename string, content []byte) ([]Result, error) {
	call := &chunkCall{
		language: language,
		filename: filename,
		content:  content,
	}

	result, err := c.pipeline.Process(ctx, call)
	if err != nil {
		return nil, err
	}

	return result.results, nil
}

// Supports returns true if the chunker sidecar is configured.
func (c *Client) Supports(language string) bool {
	return c.addr != ""
}

// Close shuts down the gRPC connection and pipeline.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error

	if c.pipeline != nil {
		if err := c.pipeline.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, err)
		}
		c.conn = nil
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
