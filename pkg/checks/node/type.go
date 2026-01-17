package nodecheck

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

// TypeCheck runs TypeScript compiler for type checking.
type TypeCheck struct {
	Config *config.NodeLanguageConfig
}

// ID returns the unique identifier for this check.
func (c *TypeCheck) ID() string {
	return "node:type"
}

// Name returns the human-readable name for this check.
func (c *TypeCheck) Name() string {
	return "TypeScript Type Check"
}

// Run executes the TypeScript type check.
func (c *TypeCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	// Check if this is a TypeScript project
	if !c.isTypeScriptProject(path) {
		return rb.Pass("Not a TypeScript project (no tsconfig.json)"), nil
	}

	// Check if npx is available
	if _, err := exec.LookPath("npx"); err != nil {
		return rb.Pass("npx not available, skipping type check"), nil
	}

	// Run tsc --noEmit for type checking
	cmd := exec.Command("npx", "tsc", "--noEmit")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		// Parse error count from output
		errorCount := c.countTypeErrors(output)
		if errorCount > 0 {
			return rb.WarnWithOutput(fmt.Sprintf("%d type %s found. Run: npx tsc --noEmit", errorCount, checkutil.Pluralize(errorCount, "error", "errors")), output), nil
		}
		return rb.WarnWithOutput("Type errors found. Run: npx tsc --noEmit", output), nil
	}

	return rb.Pass("No type errors found"), nil
}

// isTypeScriptProject checks if the project uses TypeScript.
func (c *TypeCheck) isTypeScriptProject(path string) bool {
	// Check for tsconfig.json
	if safepath.Exists(path, "tsconfig.json") {
		return true
	}

	// Check for tsconfig variants
	tsConfigs := []string{
		"tsconfig.base.json",
		"tsconfig.build.json",
		"tsconfig.app.json",
	}
	for _, cfg := range tsConfigs {
		if safepath.Exists(path, cfg) {
			return true
		}
	}

	// Check devDependencies for typescript
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		if _, ok := pkg.DevDependencies["typescript"]; ok {
			return true
		}
		if _, ok := pkg.Dependencies["typescript"]; ok {
			return true
		}
	}

	return false
}

// countTypeErrors parses tsc output to count errors.
func (c *TypeCheck) countTypeErrors(output string) int {
	// TypeScript outputs errors in format: "Found X errors."
	// or "Found X error." for single error
	re := regexp.MustCompile(`Found (\d+) errors?\.`)
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		count, err := strconv.Atoi(matches[1])
		if err == nil {
			return count
		}
	}

	// Fallback: count lines that look like errors
	// Format: path/file.ts(line,col): error TSxxxx: message
	count := 0
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "): error TS") {
			count++
		}
	}
	return count
}
