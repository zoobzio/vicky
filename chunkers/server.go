package chunkers

import (
	"context"
	"fmt"
	"net"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/chisel"
	"github.com/zoobzio/vicky/events"
	pb "github.com/zoobzio/vicky/proto/chunker"
	"google.golang.org/grpc"
)

// Server implements the ChunkerService gRPC server.
type Server struct {
	pb.UnimplementedChunkerServiceServer
	chunker *chisel.Chunker
}

// NewServer creates a new chunker gRPC server.
func NewServer(providers ...chisel.Provider) *Server {
	return &Server{
		chunker: chisel.New(providers...),
	}
}

// Chunk handles a chunking request.
func (s *Server) Chunk(ctx context.Context, req *pb.ChunkRequest) (*pb.ChunkResponse, error) {
	capitan.Debug(ctx, events.ChunkerRequestSignal,
		events.FilenameKey.Field(req.Filename),
		events.LanguageKey.Field(req.Language),
		events.BytesKey.Field(int64(len(req.Content))),
	)

	chunks, err := s.chunker.Chunk(ctx, chisel.Language(req.Language), req.Filename, req.Content)
	if err != nil {
		capitan.Error(ctx, events.ChunkerErrorSignal,
			events.FilenameKey.Field(req.Filename),
			events.ErrorKey.Field(err),
		)
		return nil, fmt.Errorf("chunk %s: %w", req.Filename, err)
	}

	results := make([]*pb.ChunkResult, len(chunks))
	for i, ch := range chunks {
		results[i] = &pb.ChunkResult{
			Content:   ch.Content,
			Symbol:    ch.Symbol,
			Kind:      string(ch.Kind),
			StartLine: int32(ch.StartLine),
			EndLine:   int32(ch.EndLine),
			Context:   ch.Context,
		}
	}

	capitan.Debug(ctx, events.ChunkerCompletedSignal,
		events.FilenameKey.Field(req.Filename),
		events.ChunksKey.Field(len(results)),
	)

	return &pb.ChunkResponse{Chunks: results}, nil
}

// ListenAndServe starts the gRPC server on the given address.
func ListenAndServe(ctx context.Context, addr string, srv *Server) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChunkerServiceServer(grpcServer, srv)

	capitan.Info(ctx, events.ChunkerListeningSignal,
		events.AddrKey.Field(addr),
	)

	return grpcServer.Serve(listener)
}
