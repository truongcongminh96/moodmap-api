package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	moodpackdomain "moodmap-api/internal/moodpack/domain"
	moodstorydomain "moodmap-api/internal/moodstory/domain"
)

type MoodStoryService struct {
	moodPackReader moodstorydomain.MoodPackReader
	now            func() time.Time
}

func NewMoodStoryService(moodPackReader moodstorydomain.MoodPackReader, now func() time.Time) *MoodStoryService {
	if now == nil {
		now = time.Now
	}

	return &MoodStoryService{
		moodPackReader: moodPackReader,
		now:            now,
	}
}

func (s *MoodStoryService) GetMoodStory(ctx context.Context, input moodstorydomain.GetMoodStoryInput) (*moodstorydomain.MoodStory, error) {
	pack, err := s.moodPackReader.GetMoodPack(ctx, moodpackdomain.GetMoodPackInput{
		City:    input.City,
		Country: input.Country,
		Units:   input.Units,
		Source:  moodpackdomain.ContentSourceAll,
	})
	if err != nil {
		return nil, err
	}

	return s.composeStory(pack, s.now()), nil
}

func (s *MoodStoryService) composeStory(pack *moodpackdomain.MoodPack, now time.Time) *moodstorydomain.MoodStory {
	timeOfDay := inferTimeOfDay(now)
	storyText, bestMoment, energyTip := composeNarrative(pack, timeOfDay)
	sources := append([]string{}, pack.Sources...)
	sources = appendUnique(sources, "system")

	return &moodstorydomain.MoodStory{
		City:     pack.Location.City,
		Country:  pack.Location.Country,
		Headline: buildHeadline(pack.Location.City, pack.Weather, pack.Mood, timeOfDay),
		Mood: moodstorydomain.Mood{
			Key:   pack.Mood.Key,
			Label: pack.Mood.Label,
			Theme: pack.Mood.Theme,
		},
		Visual: moodstorydomain.Visual{
			Gradient:  gradientFor(pack.Mood),
			TimeOfDay: timeOfDay,
		},
		Highlight: moodstorydomain.Highlight{
			Quote: pack.Quote,
			Track: pickHighlightTrack(pack.Music),
		},
		Story:      storyText,
		BestMoment: bestMoment,
		EnergyTip:  energyTip,
		Meta: moodstorydomain.Meta{
			GeneratedAt: now.UTC(),
			Sources:     sources,
		},
	}
}

func inferTimeOfDay(now time.Time) string {
	hour := now.Hour()

	switch {
	case hour >= 5 && hour < 12:
		return "morning"
	case hour >= 12 && hour < 17:
		return "afternoon"
	case hour >= 17 && hour < 21:
		return "evening"
	default:
		return "night"
	}
}

func buildHeadline(city string, weather moodpackdomain.Weather, mood moodpackdomain.Mood, timeOfDay string) string {
	moodWord := moodAdjective(mood.Key)
	weatherWord := headlineWeatherWord(weather)

	switch {
	case moodWord != "" && weatherWord != "":
		return fmt.Sprintf("A %s %s %s in %s", moodWord, weatherWord, timeOfDay, city)
	case weatherWord != "":
		return fmt.Sprintf("A %s %s in %s", weatherWord, timeOfDay, city)
	default:
		return fmt.Sprintf("A %s mood in %s", timeOfDay, city)
	}
}

func composeNarrative(pack *moodpackdomain.MoodPack, timeOfDay string) (moodstorydomain.LocalizedText, moodstorydomain.LocalizedText, moodstorydomain.LocalizedText) {
	city := pack.Location.City
	sceneEN := weatherSceneEN(pack.Weather.Description)
	sceneVI := weatherSceneVI(pack.Weather.Description)

	switch pack.Mood.Key {
	case "energetic_bright":
		return moodstorydomain.LocalizedText{
				EN: fmt.Sprintf("%s feels bright and open today %s. There is a little lift in the air, making this a lovely time to move, wander, and enjoy the momentum of the day.", city, sceneEN),
				VI: fmt.Sprintf("%s hôm nay sáng và thoáng %s. Không khí mang theo một nguồn năng lượng vừa đủ để bạn muốn bước ra ngoài, chuyển động nhẹ, và tận hưởng nhịp điệu hứng khởi của ngày mới.", city, sceneVI),
			},
			bestMomentFor(pack.Mood.Key, timeOfDay),
			energyTipFor(pack.Mood.Key)
	case "chill_reflective":
		return moodstorydomain.LocalizedText{
				EN: fmt.Sprintf("%s feels calm and reflective today %s. The day leans inward in a beautiful way, inviting warm drinks, softer notes, and a little more space for your own thoughts.", city, sceneEN),
				VI: fmt.Sprintf("%s hôm nay trầm hơn và nhiều suy tư hơn %s. Cả ngày như khẽ nghiêng vào bên trong, rất hợp cho một tách đồ uống ấm, vài giai điệu nhẹ, và một khoảng riêng cho chính mình.", city, sceneVI),
			},
			bestMomentFor(pack.Mood.Key, timeOfDay),
			energyTipFor(pack.Mood.Key)
	case "calm_soft":
		return moodstorydomain.LocalizedText{
				EN: fmt.Sprintf("%s feels calm and unhurried today %s. It is the kind of weather that invites slower thinking, light music, and a gentler pace.", city, sceneEN),
				VI: fmt.Sprintf("%s hôm nay mang một nhịp điệu chậm và dịu %s. Đây là kiểu thời tiết khiến bạn muốn suy nghĩ nhẹ nhàng hơn, nghe nhạc khẽ hơn, và sống chậm lại một chút.", city, sceneVI),
			},
			bestMomentFor(pack.Mood.Key, timeOfDay),
			energyTipFor(pack.Mood.Key)
	case "intense_moody":
		return moodstorydomain.LocalizedText{
				EN: fmt.Sprintf("%s feels moody and charged today %s. There is depth in the air, making it a strong moment for deep focus, cinematic thinking, and creative concentration.", city, sceneEN),
				VI: fmt.Sprintf("%s hôm nay có chiều sâu và một chút sắc độ mạnh mẽ %s. Bầu không khí này rất hợp cho những giờ tập trung sâu, suy nghĩ đậm hơn, và để cảm hứng sáng tạo đi xa hơn một chút.", city, sceneVI),
			},
			bestMomentFor(pack.Mood.Key, timeOfDay),
			energyTipFor(pack.Mood.Key)
	case "cozy_gentle":
		return moodstorydomain.LocalizedText{
				EN: fmt.Sprintf("%s feels cozy and unhurried today %s. Small comforts land a little more beautifully in this kind of weather, and even an ordinary moment can feel softly restorative.", city, sceneEN),
				VI: fmt.Sprintf("%s hôm nay đem lại cảm giác ấm và hiền %s. Những niềm vui nhỏ dường như chạm đến bạn rõ hơn trong kiểu thời tiết này, khiến cả một khoảnh khắc bình thường cũng trở nên dễ chịu và hồi phục hơn.", city, sceneVI),
			},
			bestMomentFor(pack.Mood.Key, timeOfDay),
			energyTipFor(pack.Mood.Key)
	case "dreamy_quiet":
		return moodstorydomain.LocalizedText{
				EN: fmt.Sprintf("%s feels hushed and dreamy today %s. The atmosphere is slightly suspended, ideal for slow walks, private thoughts, and letting the outside noise fall away.", city, sceneEN),
				VI: fmt.Sprintf("%s hôm nay yên và mơ màng hơn thường lệ %s. Không khí như chậm lại vừa đủ để bạn đi thật thong thả, nghĩ sâu hơn một chút, và để tiếng ồn xung quanh lùi ra xa.", city, sceneVI),
			},
			bestMomentFor(pack.Mood.Key, timeOfDay),
			energyTipFor(pack.Mood.Key)
	default:
		return moodstorydomain.LocalizedText{
				EN: fmt.Sprintf("%s feels balanced today %s. It is a good moment to keep things steady, soften the pressure, and move through the day with a clear head.", city, sceneEN),
				VI: fmt.Sprintf("%s hôm nay giữ một nhịp khá cân bằng %s. Đây là lúc phù hợp để sống chậm vừa phải, giảm bớt áp lực, và đi qua ngày với một tâm trí sáng rõ hơn.", city, sceneVI),
			},
			bestMomentFor(pack.Mood.Key, timeOfDay),
			energyTipFor(pack.Mood.Key)
	}
}

func bestMomentFor(moodKey, timeOfDay string) moodstorydomain.LocalizedText {
	switch moodKey {
	case "energetic_bright":
		switch timeOfDay {
		case "morning":
			return moodstorydomain.LocalizedText{EN: "Golden-hour stretch", VI: "Một nhịp vươn vai trong nắng sớm"}
		case "afternoon":
			return moodstorydomain.LocalizedText{EN: "Sunny city wander", VI: "Một vòng dạo phố đầy nắng"}
		case "evening":
			return moodstorydomain.LocalizedText{EN: "Sunset recharge walk", VI: "Một buổi đi bộ đón hoàng hôn"}
		default:
			return moodstorydomain.LocalizedText{EN: "Light late-night reset", VI: "Một nhịp reset đêm thật nhẹ"}
		}
	case "chill_reflective":
		switch timeOfDay {
		case "morning":
			return moodstorydomain.LocalizedText{EN: "Quiet window coffee", VI: "Một tách cà phê bên khung cửa sổ"}
		case "afternoon":
			return moodstorydomain.LocalizedText{EN: "Slow journaling pause", VI: "Một khoảng viết vài dòng thật chậm"}
		case "evening":
			return moodstorydomain.LocalizedText{EN: "Rainy playlist hour", VI: "Một giờ nghe playlist trong mưa"}
		default:
			return moodstorydomain.LocalizedText{EN: "Nighttime reflection", VI: "Một khoảng lặng để ngẫm lại trong đêm"}
		}
	case "calm_soft":
		switch timeOfDay {
		case "morning":
			return moodstorydomain.LocalizedText{EN: "Slow coffee break", VI: "Một tách cà phê buổi sáng thật chậm"}
		case "afternoon":
			return moodstorydomain.LocalizedText{EN: "Quiet desk reset", VI: "Một nhịp sắp lại bàn làm việc thật yên"}
		case "evening":
			return moodstorydomain.LocalizedText{EN: "Early evening walk", VI: "Một buổi dạo bộ lúc chiều muộn"}
		default:
			return moodstorydomain.LocalizedText{EN: "Soft night wind-down", VI: "Một khoảng hạ nhịp thật êm trước khi ngủ"}
		}
	case "intense_moody":
		return moodstorydomain.LocalizedText{EN: "Deep focus window", VI: "Một khoảng tập trung thật sâu"}
	case "cozy_gentle":
		return moodstorydomain.LocalizedText{EN: "Warm indoor reset", VI: "Một khoảng nghỉ ấm áp trong nhà"}
	case "dreamy_quiet":
		return moodstorydomain.LocalizedText{EN: "Solo twilight walk", VI: "Một buổi đi bộ một mình lúc chạng vạng"}
	default:
		return moodstorydomain.LocalizedText{EN: "A steady mindful pause", VI: "Một khoảng dừng vừa đủ và tỉnh táo"}
	}
}

func energyTipFor(moodKey string) moodstorydomain.LocalizedText {
	switch moodKey {
	case "energetic_bright":
		return moodstorydomain.LocalizedText{
			EN: "Let the extra brightness carry you, but keep the pace playful.",
			VI: "Hãy để nguồn sáng của ngày hôm nay nâng bạn lên, nhưng vẫn giữ nhịp thật nhẹ và vui.",
		}
	case "chill_reflective":
		return moodstorydomain.LocalizedText{
			EN: "Give yourself a little room to think, instead of filling every moment.",
			VI: "Hãy chừa cho mình một khoảng trống để nghĩ ngợi, thay vì lấp đầy mọi khoảnh khắc.",
		}
	case "calm_soft":
		return moodstorydomain.LocalizedText{
			EN: "Keep the day light and low-pressure.",
			VI: "Hãy giữ ngày hôm nay nhẹ nhàng và không áp lực.",
		}
	case "intense_moody":
		return moodstorydomain.LocalizedText{
			EN: "Channel the intensity into one thing that matters, not ten at once.",
			VI: "Hãy dồn năng lượng mạnh của hôm nay vào một điều thật sự quan trọng, thay vì cùng lúc quá nhiều thứ.",
		}
	case "cozy_gentle":
		return moodstorydomain.LocalizedText{
			EN: "Lean into comfort without feeling guilty about moving a little slower.",
			VI: "Hãy cho phép mình tìm đến sự dễ chịu mà không thấy áy náy vì đi chậm hơn một chút.",
		}
	case "dreamy_quiet":
		return moodstorydomain.LocalizedText{
			EN: "Let today move at a softer pace and protect your quiet.",
			VI: "Hãy để hôm nay trôi theo một nhịp mềm hơn và giữ gìn khoảng yên của riêng bạn.",
		}
	default:
		return moodstorydomain.LocalizedText{
			EN: "Stay steady and keep the next step uncomplicated.",
			VI: "Hãy giữ nhịp ổn định và để bước tiếp theo thật gọn, thật nhẹ.",
		}
	}
}

func gradientFor(mood moodpackdomain.Mood) string {
	if strings.TrimSpace(mood.Theme) != "" {
		return mood.Theme
	}

	switch mood.Key {
	case "energetic_bright":
		return "sunny-gold"
	case "chill_reflective":
		return "rainy-blue"
	case "calm_soft":
		return "cloudy-silver"
	case "intense_moody":
		return "storm-indigo"
	case "cozy_gentle":
		return "misty-latte"
	case "dreamy_quiet":
		return "fog-lilac"
	default:
		return "soft-neutral"
	}
}

func pickHighlightTrack(tracks []moodpackdomain.MusicTrack) *moodstorydomain.Track {
	if len(tracks) == 0 {
		return nil
	}

	for _, track := range tracks {
		if strings.TrimSpace(track.Title) == "" || strings.TrimSpace(track.Artist) == "" {
			continue
		}

		return &moodstorydomain.Track{
			Title:  track.Title,
			Artist: track.Artist,
			URL:    track.TrackURL,
		}
	}

	return nil
}

func appendUnique(values []string, value string) []string {
	for _, item := range values {
		if item == value {
			return values
		}
	}

	return append(values, value)
}

func moodAdjective(moodKey string) string {
	switch moodKey {
	case "energetic_bright":
		return "bright"
	case "chill_reflective":
		return "reflective"
	case "calm_soft":
		return "soft"
	case "intense_moody":
		return "moody"
	case "cozy_gentle":
		return "gentle"
	case "dreamy_quiet":
		return "dreamy"
	default:
		return "balanced"
	}
}

func headlineWeatherWord(weather moodpackdomain.Weather) string {
	description := strings.ToLower(strings.TrimSpace(weather.Description))

	switch {
	case strings.Contains(description, "clear"):
		return "clear"
	case strings.Contains(description, "thunder"):
		return "stormy"
	case strings.Contains(description, "drizzle"):
		return "drizzly"
	case strings.Contains(description, "rain"):
		return "rainy"
	case strings.Contains(description, "cloud"):
		return "cloudy"
	case strings.Contains(description, "mist"), strings.Contains(description, "fog"), strings.Contains(description, "haze"), strings.Contains(description, "smoke"):
		return "misty"
	default:
		switch strings.ToUpper(strings.TrimSpace(weather.Main)) {
		case "CLEAR":
			return "clear"
		case "RAIN":
			return "rainy"
		case "CLOUDS":
			return "cloudy"
		case "THUNDERSTORM":
			return "stormy"
		default:
			return ""
		}
	}
}

func weatherSceneEN(description string) string {
	description = strings.ToLower(strings.TrimSpace(description))

	switch {
	case strings.Contains(description, "broken clouds"):
		return "under broken clouds"
	case strings.Contains(description, "scattered clouds"):
		return "under scattered clouds"
	case strings.Contains(description, "few clouds"):
		return "under a light scatter of clouds"
	case strings.Contains(description, "overcast clouds"):
		return "under a muted blanket of clouds"
	case strings.Contains(description, "cloud"):
		return "under soft clouds"
	case strings.Contains(description, "clear sky"):
		return "under a clear sky"
	case strings.Contains(description, "drizzle"):
		return "through a soft drizzle"
	case strings.Contains(description, "rain"):
		return "in light rain"
	case strings.Contains(description, "thunderstorm"):
		return "under stormy skies"
	case strings.Contains(description, "mist"), strings.Contains(description, "fog"):
		return "through a quiet mist"
	case strings.Contains(description, "haze"), strings.Contains(description, "smoke"):
		return "through a hazy atmosphere"
	default:
		return "with the weather settling in gently"
	}
}

func weatherSceneVI(description string) string {
	description = strings.ToLower(strings.TrimSpace(description))

	switch {
	case strings.Contains(description, "broken clouds"):
		return "dưới những đám mây lững lờ"
	case strings.Contains(description, "scattered clouds"):
		return "dưới những cụm mây thưa"
	case strings.Contains(description, "few clouds"):
		return "dưới vài đám mây nhẹ"
	case strings.Contains(description, "overcast clouds"):
		return "dưới một bầu trời mây phủ"
	case strings.Contains(description, "cloud"):
		return "dưới bầu trời nhiều mây"
	case strings.Contains(description, "clear sky"):
		return "dưới bầu trời trong"
	case strings.Contains(description, "drizzle"):
		return "giữa màn mưa bụi nhẹ"
	case strings.Contains(description, "rain"):
		return "trong cơn mưa nhẹ"
	case strings.Contains(description, "thunderstorm"):
		return "dưới bầu trời giông"
	case strings.Contains(description, "mist"), strings.Contains(description, "fog"):
		return "giữa làn sương mỏng"
	case strings.Contains(description, "haze"), strings.Contains(description, "smoke"):
		return "giữa lớp không khí mờ nhẹ"
	default:
		return "giữa nhịp thời tiết hiện tại"
	}
}
