package http

import (
	"errors"
	"net/http"

	"moodmap-api/internal/moodpack/domain"
	"moodmap-api/internal/moodpack/presenter"
	"moodmap-api/internal/moodpack/service"
	"moodmap-api/internal/platform/apperror"
	"moodmap-api/internal/platform/httpx"
)

type MoodHandler struct {
	moodService *service.MoodService
}

func NewMoodHandler(moodService *service.MoodService) *MoodHandler {
	return &MoodHandler{moodService: moodService}
}

func (h *MoodHandler) GetMoodPack(w http.ResponseWriter, r *http.Request) {
	input := domain.GetMoodPackInput{
		City:    httpx.NormalizeQueryValue(r.URL.Query().Get("city")),
		Country: httpx.NormalizeQueryValue(r.URL.Query().Get("country")),
		Units:   domain.Units(httpx.NormalizeQueryValue(r.URL.Query().Get("units"))),
		Source:  domain.ContentSource(httpx.NormalizeQueryValue(r.URL.Query().Get("source"))),
	}

	pack, err := h.moodService.GetMoodPack(r.Context(), input)
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

	httpx.WriteJSON(w, http.StatusOK, presenter.Success(pack))
}
