package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// handleRoot returns a greeting message.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello from Ocean!")); err != nil {
		slog.Error("write response failed", "handler", "root", "error", err)
	}
}

// handleHealth returns the health status of the service.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, "health", map[string]string{
		"status": "ok",
	})
}

// handleVersion returns the current version of the service.
func handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, "version", map[string]string{
		"version": version,
	})
}

// writeJSON serializes v as JSON and writes it to w, logging any error.
func writeJSON(w http.ResponseWriter, handler string, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response failed", "handler", handler, "error", err)
	}
}
