package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const uiPath = "/usr/local/share/a2/ui"

// getUIFS returns a filesystem for serving the UI.
// In production, the UI is copied to /usr/local/share/a2/ui in the Docker image.
func getUIFS() (http.FileSystem, error) {
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		// No UI available
		return nil, os.ErrNotExist
	}

	return http.Dir(uiPath), nil
}

// newSPAHandler creates an SPA handler that serves the UI.
func newSPAHandler(fs http.FileSystem) http.Handler {
	return &spaHandler{
		fs:         fs,
		uiPath:     uiPath,
		fileServer: http.FileServer(fs),
	}
}

// spaHandler serves the SPA with proper routing.
type spaHandler struct {
	fs         http.FileSystem
	uiPath     string
	fileServer http.Handler
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Don't handle API or health endpoints
	if r.URL.Path == "/health" || strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Clean the path to prevent path traversal
	cleanPath := filepath.Clean(r.URL.Path)
	fullPath := filepath.Join(h.uiPath, cleanPath)

	// Ensure the resolved path is within the UI directory
	if !strings.HasPrefix(fullPath, h.uiPath) {
		http.NotFound(w, r)
		return
	}

	// Check if the requested path exists and is a file
	info, err := os.Stat(fullPath) // #nosec G703 -- path is sanitized above
	if err == nil && !info.IsDir() {
		// File exists, serve it
		h.fileServer.ServeHTTP(w, r)
		return
	}

	// Either doesn't exist or is a directory - serve index.html for SPA
	http.ServeFile(w, r, filepath.Join(h.uiPath, "index.html"))
}
