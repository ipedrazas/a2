// Package checkutil provides common utilities for check implementations.
package checkutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// TruncateMessage limits a message to maxLen characters.
// It trims whitespace and appends "..." if truncated.
func TruncateMessage(msg string, maxLen int) string {
	msg = strings.TrimSpace(msg)
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen] + "..."
}

// Pluralize returns the singular form if count is 1, otherwise the plural form.
// This returns just the word without the count prefix.
// Example: Pluralize(1, "file", "files") returns "file"
// Example: Pluralize(3, "file", "files") returns "files"
func Pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// PluralizeCount returns "N word" with proper singular/plural form.
// Example: PluralizeCount(1, "file", "files") returns "1 file"
// Example: PluralizeCount(3, "file", "files") returns "3 files"
func PluralizeCount(count int, singular, plural string) string {
	if count == 1 {
		return "1 " + singular
	}
	return fmt.Sprintf("%d %s", count, plural)
}

// CommandResult holds the output from running a command.
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

// Output returns stderr if non-empty, otherwise stdout (trimmed).
// This follows the common pattern where errors go to stderr.
func (r *CommandResult) Output() string {
	output := strings.TrimSpace(r.Stderr)
	if output == "" {
		output = strings.TrimSpace(r.Stdout)
	}
	return output
}

// CombinedOutput returns both stdout and stderr concatenated.
func (r *CommandResult) CombinedOutput() string {
	return strings.TrimSpace(r.Stdout + r.Stderr)
}

// Success returns true if the command exited with code 0.
func (r *CommandResult) Success() bool {
	return r.Err == nil && r.ExitCode == 0
}

// RunCommand executes a command in the specified directory and captures output.
// It returns a CommandResult with stdout, stderr, and exit information.
func RunCommand(dir, name string, args ...string) *CommandResult {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}

	// Extract exit code if available
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1 // Unknown error (e.g., command not found)
		}
	}

	return result
}

// ToolAvailable checks if a command-line tool is available in PATH.
// Returns true if the tool can be found.
func ToolAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// ToolNotFoundError checks if an error indicates a tool was not found.
func ToolNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "executable file not found")
}
