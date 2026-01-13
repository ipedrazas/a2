package nodecheck

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck verifies that Node.js dependencies can be installed.
type BuildCheck struct {
	Config *config.NodeLanguageConfig
}

// ID returns the unique identifier for this check.
func (c *BuildCheck) ID() string {
	return "node:build"
}

// Name returns the human-readable name for this check.
func (c *BuildCheck) Name() string {
	return "Node Build"
}

// Run executes the build check.
func (c *BuildCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangNode,
	}

	// Check if package.json exists first
	if !safepath.Exists(path, "package.json") {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = "package.json not found"
		return result, nil
	}

	pm := c.detectPackageManager(path)

	// Verify package manager is available
	if _, err := exec.LookPath(pm); err != nil {
		result.Status = checker.Pass
		result.Passed = true
		result.Message = fmt.Sprintf("%s not installed. Install it to verify dependencies.", pm)
		return result, nil
	}

	// Run validation command based on package manager
	var cmd *exec.Cmd
	switch pm {
	case "pnpm":
		cmd = exec.Command("pnpm", "install", "--frozen-lockfile", "--dry-run")
	case "yarn":
		cmd = exec.Command("yarn", "install", "--check-files")
	case "bun":
		cmd = exec.Command("bun", "install", "--dry-run")
	default: // npm
		// Check if package-lock.json exists for npm ci
		if safepath.Exists(path, "package-lock.json") {
			cmd = exec.Command("npm", "ci", "--dry-run")
		} else {
			cmd = exec.Command("npm", "install", "--dry-run")
		}
	}

	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		output := stderr.String()
		if output == "" {
			output = stdout.String()
		}
		result.Status = checker.Fail
		result.Passed = false
		result.Message = fmt.Sprintf("Dependency validation failed (%s): %s", pm, checkutil.TruncateMessage(output, 200))
		return result, nil
	}

	result.Status = checker.Pass
	result.Passed = true
	result.Message = fmt.Sprintf("Dependencies valid (%s)", pm)
	return result, nil
}

// detectPackageManager determines which package manager to use.
func (c *BuildCheck) detectPackageManager(path string) string {
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
