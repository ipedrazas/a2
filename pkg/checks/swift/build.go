package swiftcheck

import (
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck runs swift build to verify the project compiles.
type BuildCheck struct{}

func (c *BuildCheck) ID() string   { return "swift:build" }
func (c *BuildCheck) Name() string { return "Swift Build" }

// Run executes swift build to check compilation.
func (c *BuildCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangSwift,
	}

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No Package.swift found"
		return result, nil
	}

	// Run swift build
	cmd := exec.Command("swift", "build")
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
