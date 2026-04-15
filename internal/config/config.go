package config

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port               string
	OpenWeatherAPIKey  string
	OpenWeatherBaseURL string
	ZenQuotesBaseURL   string
	DeezerBaseURL      string
	HTTPTimeout        time.Duration
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout  time.Duration
}

func Load() Config {
	loadDotEnv(".env")

	return Config{
		Port:               getEnv("PORT", "8080"),
		OpenWeatherAPIKey:  getEnv("OPENWEATHER_API_KEY", ""),
		OpenWeatherBaseURL: getEnv("OPENWEATHER_BASE_URL", "https://api.openweathermap.org/data/2.5"),
		ZenQuotesBaseURL:   getEnv("ZENQUOTES_BASE_URL", "https://zenquotes.io/api"),
		DeezerBaseURL:      getEnv("DEEZER_BASE_URL", "https://api.deezer.com"),
		HTTPTimeout:        8 * time.Second,
		ServerReadTimeout:  5 * time.Second,
		ServerWriteTimeout: 10 * time.Second,
		ServerIdleTimeout:  30 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		_ = os.Setenv(key, value)
	}
}
