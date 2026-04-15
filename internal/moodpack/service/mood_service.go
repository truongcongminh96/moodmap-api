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
	pack.Summary = buildSummary(pack.Location.City, pack.Weather.Description, pack.Mood)

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

func buildSummary(city, weatherDescription string, mood domain.Mood) string {
	switch mood.Key {
	case "energetic_bright":
		return fmt.Sprintf(
			"%s today feels energetic and bright under %s. It is a great day to get outside, move your body, and enjoy the sunlight.",
			city,
			addArticle(weatherDescription),
		)
	case "chill_reflective":
		return fmt.Sprintf(
			"%s today feels calm and reflective with %s. It is a lovely day to slow down, sip something warm, and spend a little time with your thoughts.",
			city,
			addArticle(weatherDescription),
		)
	case "calm_soft":
		return fmt.Sprintf(
			"%s today feels soft and balanced with %s. It is a nice day to read, plan gently, and ease into your rhythm.",
			city,
			addArticle(weatherDescription),
		)
	case "intense_moody":
		return fmt.Sprintf(
			"%s today feels intense and moody with %s. It is a strong day for deep focus, creative work, and staying tucked into your own flow.",
			city,
			addArticle(weatherDescription),
		)
	case "cozy_gentle":
		return fmt.Sprintf(
			"%s today feels cozy and gentle with %s. It is a perfect day to stay comfortable, play soft music, and recharge indoors.",
			city,
			addArticle(weatherDescription),
		)
	case "dreamy_quiet":
		return fmt.Sprintf(
			"%s today feels dreamy and quiet with %s. It is a good day to wander slowly, think deeply, and keep things light.",
			city,
			addArticle(weatherDescription),
		)
	default:
		return fmt.Sprintf(
			"%s today feels balanced with %s. It is a good day to take a mindful pause and choose a steady pace.",
			city,
			addArticle(weatherDescription),
		)
	}
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

func addArticle(description string) string {
	description = strings.TrimSpace(strings.ToLower(description))
	if description == "" {
		return "the current weather"
	}

	if strings.HasPrefix(description, "a ") || strings.HasPrefix(description, "an ") || strings.HasPrefix(description, "the ") {
		return description
	}

	if strings.HasSuffix(description, "s") {
		return description
	}

	if strings.ContainsRune("aeiou", rune(description[0])) {
		return "an " + description
	}

	return "a " + description
}
