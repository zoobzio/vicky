package chunkers

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/zoobzio/chisel"
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
	log.Printf("chunk request: file=%s lang=%s bytes=%d", req.Filename, req.Language, len(req.Content))

	chunks, err := s.chunker.Chunk(ctx, chisel.Language(req.Language), req.Filename, req.Content)
	if err != nil {
		log.Printf("chunk error: file=%s err=%v", req.Filename, err)
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

	log.Printf("chunk completed: file=%s chunks=%d", req.Filename, len(results))

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

	log.Printf("chunker listening on %s", addr)

	return grpcServer.Serve(listener)
}
