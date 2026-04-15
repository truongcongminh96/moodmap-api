package zenquotes

import (
	"context"
	"encoding/json"
	"net/http"
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

func (c *Client) GetRandomQuote(ctx context.Context) (*domain.Quote, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/random", nil)
	if err != nil {
		return nil, apperror.ErrQuoteProvider
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, apperror.ErrQuoteProvider
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, apperror.ErrQuoteProvider
	}

	var payload randomQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, apperror.ErrQuoteProvider
	}

	if len(payload) == 0 || strings.TrimSpace(payload[0].Q) == "" {
		return nil, apperror.ErrQuoteProvider
	}

	return &domain.Quote{
		Text:   payload[0].Q,
		Author: payload[0].A,
	}, nil
}
