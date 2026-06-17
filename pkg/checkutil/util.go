// Package checkutil provides common utilities for check implementations.
package checkutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// ResultBuilder simplifies creating checker.Result with consistent fields.
// It eliminates the need to repeatedly set Name, ID, and Language for each result.
type ResultBuilder struct {
	name     string
	id       string
	language checker.Language
	command  string // last command recorded via RunCommand/RecordCommand, for display
}

// RunCommand runs a command via the package-level RunCommand and records the
// executed command line (with working dir) so any Result built afterwards
// reports exactly what a2 ran. Use this instead of checkutil.RunCommand inside
// a check's Run method whenever a ResultBuilder is in scope.
func (b *ResultBuilder) RunCommand(dir, name string, args ...string) *CommandResult {
	res := RunCommand(dir, name, args...)
	b.RecordCommand(res)
	return res
}

// RecordCommand records the command from an already-executed CommandResult so
// the next Result built by this builder reports it. Use this when the command
// is run indirectly (e.g. via a shared helper) rather than through b.RunCommand.
func (b *ResultBuilder) RecordCommand(res *CommandResult) {
	if res != nil {
		b.command = res.DisplayCommand()
	}
}

// NewResultBuilder creates a new ResultBuilder for a checker.
// The builder captures the checker's Name() and ID() so they don't need
// to be repeated in every result construction.
func NewResultBuilder(c checker.Checker, lang checker.Language) *ResultBuilder {
	return &ResultBuilder{
		name:     c.Name(),
		id:       c.ID(),
		language: lang,
	}
}

// ShortMessage derives a concise "what happened" summary from a longer reason.
func ShortMessage(reason string) string {
	msg := strings.TrimSpace(reason)
	if msg == "" {
		return ""
	}
	// Drop parenthetical guidance or detail.
	if idx := strings.Index(msg, " ("); idx != -1 {
		msg = strings.TrimSpace(msg[:idx])
	}
	// Drop trailing clause after a dash when it's likely extra detail.
	if idx := strings.Index(msg, " - "); idx != -1 {
		msg = strings.TrimSpace(msg[:idx])
	}
	// Drop trailing clause after a semicolon when it's likely extra detail.
	if idx := strings.Index(msg, "; "); idx != -1 && idx < 80 {
		msg = strings.TrimSpace(msg[:idx])
	}
	// Truncate very long messages.
	const maxLen = 80
	if len(msg) > maxLen {
		return msg[:maxLen] + "..."
	}
	return msg
}

// Pass creates a passing result with the given reason.
func (b *ResultBuilder) Pass(reason string) checker.Result {
	return checker.Result{
		Name:     b.name,
		ID:       b.id,
		Passed:   true,
		Status:   checker.Pass,
		Message:  ShortMessage(reason),
		Reason:   reason,
		Language: b.language,
		Command:  b.command,
	}
}

// Fail creates a failing result with the given reason.
// Fail status indicates a critical failure that may abort execution.
func (b *ResultBuilder) Fail(reason string) checker.Result {
	return checker.Result{
		Name:     b.name,
		ID:       b.id,
		Passed:   false,
		Status:   checker.Fail,
		Message:  ShortMessage(reason),
		Reason:   reason,
		Language: b.language,
		Command:  b.command,
	}
}

// Warn creates a warning result with the given reason.
// Warnings indicate issues but don't cause execution to abort.
func (b *ResultBuilder) Warn(reason string) checker.Result {
	return checker.Result{
		Name:     b.name,
		ID:       b.id,
		Passed:   false,
		Status:   checker.Warn,
		Message:  ShortMessage(reason),
		Reason:   reason,
		Language: b.language,
		Command:  b.command,
	}
}

// Info creates an informational result with the given reason.
// Info results don't affect the pass/fail status or maturity score.
func (b *ResultBuilder) Info(reason string) checker.Result {
	return checker.Result{
		Name:     b.name,
		ID:       b.id,
		Passed:   true,
		Status:   checker.Info,
		Message:  ShortMessage(reason),
		Reason:   reason,
		Language: b.language,
		Command:  b.command,
	}
}

// ToolNotInstalled creates an Info result indicating a tool is not installed.
// This standardizes the handling of missing tools across all checks.
// Use this when a check cannot run because an optional tool is not available.
// The installHint should provide installation instructions (e.g., "pip install black").
func (b *ResultBuilder) ToolNotInstalled(toolName, installHint string) checker.Result {
	message := toolName + " not installed"
	if installHint != "" {
		message += " (" + installHint + ")"
	}
	return checker.Result{
		Name:     b.name,
		ID:       b.id,
		Passed:   true,
		Status:   checker.Info,
		Message:  ShortMessage(message),
		Reason:   message,
		Language: b.language,
	}
}

// PassWithOutput creates a passing result with the given reason and raw output.
func (b *ResultBuilder) PassWithOutput(reason, rawOutput string) checker.Result {
	return checker.Result{
		Name:      b.name,
		ID:        b.id,
		Passed:    true,
		Status:    checker.Pass,
		Message:   ShortMessage(reason),
		Reason:    reason,
		Language:  b.language,
		RawOutput: rawOutput,
		Command:   b.command,
	}
}

// FailWithOutput creates a failing result with the given reason and raw output.
func (b *ResultBuilder) FailWithOutput(reason, rawOutput string) checker.Result {
	return checker.Result{
		Name:      b.name,
		ID:        b.id,
		Passed:    false,
		Status:    checker.Fail,
		Message:   ShortMessage(reason),
		Reason:    reason,
		Language:  b.language,
		RawOutput: rawOutput,
		Command:   b.command,
	}
}

// WarnWithOutput creates a warning result with the given reason and raw output.
func (b *ResultBuilder) WarnWithOutput(reason, rawOutput string) checker.Result {
	return checker.Result{
		Name:      b.name,
		ID:        b.id,
		Passed:    false,
		Status:    checker.Warn,
		Message:   ShortMessage(reason),
		Reason:    reason,
		Language:  b.language,
		RawOutput: rawOutput,
		Command:   b.command,
	}
}

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

	// Command is the command line that was executed (e.g. "govulncheck ./...").
	// Dir is the working directory it ran in (empty or "." means the repo root).
	// These let a2 surface exactly what it executed so failures are reproducible.
	Command string
	Dir     string
}

// DisplayCommand returns a copy-pasteable representation of the command,
// prefixed with a `cd` into the working directory when the command did not
// run at the repo root. Example: "cd nodeagent && govulncheck ./...".
func (r *CommandResult) DisplayCommand() string {
	if r.Command == "" {
		return ""
	}
	if r.Dir != "" && r.Dir != "." {
		return fmt.Sprintf("cd %s && %s", r.Dir, r.Command)
	}
	return r.Command
}

// formatCommandLine renders a command and its arguments as a single shell-like
// string, quoting any token that contains whitespace so it stays copy-pasteable.
func formatCommandLine(name string, args ...string) string {
	parts := make([]string, 0, len(args)+1)
	parts = append(parts, name)
	for _, a := range args {
		if strings.ContainsAny(a, " \t") {
			parts = append(parts, fmt.Sprintf("%q", a))
		} else {
			parts = append(parts, a)
		}
	}
	return strings.Join(parts, " ")
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
	cmd := exec.Command(name, args...) // #nosec G204 -- command execution is by design
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &CommandResult{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Err:     err,
		Command: formatCommandLine(name, args...),
		Dir:     dir,
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
