package presenter

import (
	"time"

	"moodmap-api/internal/moodpack/domain"
)

type MoodPackResponse struct {
	Success bool             `json:"success"`
	Data    *domain.MoodPack `json:"data,omitempty"`
	Meta    *Meta            `json:"meta,omitempty"`
	Error   *ErrorPayload    `json:"error,omitempty"`
}

type Meta struct {
	RequestedAt time.Time `json:"requestedAt"`
	Sources     []string  `json:"sources"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(pack *domain.MoodPack) MoodPackResponse {
	return MoodPackResponse{
		Success: true,
		Data:    pack,
		Meta: &Meta{
			RequestedAt: pack.RequestedAt,
			Sources:     pack.Sources,
		},
	}
}

func Failure(code, message string) MoodPackResponse {
	return MoodPackResponse{
		Success: false,
		Error: &ErrorPayload{
			Code:    code,
			Message: message,
		},
	}
}
