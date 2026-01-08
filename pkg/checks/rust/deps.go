package rustcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsCheck runs cargo audit to check for vulnerable dependencies.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "rust:deps" }
func (c *DepsCheck) Name() string { return "Rust Vulnerabilities" }

// Run checks for dependency vulnerabilities using cargo audit.
func (c *DepsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangRust,
	}

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No Cargo.toml found"
		return result, nil
	}

	// Check for Cargo.lock (required for audit)
	if !safepath.Exists(path, "Cargo.lock") {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "No Cargo.lock found (run cargo build to generate)"
		return result, nil
	}

	// Check for deny.toml (cargo-deny config)
	hasDeny := safepath.Exists(path, "deny.toml")

	// Try cargo audit first
	cmd := exec.Command("cargo", "audit", "--quiet")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// If cargo audit is not installed, check for other tools
	if err != nil && strings.Contains(outputStr, "no such subcommand") {
		// Check if cargo-deny is configured
		if hasDeny {
			result.Passed = true
			result.Status = checker.Pass
			result.Message = "cargo-deny configured for security scanning"
			return result, nil
		}

		// Check for audit in CI
		ciConfigured := c.checkCIForAudit(path)
		if ciConfigured {
			result.Passed = true
			result.Status = checker.Pass
			result.Message = "Security audit configured in CI"
			return result, nil
		}

		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No security audit tool found (consider cargo-audit or cargo-deny)"
		return result, nil
	}

	if err != nil {
		// Vulnerabilities found
		vulnRe := regexp.MustCompile(`(\d+) vulnerabilit`)
		matches := vulnRe.FindStringSubmatch(outputStr)
		if len(matches) > 1 {
			result.Passed = false
			result.Status = checker.Warn
			result.Message = matches[1] + " vulnerabilities found"
		} else {
			result.Passed = false
			result.Status = checker.Warn
			result.Message = "Vulnerabilities found in dependencies"
		}
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "No known vulnerabilities found"

	return result, nil
}

// checkCIForAudit looks for audit configuration in CI files.
func (c *DepsCheck) checkCIForAudit(path string) bool {
	ciFiles := []string{
		".github/workflows/ci.yml",
		".github/workflows/ci.yaml",
		".github/workflows/rust.yml",
		".github/workflows/rust.yaml",
		".github/workflows/security.yml",
		".gitlab-ci.yml",
	}

	for _, ciFile := range ciFiles {
		if content, err := safepath.ReadFile(path, ciFile); err == nil {
			contentLower := strings.ToLower(string(content))
			if strings.Contains(contentLower, "cargo audit") ||
				strings.Contains(contentLower, "cargo-audit") ||
				strings.Contains(contentLower, "cargo deny") ||
				strings.Contains(contentLower, "cargo-deny") {
				return true
			}
		}
	}

	return false
}
