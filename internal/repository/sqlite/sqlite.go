package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetThumbnail(ctx context.Context, videoURL string) ([]byte, error) {
	query := "SELECT thumbnail_data FROM thumbnails WHERE video_url = ?;"
	row := r.db.QueryRowContext(ctx, query, videoURL)

	var thumbnailData []byte

	err := row.Scan(&thumbnailData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get thumbnail: %w", err)
	}

	return thumbnailData, nil
}

func (r *Repository) SaveThumbnail(ctx context.Context, videoURL string, thumbnailData []byte) error {
	query := "INSERT OR REPLACE INTO thumbnails (video_url, thumbnail_data) VALUES (?, ?);"

	_, err := r.db.ExecContext(ctx, query, videoURL, thumbnailData)
	if err != nil {
		return fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return nil
}
