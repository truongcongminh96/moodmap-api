package presenter

import "moodmap-api/internal/moodstory/domain"

type MoodStoryResponse struct {
	Success bool              `json:"success"`
	Data    *domain.MoodStory `json:"data,omitempty"`
	Error   *ErrorPayload     `json:"error,omitempty"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(story *domain.MoodStory) MoodStoryResponse {
	return MoodStoryResponse{
		Success: true,
		Data:    story,
	}
}

func Failure(code, message string) MoodStoryResponse {
	return MoodStoryResponse{
		Success: false,
		Error: &ErrorPayload{
			Code:    code,
			Message: message,
		},
	}
}
