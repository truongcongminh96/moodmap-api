# moodmap-api

API for MoodMap that transforms live weather data into a frontend-ready mood pack with quotes, music, and activity suggestions.

## Stack

- Go
- OpenWeather for current weather
- ZenQuotes for quotes
- Deezer for music recommendations

## Run

1. Copy `.env.example` values into your local shell or `.env`.
2. Set at least `OPENWEATHER_API_KEY`.
3. Start the API:

```bash
go run ./cmd/api
```

The server runs on `http://localhost:8080` by default.
If a local `.env` file exists, the app now auto-loads it on startup for local development.

## Main Endpoint

```bash
curl "http://localhost:8080/api/v1/mood-pack?city=Hanoi&source=all"
```

Supported query params:

- `city` required
- `country` optional
- `units` optional: `metric` or `imperial`
- `source` optional: `quotes`, `music`, or `all`

## Notes

- `source=all` degrades gracefully if quote or music lookup fails.
- `source=quotes` and `source=music` return provider errors if the requested provider fails.
- Weather-to-mood mapping and activity suggestions live in backend business logic.
- Deezer search is used directly for the music MVP, so no extra music API key is required.
