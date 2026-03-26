package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handleRoot returns a greeting message.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, GitHub Actions!")
}

// handleHealth returns the health status of the service.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// handleVersion returns the current version of the service.
func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": version,
	})
}
