package app

import (
	"log"

	"github.com/realPointer/YouTube-thumbnails-downloader/config"
	grpcserver "github.com/realPointer/YouTube-thumbnails-downloader/pkg/grpc_server"
	pb "github.com/realPointer/YouTube-thumbnails-downloader/pkg/thumbnail_v1"
)

func Run() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// GRPC
	grpcServer, err := grpcserver.New(grpcserver.WithPort(cfg.GRPC.Port))
	if err != nil {
		log.Printf("app - Run - grpcserver.New: %v", err)
	}

	err = grpcServer.Start(pb.UnimplementedThumbnailServiceServer{})
	if err != nil {
		log.Printf("app - Run - grpcServer.Start: %v", err)
	}
}
