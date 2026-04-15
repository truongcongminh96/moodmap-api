package service

import (
	"context"
	"testing"
	"time"

	moodpackdomain "moodmap-api/internal/moodpack/domain"
	moodstorydomain "moodmap-api/internal/moodstory/domain"
)

type fakeMoodPackReader struct {
	pack *moodpackdomain.MoodPack
	err  error
}

func (f fakeMoodPackReader) GetMoodPack(_ context.Context, _ moodpackdomain.GetMoodPackInput) (*moodpackdomain.MoodPack, error) {
	return f.pack, f.err
}

func TestInferTimeOfDay(t *testing.T) {
	testCases := []struct {
		name string
		hour int
		want string
	}{
		{name: "morning", hour: 9, want: "morning"},
		{name: "afternoon", hour: 14, want: "afternoon"},
		{name: "evening", hour: 18, want: "evening"},
		{name: "night", hour: 23, want: "night"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			now := time.Date(2026, 4, 16, testCase.hour, 0, 0, 0, time.UTC)
			if got := inferTimeOfDay(now); got != testCase.want {
				t.Fatalf("inferTimeOfDay() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestBuildHeadline(t *testing.T) {
	headline := buildHeadline(
		"Hanoi",
		moodpackdomain.Weather{Main: "Clouds", Description: "broken clouds"},
		moodpackdomain.Mood{Key: "calm_soft", Label: "Calm & Soft", Theme: "cloudy-silver"},
		"evening",
	)

	if headline != "A soft cloudy evening in Hanoi" {
		t.Fatalf("buildHeadline() = %q", headline)
	}
}

func TestComposeNarrativeCalmSoft(t *testing.T) {
	story, bestMoment, energyTip := composeNarrative(&moodpackdomain.MoodPack{
		Location: moodpackdomain.Location{City: "Hanoi", Country: "VN"},
		Weather:  moodpackdomain.Weather{Main: "Clouds", Description: "broken clouds"},
		Mood:     moodpackdomain.Mood{Key: "calm_soft", Label: "Calm & Soft", Theme: "cloudy-silver"},
	}, "evening")

	if story.EN == "" || story.VI == "" {
		t.Fatalf("expected localized story, got %+v", story)
	}

	if bestMoment.EN != "Early evening walk" {
		t.Fatalf("bestMoment.EN = %q", bestMoment.EN)
	}

	if energyTip.VI != "Hãy giữ ngày hôm nay nhẹ nhàng và không áp lực." {
		t.Fatalf("energyTip.VI = %q", energyTip.VI)
	}
}

func TestMoodStoryServiceGetMoodStoryGracefulFallback(t *testing.T) {
	svc := NewMoodStoryService(fakeMoodPackReader{
		pack: &moodpackdomain.MoodPack{
			Location:    moodpackdomain.Location{City: "Hanoi", Country: "VN"},
			Weather:     moodpackdomain.Weather{Main: "Clouds", Description: "broken clouds"},
			Mood:        moodpackdomain.Mood{Key: "calm_soft", Label: "Calm & Soft", Theme: "cloudy-silver"},
			RequestedAt: time.Date(2026, 4, 16, 12, 0, 0, 0, time.UTC),
			Sources:     []string{"openweather"},
		},
	}, func() time.Time {
		return time.Date(2026, 4, 16, 18, 0, 0, 0, time.UTC)
	})

	story, err := svc.GetMoodStory(context.Background(), moodstorydomain.GetMoodStoryInput{
		City: "Hanoi",
	})
	if err != nil {
		t.Fatalf("GetMoodStory() error = %v", err)
	}

	if story.Highlight.Quote != nil {
		t.Fatalf("expected nil quote highlight")
	}

	if story.Highlight.Track != nil {
		t.Fatalf("expected nil track highlight")
	}

	if len(story.Meta.Sources) != 2 || story.Meta.Sources[1] != "system" {
		t.Fatalf("unexpected sources: %+v", story.Meta.Sources)
	}
}
