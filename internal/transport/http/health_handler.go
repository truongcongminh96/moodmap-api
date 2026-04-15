package http

import (
	"net/http"

	"moodmap-api/internal/platform/httpx"
)

type healthResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, healthResponse{
		OK:      true,
		Service: "moodmap-api",
	})
}
