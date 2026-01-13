package nodecheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsCheck scans for security vulnerabilities in Node.js dependencies.
type DepsCheck struct {
	Config *config.NodeLanguageConfig
}

// ID returns the unique identifier for this check.
func (c *DepsCheck) ID() string {
	return "node:deps"
}

// Name returns the human-readable name for this check.
func (c *DepsCheck) Name() string {
	return "Node Vulnerabilities"
}

// Run executes the dependency vulnerability check.
func (c *DepsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	pm := c.detectPackageManager(path)

	// Bun doesn't have built-in audit
	if pm == "bun" {
		return rb.Pass("Bun does not have built-in security audit"), nil
	}

	// Run security audit
	vulnCount, err := c.runAudit(path, pm)
	if err != nil {
		return rb.Warn(fmt.Sprintf("Security audit failed: %v", err)), nil
	}

	if vulnCount > 0 {
		return rb.Warn(fmt.Sprintf("%d %s found. Run: %s audit for details", vulnCount, checkutil.Pluralize(vulnCount, "vulnerability", "vulnerabilities"), pm)), nil
	}

	return rb.Pass(fmt.Sprintf("No known vulnerabilities found (%s audit)", pm)), nil
}

// detectPackageManager determines which package manager to use.
func (c *DepsCheck) detectPackageManager(path string) string {
	// Check config override first
	if c.Config != nil && c.Config.PackageManager != "auto" && c.Config.PackageManager != "" {
		return c.Config.PackageManager
	}

	// Auto-detect from lock files
	if safepath.Exists(path, "pnpm-lock.yaml") {
		return "pnpm"
	}
	if safepath.Exists(path, "yarn.lock") {
		return "yarn"
	}
	if safepath.Exists(path, "bun.lockb") {
		return "bun"
	}

	return "npm"
}

// runAudit runs the security audit and returns the vulnerability count.
func (c *DepsCheck) runAudit(path, pm string) (int, error) {
	var cmd *exec.Cmd

	switch pm {
	case "npm":
		cmd = exec.Command("npm", "audit", "--json")
	case "yarn":
		// Yarn 1.x uses different format than yarn 2+
		cmd = exec.Command("yarn", "audit", "--json")
	case "pnpm":
		cmd = exec.Command("pnpm", "audit", "--json")
	default:
		cmd = exec.Command("npm", "audit", "--json")
	}

	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// npm audit returns non-zero exit code when vulnerabilities are found
	// so we ignore the error and parse the output
	_ = cmd.Run()

	output := stdout.String()
	if output == "" {
		output = stderr.String()
	}

	return parseAuditOutput(output, pm)
}

// AuditResult represents the npm audit JSON output.
type AuditResult struct {
	Metadata struct {
		Vulnerabilities struct {
			Total    int `json:"total"`
			Low      int `json:"low"`
			Moderate int `json:"moderate"`
			High     int `json:"high"`
			Critical int `json:"critical"`
		} `json:"vulnerabilities"`
	} `json:"metadata"`
	// npm 7+ format
	Vulnerabilities map[string]interface{} `json:"vulnerabilities"`
}

// parseAuditOutput parses the audit output and returns the vulnerability count.
func parseAuditOutput(output, pm string) (int, error) {
	if output == "" {
		return 0, nil
	}

	// Try to parse npm/pnpm JSON format
	var auditResult AuditResult
	if err := json.Unmarshal([]byte(output), &auditResult); err == nil {
		// npm 6 format with metadata
		if auditResult.Metadata.Vulnerabilities.Total > 0 {
			return auditResult.Metadata.Vulnerabilities.Total, nil
		}
		// npm 7+ format with vulnerabilities map
		if len(auditResult.Vulnerabilities) > 0 {
			return len(auditResult.Vulnerabilities), nil
		}
		return 0, nil
	}

	// For yarn, try line-by-line JSON parsing (NDJSON format)
	if pm == "yarn" {
		return parseYarnAudit(output)
	}

	return 0, nil
}

// YarnAuditLine represents a line in yarn audit NDJSON output.
type YarnAuditLine struct {
	Type string `json:"type"`
	Data struct {
		Advisory struct {
			ID int `json:"id"`
		} `json:"advisory"`
		Resolution struct {
			ID int `json:"id"`
		} `json:"resolution"`
	} `json:"data"`
}

// parseYarnAudit parses yarn audit NDJSON output.
func parseYarnAudit(output string) (int, error) {
	advisories := make(map[int]bool)
	lines := bytes.Split([]byte(output), []byte("\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var auditLine YarnAuditLine
		if err := json.Unmarshal(line, &auditLine); err != nil {
			continue
		}

		if auditLine.Type == "auditAdvisory" {
			if auditLine.Data.Advisory.ID > 0 {
				advisories[auditLine.Data.Advisory.ID] = true
			}
		}
	}

	return len(advisories), nil
}
