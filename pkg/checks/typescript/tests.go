package typescriptcheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
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
func (c *TestsCheck) Run(path string) (checker.Result, error) {
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

	// Detect test runner
	runner := c.detectTestRunner(path)
	pm := c.detectPackageManager(path)

	switch runner {
	case "jest":
		return c.runJest(path, pm)
	case "vitest":
		return c.runVitest(path, pm)
	case "mocha":
		return c.runMocha(path, pm)
	case "npm-test":
		return c.runNpmTest(path, pm)
	default:
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "No test runner detected"
		return result, nil
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
func (c *TestsCheck) runJest(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "jest", "--passWithNoTests")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "jest", "--passWithNoTests")
	case "bun":
		cmd = exec.Command("bun", "run", "jest", "--passWithNoTests")
	default:
		cmd = exec.Command("npx", "jest", "--passWithNoTests")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "Jest tests failed"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Jest tests passed"
	return result, nil
}

// runVitest runs Vitest tests.
func (c *TestsCheck) runVitest(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "vitest", "run")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "vitest", "run")
	case "bun":
		cmd = exec.Command("bun", "run", "vitest", "run")
	default:
		cmd = exec.Command("npx", "vitest", "run")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "Vitest tests failed"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Vitest tests passed"
	return result, nil
}

// runMocha runs Mocha tests.
func (c *TestsCheck) runMocha(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "mocha")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "mocha")
	case "bun":
		cmd = exec.Command("bun", "run", "mocha")
	default:
		cmd = exec.Command("npx", "mocha")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "Mocha tests failed"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Mocha tests passed"
	return result, nil
}

// runNpmTest runs npm test script.
func (c *TestsCheck) runNpmTest(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "test")
	case "pnpm":
		cmd = exec.Command("pnpm", "test")
	case "bun":
		cmd = exec.Command("bun", "test")
	default:
		cmd = exec.Command("npm", "test")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())
		// Check if it's just "no tests found" type message
		if strings.Contains(output, "no test") || strings.Contains(stdout.String(), "no test") {
			result.Passed = true
			result.Status = checker.Pass
			result.Message = "No tests found"
			return result, nil
		}
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "Tests failed"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Tests passed"
	return result, nil
}
