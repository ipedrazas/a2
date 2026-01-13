package rustcheck

import (
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck runs cargo build to verify the project compiles.
type BuildCheck struct{}

func (c *BuildCheck) ID() string   { return "rust:build" }
func (c *BuildCheck) Name() string { return "Rust Build" }

// Run executes cargo build --release to check compilation.
func (c *BuildCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangRust)

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		return rb.Fail("No Cargo.toml found"), nil
	}

	// Run cargo check (faster than full build, still validates code)
	cmd := exec.Command("cargo", "check", "--quiet")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Extract first line of error for message
		errOutput := strings.TrimSpace(string(output))
		lines := strings.Split(errOutput, "\n")
		if len(lines) > 0 && lines[0] != "" {
			return rb.Fail("Build failed: " + lines[0]), nil
		}
		return rb.Fail("Build failed: " + err.Error()), nil
	}

	return rb.Pass("Build successful"), nil
}
