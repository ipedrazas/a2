package swiftcheck

import (
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck runs swift build to verify the project compiles.
type BuildCheck struct{}

func (c *BuildCheck) ID() string   { return "swift:build" }
func (c *BuildCheck) Name() string { return "Swift Build" }

// Run executes swift build to check compilation.
func (c *BuildCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangSwift)

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		return rb.Fail("No Package.swift found"), nil
	}

	// Run swift build
	cmd := exec.Command("swift", "build")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// Extract first line of error for message
		errOutput := strings.TrimSpace(outputStr)
		lines := strings.Split(errOutput, "\n")
		if len(lines) > 0 && lines[0] != "" {
			return rb.FailWithOutput("Build failed: "+lines[0], outputStr), nil
		}
		return rb.FailWithOutput("Build failed: "+err.Error(), outputStr), nil
	}

	return rb.PassWithOutput("Build successful", outputStr), nil
}
