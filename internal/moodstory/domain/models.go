package domain

import (
	"context"
	"time"

	moodpackdomain "moodmap-api/internal/moodpack/domain"
)

type GetMoodStoryInput struct {
	City    string
	Country string
	Units   moodpackdomain.Units
}

type MoodPackReader interface {
	GetMoodPack(ctx context.Context, input moodpackdomain.GetMoodPackInput) (*moodpackdomain.MoodPack, error)
}

type MoodStory struct {
	City       string        `json:"city"`
	Country    string        `json:"country"`
	Headline   string        `json:"headline"`
	Mood       Mood          `json:"mood"`
	Visual     Visual        `json:"visual"`
	Highlight  Highlight     `json:"highlight"`
	Story      LocalizedText `json:"story"`
	BestMoment LocalizedText `json:"bestMoment"`
	EnergyTip  LocalizedText `json:"energyTip"`
	Meta       Meta          `json:"meta"`
}

type Mood struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Theme string `json:"theme"`
}

type Visual struct {
	Gradient  string `json:"gradient"`
	TimeOfDay string `json:"timeOfDay"`
}

type Highlight struct {
	Quote *moodpackdomain.Quote `json:"quote,omitempty"`
	Track *Track                `json:"track,omitempty"`
}

type Track struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	URL    string `json:"url"`
}

type LocalizedText struct {
	EN string `json:"en"`
	VI string `json:"vi"`
}

type Meta struct {
	GeneratedAt time.Time `json:"generatedAt"`
	Sources     []string  `json:"sources"`
}
