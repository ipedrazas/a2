package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// IntegrationCheck verifies that integration tests exist.
type IntegrationCheck struct{}

func (c *IntegrationCheck) ID() string   { return "common:integration" }
func (c *IntegrationCheck) Name() string { return "Integration Tests" }

// Run checks for integration test directories and files.
func (c *IntegrationCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var found []string

	// Check for integration test directories
	integrationDirs := []string{
		"tests/integration",
		"test/integration",
		"integration_tests",
		"integration-tests",
		"tests/e2e",
		"test/e2e",
		"e2e",
		"e2e-tests",
	}

	for _, dir := range integrationDirs {
		if safepath.IsDir(path, dir) {
			found = append(found, dir+"/")
		}
	}

	// Check for Go integration test files
	goIntegrationPatterns := []string{
		"*_integration_test.go",
		"*_integ_test.go",
		"integration_test.go",
	}
	for _, pattern := range goIntegrationPatterns {
		if files, err := safepath.Glob(path, pattern); err == nil && len(files) > 0 {
			found = append(found, "Go integration tests")
			break
		}
	}

	// Check in subdirectories for Go
	goSubdirs := []string{"pkg", "internal", "cmd"}
	for _, subdir := range goSubdirs {
		for _, pattern := range goIntegrationPatterns {
			if files, err := safepath.Glob(path+"/"+subdir, "**/"+pattern); err == nil && len(files) > 0 {
				found = append(found, "Go integration tests")
				break
			}
		}
	}

	// Check for Python integration test files
	pythonIntegrationPatterns := []string{
		"test_integration_*.py",
		"test_*_integration.py",
		"*_integration_test.py",
	}
	for _, pattern := range pythonIntegrationPatterns {
		if files, err := safepath.Glob(path, pattern); err == nil && len(files) > 0 {
			found = append(found, "Python integration tests")
			break
		}
		if files, err := safepath.Glob(path+"/tests", pattern); err == nil && len(files) > 0 {
			found = append(found, "Python integration tests")
			break
		}
	}

	// Check for Node.js/TypeScript integration test files
	nodeIntegrationPatterns := []string{
		"*.integration.test.ts",
		"*.integration.test.js",
		"*.integration.spec.ts",
		"*.integration.spec.js",
		"*.e2e.test.ts",
		"*.e2e.test.js",
		"*.e2e.spec.ts",
		"*.e2e.spec.js",
	}
	for _, pattern := range nodeIntegrationPatterns {
		if files, err := safepath.Glob(path, pattern); err == nil && len(files) > 0 {
			found = append(found, "Node.js integration tests")
			break
		}
		testDirs := []string{"test", "tests", "__tests__", "src"}
		for _, dir := range testDirs {
			if files, err := safepath.Glob(path+"/"+dir, "**/"+pattern); err == nil && len(files) > 0 {
				found = append(found, "Node.js integration tests")
				break
			}
		}
	}

	// Check for test infrastructure (docker-compose for tests)
	testInfraFiles := []string{
		"docker-compose.test.yml",
		"docker-compose.test.yaml",
		"docker-compose.e2e.yml",
		"docker-compose.e2e.yaml",
		"docker-compose.integration.yml",
		"docker-compose.integration.yaml",
		"tests/docker-compose.yml",
		"test/docker-compose.yml",
	}
	for _, file := range testInfraFiles {
		if safepath.Exists(path, file) {
			found = append(found, "test infrastructure ("+file+")")
		}
	}

	// Check for testcontainers usage
	if c.hasTestcontainers(path) {
		found = append(found, "testcontainers")
	}

	// Check for E2E testing frameworks
	e2eFrameworks := c.detectE2EFrameworks(path)
	found = append(found, e2eFrameworks...)

	// Build result
	found = unique(found)
	if len(found) > 0 {
		return rb.Pass("Integration tests found: " + strings.Join(found, ", ")), nil
	}
	return rb.Warn("No integration tests found (consider adding tests/integration/ directory)"), nil
}

func (c *IntegrationCheck) hasTestcontainers(path string) bool {
	// Go
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			if strings.Contains(string(content), "testcontainers-go") {
				return true
			}
		}
	}

	// Python
	pythonConfigs := []string{"pyproject.toml", "requirements.txt", "requirements-dev.txt"}
	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				if strings.Contains(strings.ToLower(string(content)), "testcontainers") {
					return true
				}
			}
		}
	}

	// Node.js
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			if strings.Contains(string(content), "testcontainers") {
				return true
			}
		}
	}

	return false
}

func (c *IntegrationCheck) detectE2EFrameworks(path string) []string {
	var frameworks []string

	// Cypress
	cypressConfigs := []string{"cypress.config.js", "cypress.config.ts", "cypress.json"}
	for _, cfg := range cypressConfigs {
		if safepath.Exists(path, cfg) {
			frameworks = append(frameworks, "Cypress")
			break
		}
	}

	// Playwright
	playwrightConfigs := []string{"playwright.config.js", "playwright.config.ts"}
	for _, cfg := range playwrightConfigs {
		if safepath.Exists(path, cfg) {
			frameworks = append(frameworks, "Playwright")
			break
		}
	}

	// WebdriverIO
	wdioConfigs := []string{"wdio.conf.js", "wdio.conf.ts"}
	for _, cfg := range wdioConfigs {
		if safepath.Exists(path, cfg) {
			frameworks = append(frameworks, "WebdriverIO")
			break
		}
	}

	// Selenium (check in dependencies)
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			if strings.Contains(string(content), "selenium-webdriver") {
				frameworks = append(frameworks, "Selenium")
			}
		}
	}

	return frameworks
}
