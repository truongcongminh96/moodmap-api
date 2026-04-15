package http

import (
	"net/http"

	"moodmap-api/internal/config"
	"moodmap-api/internal/moodpack/provider/deezer"
	"moodmap-api/internal/moodpack/provider/openweather"
	"moodmap-api/internal/moodpack/provider/zenquotes"
	"moodmap-api/internal/moodpack/service"
	moodhttp "moodmap-api/internal/moodpack/transport/http"
)

func NewRouter(cfg config.Config) http.Handler {
	client := &http.Client{Timeout: cfg.HTTPTimeout}

	weatherClient := openweather.NewClient(cfg.OpenWeatherBaseURL, cfg.OpenWeatherAPIKey, client)
	quoteClient := zenquotes.NewClient(cfg.ZenQuotesBaseURL, client)
	musicClient := deezer.NewClient(cfg.DeezerBaseURL, client)

	moodService := service.NewMoodService(weatherClient, quoteClient, musicClient)
	moodHandler := moodhttp.NewMoodHandler(moodService)

	mux := http.NewServeMux()
	mux.HandleFunc("/kaithhealth", HealthCheck)
	mux.HandleFunc("/healthz", HealthCheck)
	mux.HandleFunc("/api/v1/mood-pack", moodHandler.GetMoodPack)

	return mux
}
