package main

import (
	"log/slog"
	"net/http"
	"os"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /version", handleVersion)

	slog.Info("server starting", "port", port, "version", version)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
