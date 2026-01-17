package nodecheck

import (
	"fmt"

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
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	// Check if package.json exists first
	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	pm := c.detectPackageManager(path)

	// Verify package manager is available
	if !checkutil.ToolAvailable(pm) {
		return rb.ToolNotInstalled(pm, ""), nil
	}

	// Run validation command based on package manager
	var result *checkutil.CommandResult
	switch pm {
	case "pnpm":
		result = checkutil.RunCommand(path, "pnpm", "install", "--frozen-lockfile", "--dry-run")
	case "yarn":
		result = checkutil.RunCommand(path, "yarn", "install", "--check-files")
	case "bun":
		result = checkutil.RunCommand(path, "bun", "install", "--dry-run")
	default: // npm
		if safepath.Exists(path, "package-lock.json") {
			result = checkutil.RunCommand(path, "npm", "ci", "--dry-run")
		} else {
			result = checkutil.RunCommand(path, "npm", "install", "--dry-run")
		}
	}

	output := result.CombinedOutput()
	if !result.Success() {
		return rb.FailWithOutput(fmt.Sprintf("Dependency validation failed (%s): %s", pm, checkutil.TruncateMessage(result.Output(), 200)), output), nil
	}

	return rb.PassWithOutput(fmt.Sprintf("Dependencies valid (%s)", pm), output), nil
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
