package rustcheck

import (
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// FormatCheck verifies Rust code is formatted with rustfmt.
type FormatCheck struct{}

func (c *FormatCheck) ID() string   { return "rust:format" }
func (c *FormatCheck) Name() string { return "Rust Format" }

// Run checks if code is properly formatted using rustfmt.
func (c *FormatCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangRust)

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		return rb.Fail("No Cargo.toml found"), nil
	}

	// Check for rustfmt config
	hasConfig := safepath.Exists(path, "rustfmt.toml") || safepath.Exists(path, ".rustfmt.toml")

	// Run cargo fmt --check to see if code is formatted
	cmd := exec.Command("cargo", "fmt", "--check")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Exit code 1 means unformatted code found
		outputStr := strings.TrimSpace(string(output))
		if strings.Contains(outputStr, "Diff in") || strings.Contains(outputStr, "would be reformatted") {
			// Count files that need formatting
			lines := strings.Split(outputStr, "\n")
			count := 0
			for _, line := range lines {
				if strings.Contains(line, "Diff in") || strings.HasSuffix(line, ".rs") {
					count++
				}
			}
			if count > 0 {
				return rb.Warn("Code not formatted: " + string(rune(count)) + " file(s) need formatting"), nil
			}
			return rb.Warn("Code not formatted (run 'cargo fmt')"), nil
		}
		// Some other error (rustfmt not installed, etc.)
		return rb.Warn("Cannot check format: " + err.Error()), nil
	}

	if hasConfig {
		return rb.Pass("Code formatted (custom config)"), nil
	}
	return rb.Pass("Code formatted"), nil
}
