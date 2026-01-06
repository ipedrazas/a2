package pythoncheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// DepsCheck scans for known vulnerabilities in Python dependencies.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "python:deps" }
func (c *DepsCheck) Name() string { return "Python Vulnerabilities" }

func (c *DepsCheck) Run(path string) (checker.Result, error) {
	// Try pip-audit first, then safety
	var cmd *exec.Cmd
	var cmdDesc string

	if _, err := exec.LookPath("pip-audit"); err == nil {
		cmd = exec.Command("pip-audit")
		cmdDesc = "pip-audit"
	} else if _, err := exec.LookPath("safety"); err == nil {
		cmd = exec.Command("safety", "check")
		cmdDesc = "safety"
	} else {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Pass,
			Message:  "No vulnerability scanner installed (install pip-audit or safety)",
			Language: checker.LangPython,
		}, nil
	}

	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := strings.TrimSpace(stdout.String())

	if err != nil {
		// pip-audit and safety exit with non-zero when vulnerabilities are found
		vulnCount := countPythonVulnerabilities(output, cmdDesc)

		if vulnCount > 0 {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   false,
				Status:   checker.Warn,
				Message:  fmt.Sprintf("%d vulnerabilities found. Run '%s' for details.", vulnCount, cmdDesc),
				Language: checker.LangPython,
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
				Message:  cmdDesc + " error: " + truncateMessage(errOutput, 150),
				Language: checker.LangPython,
			}, nil
		}
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "No known vulnerabilities found",
		Language: checker.LangPython,
	}, nil
}

// countPythonVulnerabilities counts vulnerabilities in pip-audit or safety output.
func countPythonVulnerabilities(output, tool string) int {
	switch tool {
	case "pip-audit":
		// pip-audit outputs one line per vulnerability
		lines := strings.Split(output, "\n")
		count := 0
		for _, line := range lines {
			if strings.Contains(line, "vulnerability") || strings.Contains(line, "PYSEC-") || strings.Contains(line, "CVE-") {
				count++
			}
		}
		if count == 0 {
			// Alternative: count non-empty lines that look like vulnerability entries
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" && !strings.HasPrefix(trimmed, "Name") && !strings.HasPrefix(trimmed, "---") {
					if strings.Contains(trimmed, "fix") || strings.Contains(trimmed, "Fixed") {
						count++
					}
				}
			}
		}
		return count
	case "safety":
		// safety outputs vulnerability count
		count := strings.Count(output, "vulnerability found")
		if count == 0 {
			count = strings.Count(output, "->")
		}
		return count
	}
	return 0
}
