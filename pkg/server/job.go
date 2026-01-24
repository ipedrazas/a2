package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ipedrazas/a2/pkg/output"
)

// JobStore manages job storage with thread-safe access.
type JobStore struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

// NewJobStore creates a new job store.
func NewJobStore() *JobStore {
	return &JobStore{
		jobs: make(map[string]*Job),
	}
}

// CreateJob creates a new job with a unique ID.
func (s *JobStore) CreateJob(githubURL string, req CheckRequest, workspaceDir string) *Job {
	id := uuid.New().String()
	now := time.Now()

	job := &Job{
		ID:           id,
		Status:       JobStatusPending,
		GitHubURL:    githubURL,
		WorkspaceDir: workspaceDir,
		Request:      req,
		SubmittedAt:  now,
	}

	s.mu.Lock()
	s.jobs[id] = job
	s.mu.Unlock()

	return job
}

// Get retrieves a job by ID.
func (s *JobStore) Get(id string) (*Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	return job, ok
}

// UpdateStatus updates the status of a job.
func (s *JobStore) UpdateStatus(id string, status JobStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return
	}

	job.Status = status

	if status == JobStatusRunning && job.StartedAt == nil {
		now := time.Now()
		job.StartedAt = &now
	}

	if status == JobStatusCompleted || status == JobStatusFailed {
		now := time.Now()
		job.CompletedAt = &now
	}
}

// SetResult sets the result of a job and marks it as completed.
func (s *JobStore) SetResult(id string, result *output.JSONOutput) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return
	}

	job.Result = result
	job.Status = JobStatusCompleted
	now := time.Now()
	job.CompletedAt = &now
}

// SetError sets an error on a job and marks it as failed.
func (s *JobStore) SetError(id string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return
	}

	job.Error = err.Error()
	job.Status = JobStatusFailed
	now := time.Now()
	job.CompletedAt = &now
}

// Delete removes a job from the store.
func (s *JobStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
}

// List returns all jobs (for admin/debugging purposes).
func (s *JobStore) List() []*Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// CleanupOldJobs removes jobs older than the specified duration.
// Returns the number of jobs cleaned up.
func (s *JobStore) CleanupOldJobs(maxAge time.Duration) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	var count int

	for id, job := range s.jobs {
		// Only clean up completed/failed jobs
		if job.Status != JobStatusCompleted && job.Status != JobStatusFailed {
			continue
		}

		if job.CompletedAt != nil && job.CompletedAt.Before(cutoff) {
			delete(s.jobs, id)
			count++
		}
	}

	return count
}

// ValidateJobID checks if the job ID format is valid.
func ValidateJobID(id string) error {
	if id == "" {
		return fmt.Errorf("job ID cannot be empty")
	}
	// Check if it's a valid UUID
	_, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid job ID format: %w", err)
	}
	return nil
}
