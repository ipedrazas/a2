package pythoncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// DepsCheck scans for known vulnerabilities in Python dependencies.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "python:deps" }
func (c *DepsCheck) Name() string { return "Python Vulnerabilities" }

func (c *DepsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)

	// Try pip-audit first, then safety
	var result *checkutil.CommandResult
	var cmdDesc string

	if checkutil.ToolAvailable("pip-audit") {
		result = checkutil.RunCommand(path, "pip-audit")
		cmdDesc = "pip-audit"
	} else if checkutil.ToolAvailable("safety") {
		result = checkutil.RunCommand(path, "safety", "check")
		cmdDesc = "safety"
	} else {
		return rb.Pass("No vulnerability scanner installed (install pip-audit or safety)"), nil
	}

	if !result.Success() {
		// pip-audit and safety exit with non-zero when vulnerabilities are found
		output := strings.TrimSpace(result.Stdout)
		vulnCount := countPythonVulnerabilities(output, cmdDesc)

		if vulnCount > 0 {
			return rb.Warn(fmt.Sprintf("%d vulnerabilities found. Run '%s' for details.", vulnCount, cmdDesc)), nil
		}

		// Some other error
		if result.Stderr != "" {
			return rb.Warn(cmdDesc + " error: " + checkutil.TruncateMessage(strings.TrimSpace(result.Stderr), 150)), nil
		}
	}

	return rb.Pass("No known vulnerabilities found"), nil
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
