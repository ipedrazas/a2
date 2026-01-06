package gocheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// DepsCheck scans for known vulnerabilities using govulncheck.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "go:deps" }
func (c *DepsCheck) Name() string { return "Go Vulnerabilities" }

func (c *DepsCheck) Run(path string) (checker.Result, error) {
	// Check if govulncheck is available
	if _, err := exec.LookPath("govulncheck"); err != nil {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Pass,
			Message:  "govulncheck not installed (run: go install golang.org/x/vuln/cmd/govulncheck@latest)",
			Language: checker.LangGo,
		}, nil
	}

	cmd := exec.Command("govulncheck", "./...")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := strings.TrimSpace(stdout.String())

	if err != nil {
		// govulncheck exits with non-zero when vulnerabilities are found
		// Count vulnerabilities from output
		vulnCount := countVulnerabilities(output)

		if vulnCount > 0 {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   false,
				Status:   checker.Warn,
				Message:  formatVulnMessage(vulnCount),
				Language: checker.LangGo,
			}, nil
		}

		// Some other error
		errOutput := strings.TrimSpace(stderr.String())
		if errOutput != "" {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   false,
				Status:   checker.Warn,
				Message:  "govulncheck error: " + errOutput,
				Language: checker.LangGo,
			}, nil
		}
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "No known vulnerabilities found",
		Language: checker.LangGo,
	}, nil
}

// countVulnerabilities counts the number of vulnerabilities in govulncheck output.
func countVulnerabilities(output string) int {
	// govulncheck outputs "Vulnerability #N" for each finding
	count := strings.Count(output, "Vulnerability #")
	if count == 0 {
		// Alternative format check
		count = strings.Count(output, "GO-")
	}
	return count
}

// formatVulnMessage creates a summary message.
func formatVulnMessage(count int) string {
	if count == 1 {
		return "1 vulnerability found. Run 'govulncheck ./...' for details."
	}
	return fmt.Sprintf("%d vulnerabilities found. Run 'govulncheck ./...' for details.", count)
}
