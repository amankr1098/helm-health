package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/amankr1098/helm-health/internal/output"
	rel "github.com/amankr1098/helm-health/internal/release"
)

// Start runs the helm-health HTTP API server on the given address (e.g. ":8080").
func Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/releases", handleReleases)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("helm-health server listening on %s", addr)
	return srv.ListenAndServe()
}

// handleHealth returns the full health report for a single release.
// GET /api/health?release=<name>&namespace=<ns>
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	releaseName := r.URL.Query().Get("release")
	if releaseName == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'release' is required")
		return
	}
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	result, err := rel.FetchHelmRelease(releaseName, namespace)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// A missing release is a 404; a found-but-unhealthy release is still a
	// successful health check and returns 200 with details in the body.
	status := http.StatusOK
	if result.Status == output.StatusNotFound {
		status = http.StatusNotFound
	}
	writeJSON(w, status, result)
}

// handleReleases lists Helm releases the dashboard can pick from.
// GET /api/releases?namespace=<ns>  (namespace omitted lists all namespaces)
func handleReleases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	namespace := r.URL.Query().Get("namespace")
	releases, err := rel.ListReleases(namespace)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"releases": releases})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
