package typescriptcheck

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
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
func (c *BuildCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangTypeScript)

	// Check for tsconfig.json
	if !safepath.Exists(path, "tsconfig.json") && !safepath.Exists(path, "tsconfig.base.json") {
		return rb.Fail("No tsconfig.json found"), nil
	}

	// Try to find package manager and run build script
	pm := c.detectPackageManager(path)
	pkg, _ := ParsePackageJSON(path)

	// Check if there's a build script in package.json
	if pkg != nil {
		if _, hasBuild := pkg.Scripts["build"]; hasBuild {
			return c.runBuildScript(ctx, path, pm, rb)
		}
	}

	// Fall back to tsc --noEmit for type checking only
	return c.runTscNoEmit(ctx, path, rb)
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
func (c *BuildCheck) runBuildScript(ctx context.Context, path, pm string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = checkutil.Command(ctx, "yarn", "build")
	case "pnpm":
		cmd = checkutil.Command(ctx, "pnpm", "run", "build")
	case "bun":
		cmd = checkutil.Command(ctx, "bun", "run", "build")
	default:
		cmd = checkutil.Command(ctx, "npm", "run", "build")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return rb.FailWithOutput("Build failed: "+checkutil.TruncateMessage(errMsg, 200), output), nil
		}
		return rb.FailWithOutput("Build failed", output), nil
	}

	return rb.PassWithOutput("Build successful", output), nil
}

// runTscNoEmit runs tsc --noEmit for type checking.
func (c *BuildCheck) runTscNoEmit(ctx context.Context, path string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	// Check if npx is available
	if _, err := exec.LookPath("npx"); err != nil {
		return rb.ToolNotInstalled("npx", "install Node.js"), nil
	}

	cmd := checkutil.Command(ctx, "npx", "tsc", "--noEmit")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()
	if err != nil {
		return rb.FailWithOutput("TypeScript compilation failed", output), nil
	}

	return rb.PassWithOutput("TypeScript compiles successfully", output), nil
}
