package apperror

import "net/http"

type Error struct {
	Code       string
	Message    string
	StatusCode int
}

func (e *Error) Error() string {
	return e.Message
}

func New(code, message string, statusCode int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

var (
	ErrCityRequired      = New("CITY_REQUIRED", "City is required.", http.StatusBadRequest)
	ErrCityNotFound      = New("CITY_NOT_FOUND", "We could not find weather data for this city.", http.StatusNotFound)
	ErrWeatherProvider   = New("WEATHER_PROVIDER_ERROR", "Weather data is temporarily unavailable.", http.StatusBadGateway)
	ErrQuoteProvider     = New("QUOTE_PROVIDER_ERROR", "Quote data is temporarily unavailable.", http.StatusBadGateway)
	ErrMusicProvider     = New("MUSIC_PROVIDER_ERROR", "Music recommendations are temporarily unavailable.", http.StatusBadGateway)
	ErrMoodMappingFailed = New("MOOD_MAPPING_FAILED", "We could not map the current weather to a mood.", http.StatusUnprocessableEntity)
	ErrInternalServer    = New("INTERNAL_SERVER_ERROR", "Something went wrong while building the mood pack.", http.StatusInternalServerError)
)
