package grpccontroller

import (
	"context"

	"github.com/realPointer/YouTube-thumbnails-downloader/internal/service"
	pb "github.com/realPointer/YouTube-thumbnails-downloader/pkg/thumbnail_v1"
)

type Service struct {
	pb.UnimplementedThumbnailServiceServer
	thumbnailService service.ThumbnailService
}

func NewService(thumbnailService service.ThumbnailService) *Service {
	return &Service{
		thumbnailService: thumbnailService,
	}
}

func (s *Service) DownloadThumbnail(ctx context.Context, req *pb.DownloadThumbnailRequest) (*pb.DownloadThumbnailResponse, error) {
	thumbnailData, err := s.thumbnailService.DownloadThumbnail(ctx, req.VideoUrl)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadThumbnailResponse{ThumbnailData: thumbnailData}, nil
}
