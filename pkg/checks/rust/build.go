package rustcheck

import (
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck runs cargo build to verify the project compiles.
type BuildCheck struct{}

func (c *BuildCheck) ID() string   { return "rust:build" }
func (c *BuildCheck) Name() string { return "Rust Build" }

// Run executes cargo build --release to check compilation.
func (c *BuildCheck) Run(path string) (checker.Result, error) {
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

	// Run cargo check (faster than full build, still validates code)
	cmd := exec.Command("cargo", "check", "--quiet")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Passed = false
		result.Status = checker.Fail
		// Extract first line of error for message
		errOutput := strings.TrimSpace(string(output))
		lines := strings.Split(errOutput, "\n")
		if len(lines) > 0 && lines[0] != "" {
			result.Message = "Build failed: " + lines[0]
		} else {
			result.Message = "Build failed: " + err.Error()
		}
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Build successful"

	return result, nil
}
