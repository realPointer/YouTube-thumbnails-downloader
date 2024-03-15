package services

import (
	"context"

	"github.com/realPointer/YouTube-thumbnails-downloader/internal/repository"
	"github.com/realPointer/YouTube-thumbnails-downloader/internal/youtube"
)

type ThumbnailService struct {
	thumbRepo repository.ThumbnailRepository
	yt        youtube.YouTube
}

func NewThumbnailService(thumbRepo repository.ThumbnailRepository, yt youtube.YouTube) *ThumbnailService {
	return &ThumbnailService{
		thumbRepo: thumbRepo,
		yt:        yt,
	}
}

func (s *ThumbnailService) DownloadThumbnail(ctx context.Context, videoURL string) ([]byte, error) {
	thumbnailData, err := s.thumbRepo.GetThumbnail(ctx, videoURL)
	if err != nil {
		return nil, err
	}

	if thumbnailData != nil {
		return thumbnailData, nil
	}

	thumbnailData, err = s.yt.DownloadThumbnail(ctx, videoURL)
	if err != nil {
		return nil, err
	}

	err = s.thumbRepo.SaveThumbnail(ctx, videoURL, thumbnailData)
	if err != nil {
		return nil, err
	}

	return thumbnailData, nil
}
