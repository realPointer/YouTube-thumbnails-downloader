package grpcserver

import (
	"fmt"
	"net"

	pb "github.com/realPointer/YouTube-thumbnails-downloader/pkg/thumbnail_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	_defaultPort = 50051
)

type Server struct {
	GrpcServer *grpc.Server
	port       int
}

func New(opts ...Option) (*Server, error) {
	server := &Server{
		GrpcServer: grpc.NewServer(),
		port:       _defaultPort,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

func (s *Server) Start(thumbnailService pb.ThumbnailServiceServer) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	reflection.Register(s.GrpcServer)
	pb.RegisterThumbnailServiceServer(s.GrpcServer, thumbnailService)

	if err := s.GrpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
