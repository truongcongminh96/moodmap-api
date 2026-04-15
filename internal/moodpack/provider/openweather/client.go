package openweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"moodmap-api/internal/moodpack/domain"
	"moodmap-api/internal/platform/apperror"
)

type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewClient(baseURL, apiKey string, httpClient *http.Client) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  strings.TrimSpace(apiKey),
		client:  httpClient,
	}
}

func (c *Client) GetCurrentWeather(ctx context.Context, city, country string, units domain.Units) (domain.Location, domain.Weather, error) {
	if c.apiKey == "" {
		return domain.Location{}, domain.Weather{}, apperror.ErrWeatherProvider
	}

	queryCity := city
	if country != "" {
		queryCity = fmt.Sprintf("%s,%s", city, country)
	}

	reqURL := fmt.Sprintf(
		"%s/weather?q=%s&appid=%s&units=%s",
		c.baseURL,
		url.QueryEscape(queryCity),
		url.QueryEscape(c.apiKey),
		url.QueryEscape(string(units)),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return domain.Location{}, domain.Weather{}, apperror.ErrWeatherProvider
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return domain.Location{}, domain.Weather{}, apperror.ErrWeatherProvider
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Location{}, domain.Weather{}, apperror.ErrCityNotFound
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return domain.Location{}, domain.Weather{}, apperror.ErrWeatherProvider
	}

	var payload currentWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return domain.Location{}, domain.Weather{}, apperror.ErrWeatherProvider
	}

	if len(payload.Weather) == 0 || strings.TrimSpace(payload.Weather[0].Main) == "" {
		return domain.Location{}, domain.Weather{}, apperror.ErrWeatherProvider
	}

	location := domain.Location{
		City:    payload.Name,
		Country: payload.Sys.Country,
	}

	weather := domain.Weather{
		Main:        payload.Weather[0].Main,
		Description: payload.Weather[0].Description,
		Temperature: payload.Main.Temp,
		FeelsLike:   payload.Main.FeelsLike,
		Humidity:    payload.Main.Humidity,
		Icon:        payload.Weather[0].Icon,
	}

	return location, weather, nil
}

func IsCityNotFound(err error) bool {
	return errors.Is(err, apperror.ErrCityNotFound) || err == apperror.ErrCityNotFound
}
