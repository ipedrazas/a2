package typescriptcheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsCheck checks for dependency vulnerabilities.
type DepsCheck struct {
	Config *config.TypeScriptLanguageConfig
}

func (c *DepsCheck) ID() string   { return "typescript:deps" }
func (c *DepsCheck) Name() string { return "TypeScript Vulnerabilities" }

// Run checks for dependency vulnerabilities.
func (c *DepsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	// Check for package.json
	if !safepath.Exists(path, "package.json") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No package.json found"
		return result, nil
	}

	// Detect package manager
	pm := c.detectPackageManager(path)

	// Check for vulnerability scanning tools
	var tools []string

	// Check for Snyk
	if c.hasSnyk(path) {
		tools = append(tools, "Snyk")
	}

	// Check for Dependabot
	if safepath.Exists(path, ".github/dependabot.yml") || safepath.Exists(path, ".github/dependabot.yaml") {
		tools = append(tools, "Dependabot")
	}

	// Check for Renovate
	if safepath.Exists(path, "renovate.json") || safepath.Exists(path, ".renovaterc") ||
		safepath.Exists(path, ".renovaterc.json") {
		tools = append(tools, "Renovate")
	}

	// Run audit command
	auditResult := c.runAudit(path, pm)
	if auditResult != "" {
		if len(tools) > 0 {
			result.Passed = true
			result.Status = checker.Pass
			result.Message = auditResult + "; configured: " + strings.Join(tools, ", ")
		} else {
			result.Passed = true
			result.Status = checker.Pass
			result.Message = auditResult
		}
		return result, nil
	}

	// If audit failed or not available, check for tools
	if len(tools) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Vulnerability scanning configured: " + strings.Join(tools, ", ")
		return result, nil
	}

	result.Passed = false
	result.Status = checker.Warn
	result.Message = "No vulnerability scanning configured (consider npm audit or Snyk)"
	return result, nil
}

// detectPackageManager determines which package manager to use.
func (c *DepsCheck) detectPackageManager(path string) string {
	if c.Config != nil && c.Config.PackageManager != "" && c.Config.PackageManager != "auto" {
		return c.Config.PackageManager
	}

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

// hasSnyk checks if Snyk is configured.
func (c *DepsCheck) hasSnyk(path string) bool {
	if safepath.Exists(path, ".snyk") {
		return true
	}

	// Check CI configs for snyk
	ciFiles := []string{
		".github/workflows/ci.yml", ".github/workflows/ci.yaml",
		".github/workflows/security.yml", ".github/workflows/security.yaml",
	}
	for _, ciFile := range ciFiles {
		if content, err := safepath.ReadFile(path, ciFile); err == nil {
			if strings.Contains(strings.ToLower(string(content)), "snyk") {
				return true
			}
		}
	}

	return false
}

// runAudit runs the package manager's audit command.
func (c *DepsCheck) runAudit(path, pm string) string {
	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		// yarn audit exits with non-zero if vulnerabilities found
		cmd = exec.Command("yarn", "audit", "--summary")
	case "pnpm":
		cmd = exec.Command("pnpm", "audit")
	case "bun":
		// bun doesn't have audit yet, skip
		return ""
	default:
		cmd = exec.Command("npm", "audit", "--audit-level=high")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()

	// Parse output for vulnerability info
	if err == nil {
		return "No known vulnerabilities found"
	}

	// Check if there are actual vulnerabilities or just an error
	outputLower := strings.ToLower(output + stderr.String())
	if strings.Contains(outputLower, "vulnerabilit") {
		if strings.Contains(outputLower, "critical") || strings.Contains(outputLower, "high") {
			return ""
		}
		return "Low/moderate vulnerabilities found (run audit for details)"
	}

	// Audit command failed (possibly no lock file)
	return ""
}
