package app

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/realPointer/YouTube-thumbnails-downloader/config"
	grpccontroller "github.com/realPointer/YouTube-thumbnails-downloader/internal/controller/grpc"
	"github.com/realPointer/YouTube-thumbnails-downloader/internal/repository"
	"github.com/realPointer/YouTube-thumbnails-downloader/internal/service"
	"github.com/realPointer/YouTube-thumbnails-downloader/internal/youtube"
	grpcserver "github.com/realPointer/YouTube-thumbnails-downloader/pkg/grpc_server"

	// SQLite
	_ "modernc.org/sqlite"
)

func Run() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// SQLite
	db, err := sql.Open("sqlite", "/data/thumbnails.db")
	if err != nil {
		log.Fatal("Failed to open SQLite database:", err)
	}
	defer db.Close()

	query := `
			CREATE TABLE IF NOT EXISTS thumbnails (
					video_url TEXT PRIMARY KEY,
					thumbnail_data BLOB
			);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Failed to create tables: %v", err)

		return
	}

	// Repositories
	repositories := repository.NewRepositories(db)

	// Services dependencies
	deps := service.ServicesDependencies{
		Repos:   repositories,
		YouTube: youtube.New(cfg.YouTube.APIKey),
	}
	services := service.NewServices(deps)

	// GRPC
	grpcServer, err := grpcserver.New(grpcserver.WithPort(cfg.GRPC.Port))
	if err != nil {
		log.Printf("app - Run - grpcserver.New: %v", err)
	}

	grpcSerivce := grpccontroller.NewService(services)

	go func() {
		err = grpcServer.Start(grpcSerivce)
		if err != nil {
			log.Printf("app - Run - grpcServer.Start: %v", err)
		}
	}()

	log.Printf("gRPC server started on port: %v", cfg.GRPC.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Printf("app - Run - signal: " + s.String())

	// Shutdown
	grpcServer.GrpcServer.GracefulStop()
}
