package indexers

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/zoobzio/vicky/proto/indexer"
	"google.golang.org/grpc"
)

// Server implements the IndexerService gRPC server.
type Server struct {
	pb.UnimplementedIndexerServiceServer
	storage  *Storage
	executor Executor
	language string
}

// NewServer creates a new indexer gRPC server.
func NewServer(storage *Storage, executor Executor, language string) *Server {
	return &Server{
		storage:  storage,
		executor: executor,
		language: language,
	}
}

// Index handles an indexing request.
func (s *Server) Index(ctx context.Context, req *pb.IndexRequest) (*pb.IndexResult, error) {
	log.Printf("index request: job=%d owner=%s repo=%s tag=%s lang=%s",
		req.JobId, req.Owner, req.RepoName, req.Tag, req.Language)

	result := &pb.IndexResult{
		JobId:     req.JobId,
		VersionId: req.VersionId,
	}

	// Create temp work directory
	workDir, err := os.MkdirTemp("", fmt.Sprintf("indexer-%d-*", req.JobId))
	if err != nil {
		log.Printf("index error: job=%d err=%v", req.JobId, err)
		result.Error = fmt.Sprintf("create work dir: %v", err)
		return result, nil
	}
	defer os.RemoveAll(workDir)

	// Fetch blobs from minio
	count, err := s.storage.FetchBlobs(ctx, req.UserId, req.Owner, req.RepoName, req.Tag, workDir)
	if err != nil {
		log.Printf("index error: job=%d err=%v", req.JobId, err)
		result.Error = fmt.Sprintf("fetch blobs: %v", err)
		return result, nil
	}

	log.Printf("index fetched: job=%d count=%d workdir=%s", req.JobId, count, workDir)

	if count == 0 {
		result.Error = "no files found in blob storage"
		return result, nil
	}

	// Run SCIP indexer
	indexData, err := s.executor.Execute(ctx, workDir)
	if err != nil {
		log.Printf("index error: job=%d err=%v", req.JobId, err)
		result.Error = fmt.Sprintf("execute indexer: %v", err)
		return result, nil
	}

	log.Printf("index completed: job=%d bytes=%d", req.JobId, len(indexData))

	result.IndexData = indexData
	return result, nil
}

// ListenAndServe starts the gRPC server on the given address.
func ListenAndServe(ctx context.Context, addr string, srv *Server) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIndexerServiceServer(grpcServer, srv)

	log.Printf("indexer listening on %s (lang=%s)", addr, srv.language)

	return grpcServer.Serve(listener)
}
