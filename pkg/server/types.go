package server

import (
	"time"

	"github.com/ipedrazas/a2/pkg/output"
)

// CheckRequest is the request payload for submitting a check.
type CheckRequest struct {
	URL         string   `json:"url"`                    // GitHub URL (required)
	Languages   []string `json:"languages,omitempty"`    // Override auto-detection
	Profile     string   `json:"profile,omitempty"`      // Application profile (cli, api, library, desktop)
	Target      string   `json:"target,omitempty"`       // Maturity target (poc, production)
	SkipChecks  []string `json:"skip_checks,omitempty"`  // Checks to skip
	TimeoutSecs int      `json:"timeout_secs,omitempty"` // Per-check timeout (0 = no timeout)
	Verbose     bool     `json:"verbose,omitempty"`      // Show command output for failed/warning checks
}

// CheckResponse is the response when a check is submitted.
type CheckResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// JobStatus represents the current status of a job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// Job represents a single check job.
type Job struct {
	ID           string             `json:"id"`
	Status       JobStatus          `json:"status"`
	GitHubURL    string             `json:"github_url"`
	WorkspaceDir string             `json:"-"` // Not exposed in JSON
	Request      CheckRequest       `json:"request"`
	SubmittedAt  time.Time          `json:"submitted_at"`
	StartedAt    *time.Time         `json:"started_at,omitempty"`
	CompletedAt  *time.Time         `json:"completed_at,omitempty"`
	Result       *output.JSONOutput `json:"result,omitempty"`
	Error        string             `json:"error,omitempty"`
}

// JobResponse is the response for GET /api/check/{id}.
type JobResponse struct {
	JobID       string             `json:"job_id"`
	Status      JobStatus          `json:"status"`
	GitHubURL   string             `json:"github_url"`
	SubmittedAt time.Time          `json:"submitted_at"`
	StartedAt   *time.Time         `json:"started_at,omitempty"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
	Request     CheckRequest       `json:"request"`
	Result      *output.JSONOutput `json:"result,omitempty"`
	Error       string             `json:"error,omitempty"`
}

// ToJobResponse converts a Job to a JobResponse (excludes WorkspaceDir).
func (j *Job) ToJobResponse() JobResponse {
	return JobResponse{
		JobID:       j.ID,
		Status:      j.Status,
		GitHubURL:   j.GitHubURL,
		SubmittedAt: j.SubmittedAt,
		StartedAt:   j.StartedAt,
		CompletedAt: j.CompletedAt,
		Request:     j.Request,
		Result:      j.Result,
		Error:       j.Error,
	}
}

// HealthResponse is the response for GET /health.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// GitHubRepo represents a parsed GitHub repository URL.
type GitHubRepo struct {
	Owner    string
	Repo     string
	Branch   string // Empty if using default branch
	IsSSH    bool
	Original string // Original URL for reference
}
