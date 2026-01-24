package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server.
type Server struct {
	router       *mux.Router
	httpServer   *http.Server
	jobStore     *JobStore
	queue        *JobQueue
	workspaceDir string
	cleanupAfter bool
}

// NewServer creates a new HTTP server.
func NewServer(host string, port int, jobStore *JobStore, queue *JobQueue, workspaceDir string, cleanupAfter bool) *Server {
	router := mux.NewRouter()
	s := &Server{
		router:       router,
		jobStore:     jobStore,
		queue:        queue,
		workspaceDir: workspaceDir,
		cleanupAfter: cleanupAfter,
	}

	// Set up routes
	s.setupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// setupRoutes configures all HTTP routes.
func (s *Server) setupRoutes() {
	api := s.router.PathPrefix("/api").Subrouter()

	// API endpoints
	api.HandleFunc("/check", s.handleSubmitCheck).Methods("POST")
	api.HandleFunc("/check/{id}", s.handleGetCheck).Methods("GET")
	api.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Serve static files (UI) at root
	s.router.PathPrefix("/").Handler(s.fileServerHandler())
}

// ListenAndServe starts the server.
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// handleHealth returns the health status.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(HealthResponse{
		Status: "healthy",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleSubmitCheck handles POST /api/check.
func (s *Server) handleSubmitCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Validate URL
	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "URL is required",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse and validate GitHub URL
	if _, err := ParseGitHubURL(req.URL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid GitHub URL: " + err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create workspace
	wm := NewWorkspaceManager(s.workspaceDir, s.cleanupAfter)
	workspaceDir, err := wm.CreateWorkspace("temp")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create workspace: " + err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create job
	job := s.jobStore.CreateJob(req.URL, req, workspaceDir)

	// Enqueue job
	if err := s.queue.Enqueue(job); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to queue job: " + err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return response
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(CheckResponse{
		JobID:   job.ID,
		Status:  string(JobStatusPending),
		Message: "Job queued successfully",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetCheck handles GET /api/check/{id}.
func (s *Server) handleGetCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get job ID from URL
	vars := mux.Vars(r)
	jobID := vars["id"]

	// Validate job ID
	if err := ValidateJobID(jobID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid job ID: " + err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get job from store
	job, ok := s.jobStore.Get(jobID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Job not found",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return job response
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(job.ToJobResponse())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// fileServerHandler serves static files for the UI.
// Falls back to a simple HTML placeholder if UI is not available.
func (s *Server) fileServerHandler() http.Handler {
	// Try to get UI filesystem
	uiFS, err := getUIFS()
	if err == nil && uiFS != nil {
		// Use SPA handler for proper routing
		return newSPAHandler(uiFS)
	}

	// Fallback to simple placeholder HTML
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
	<title>A2 Server</title>
	<style>
		body { font-family: system-ui, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
		h1 { color: #333; }
		.info { background: #f0f0f0; padding: 15px; border-radius: 5px; }
		.endpoint { background: #e8f4f8; padding: 10px; margin: 10px 0; border-left: 3px solid #007bff; }
		code { background: #f4f4f4; padding: 2px 5px; border-radius: 3px; }
	</style>
</head>
<body>
	<h1>A2 Server</h1>
	<div class="info">
		<p>A2 code quality checker server is running!</p>
		<p><strong>Note:</strong> UI is not available. Build the UI with <code>cd ui && npm run build</code></p>
	</div>
	<h2>API Endpoints</h2>
	<div class="endpoint">
		<strong>POST /api/check</strong>
		<p>Submit a GitHub URL for checking</p>
		<code>curl -X POST http://localhost:8080/api/check -H "Content-Type: application/json" -d '{"url":"https://github.com/user/repo"}'</code>
	</div>
	<div class="endpoint">
		<strong>GET /api/check/{id}</strong>
		<p>Get check status and results</p>
		<code>curl http://localhost:8080/api/check/{job_id}</code>
	</div>
	<div class="endpoint">
		<strong>GET /health</strong>
		<p>Health check endpoint</p>
		<code>curl http://localhost:8080/health</code>
	</div>
</body>
</html>`)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.NotFound(w, r)
	})
}
