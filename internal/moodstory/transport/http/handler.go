package http

import (
	"errors"
	"net/http"

	moodpackdomain "moodmap-api/internal/moodpack/domain"
	moodstorydomain "moodmap-api/internal/moodstory/domain"
	"moodmap-api/internal/moodstory/presenter"
	"moodmap-api/internal/moodstory/service"
	"moodmap-api/internal/platform/apperror"
	"moodmap-api/internal/platform/httpx"
)

type MoodStoryHandler struct {
	moodStoryService *service.MoodStoryService
}

func NewMoodStoryHandler(moodStoryService *service.MoodStoryService) *MoodStoryHandler {
	return &MoodStoryHandler{moodStoryService: moodStoryService}
}

func (h *MoodStoryHandler) GetMoodStory(w http.ResponseWriter, r *http.Request) {
	input := moodstorydomain.GetMoodStoryInput{
		City:    httpx.NormalizeQueryValue(r.URL.Query().Get("city")),
		Country: httpx.NormalizeQueryValue(r.URL.Query().Get("country")),
		Units:   moodpackdomain.Units(httpx.NormalizeQueryValue(r.URL.Query().Get("units"))),
	}

	story, err := h.moodStoryService.GetMoodStory(r.Context(), input)
	if err != nil {
		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			httpx.WriteJSON(w, appErr.StatusCode, presenter.Failure(appErr.Code, appErr.Message))
			return
		}

		httpx.WriteJSON(
			w,
			apperror.ErrInternalServer.StatusCode,
			presenter.Failure(apperror.ErrInternalServer.Code, apperror.ErrInternalServer.Message),
		)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, presenter.Success(story))
}
