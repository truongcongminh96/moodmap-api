package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"moodmap-api/internal/moodpack/domain"
	"moodmap-api/internal/platform/apperror"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  httpClient,
	}
}

func (c *Client) GetRecommendations(ctx context.Context, mood domain.Mood) ([]domain.MusicTrack, error) {
	query := queryForMood(mood.Key)
	reqURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, apperror.ErrMusicProvider
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, apperror.ErrMusicProvider
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, apperror.ErrMusicProvider
	}

	var payload searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, apperror.ErrMusicProvider
	}

	tracks := make([]domain.MusicTrack, 0, 5)
	for _, item := range payload.Data {
		if item.Title == "" || item.Artist.Name == "" {
			continue
		}

		tracks = append(tracks, domain.MusicTrack{
			Title:  item.Title,
			Artist: item.Artist.Name,
			URL:    item.Link,
		})

		if len(tracks) == 5 {
			break
		}
	}

	if len(tracks) == 0 {
		return nil, apperror.ErrMusicProvider
	}

	return tracks, nil
}

func queryForMood(moodKey string) string {
	switch moodKey {
	case "chill_reflective":
		return "lofi chill"
	case "energetic_bright":
		return "upbeat indie pop"
	case "calm_soft":
		return "ambient soft"
	case "intense_moody":
		return "alternative moody"
	case "cozy_gentle":
		return "acoustic cozy"
	case "dreamy_quiet":
		return "dream pop"
	default:
		return "focus instrumental"
	}
}
