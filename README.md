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

Health checks:

```bash
curl "http://localhost:8080/kaithhealth"
curl "http://localhost:8080/healthz"
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

## Deploy To Leapcell

The current app shape works well on Leapcell as a Go service.

Use these values when creating the service in the Leapcell dashboard:

- Runtime: `Go`
- Root Directory: repository root
- Build Command: `sh ./scripts/leapcell-build.sh`
- Start Command: `./app`
- Port: `8080`

Set these environment variables in the Leapcell service settings:

- `OPENWEATHER_API_KEY` required
- `OPENWEATHER_BASE_URL=https://api.openweathermap.org/data/2.5`
- `ZENQUOTES_BASE_URL=https://zenquotes.io/api`
- `DEEZER_BASE_URL=https://api.deezer.com`

Deployment notes:

- Leapcell serverless startup polls `/kaithhealth`, so this repo exposes that endpoint.
- The service filesystem is effectively read-only except for `/tmp`, which is fine for this API because it is stateless.
- Do not rely on local `.env` files in production. Configure secrets in the Leapcell dashboard instead.
- If the build still fails, open the detailed build log and look for the first `go` error line, not just the final `failed to build image` summary.
