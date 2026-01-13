package typescriptcheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LintCheck runs linting on TypeScript code.
type LintCheck struct {
	Config *config.TypeScriptLanguageConfig
}

func (c *LintCheck) ID() string   { return "typescript:lint" }
func (c *LintCheck) Name() string { return "TypeScript Lint" }

// Run executes the linter.
func (c *LintCheck) Run(path string) (checker.Result, error) {
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

	// Detect linter
	linter := c.detectLinter(path)
	pm := c.detectPackageManager(path)

	switch linter {
	case "eslint":
		return c.runESLint(path, pm)
	case "biome":
		return c.runBiomeLint(path, pm)
	case "oxlint":
		return c.runOxlint(path, pm)
	default:
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No linter configured (consider ESLint or Biome)"
		return result, nil
	}
}

// detectLinter identifies which linter is configured.
func (c *LintCheck) detectLinter(path string) string {
	if c.Config != nil && c.Config.Linter != "" && c.Config.Linter != "auto" {
		return c.Config.Linter
	}

	// Check for ESLint config files
	eslintConfigs := []string{
		".eslintrc.js", ".eslintrc.cjs", ".eslintrc.json", ".eslintrc.yaml", ".eslintrc.yml",
		".eslintrc", "eslint.config.js", "eslint.config.mjs", "eslint.config.cjs",
	}
	for _, cfg := range eslintConfigs {
		if safepath.Exists(path, cfg) {
			return "eslint"
		}
	}

	// Check for Biome config
	if safepath.Exists(path, "biome.json") || safepath.Exists(path, "biome.jsonc") {
		return "biome"
	}

	// Check package.json dependencies
	pkg, err := ParsePackageJSON(path)
	if err != nil {
		return ""
	}

	if _, ok := pkg.DevDependencies["eslint"]; ok {
		return "eslint"
	}
	if _, ok := pkg.DevDependencies["@biomejs/biome"]; ok {
		return "biome"
	}
	if _, ok := pkg.DevDependencies["oxlint"]; ok {
		return "oxlint"
	}

	return ""
}

// detectPackageManager determines which package manager to use.
func (c *LintCheck) detectPackageManager(path string) string {
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

// runESLint runs ESLint.
func (c *LintCheck) runESLint(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "eslint", ".", "--ext", ".ts,.tsx")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "eslint", ".", "--ext", ".ts,.tsx")
	case "bun":
		cmd = exec.Command("bun", "run", "eslint", ".", "--ext", ".ts,.tsx")
	default:
		cmd = exec.Command("npx", "eslint", ".", "--ext", ".ts,.tsx")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := stdout.String() + stderr.String()
		errors, warnings := parseESLintOutput(output)
		if errors > 0 || warnings > 0 {
			result.Passed = false
			result.Status = checker.Warn
			result.Message = fmt.Sprintf("ESLint: %d %s, %d %s",
				errors, checkutil.Pluralize(errors, "error", "errors"),
				warnings, checkutil.Pluralize(warnings, "warning", "warnings"))
		} else {
			result.Passed = false
			result.Status = checker.Warn
			result.Message = "ESLint found issues"
		}
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "ESLint: No issues found"
	return result, nil
}

// runBiomeLint runs Biome linter.
func (c *LintCheck) runBiomeLint(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "biome", "lint", ".")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "biome", "lint", ".")
	case "bun":
		cmd = exec.Command("bun", "run", "biome", "lint", ".")
	default:
		cmd = exec.Command("npx", "@biomejs/biome", "lint", ".")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Biome lint found issues"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Biome lint: No issues found"
	return result, nil
}

// runOxlint runs oxlint.
func (c *LintCheck) runOxlint(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "oxlint")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "oxlint")
	case "bun":
		cmd = exec.Command("bun", "run", "oxlint")
	default:
		cmd = exec.Command("npx", "oxlint")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "oxlint found issues"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "oxlint: No issues found"
	return result, nil
}

// parseESLintOutput extracts error and warning counts from ESLint output.
func parseESLintOutput(output string) (errors, warnings int) {
	// ESLint summary format: "X problems (Y errors, Z warnings)"
	re := regexp.MustCompile(`(\d+) problems? \((\d+) errors?, (\d+) warnings?\)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 4 {
		errors, _ = strconv.Atoi(matches[2])
		warnings, _ = strconv.Atoi(matches[3])
		return
	}

	// Alternative format: count lines with error/warning
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "error") {
			errors++
		}
		if strings.Contains(line, "warning") {
			warnings++
		}
	}
	return
}
