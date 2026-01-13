package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// E2ECheck verifies end-to-end tests exist.
type E2ECheck struct{}

func (c *E2ECheck) ID() string   { return "common:e2e" }
func (c *E2ECheck) Name() string { return "E2E Tests" }

// Run checks for end-to-end test configuration.
func (c *E2ECheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var found []string

	// Check for e2e frameworks config files
	frameworkConfigs := map[string]string{
		"cypress.config.js":       "Cypress",
		"cypress.config.ts":       "Cypress",
		"cypress.config.mjs":      "Cypress",
		"cypress.json":            "Cypress",
		"playwright.config.ts":    "Playwright",
		"playwright.config.js":    "Playwright",
		"wdio.conf.js":            "WebdriverIO",
		"wdio.conf.ts":            "WebdriverIO",
		"nightwatch.conf.js":      "Nightwatch",
		"codecept.conf.js":        "CodeceptJS",
		"codecept.conf.ts":        "CodeceptJS",
		"testcafe.config.js":      "TestCafe",
		".testcaferc.json":        "TestCafe",
		"detox.config.js":         "Detox",
		".detoxrc.json":           "Detox",
		"appium.conf.js":          "Appium",
		"chimp.js":                "Chimp",
		"gauge.json":              "Gauge",
		"behave.ini":              "Behave",
		"features/support/env.rb": "Cucumber",
	}

	for file, framework := range frameworkConfigs {
		if safepath.Exists(path, file) {
			if !containsString(found, framework) {
				found = append(found, framework)
			}
		}
	}

	// Check for e2e directories
	e2eDirs := []string{
		"e2e",
		"cypress",
		"playwright",
		"tests/e2e",
		"test/e2e",
		"e2e-tests",
		"acceptance",
		"features", // Cucumber/Behave
		"specs",
	}
	for _, dir := range e2eDirs {
		if safepath.Exists(path, dir) {
			// Check if it's a directory with test files
			found = append(found, dir+" directory")
			break
		}
	}

	// Check for e2e dependencies in package.json
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			contentStr := string(content)
			e2eDeps := map[string]string{
				"cypress":            "Cypress",
				"playwright":         "Playwright",
				"@playwright/test":   "Playwright",
				"webdriverio":        "WebdriverIO",
				"puppeteer":          "Puppeteer",
				"selenium-webdriver": "Selenium",
				"nightwatch":         "Nightwatch",
				"testcafe":           "TestCafe",
				"codeceptjs":         "CodeceptJS",
				"detox":              "Detox",
			}
			for dep, name := range e2eDeps {
				if strings.Contains(contentStr, `"`+dep+`"`) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check for Python e2e tools
	if safepath.Exists(path, "pyproject.toml") || safepath.Exists(path, "requirements.txt") {
		pythonE2E := map[string]string{
			"selenium":   "Selenium",
			"playwright": "Playwright",
			"splinter":   "Splinter",
			"robot":      "Robot",
		}
		for tool, name := range pythonE2E {
			if c.hasPythonDep(path, tool) {
				if !containsString(found, name) {
					found = append(found, name)
				}
			}
		}
	}

	// Build result
	if len(found) > 0 {
		return rb.Pass("E2E testing: " + strings.Join(found, ", ")), nil
	}
	return rb.Warn("No E2E tests found (consider Playwright or Cypress)"), nil
}

// hasPythonDep checks if a Python dependency is present.
func (c *E2ECheck) hasPythonDep(path, dep string) bool {
	files := []string{"pyproject.toml", "requirements.txt", "requirements-dev.txt", "setup.py"}
	for _, file := range files {
		if content, err := safepath.ReadFile(path, file); err == nil {
			if strings.Contains(strings.ToLower(string(content)), dep) {
				return true
			}
		}
	}
	return false
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
