package gocheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// DeadcodeCheck detects unreachable code using the deadcode tool.
type DeadcodeCheck struct{}

func (c *DeadcodeCheck) ID() string   { return "go:deadcode" }
func (c *DeadcodeCheck) Name() string { return "Go Dead Code" }

func (c *DeadcodeCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	if !checkutil.ToolAvailable("deadcode") {
		return rb.ToolNotInstalled("deadcode", "go install golang.org/x/tools/cmd/deadcode@latest"), nil
	}

	result := checkutil.RunCommand(path, "deadcode", "-test", "./...")
	output := result.CombinedOutput()

	if !result.Success() {
		// deadcode may fail for build issues
		if strings.Contains(output, "no Go files") || strings.Contains(output, "build constraints") {
			return rb.Pass("No Go files to analyze"), nil
		}
		return rb.WarnWithOutput("deadcode error: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	// Count dead functions from output — each line is a dead function
	lines := countNonEmptyLines(result.Stdout)

	if lines == 0 {
		return rb.PassWithOutput("No unreachable functions detected", output), nil
	}

	msg := fmt.Sprintf("%d unreachable %s found",
		lines,
		checkutil.Pluralize(lines, "function", "functions"),
	)

	return rb.WarnWithOutput(msg, output), nil
}

// countNonEmptyLines counts non-empty lines in output.
func countNonEmptyLines(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}
