# Global Hack Week: API .MoodMap

MoodMap API turns live weather data into a frontend-ready emotional experience. It combines current weather with a mapped mood, a quote, music suggestions, activity ideas, and a story-style summary that can be used directly in a web or mobile app.

## Features

- Get a `mood-pack` based on live weather data for a city
- Add optional quote and music recommendations
- Generate a richer `mood-story` response for storytelling UI
- Return consistent JSON success and error payloads
- Gracefully degrade when optional content providers fail

## Tech Stack

- Go
- OpenWeather for current weather
- ZenQuotes for quotes
- Deezer for music recommendations

## Project Structure

```text
cmd/api/                     # API entrypoint
internal/moodpack/          # Mood pack domain, services, providers, HTTP handler
internal/moodstory/         # Mood story domain, services, HTTP handler
internal/transport/http/    # Router and health endpoints
internal/config/            # Environment-based configuration
scripts/                    # Deployment/build scripts
```

## Environment Variables

The app auto-loads a local `.env` file if it exists.

Required:

- `OPENWEATHER_API_KEY`

Optional:

- `PORT` default: `8080`
- `OPENWEATHER_BASE_URL` default: `https://api.openweathermap.org/data/2.5`
- `ZENQUOTES_BASE_URL` default: `https://zenquotes.io/api`
- `DEEZER_BASE_URL` default: `https://api.deezer.com`

## Run Locally

1. Copy `.env.example` to `.env`.
2. Add your `OPENWEATHER_API_KEY`.
3. Start the API:

```bash
go run ./cmd/api
```

The server runs on `http://localhost:8080` by default.

## API Base URL

Local:

```text
http://localhost:8080
```

## API Reference

### Health Check

`GET /healthz`

`GET /kaithhealth`

Example:

```bash
curl "http://localhost:8080/healthz"
```

Success response:

```json
{
  "ok": true,
  "service": "moodmap-api"
}
```

### Get Mood Pack

`GET /api/v1/mood-pack`

Returns the current weather mood pack for a city, including mood data, activity suggestions, and optional quote/music content.

Query parameters:

- `city` required. Example: `Hanoi`
- `country` optional. Country code or country text
- `units` optional. Supported values: `metric`, `imperial`. Default: `metric`
- `source` optional. Supported values: `quotes`, `music`, `all`. Default: `all`

Example requests:

```bash
curl "http://localhost:8080/api/v1/mood-pack?city=Hanoi"
```

```bash
curl "http://localhost:8080/api/v1/mood-pack?city=Ho%20Chi%20Minh&units=metric&source=all"
```

Example success response:

```json
{
  "success": true,
  "data": {
    "location": {
      "city": "Hanoi",
      "country": "VN"
    },
    "weather": {
      "main": "Clouds",
      "description": "broken clouds",
      "temperature": 29.4,
      "feelsLike": 33.1,
      "humidity": 74,
      "icon": "04d"
    },
    "mood": {
      "key": "calm_soft",
      "label": "Calm Soft",
      "theme": "cloudy-silver",
      "confidence": 0.82
    },
    "quote": {
      "text": "Keep going. Everything you need will come to you at the perfect time.",
      "author": "Unknown"
    },
    "music": [
      {
        "title": "Sunset Lover",
        "artist": "Petit Biscuit",
        "trackUrl": "https://www.deezer.com/track/...",
        "source": "deezer"
      }
    ],
    "activities": [
      "Take a slow walk",
      "Journal for ten minutes",
      "Listen to a soft playlist"
    ],
    "summary": "Hanoi today feels soft and balanced with broken clouds. It is a nice day to read, plan gently, and ease into your rhythm."
  },
  "meta": {
    "requestedAt": "2026-04-16T03:00:00Z",
    "sources": [
      "openweather",
      "zenquotes",
      "deezer"
    ]
  }
}
```

Notes:

- `source=all` will still return a successful response even if quote or music lookup fails.
- `source=quotes` returns an error if the quote provider fails.
- `source=music` returns an error if the music provider fails.

### Get Mood Story

`GET /api/v1/mood-story`

Returns a more narrative response built on top of the mood pack. This endpoint is useful for cards, hero sections, and emotionally styled frontend experiences.

Query parameters:

- `city` required
- `country` optional
- `units` optional. Supported values: `metric`, `imperial`. Default: `metric`

Example request:

```bash
curl "http://localhost:8080/api/v1/mood-story?city=Hanoi"
```

Example success response:

```json
{
  "success": true,
  "data": {
    "city": "Hanoi",
    "country": "VN",
    "headline": "A calm cloudy afternoon in Hanoi",
    "mood": {
      "key": "calm_soft",
      "label": "Calm Soft",
      "theme": "cloudy-silver"
    },
    "visual": {
      "gradient": "cloudy-silver",
      "timeOfDay": "afternoon"
    },
    "highlight": {
      "quote": {
        "text": "Keep going. Everything you need will come to you at the perfect time.",
        "author": "Unknown"
      },
      "track": {
        "title": "Sunset Lover",
        "artist": "Petit Biscuit",
        "url": "https://www.deezer.com/track/..."
      }
    },
    "story": {
      "en": "Hanoi feels calm and unhurried today with broken clouds. It is the kind of weather that invites slower thinking, light music, and a gentler pace.",
      "vi": "Hanoi hom nay mang mot nhip dieu cham va diu voi troi nhieu may. Day la kieu thoi tiet khien ban muon suy nghi nhe nhang hon, nghe nhac khe hon, va song cham lai mot chut."
    },
    "bestMoment": {
      "en": "Quiet desk reset",
      "vi": "Mot nhip sap lai ban lam viec that yen"
    },
    "energyTip": {
      "en": "Keep the day light and low-pressure.",
      "vi": "Hay giu ngay hom nay nhe nhang va khong ap luc."
    },
    "meta": {
      "generatedAt": "2026-04-16T03:00:00Z",
      "sources": [
        "openweather",
        "zenquotes",
        "deezer",
        "system"
      ]
    }
  }
}
```

## Error Format

Both main endpoints return errors in this structure:

```json
{
  "success": false,
  "error": {
    "code": "CITY_REQUIRED",
    "message": "City is required."
  }
}
```

Known error codes:

- `CITY_REQUIRED` HTTP `400`
- `CITY_NOT_FOUND` HTTP `404`
- `WEATHER_PROVIDER_ERROR` HTTP `502`
- `QUOTE_PROVIDER_ERROR` HTTP `502`
- `MUSIC_PROVIDER_ERROR` HTTP `502`
- `MOOD_MAPPING_FAILED` HTTP `422`
- `INTERNAL_SERVER_ERROR` HTTP `500`

## Example cURL Commands

```bash
curl "http://localhost:8080/api/v1/mood-pack?city=Hanoi&source=all"
```

```bash
curl "http://localhost:8080/api/v1/mood-pack?city=Tokyo&units=imperial&source=music"
```

```bash
curl "http://localhost:8080/api/v1/mood-story?city=Da%20Nang"
```

## Deployment

### Deploy to Leapcell

Use these values when creating the service in the Leapcell dashboard:

- Runtime: `Go`
- Root Directory: repository root
- Build Command: `sh ./scripts/leapcell-build.sh`
- Start Command: `./app`
- Port: `8080`

Set these environment variables in Leapcell:

- `OPENWEATHER_API_KEY` required
- `OPENWEATHER_BASE_URL=https://api.openweathermap.org/data/2.5`
- `ZENQUOTES_BASE_URL=https://zenquotes.io/api`
- `DEEZER_BASE_URL=https://api.deezer.com`

Deployment notes:

- Leapcell serverless startup polls `/kaithhealth`, so this endpoint is included.
- The app is stateless and does not rely on writable local storage.
- Do not rely on `.env` in production. Configure secrets in the deployment platform instead.

### Docker

This repository also includes a production-ready `Dockerfile`.

## License

MIT
