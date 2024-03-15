package service

import (
	"context"

	"github.com/realPointer/YouTube-thumbnails-downloader/internal/repository"
	"github.com/realPointer/YouTube-thumbnails-downloader/internal/service/services"
	"github.com/realPointer/YouTube-thumbnails-downloader/internal/youtube"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type ThumbnailService interface {
	DownloadThumbnail(ctx context.Context, videoURL string) ([]byte, error)
}

type Services struct {
	ThumbnailService
}

type ServicesDependencies struct {
	Repos   *repository.Repositories
	YouTube *youtube.Service
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		ThumbnailService: services.NewThumbnailService(deps.Repos, deps.YouTube),
	}
}
