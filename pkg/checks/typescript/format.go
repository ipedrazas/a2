package typescriptcheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// FormatCheck verifies code formatting in TypeScript projects.
type FormatCheck struct {
	Config *config.TypeScriptLanguageConfig
}

func (c *FormatCheck) ID() string   { return "typescript:format" }
func (c *FormatCheck) Name() string { return "TypeScript Format" }

// Run checks code formatting.
func (c *FormatCheck) Run(path string) (checker.Result, error) {
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

	// Detect formatter
	formatter := c.detectFormatter(path)
	pm := c.detectPackageManager(path)

	switch formatter {
	case "prettier":
		return c.runPrettier(path, pm)
	case "biome":
		return c.runBiome(path, pm)
	case "dprint":
		return c.runDprint(path, pm)
	default:
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No formatter configured (consider Prettier or Biome)"
		return result, nil
	}
}

// detectFormatter identifies which formatter is configured.
func (c *FormatCheck) detectFormatter(path string) string {
	if c.Config != nil && c.Config.Formatter != "" && c.Config.Formatter != "auto" {
		return c.Config.Formatter
	}

	// Check for Prettier config files
	prettierConfigs := []string{
		".prettierrc", ".prettierrc.json", ".prettierrc.yaml", ".prettierrc.yml",
		".prettierrc.js", ".prettierrc.cjs", ".prettierrc.mjs", "prettier.config.js",
		"prettier.config.cjs", "prettier.config.mjs",
	}
	for _, cfg := range prettierConfigs {
		if safepath.Exists(path, cfg) {
			return "prettier"
		}
	}

	// Check for Biome config
	if safepath.Exists(path, "biome.json") || safepath.Exists(path, "biome.jsonc") {
		return "biome"
	}

	// Check for dprint config
	if safepath.Exists(path, "dprint.json") || safepath.Exists(path, ".dprint.json") {
		return "dprint"
	}

	// Check package.json dependencies
	pkg, err := ParsePackageJSON(path)
	if err != nil {
		return ""
	}

	if _, ok := pkg.DevDependencies["prettier"]; ok {
		return "prettier"
	}
	if _, ok := pkg.DevDependencies["@biomejs/biome"]; ok {
		return "biome"
	}
	if _, ok := pkg.DevDependencies["dprint"]; ok {
		return "dprint"
	}

	return ""
}

// detectPackageManager determines which package manager to use.
func (c *FormatCheck) detectPackageManager(path string) string {
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

// runPrettier checks formatting with Prettier.
func (c *FormatCheck) runPrettier(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "prettier", "--check", ".")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "prettier", "--check", ".")
	case "bun":
		cmd = exec.Command("bun", "run", "prettier", "--check", ".")
	default:
		cmd = exec.Command("npx", "prettier", "--check", ".")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := stdout.String() + stderr.String()
		unformatted := countUnformattedFiles(output)
		if unformatted > 0 {
			result.Passed = false
			result.Status = checker.Warn
			result.Message = "Formatting issues found. Run: npx prettier --write ."
		} else {
			result.Passed = false
			result.Status = checker.Warn
			result.Message = "Prettier check failed"
		}
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "All files formatted correctly (Prettier)"
	return result, nil
}

// runBiome checks formatting with Biome.
func (c *FormatCheck) runBiome(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "biome", "format", "--error-on-warnings", ".")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "biome", "format", "--error-on-warnings", ".")
	case "bun":
		cmd = exec.Command("bun", "run", "biome", "format", "--error-on-warnings", ".")
	default:
		cmd = exec.Command("npx", "@biomejs/biome", "format", "--error-on-warnings", ".")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Formatting issues found. Run: npx @biomejs/biome format --write ."
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "All files formatted correctly (Biome)"
	return result, nil
}

// runDprint checks formatting with dprint.
func (c *FormatCheck) runDprint(path, pm string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	var cmd *exec.Cmd
	switch pm {
	case "yarn":
		cmd = exec.Command("yarn", "dprint", "check")
	case "pnpm":
		cmd = exec.Command("pnpm", "exec", "dprint", "check")
	case "bun":
		cmd = exec.Command("bun", "run", "dprint", "check")
	default:
		cmd = exec.Command("npx", "dprint", "check")
	}
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Formatting issues found. Run: npx dprint fmt"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "All files formatted correctly (dprint)"
	return result, nil
}

// countUnformattedFiles counts files that need formatting from Prettier output.
func countUnformattedFiles(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, ".ts") || strings.HasSuffix(line, ".tsx") ||
			strings.HasSuffix(line, ".js") || strings.HasSuffix(line, ".jsx") {
			count++
		}
	}
	return count
}
