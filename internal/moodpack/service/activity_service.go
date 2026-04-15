package service

func activitiesForMood(moodKey string) []string {
	switch moodKey {
	case "chill_reflective":
		return []string{
			"Drink coffee near the window",
			"Write a short journal entry",
			"Put on a lo-fi playlist and slow down",
		}
	case "energetic_bright":
		return []string{
			"Take a sunny outdoor walk",
			"Try a quick energizing workout",
			"Capture a few bright photos around the city",
		}
	case "calm_soft":
		return []string{
			"Read a few pages of a book",
			"Plan the rest of your week",
			"Play light background music while you work",
		}
	case "intense_moody":
		return []string{
			"Set up a focused deep work block",
			"Sketch or brainstorm without distractions",
			"End the day with a movie night",
		}
	case "cozy_gentle":
		return []string{
			"Make tea and reset your pace",
			"Listen to acoustic music indoors",
			"Tidy a small corner of your space",
		}
	case "dreamy_quiet":
		return []string{
			"Take a slow solo walk",
			"Write down a few ideas or daydreams",
			"Keep the evening quiet with soft music",
		}
	default:
		return []string{
			"Take a mindful break",
			"Check in with your energy for the day",
			"Pick one small thing to finish well",
		}
	}
}
