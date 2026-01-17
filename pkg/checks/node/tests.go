package nodecheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TestsCheck runs the project's test suite.
type TestsCheck struct {
	Config *config.NodeLanguageConfig
}

// ID returns the unique identifier for this check.
func (c *TestsCheck) ID() string {
	return "node:tests"
}

// Name returns the human-readable name for this check.
func (c *TestsCheck) Name() string {
	return "Node Tests"
}

// Run executes the tests check.
func (c *TestsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	// Parse package.json to check for test script
	pkg, err := ParsePackageJSON(path)
	if err != nil {
		return rb.Fail(fmt.Sprintf("Failed to parse package.json: %v", err)), nil
	}

	// Check if test script exists
	testScript, hasTest := pkg.Scripts["test"]
	if !hasTest || testScript == "" {
		return rb.Pass("No test script defined in package.json"), nil
	}

	// Check for default "no test specified" script
	if strings.Contains(testScript, "no test specified") {
		return rb.Pass("No tests configured (default npm init script)"), nil
	}

	// Detect test runner and run tests
	pm := c.detectPackageManager(path)
	runner := c.detectTestRunner(path, pkg)

	var cmd *exec.Cmd
	switch runner {
	case "jest":
		if _, err := exec.LookPath("npx"); err == nil {
			cmd = exec.Command("npx", "jest", "--passWithNoTests", "--ci")
		} else {
			cmd = c.createTestCommand(pm)
		}
	case "vitest":
		if _, err := exec.LookPath("npx"); err == nil {
			cmd = exec.Command("npx", "vitest", "run")
		} else {
			cmd = c.createTestCommand(pm)
		}
	case "mocha":
		if _, err := exec.LookPath("npx"); err == nil {
			cmd = exec.Command("npx", "mocha")
		} else {
			cmd = c.createTestCommand(pm)
		}
	default:
		cmd = c.createTestCommand(pm)
	}

	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	combinedOutput := stdout.String() + stderr.String()
	if err != nil {
		errOutput := stderr.String()
		if errOutput == "" {
			errOutput = stdout.String()
		}
		return rb.FailWithOutput(fmt.Sprintf("Tests failed: %s", checkutil.TruncateMessage(errOutput, 200)), combinedOutput), nil
	}

	return rb.PassWithOutput("All tests passed", combinedOutput), nil
}

// createTestCommand creates a test command for the given package manager.
// This validates the package manager to avoid command injection.
func (c *TestsCheck) createTestCommand(pm string) *exec.Cmd {
	switch pm {
	case "npm":
		return exec.Command("npm", "test")
	case "yarn":
		return exec.Command("yarn", "test")
	case "pnpm":
		return exec.Command("pnpm", "test")
	case "bun":
		return exec.Command("bun", "test")
	default:
		return exec.Command("npm", "test")
	}
}

// detectPackageManager determines which package manager to use.
func (c *TestsCheck) detectPackageManager(path string) string {
	if c.Config != nil && c.Config.PackageManager != "auto" && c.Config.PackageManager != "" {
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

// detectTestRunner determines which test runner to use.
func (c *TestsCheck) detectTestRunner(path string, pkg *PackageJSON) string {
	// Check config override first
	if c.Config != nil && c.Config.TestRunner != "auto" && c.Config.TestRunner != "" {
		return c.Config.TestRunner
	}

	// Check for Jest config files
	jestConfigs := []string{"jest.config.js", "jest.config.ts", "jest.config.mjs", "jest.config.cjs", "jest.config.json"}
	for _, cfg := range jestConfigs {
		if safepath.Exists(path, cfg) {
			return "jest"
		}
	}

	// Check for Vitest config files
	vitestConfigs := []string{"vitest.config.js", "vitest.config.ts", "vitest.config.mjs", "vitest.config.mts"}
	for _, cfg := range vitestConfigs {
		if safepath.Exists(path, cfg) {
			return "vitest"
		}
	}

	// Check for Mocha config files
	mochaConfigs := []string{".mocharc.js", ".mocharc.json", ".mocharc.yaml", ".mocharc.yml", "mocha.opts"}
	for _, cfg := range mochaConfigs {
		if safepath.Exists(path, cfg) {
			return "mocha"
		}
	}

	// Check devDependencies
	if pkg != nil {
		if _, ok := pkg.DevDependencies["jest"]; ok {
			return "jest"
		}
		if _, ok := pkg.DevDependencies["vitest"]; ok {
			return "vitest"
		}
		if _, ok := pkg.DevDependencies["mocha"]; ok {
			return "mocha"
		}
	}

	// Default to using npm test script
	return "npm-test"
}
