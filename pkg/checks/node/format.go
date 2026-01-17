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

// FormatCheck verifies code formatting with prettier or biome.
type FormatCheck struct {
	Config *config.NodeLanguageConfig
}

// ID returns the unique identifier for this check.
func (c *FormatCheck) ID() string {
	return "node:format"
}

// Name returns the human-readable name for this check.
func (c *FormatCheck) Name() string {
	return "Node Format"
}

// Run executes the format check.
func (c *FormatCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	formatter := c.detectFormatter(path)

	// Handle auto mode - try available formatters
	if formatter == "auto" {
		if _, err := exec.LookPath("npx"); err != nil {
			return rb.Pass("npx not available, skipping format check"), nil
		}

		// Try prettier first, then biome
		prettierResult := c.runPrettier(path, rb)
		if prettierResult != nil {
			return *prettierResult, nil
		}

		biomeResult := c.runBiome(path, rb)
		if biomeResult != nil {
			return *biomeResult, nil
		}

		return rb.Pass("No formatter configured (prettier or biome)"), nil
	}

	// Run specific formatter
	switch formatter {
	case "prettier":
		if prettierResult := c.runPrettier(path, rb); prettierResult != nil {
			return *prettierResult, nil
		}
		return rb.ToolNotInstalled("Prettier", "npm install prettier"), nil
	case "biome":
		if biomeResult := c.runBiome(path, rb); biomeResult != nil {
			return *biomeResult, nil
		}
		return rb.ToolNotInstalled("Biome", "npm install @biomejs/biome"), nil
	default:
		return rb.Pass(fmt.Sprintf("Unknown formatter: %s", formatter)), nil
	}
}

// runPrettier runs prettier --check and returns a result or nil if prettier is not available.
func (c *FormatCheck) runPrettier(path string, rb *checkutil.ResultBuilder) *checker.Result {
	// Check if prettier config exists
	if !c.hasPrettierConfig(path) {
		return nil
	}

	cmd := exec.Command("npx", "prettier", "--check", ".")
	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// Count files needing formatting
	unformattedCount := 0
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(strings.ToLower(line), "would reformat") ||
			strings.Contains(line, "[warn]") {
			unformattedCount++
		}
	}

	if err != nil {
		if unformattedCount > 0 {
			result := rb.WarnWithOutput(fmt.Sprintf("%d %s need formatting. Run: npx prettier --write .", unformattedCount, checkutil.Pluralize(unformattedCount, "file", "files")), output)
			return &result
		}
		result := rb.WarnWithOutput("Files need formatting. Run: npx prettier --write .", output)
		return &result
	}

	result := rb.Pass("All files properly formatted (prettier)")
	return &result
}

// runBiome runs biome format --check and returns a result or nil if biome is not available.
func (c *FormatCheck) runBiome(path string, rb *checkutil.ResultBuilder) *checker.Result {
	// Check if biome config exists
	if !c.hasBiomeConfig(path) {
		return nil
	}

	cmd := exec.Command("npx", "@biomejs/biome", "format", "--check", ".")
	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		result := rb.WarnWithOutput("Files need formatting. Run: npx @biomejs/biome format --write .", output)
		return &result
	}

	result := rb.Pass("All files properly formatted (biome)")
	return &result
}

// detectFormatter determines which formatter to use.
func (c *FormatCheck) detectFormatter(path string) string {
	// Check config override first
	if c.Config != nil && c.Config.Formatter != "auto" && c.Config.Formatter != "" {
		return c.Config.Formatter
	}

	// Check for Prettier config files
	if c.hasPrettierConfig(path) {
		return "prettier"
	}

	// Check for Biome config files
	if c.hasBiomeConfig(path) {
		return "biome"
	}

	// Check devDependencies
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		if _, ok := pkg.DevDependencies["prettier"]; ok {
			return "prettier"
		}
		if _, ok := pkg.DevDependencies["@biomejs/biome"]; ok {
			return "biome"
		}
	}

	return "auto"
}

// hasPrettierConfig checks if prettier configuration exists.
func (c *FormatCheck) hasPrettierConfig(path string) bool {
	prettierConfigs := []string{
		".prettierrc",
		".prettierrc.json",
		".prettierrc.yaml",
		".prettierrc.yml",
		".prettierrc.js",
		".prettierrc.cjs",
		".prettierrc.mjs",
		"prettier.config.js",
		"prettier.config.cjs",
		"prettier.config.mjs",
	}
	for _, cfg := range prettierConfigs {
		if safepath.Exists(path, cfg) {
			return true
		}
	}

	// Also check devDependencies
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		if _, ok := pkg.DevDependencies["prettier"]; ok {
			return true
		}
	}

	return false
}

// hasBiomeConfig checks if biome configuration exists.
func (c *FormatCheck) hasBiomeConfig(path string) bool {
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
