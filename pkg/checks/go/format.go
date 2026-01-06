package gocheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// FormatCheck verifies that all Go code is properly formatted.
type FormatCheck struct{}

func (c *FormatCheck) ID() string   { return "go:format" }
func (c *FormatCheck) Name() string { return "Go Format" }

func (c *FormatCheck) Run(path string) (checker.Result, error) {
	// Run gofmt -l to list files that need formatting
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// gofmt returns non-zero if there's an error (not just unformatted files)
		if stderr.Len() > 0 {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   false,
				Status:   checker.Warn,
				Message:  "gofmt error: " + stderr.String(),
				Language: checker.LangGo,
			}, nil
		}
	}

	output := strings.TrimSpace(stdout.String())
	if output != "" {
		// Count unformatted files
		files := strings.Split(output, "\n")
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  "Unformatted files: " + strings.Join(files, ", ") + ". Run 'gofmt -w .' to fix.",
			Language: checker.LangGo,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "All Go files are properly formatted",
		Language: checker.LangGo,
	}, nil
}
