package youtube

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Service struct {
	apiKey string
}

type YouTube interface {
	DownloadThumbnail(ctx context.Context, videoID string) ([]byte, error)
}

type VideoResponse struct {
	Items []struct {
		Snippet struct {
			Thumbnails struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
				Standard struct {
					URL string `json:"url"`
				} `json:"standard"`
				Maxres struct {
					URL string `json:"url"`
				} `json:"maxres"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

func New(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
	}
}

var ErrNoThumbnail = fmt.Errorf("no thumbnail found")

func (s Service) DownloadThumbnail(ctx context.Context, videoID string) ([]byte, error) {
	videoResponse, err := s.getVideoResponse(ctx, videoID)
	if err != nil {
		return nil, err
	}

	thumbnailURL, err := s.selectThumbnailURL(videoResponse, videoID)
	if err != nil {
		return nil, err
	}

	thumbnailData, err := s.getThumbnailData(ctx, thumbnailURL)
	if err != nil {
		return nil, err
	}

	return thumbnailData, nil
}

func (s Service) getVideoResponse(ctx context.Context, videoID string) (VideoResponse, error) {
	// #nosec G402
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s", videoID, s.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return VideoResponse{}, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return VideoResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return VideoResponse{}, err
	}

	var videoResponse VideoResponse

	err = json.Unmarshal(body, &videoResponse)
	if err != nil {
		return VideoResponse{}, err
	}

	return videoResponse, nil
}

func (s Service) selectThumbnailURL(videoResponse VideoResponse, videoID string) (string, error) {
	if len(videoResponse.Items) == 0 {
		return "", fmt.Errorf("%w: video ID: %s", ErrNoThumbnail, videoID)
	}

	thumbnails := videoResponse.Items[0].Snippet.Thumbnails

	var thumbnailURL string

	switch {
	case thumbnails.Maxres.URL != "":
		thumbnailURL = thumbnails.Maxres.URL
	case thumbnails.Standard.URL != "":
		thumbnailURL = thumbnails.Standard.URL
	case thumbnails.High.URL != "":
		thumbnailURL = thumbnails.High.URL
	case thumbnails.Medium.URL != "":
		thumbnailURL = thumbnails.Medium.URL
	case thumbnails.Default.URL != "":
		thumbnailURL = thumbnails.Default.URL
	default:
		return "", fmt.Errorf("%w: video ID: %s", ErrNoThumbnail, videoID)
	}

	return thumbnailURL, nil
}

func (s Service) getThumbnailData(ctx context.Context, thumbnailURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, thumbnailURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	thumbnailData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return thumbnailData, nil
}
