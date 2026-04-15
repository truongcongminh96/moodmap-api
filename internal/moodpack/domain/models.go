package domain

import "time"

type ContentSource string

const (
	ContentSourceQuotes ContentSource = "quotes"
	ContentSourceMusic  ContentSource = "music"
	ContentSourceAll    ContentSource = "all"
)

type Units string

const (
	UnitsMetric   Units = "metric"
	UnitsImperial Units = "imperial"
)

type GetMoodPackInput struct {
	City    string
	Country string
	Units   Units
	Source  ContentSource
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type Weather struct {
	Main        string  `json:"main"`
	Description string  `json:"description"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feelsLike"`
	Humidity    int     `json:"humidity"`
	Icon        string  `json:"icon,omitempty"`
}

type Mood struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Theme string `json:"theme"`
}

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

type MusicTrack struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	URL    string `json:"url,omitempty"`
}

type MoodPack struct {
	Location    Location     `json:"location"`
	Weather     Weather      `json:"weather"`
	Mood        Mood         `json:"mood"`
	Quote       *Quote       `json:"quote,omitempty"`
	Music       []MusicTrack `json:"music,omitempty"`
	Activities  []string     `json:"activities"`
	Summary     string       `json:"summary"`
	RequestedAt time.Time    `json:"-"`
	Sources     []string     `json:"-"`
}
