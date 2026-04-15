package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"moodmap-api/internal/moodpack/domain"
	"moodmap-api/internal/moodpack/mapper"
	"moodmap-api/internal/platform/apperror"
)

type MoodService struct {
	weatherProvider domain.WeatherProvider
	quoteProvider   domain.QuoteProvider
	musicProvider   domain.MusicProvider
}

func NewMoodService(
	weatherProvider domain.WeatherProvider,
	quoteProvider domain.QuoteProvider,
	musicProvider domain.MusicProvider,
) *MoodService {
	return &MoodService{
		weatherProvider: weatherProvider,
		quoteProvider:   quoteProvider,
		musicProvider:   musicProvider,
	}
}

func (s *MoodService) GetMoodPack(ctx context.Context, input domain.GetMoodPackInput) (*domain.MoodPack, error) {
	normalizedInput, err := normalizeInput(input)
	if err != nil {
		return nil, err
	}

	location, weather, err := s.weatherProvider.GetCurrentWeather(
		ctx,
		normalizedInput.City,
		normalizedInput.Country,
		normalizedInput.Units,
	)
	if err != nil {
		return nil, err
	}

	mood, err := mapper.ResolveMood(weather.Main)
	if err != nil {
		return nil, err
	}

	pack := &domain.MoodPack{
		Location:    location,
		Weather:     weather,
		Mood:        mood,
		Activities:  activitiesForMood(mood.Key),
		RequestedAt: time.Now().UTC(),
		Sources:     []string{"openweather"},
	}
	pack.Summary = buildSummary(pack.Location.City, pack.Weather.Description, pack.Mood.Label, pack.Activities)

	switch normalizedInput.Source {
	case domain.ContentSourceQuotes:
		quote, quoteErr := s.quoteProvider.GetRandomQuote(ctx)
		if quoteErr != nil {
			return nil, apperror.ErrQuoteProvider
		}
		pack.Quote = quote
		pack.Sources = append(pack.Sources, "zenquotes")
	case domain.ContentSourceMusic:
		tracks, musicErr := s.musicProvider.GetRecommendations(ctx, mood)
		if musicErr != nil {
			return nil, apperror.ErrMusicProvider
		}
		pack.Music = tracks
		pack.Sources = append(pack.Sources, "deezer")
	case domain.ContentSourceAll:
		if quote, quoteErr := s.quoteProvider.GetRandomQuote(ctx); quoteErr == nil && quote != nil {
			pack.Quote = quote
			pack.Sources = append(pack.Sources, "zenquotes")
		}

		if tracks, musicErr := s.musicProvider.GetRecommendations(ctx, mood); musicErr == nil && len(tracks) > 0 {
			pack.Music = tracks
			pack.Sources = append(pack.Sources, "deezer")
		}
	}

	return pack, nil
}

func normalizeInput(input domain.GetMoodPackInput) (domain.GetMoodPackInput, error) {
	input.City = strings.TrimSpace(input.City)
	input.Country = strings.TrimSpace(input.Country)

	if input.City == "" {
		return domain.GetMoodPackInput{}, apperror.ErrCityRequired
	}

	if input.Units == "" {
		input.Units = domain.UnitsMetric
	}

	if input.Units != domain.UnitsMetric && input.Units != domain.UnitsImperial {
		input.Units = domain.UnitsMetric
	}

	if input.Source == "" {
		input.Source = domain.ContentSourceAll
	}

	switch input.Source {
	case domain.ContentSourceQuotes, domain.ContentSourceMusic, domain.ContentSourceAll:
	default:
		input.Source = domain.ContentSourceAll
	}

	input.City = titleCase(input.City)
	input.Country = strings.ToUpper(input.Country)

	return input, nil
}

func buildSummary(city, weatherDescription, moodLabel string, activities []string) string {
	activityPhrase := "to take a mindful pause"
	if len(activities) > 0 {
		activityPhrase = "to " + strings.ToLower(activities[0])
	}

	return fmt.Sprintf(
		"%s today feels %s with %s. A good day %s.",
		city,
		strings.ToLower(moodLabel),
		weatherDescription,
		activityPhrase,
	)
}

func titleCase(value string) string {
	parts := strings.Fields(strings.ToLower(value))
	for index, part := range parts {
		if part == "" {
			continue
		}

		parts[index] = strings.ToUpper(part[:1]) + part[1:]
	}

	return strings.Join(parts, " ")
}
