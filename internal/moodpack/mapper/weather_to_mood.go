package mapper

import (
	"strings"

	"moodmap-api/internal/moodpack/domain"
	"moodmap-api/internal/platform/apperror"
)

var moodByWeather = map[string]domain.Mood{
	"RAIN":         {Key: "chill_reflective", Label: "Chill & Reflective", Theme: "rainy-blue"},
	"CLEAR":        {Key: "energetic_bright", Label: "Energetic & Bright", Theme: "sunny-gold"},
	"CLOUDS":       {Key: "calm_soft", Label: "Calm & Soft", Theme: "cloudy-silver"},
	"THUNDERSTORM": {Key: "intense_moody", Label: "Intense & Moody", Theme: "storm-indigo"},
	"DRIZZLE":      {Key: "cozy_gentle", Label: "Cozy & Gentle", Theme: "misty-latte"},
	"MIST":         {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac"},
	"FOG":          {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac"},
	"HAZE":         {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac"},
	"SMOKE":        {Key: "dreamy_quiet", Label: "Dreamy & Quiet", Theme: "fog-lilac"},
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
		Key:   "balanced_neutral",
		Label: "Balanced & Neutral",
		Theme: "soft-neutral",
	}, nil
}
