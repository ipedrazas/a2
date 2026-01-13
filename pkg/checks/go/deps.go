package gocheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// DepsCheck scans for known vulnerabilities using govulncheck.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "go:deps" }
func (c *DepsCheck) Name() string { return "Go Vulnerabilities" }

func (c *DepsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	// Check if govulncheck is available
	if !checkutil.ToolAvailable("govulncheck") {
		return rb.ToolNotInstalled("govulncheck", "go install golang.org/x/vuln/cmd/govulncheck@latest"), nil
	}

	result := checkutil.RunCommand(path, "govulncheck", "./...")

	if !result.Success() {
		// govulncheck exits with non-zero when vulnerabilities are found
		output := strings.TrimSpace(result.Stdout)
		vulnCount := countVulnerabilities(output)

		if vulnCount > 0 {
			return rb.Warn(formatVulnMessage(vulnCount)), nil
		}

		// Some other error
		if result.Stderr != "" {
			return rb.Warn("govulncheck error: " + checkutil.TruncateMessage(strings.TrimSpace(result.Stderr), 150)), nil
		}
	}

	return rb.Pass("No known vulnerabilities found"), nil
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
