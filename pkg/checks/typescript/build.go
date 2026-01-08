package typescriptcheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck verifies TypeScript compilation.
type BuildCheck struct {
	Config *config.TypeScriptLanguageConfig
}

func (c *BuildCheck) ID() string   { return "typescript:build" }
func (c *BuildCheck) Name() string { return "TypeScript Build" }

// Run checks if TypeScript compiles successfully.
func (c *BuildCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	// Check for tsconfig.json
	if !safepath.Exists(path, "tsconfig.json") && !safepath.Exists(path, "tsconfig.base.json") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No tsconfig.json found"
		return result, nil
	}

	// Try to find package manager and run build script
	pm := c.detectPackageManager(path)
	pkg, _ := ParsePackageJSON(path)

	// Check if there's a build script in package.json
	if pkg != nil {
		if _, hasBuild := pkg.Scripts["build"]; hasBuild {
			return c.runBuildScript(path, pm)
		}
	}

	// Fall back to tsc --noEmit for type checking only
	return c.runTscNoEmit(path)
}

// detectPackageManager determines which package manager to use.
func (c *BuildCheck) detectPackageManager(path string) string {
	if c.Config != nil && c.Config.PackageManager != "" && c.Config.PackageManager != "auto" {
		return c.Config.PackageManager
	}

	// Detect from lock files
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

// runBuildScript runs the package.json build script.
func (c *BuildCheck) runBuildScript(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "build")
	case "pnpm":
		cmd = exec.Command("pnpm", "run", "build")
	case "bun":
		cmd = exec.Command("bun", "run", "build")
	default:
		cmd = exec.Command("npm", "run", "build")
	}
	cmd.Dir = path

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Fail
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			// Truncate long error messages
			if len(errMsg) > 200 {
				errMsg = errMsg[:200] + "..."
			}
			result.Message = "Build failed: " + errMsg
		} else {
			result.Message = "Build failed"
		}
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Build successful"
	return result, nil
}

// runTscNoEmit runs tsc --noEmit for type checking.
func (c *BuildCheck) runTscNoEmit(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	// Check if npx is available
	if _, err := exec.LookPath("npx"); err != nil {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "npx not available, skipping build check"
		return result, nil
	}

	cmd := exec.Command("npx", "tsc", "--noEmit")
	cmd.Dir = path

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "TypeScript compilation failed"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "TypeScript compiles successfully"
	return result, nil
}
