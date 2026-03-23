package rustcheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsFreshnessCheck detects outdated Rust crate dependencies.
type DepsFreshnessCheck struct{}

func (c *DepsFreshnessCheck) ID() string   { return "rust:deps_freshness" }
func (c *DepsFreshnessCheck) Name() string { return "Rust Dependency Freshness" }

func (c *DepsFreshnessCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangRust)

	if !safepath.Exists(path, "Cargo.toml") {
		return rb.Fail("No Cargo.toml found"), nil
	}

	if !checkutil.ToolAvailable("cargo-outdated") {
		// Fallback: check if Cargo.lock exists and is reasonably fresh
		if !safepath.Exists(path, "Cargo.lock") {
			return rb.Pass("No Cargo.lock found (run cargo build to generate)"), nil
		}
		return rb.ToolNotInstalled("cargo-outdated", "cargo install cargo-outdated"), nil
	}

	result := checkutil.RunCommand(path, "cargo", "outdated", "--root-deps-only")
	output := result.CombinedOutput()

	if !result.Success() {
		return rb.WarnWithOutput("cargo outdated failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	outdated := countOutdatedCrates(result.Stdout)

	if outdated == 0 {
		return rb.PassWithOutput("All dependencies are up to date", output), nil
	}

	msg := fmt.Sprintf("%d outdated %s",
		outdated,
		checkutil.Pluralize(outdated, "crate", "crates"),
	)

	if outdated > 20 {
		return rb.WarnWithOutput(msg+". Run: cargo outdated", output), nil
	}

	return rb.PassWithOutput(msg+". Run: cargo update", output), nil
}

// countOutdatedCrates counts outdated crates from cargo outdated output.
// Output format:
// Name             Project  Compat   Latest   Kind
// ----             -------  ------   ------   ----
// serde            1.0.100  1.0.200  1.0.200  Normal
func countOutdatedCrates(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "Name") || strings.HasPrefix(trimmed, "---") || strings.HasPrefix(trimmed, "====") {
			continue
		}
		if len(strings.Fields(trimmed)) >= 4 {
			count++
		}
	}
	return count
}
