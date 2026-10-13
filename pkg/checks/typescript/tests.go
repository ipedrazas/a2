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

// TestsCheck runs TypeScript project tests.
type TestsCheck struct {
	Config *config.TypeScriptLanguageConfig
}

func (c *TestsCheck) ID() string   { return "typescript:tests" }
func (c *TestsCheck) Name() string { return "TypeScript Tests" }

// Run executes the test suite.
func (c *TestsCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangTypeScript)

	// Check for tsconfig.json
	if !safepath.Exists(path, "tsconfig.json") && !safepath.Exists(path, "tsconfig.base.json") {
		return rb.Fail("No tsconfig.json found"), nil
	}

	// Detect test runner
	runner := c.detectTestRunner(path)
	pm := c.detectPackageManager(path)

	switch runner {
	case "jest":
		return c.runJest(ctx, path, pm, rb)
	case "vitest":
		return c.runVitest(ctx, path, pm, rb)
	case "mocha":
		return c.runMocha(ctx, path, pm, rb)
	case "npm-test":
		return c.runNpmTest(ctx, path, pm, rb)
	default:
		return rb.Pass("No test runner detected"), nil
	}
}

// detectTestRunner identifies which test framework is used.
func (c *TestsCheck) detectTestRunner(path string) string {
	if c.Config != nil && c.Config.TestRunner != "" && c.Config.TestRunner != "auto" {
		return c.Config.TestRunner
	}

	// Check for config files
	if safepath.Exists(path, "jest.config.js") || safepath.Exists(path, "jest.config.ts") ||
		safepath.Exists(path, "jest.config.mjs") || safepath.Exists(path, "jest.config.cjs") {
		return "jest"
	}

	if safepath.Exists(path, "vitest.config.ts") || safepath.Exists(path, "vitest.config.js") ||
		safepath.Exists(path, "vitest.config.mts") {
		return "vitest"
	}

	if safepath.Exists(path, ".mocharc.js") || safepath.Exists(path, ".mocharc.json") ||
		safepath.Exists(path, ".mocharc.yml") || safepath.Exists(path, ".mocharc.yaml") {
		return "mocha"
	}

	// Check package.json dependencies
	pkg, err := ParsePackageJSON(path)
	if err != nil {
		return ""
	}

	if _, ok := pkg.DevDependencies["vitest"]; ok {
		return "vitest"
	}
	if _, ok := pkg.DevDependencies["jest"]; ok {
		return "jest"
	}
	if _, ok := pkg.DevDependencies["@jest/core"]; ok {
		return "jest"
	}
	if _, ok := pkg.DevDependencies["mocha"]; ok {
		return "mocha"
	}

	// Check if there's a test script
	if _, hasTest := pkg.Scripts["test"]; hasTest {
		return "npm-test"
	}

	return ""
}

// detectPackageManager determines which package manager to use.
func (c *TestsCheck) detectPackageManager(path string) string {
	if c.Config != nil && c.Config.PackageManager != "" && c.Config.PackageManager != "auto" {
		return c.Config.PackageManager
	}

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

// runJest runs Jest tests.
func (c *TestsCheck) runJest(ctx context.Context, path, pm string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = checkutil.Command(ctx, "yarn", "jest", "--passWithNoTests")
	case "pnpm":
		cmd = checkutil.Command(ctx, "pnpm", "exec", "jest", "--passWithNoTests")
	case "bun":
		cmd = checkutil.Command(ctx, "bun", "run", "jest", "--passWithNoTests")
	default:
		cmd = checkutil.Command(ctx, "npx", "jest", "--passWithNoTests")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()
	if err != nil {
		return rb.FailWithOutput("Jest tests failed", output), nil
	}

	return rb.PassWithOutput("Jest tests passed", output), nil
}

// runVitest runs Vitest tests.
func (c *TestsCheck) runVitest(ctx context.Context, path, pm string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = checkutil.Command(ctx, "yarn", "vitest", "run")
	case "pnpm":
		cmd = checkutil.Command(ctx, "pnpm", "exec", "vitest", "run")
	case "bun":
		cmd = checkutil.Command(ctx, "bun", "run", "vitest", "run")
	default:
		cmd = checkutil.Command(ctx, "npx", "vitest", "run")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()
	if err != nil {
		return rb.FailWithOutput("Vitest tests failed", output), nil
	}

	return rb.PassWithOutput("Vitest tests passed", output), nil
}

// runMocha runs Mocha tests.
func (c *TestsCheck) runMocha(ctx context.Context, path, pm string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = checkutil.Command(ctx, "yarn", "mocha")
	case "pnpm":
		cmd = checkutil.Command(ctx, "pnpm", "exec", "mocha")
	case "bun":
		cmd = checkutil.Command(ctx, "bun", "run", "mocha")
	default:
		cmd = checkutil.Command(ctx, "npx", "mocha")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()
	if err != nil {
		return rb.FailWithOutput("Mocha tests failed", output), nil
	}

	return rb.PassWithOutput("Mocha tests passed", output), nil
}

// runNpmTest runs npm test script.
func (c *TestsCheck) runNpmTest(ctx context.Context, path, pm string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = checkutil.Command(ctx, "yarn", "test")
	case "pnpm":
		cmd = checkutil.Command(ctx, "pnpm", "test")
	case "bun":
		cmd = checkutil.Command(ctx, "bun", "test")
	default:
		cmd = checkutil.Command(ctx, "npm", "test")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	combinedOutput := stdout.String() + stderr.String()
	if err != nil {
		// Check if it's just "no tests found" type message
		if strings.Contains(combinedOutput, "no test") {
			return rb.Pass("No tests found"), nil
		}
		return rb.FailWithOutput("Tests failed", combinedOutput), nil
	}

	return rb.PassWithOutput("Tests passed", combinedOutput), nil
}
