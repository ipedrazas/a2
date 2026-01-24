package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WorkspaceManager manages temporary workspaces for cloning repositories.
type WorkspaceManager struct {
	baseDir      string
	cleanupAfter bool
}

// NewWorkspaceManager creates a new workspace manager.
func NewWorkspaceManager(baseDir string, cleanupAfter bool) *WorkspaceManager {
	return &WorkspaceManager{
		baseDir:      baseDir,
		cleanupAfter: cleanupAfter,
	}
}

// CreateWorkspace creates a new workspace directory for a job.
func (wm *WorkspaceManager) CreateWorkspace(jobID string) (string, error) {
	if err := os.MkdirAll(wm.baseDir, 0750); err != nil {
		return "", fmt.Errorf("failed to create base directory: %w", err)
	}

	// Create a unique workspace directory
	timestamp := time.Now().Unix()
	workspaceName := fmt.Sprintf("job-%s-%d", jobID, timestamp)
	workspaceDir := filepath.Join(wm.baseDir, workspaceName)

	if err := os.MkdirAll(workspaceDir, 0750); err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	return workspaceDir, nil
}

// Cleanup removes a workspace directory.
func (wm *WorkspaceManager) Cleanup(workspaceDir string) error {
	if !wm.cleanupAfter {
		return nil
	}

	// Security check: ensure path is within baseDir
	if !strings.HasPrefix(workspaceDir, wm.baseDir) {
		return fmt.Errorf("security: workspace path outside base directory")
	}

	if err := os.RemoveAll(workspaceDir); err != nil {
		return fmt.Errorf("failed to cleanup workspace %s: %w", workspaceDir, err)
	}

	return nil
}

// CleanupOld removes workspace directories older than the specified duration.
func (wm *WorkspaceManager) CleanupOld(maxAge time.Duration) (int, error) {
	entries, err := os.ReadDir(wm.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read base directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)
	var count int

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			path := filepath.Join(wm.baseDir, entry.Name())
			if err := os.RemoveAll(path); err != nil {
				continue
			}
			count++
		}
	}

	return count, nil
}

// ValidateWorkspacePath performs security validation of a workspace path.
func ValidateWorkspacePath(baseDir, workspaceDir string) error {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("failed to resolve base directory: %w", err)
	}

	absWorkspace, err := filepath.Abs(workspaceDir)
	if err != nil {
		return fmt.Errorf("failed to resolve workspace directory: %w", err)
	}

	rel, err := filepath.Rel(absBase, absWorkspace)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}

	// Check if the workspace is within the base directory
	if strings.HasPrefix(rel, "..") {
		return fmt.Errorf("security: workspace path escapes base directory")
	}

	return nil
}
