package nodecheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LintCheck runs linting on Node.js code.
type LintCheck struct {
	Config *config.NodeLanguageConfig
}

// ID returns the unique identifier for this check.
func (c *LintCheck) ID() string {
	return "node:lint"
}

// Name returns the human-readable name for this check.
func (c *LintCheck) Name() string {
	return "Node Lint"
}

// Run executes the lint check.
func (c *LintCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangNode,
	}

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = "package.json not found"
		return result, nil
	}

	linter := c.detectLinter(path)

	// Handle auto mode - try available linters
	if linter == "auto" {
		if _, err := exec.LookPath("npx"); err != nil {
			result.Status = checker.Pass
			result.Passed = true
			result.Message = "npx not available, skipping lint check"
			return result, nil
		}

		// Try eslint first, then biome, then oxlint
		eslintResult := c.runESLint(path)
		if eslintResult != nil {
			return *eslintResult, nil
		}

		biomeResult := c.runBiome(path)
		if biomeResult != nil {
			return *biomeResult, nil
		}

		oxlintResult := c.runOxlint(path)
		if oxlintResult != nil {
			return *oxlintResult, nil
		}

		result.Status = checker.Pass
		result.Passed = true
		result.Message = "No linter configured (eslint, biome, or oxlint)"
		return result, nil
	}

	// Run specific linter
	switch linter {
	case "eslint":
		if eslintResult := c.runESLint(path); eslintResult != nil {
			return *eslintResult, nil
		}
		result.Status = checker.Pass
		result.Passed = true
		result.Message = "ESLint not installed"
		return result, nil
	case "biome":
		if biomeResult := c.runBiome(path); biomeResult != nil {
			return *biomeResult, nil
		}
		result.Status = checker.Pass
		result.Passed = true
		result.Message = "Biome not installed"
		return result, nil
	case "oxlint":
		if oxlintResult := c.runOxlint(path); oxlintResult != nil {
			return *oxlintResult, nil
		}
		result.Status = checker.Pass
		result.Passed = true
		result.Message = "Oxlint not installed"
		return result, nil
	default:
		result.Status = checker.Pass
		result.Passed = true
		result.Message = fmt.Sprintf("Unknown linter: %s", linter)
		return result, nil
	}
}

// runESLint runs eslint and returns a result or nil if eslint is not available.
func (c *LintCheck) runESLint(path string) *checker.Result {
	// Check if eslint config exists
	if !c.hasESLintConfig(path) {
		return nil
	}

	cmd := exec.Command("npx", "eslint", ".", "--max-warnings=0")
	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	result := &checker.Result{
		Name:     "Node Lint",
		ID:       "node:lint",
		Language: checker.LangNode,
	}

	if err != nil {
		issueCount := countLintIssues(output)
		if issueCount > 0 {
			result.Status = checker.Warn
			result.Passed = false
			result.Message = fmt.Sprintf("%d linting %s found. Run: npx eslint . --fix", issueCount, pluralize(issueCount, "issue", "issues"))
		} else {
			result.Status = checker.Warn
			result.Passed = false
			result.Message = "Linting issues found. Run: npx eslint . --fix"
		}
		return result
	}

	result.Status = checker.Pass
	result.Passed = true
	result.Message = "No linting issues found (eslint)"
	return result
}

// runBiome runs biome lint and returns a result or nil if biome is not available.
func (c *LintCheck) runBiome(path string) *checker.Result {
	// Check if biome config exists
	if !c.hasBiomeConfig(path) {
		return nil
	}

	cmd := exec.Command("npx", "@biomejs/biome", "lint", ".")
	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &checker.Result{
		Name:     "Node Lint",
		ID:       "node:lint",
		Language: checker.LangNode,
	}

	if err != nil {
		result.Status = checker.Warn
		result.Passed = false
		result.Message = "Linting issues found. Run: npx @biomejs/biome lint --apply ."
		return result
	}

	result.Status = checker.Pass
	result.Passed = true
	result.Message = "No linting issues found (biome)"
	return result
}

// runOxlint runs oxlint and returns a result or nil if oxlint is not available.
func (c *LintCheck) runOxlint(path string) *checker.Result {
	// Check if oxlint config exists or is in devDependencies
	if !c.hasOxlintConfig(path) {
		return nil
	}

	cmd := exec.Command("npx", "oxlint", ".")
	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &checker.Result{
		Name:     "Node Lint",
		ID:       "node:lint",
		Language: checker.LangNode,
	}

	if err != nil {
		result.Status = checker.Warn
		result.Passed = false
		result.Message = "Linting issues found (oxlint)"
		return result
	}

	result.Status = checker.Pass
	result.Passed = true
	result.Message = "No linting issues found (oxlint)"
	return result
}

// detectLinter determines which linter to use.
func (c *LintCheck) detectLinter(path string) string {
	// Check config override first
	if c.Config != nil && c.Config.Linter != "auto" && c.Config.Linter != "" {
		return c.Config.Linter
	}

	// Check for ESLint config files
	if c.hasESLintConfig(path) {
		return "eslint"
	}

	// Check for Biome config files
	if c.hasBiomeConfig(path) {
		return "biome"
	}

	// Check for Oxlint config files
	if c.hasOxlintConfig(path) {
		return "oxlint"
	}

	// Check devDependencies
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		if _, ok := pkg.DevDependencies["eslint"]; ok {
			return "eslint"
		}
		if _, ok := pkg.DevDependencies["@biomejs/biome"]; ok {
			return "biome"
		}
		if _, ok := pkg.DevDependencies["oxlint"]; ok {
			return "oxlint"
		}
	}

	return "auto"
}

// hasESLintConfig checks if ESLint configuration exists.
func (c *LintCheck) hasESLintConfig(path string) bool {
	eslintConfigs := []string{
		".eslintrc",
		".eslintrc.json",
		".eslintrc.js",
		".eslintrc.cjs",
		".eslintrc.yaml",
		".eslintrc.yml",
		"eslint.config.js",
		"eslint.config.mjs",
		"eslint.config.cjs",
	}
	for _, cfg := range eslintConfigs {
		if safepath.Exists(path, cfg) {
			return true
		}
	}

	// Also check devDependencies
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		if _, ok := pkg.DevDependencies["eslint"]; ok {
			return true
		}
	}

	return false
}

// hasBiomeConfig checks if Biome configuration exists.
func (c *LintCheck) hasBiomeConfig(path string) bool {
	biomeConfigs := []string{"biome.json", "biome.jsonc"}
	for _, cfg := range biomeConfigs {
		if safepath.Exists(path, cfg) {
			return true
		}
	}

	// Also check devDependencies
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		if _, ok := pkg.DevDependencies["@biomejs/biome"]; ok {
			return true
		}
	}

	return false
}

// hasOxlintConfig checks if Oxlint configuration exists.
func (c *LintCheck) hasOxlintConfig(path string) bool {
	oxlintConfigs := []string{"oxlint.json", ".oxlintrc.json", "oxlintrc.json"}
	for _, cfg := range oxlintConfigs {
		if safepath.Exists(path, cfg) {
			return true
		}
	}

	// Also check devDependencies
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		if _, ok := pkg.DevDependencies["oxlint"]; ok {
			return true
		}
	}

	return false
}

// countLintIssues counts lint issues in the output.
func countLintIssues(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Count lines that look like lint errors/warnings
		if strings.Contains(line, "error") || strings.Contains(line, "warning") {
			count++
		}
	}
	return count
}
