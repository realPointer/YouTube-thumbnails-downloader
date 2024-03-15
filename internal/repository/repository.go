package repository

import (
	"context"
	"database/sql"

	"github.com/realPointer/YouTube-thumbnails-downloader/internal/repository/sqlite"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type ThumbnailRepository interface {
	GetThumbnail(ctx context.Context, videoURL string) ([]byte, error)
	SaveThumbnail(ctx context.Context, videoURL string, thumbnailData []byte) error
}

type Repositories struct {
	ThumbnailRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ThumbnailRepository: sqlite.NewRepository(db),
	}
}
