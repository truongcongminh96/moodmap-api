package mapper

import (
	"strings"

	"moodmap-api/internal/moodpack/domain"
	"moodmap-api/internal/platform/apperror"
)

var moodByWeather = map[string]domain.Mood{
	"RAIN":         {Key: "chill_reflective", Label: "Chill & Reflective", Theme: "rainy-blue", Confidence: 0.91},
	"CLEAR":        {Key: "energetic_bright", Label: "Energetic & Bright", Theme: "sunny-gold", Confidence: 0.89},
	"CLOUDS":       {Key: "calm_soft", Label: "Calm & Soft", Theme: "cloudy-silver", Confidence: 0.84},
	"THUNDERSTORM": {Key: "intense_moody", Label: "Intense & Moody", Theme: "storm-indigo", Confidence: 0.93},
	"DRIZZLE":      {Key: "cozy_gentle", Label: "Cozy & Gentle", Theme: "misty-latte", Confidence: 0.87},
	"MIST":         {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac", Confidence: 0.82},
	"FOG":          {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac", Confidence: 0.82},
	"HAZE":         {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac", Confidence: 0.8},
	"SMOKE":        {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac", Confidence: 0.78},
}

func ResolveMood(weatherMain string) (domain.Mood, error) {
	key := strings.ToUpper(strings.TrimSpace(weatherMain))
	if mood, ok := moodByWeather[key]; ok {
		return mood, nil
	}

	if key == "" {
		return domain.Mood{}, apperror.ErrMoodMappingFailed
	}

	return domain.Mood{
		Key:        "balanced_neutral",
		Label:      "Balanced & Neutral",
		Theme:      "soft-neutral",
		Confidence: 0.64,
	}, nil
}
