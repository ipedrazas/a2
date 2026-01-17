package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DockerfileCheck verifies that a Dockerfile or Containerfile exists and scans for issues.
type DockerfileCheck struct{}

func (c *DockerfileCheck) ID() string   { return "common:dockerfile" }
func (c *DockerfileCheck) Name() string { return "Container Ready" }

func (c *DockerfileCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Find Dockerfile
	dockerfile := c.findDockerfile(path)
	if dockerfile == "" {
		return rb.Warn("No Dockerfile or Containerfile found"), nil
	}

	// Check if trivy is installed
	trivyInstalled := checkutil.ToolAvailable("trivy")

	// Priority 1: If trivy is installed, scan the Dockerfile
	if trivyInstalled {
		return c.runTrivy(path, dockerfile, rb)
	}

	// Priority 2: Fall back to basic file existence check
	return c.basicCheck(path, dockerfile, rb), nil
}

// findDockerfile returns the path to a Dockerfile if one exists.
func (c *DockerfileCheck) findDockerfile(path string) string {
	dockerfiles := []string{
		"Dockerfile",
		"dockerfile",
		"Containerfile",
		"containerfile",
	}

	for _, df := range dockerfiles {
		if safepath.Exists(path, df) {
			return df
		}
	}
	return ""
}

// runTrivy executes trivy to scan the Dockerfile for misconfigurations.
func (c *DockerfileCheck) runTrivy(path, dockerfile string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	// Run trivy config scan on the Dockerfile
	args := []string{"config", "--exit-code", "0", "--format", "json", dockerfile}

	result := checkutil.RunCommand(path, "trivy", args...)
	output := result.CombinedOutput()

	// Check for .dockerignore
	hasIgnore := safepath.Exists(path, ".dockerignore")

	// Count issues from trivy output
	issueCount := c.countTrivyIssues(output)

	if issueCount > 0 {
		if hasIgnore {
			return rb.WarnWithOutput(fmt.Sprintf("trivy: %d %s in %s (with .dockerignore)",
				issueCount, pluralize(issueCount, "issue", "issues"), dockerfile), output), nil
		}
		return rb.WarnWithOutput(fmt.Sprintf("trivy: %d %s in %s, consider adding .dockerignore",
			issueCount, pluralize(issueCount, "issue", "issues"), dockerfile), output), nil
	}

	// No issues found
	if hasIgnore {
		return rb.Pass(fmt.Sprintf("trivy: %s scanned, no issues (with .dockerignore)", dockerfile)), nil
	}
	return rb.Pass(fmt.Sprintf("trivy: %s scanned, no issues, consider adding .dockerignore", dockerfile)), nil
}

// countTrivyIssues counts misconfigurations from trivy JSON output.
func (c *DockerfileCheck) countTrivyIssues(output string) int {
	// Look for "Misconfigurations" in JSON output
	// Pattern: "Misconfigurations": [...] with entries
	misconfigPattern := regexp.MustCompile(`"Misconfigurations"\s*:\s*\[`)
	if !misconfigPattern.MatchString(output) {
		return 0
	}

	// Count individual misconfigurations by looking for "ID" entries within Misconfigurations
	// Each misconfiguration has an "ID" field like "DS001", "DS002", etc.
	idCount := len(regexp.MustCompile(`"ID"\s*:\s*"DS\d+"`).FindAllString(output, -1))
	if idCount > 0 {
		return idCount
	}

	// Alternative: count "Severity" entries which indicate findings
	severityCount := strings.Count(output, `"Severity"`)
	return severityCount
}

// basicCheck performs basic Dockerfile existence check without trivy.
func (c *DockerfileCheck) basicCheck(path, dockerfile string, rb *checkutil.ResultBuilder) checker.Result {
	hasIgnore := safepath.Exists(path, ".dockerignore")

	if hasIgnore {
		return rb.Pass(dockerfile + " found with .dockerignore")
	}
	return rb.Pass(dockerfile + " found (consider adding .dockerignore)")
}
