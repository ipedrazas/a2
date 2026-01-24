package server

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// githubURLPatterns matches various GitHub URL formats.
var (
	githubHTTPPattern  = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+?)(\.git)?/?$`)
	githubHTTPBranch   = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+?)/tree/([^/]+)/?.*$`)
	githubSSHPattern   = regexp.MustCompile(`^git@github\.com:([^/]+)/([^/]+?)(\.git)?$`)
	githubShortPattern = regexp.MustCompile(`^([^/]+)/([^/]+)$`)
)

// ParseGitHubURL parses a GitHub URL and extracts owner, repo, and branch.
func ParseGitHubURL(rawURL string) (*GitHubRepo, error) {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	// Try HTTP/HTTPS pattern with branch
	if matches := githubHTTPBranch.FindStringSubmatch(rawURL); len(matches) >= 4 {
		return &GitHubRepo{
			Owner:    matches[1],
			Repo:     strings.TrimSuffix(matches[2], ".git"),
			Branch:   matches[3],
			IsSSH:    false,
			Original: rawURL,
		}, nil
	}

	// Try HTTP/HTTPS pattern
	if matches := githubHTTPPattern.FindStringSubmatch(rawURL); len(matches) >= 3 {
		return &GitHubRepo{
			Owner:    matches[1],
			Repo:     strings.TrimSuffix(matches[2], ".git"),
			Branch:   "",
			IsSSH:    false,
			Original: rawURL,
		}, nil
	}

	// Try SSH pattern
	if matches := githubSSHPattern.FindStringSubmatch(rawURL); len(matches) >= 3 {
		return &GitHubRepo{
			Owner:    matches[1],
			Repo:     strings.TrimSuffix(matches[2], ".git"),
			Branch:   "",
			IsSSH:    true,
			Original: rawURL,
		}, nil
	}

	// Try short pattern (owner/repo)
	if matches := githubShortPattern.FindStringSubmatch(rawURL); len(matches) >= 3 {
		return &GitHubRepo{
			Owner:    matches[1],
			Repo:     matches[2],
			Branch:   "",
			IsSSH:    false,
			Original: rawURL,
		}, nil
	}

	return nil, fmt.Errorf("invalid GitHub URL format: %s (expected: https://github.com/owner/repo or owner/repo)", rawURL)
}

// CloneURL returns the clone URL for the repo.
func (gr *GitHubRepo) CloneURL() string {
	if gr.IsSSH {
		return fmt.Sprintf("git@github.com:%s/%s.git", gr.Owner, gr.Repo)
	}
	return fmt.Sprintf("https://github.com/%s/%s.git", gr.Owner, gr.Repo)
}

// CloneBranch clones a specific branch if specified, otherwise default.
func (gr *GitHubRepo) CloneBranch() string {
	return gr.Branch
}

// Validate validates the GitHubRepo fields.
func (gr *GitHubRepo) Validate() error {
	if gr.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if gr.Repo == "" {
		return fmt.Errorf("repo cannot be empty")
	}
	// Check for path traversal
	if strings.Contains(gr.Owner, "..") || strings.Contains(gr.Repo, "..") {
		return fmt.Errorf("invalid characters in owner or repo")
	}
	return nil
}

// CloneRepository clones a GitHub repository to the specified workspace directory.
// It performs a shallow clone for speed and minimal disk usage.
func CloneRepository(repo *GitHubRepo, workspaceDir string) error {
	if err := repo.Validate(); err != nil {
		return fmt.Errorf("invalid repository: %w", err)
	}

	// Build git clone command
	args := []string{"clone", "--depth", "1", "--single-branch"}

	// Add branch if specified
	if repo.Branch != "" {
		args = append(args, "--branch", repo.Branch)
	}

	// Add repository URL
	args = append(args, repo.CloneURL())

	// Add destination directory
	args = append(args, workspaceDir)

	// Execute git clone
	cmd := exec.Command("git", args...)
	cmd.Stdout = nil // Suppress stdout
	cmd.Stderr = nil // Suppress stderr (will return error if failed)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

// ValidateGitHubURL performs additional validation on parsed URLs.
func ValidateGitHubURL(u *url.URL) error {
	if u.Host != "github.com" {
		return fmt.Errorf("only github.com is supported")
	}

	// Check for suspicious characters in path
	if strings.Contains(u.Path, "..") {
		return fmt.Errorf("path traversal detected")
	}

	return nil
}

// GetRepoPath returns the path where the repository was cloned.
func GetRepoPath(workspaceDir string) string {
	// For shallow clones, the repo is directly in the workspace directory
	return workspaceDir
}

// CleanWorkspace removes the workspace directory.
func CleanWorkspace(workspaceDir string) error {
	// Security check: ensure we're not deleting something important
	if workspaceDir == "" || workspaceDir == "/" || workspaceDir == "." {
		return fmt.Errorf("refusing to delete unsafe path: %s", workspaceDir)
	}

	absPath, err := filepath.Abs(workspaceDir)
	if err != nil {
		return fmt.Errorf("failed to resolve workspace path: %w", err)
	}

	// Additional safety check
	if !strings.Contains(absPath, "job-") {
		return fmt.Errorf("workspace directory does not appear to be a job directory: %s", absPath)
	}

	return os.RemoveAll(absPath)
}
