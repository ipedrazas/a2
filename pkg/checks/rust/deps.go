package rustcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsCheck runs cargo audit to check for vulnerable dependencies.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "rust:deps" }
func (c *DepsCheck) Name() string { return "Rust Vulnerabilities" }

// Run checks for dependency vulnerabilities using cargo audit.
func (c *DepsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangRust)

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		return rb.Fail("No Cargo.toml found"), nil
	}

	// Check for Cargo.lock (required for audit)
	if !safepath.Exists(path, "Cargo.lock") {
		return rb.Pass("No Cargo.lock found (run cargo build to generate)"), nil
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
			return rb.Pass("cargo-deny configured for security scanning"), nil
		}

		// Check for audit in CI
		ciConfigured := c.checkCIForAudit(path)
		if ciConfigured {
			return rb.Pass("Security audit configured in CI"), nil
		}

		return rb.Warn("No security audit tool found (consider cargo-audit or cargo-deny)"), nil
	}

	if err != nil {
		// Vulnerabilities found
		vulnRe := regexp.MustCompile(`(\d+) vulnerabilit`)
		matches := vulnRe.FindStringSubmatch(outputStr)
		if len(matches) > 1 {
			return rb.Warn(matches[1] + " vulnerabilities found"), nil
		}
		return rb.Warn("Vulnerabilities found in dependencies"), nil
	}

	return rb.Pass("No known vulnerabilities found"), nil
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
