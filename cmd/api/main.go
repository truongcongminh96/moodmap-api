package main

import (
	"log"
	"net/http"

	"moodmap-api/internal/config"
	httptransport "moodmap-api/internal/transport/http"
)

func main() {
	cfg := config.Load()

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httptransport.NewRouter(cfg),
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
	}

	log.Printf("moodmap-api listening on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
