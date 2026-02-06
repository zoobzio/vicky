package indexers

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/vicky/events"
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
	capitan.Info(ctx, events.IndexerRequestSignal,
		events.JobIDKey.Field(req.JobId),
		events.OwnerKey.Field(req.Owner),
		events.RepoNameKey.Field(req.RepoName),
		events.TagKey.Field(req.Tag),
		events.LanguageKey.Field(req.Language),
	)

	result := &pb.IndexResult{
		JobId:     req.JobId,
		VersionId: req.VersionId,
	}

	// Create temp work directory
	workDir, err := os.MkdirTemp("", fmt.Sprintf("indexer-%d-*", req.JobId))
	if err != nil {
		capitan.Error(ctx, events.IndexerErrorSignal,
			events.JobIDKey.Field(req.JobId),
			events.ErrorKey.Field(err),
		)
		result.Error = fmt.Sprintf("create work dir: %v", err)
		return result, nil
	}
	defer os.RemoveAll(workDir)

	// Fetch blobs from minio
	count, err := s.storage.FetchBlobs(ctx, req.UserId, req.Owner, req.RepoName, req.Tag, workDir)
	if err != nil {
		capitan.Error(ctx, events.IndexerErrorSignal,
			events.JobIDKey.Field(req.JobId),
			events.ErrorKey.Field(err),
		)
		result.Error = fmt.Sprintf("fetch blobs: %v", err)
		return result, nil
	}

	capitan.Debug(ctx, events.IndexerFetchedSignal,
		events.JobIDKey.Field(req.JobId),
		events.CountKey.Field(count),
		events.WorkDirKey.Field(workDir),
	)

	if count == 0 {
		result.Error = "no files found in blob storage"
		return result, nil
	}

	// Run SCIP indexer
	indexData, err := s.executor.Execute(ctx, workDir)
	if err != nil {
		capitan.Error(ctx, events.IndexerErrorSignal,
			events.JobIDKey.Field(req.JobId),
			events.ErrorKey.Field(err),
		)
		result.Error = fmt.Sprintf("execute indexer: %v", err)
		return result, nil
	}

	capitan.Info(ctx, events.IndexerCompletedSignal,
		events.JobIDKey.Field(req.JobId),
		events.IndexBytesKey.Field(len(indexData)),
	)

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

	capitan.Info(ctx, events.IndexerListeningSignal,
		events.AddrKey.Field(addr),
		events.LanguageKey.Field(srv.language),
	)

	return grpcServer.Serve(listener)
}
