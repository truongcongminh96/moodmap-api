package domain

import "context"

type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, city, country string, units Units) (Location, Weather, error)
}

type QuoteProvider interface {
	GetRandomQuote(ctx context.Context) (*Quote, error)
}

type MusicProvider interface {
	GetRecommendations(ctx context.Context, mood Mood) ([]MusicTrack, error)
}
