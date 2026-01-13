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

// TypeCheck runs TypeScript compiler for type checking.
type TypeCheck struct {
	Config *config.TypeScriptLanguageConfig
}

func (c *TypeCheck) ID() string   { return "typescript:type" }
func (c *TypeCheck) Name() string { return "TypeScript Type Check" }

// Run executes the TypeScript type check.
func (c *TypeCheck) Run(path string) (checker.Result, error) {
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

	// Check if npx is available
	if _, err := exec.LookPath("npx"); err != nil {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "npx not available, skipping type check"
		return result, nil
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
		errorCount := countTypeErrors(output)
		if errorCount > 0 {
			result.Status = checker.Fail
			result.Passed = false
			result.Message = fmt.Sprintf("%d type %s found. Run: npx tsc --noEmit",
				errorCount, checkutil.Pluralize(errorCount, "error", "errors"))
		} else {
			result.Status = checker.Fail
			result.Passed = false
			result.Message = "Type errors found. Run: npx tsc --noEmit"
		}
		return result, nil
	}

	result.Status = checker.Pass
	result.Passed = true
	result.Message = "No type errors found"
	return result, nil
}

// countTypeErrors parses tsc output to count errors.
func countTypeErrors(output string) int {
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
