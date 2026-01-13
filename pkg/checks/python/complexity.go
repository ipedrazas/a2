package pythoncheck

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ComplexityCheck measures cyclomatic complexity of Python functions using radon.
type ComplexityCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *ComplexityCheck) ID() string   { return "python:complexity" }
func (c *ComplexityCheck) Name() string { return "Python Complexity" }

// ComplexFunction holds complexity info for a single function.
type ComplexFunction struct {
	Name       string
	File       string
	Line       int
	Complexity int
	Grade      string // A, B, C, D, E, F
}

func (c *ComplexityCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)

	// Check if Python project exists
	hasPython := safepath.Exists(path, "pyproject.toml") ||
		safepath.Exists(path, "setup.py") ||
		safepath.Exists(path, "requirements.txt")

	if !hasPython {
		return rb.Fail("Python project not found"), nil
	}

	// Get threshold
	threshold := 15
	if c.Config != nil && c.Config.CyclomaticThreshold > 0 {
		threshold = c.Config.CyclomaticThreshold
	}

	// Check if radon is installed
	if !checkutil.ToolAvailable("radon") {
		return rb.ToolNotInstalled("radon", "pip install radon"), nil
	}

	// Run radon cc with show-complexity flag
	result := checkutil.RunCommand(path, "radon", "cc", "-s", ".")

	if !result.Success() {
		// radon may fail if no Python files found
		if strings.Contains(result.Stderr, "No such file") ||
			strings.Contains(result.Stderr, "Invalid") {
			return rb.Pass("No Python files to analyze"), nil
		}
		// Continue even if there's an error, try to parse output
	}

	// Parse radon output
	complexFunctions := parseRadonOutput(result.Stdout, threshold)

	if len(complexFunctions) == 0 {
		return rb.Pass(fmt.Sprintf("No functions exceed complexity threshold (%d)", threshold)), nil
	}

	// Sort by complexity descending
	sort.Slice(complexFunctions, func(i, j int) bool {
		return complexFunctions[i].Complexity > complexFunctions[j].Complexity
	})

	// Build message with top offenders
	msg := checkutil.PluralizeCount(len(complexFunctions), "function exceeds", "functions exceed") +
		fmt.Sprintf(" complexity threshold (%d)", threshold)

	// Show top 3 offenders
	showCount := 3
	if len(complexFunctions) < showCount {
		showCount = len(complexFunctions)
	}

	for i := 0; i < showCount; i++ {
		f := complexFunctions[i]
		msg += fmt.Sprintf("\n  • %s (%s:%d) = %d [%s]", f.Name, f.File, f.Line, f.Complexity, f.Grade)
	}

	if len(complexFunctions) > showCount {
		msg += fmt.Sprintf("\n  ... and %d more", len(complexFunctions)-showCount)
	}

	return rb.Warn(msg), nil
}

// parseRadonOutput parses radon cc output and returns functions exceeding threshold.
// Radon output format:
// path/to/file.py
//
//	F 10:0 function_name - A (1)
//	M 20:4 ClassName.method - B (6)
func parseRadonOutput(output string, threshold int) []ComplexFunction {
	var complexFunctions []ComplexFunction

	// Pattern to match radon output lines
	// F/M/C for function/method/class, line:col, name, grade (complexity)
	linePattern := regexp.MustCompile(`^\s+([FMC])\s+(\d+):\d+\s+(\S+)\s+-\s+([A-F])\s+\((\d+)\)`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	var currentFile string

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a file path line (no leading whitespace, ends with .py)
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(line, " ") && strings.HasSuffix(trimmed, ".py") {
			currentFile = trimmed
			continue
		}

		// Try to match function/method line
		matches := linePattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		lineNum, err := strconv.Atoi(matches[2])
		if err != nil {
			continue
		}

		complexity, err := strconv.Atoi(matches[5])
		if err != nil {
			continue
		}

		if complexity > threshold {
			complexFunctions = append(complexFunctions, ComplexFunction{
				Name:       matches[3],
				File:       currentFile,
				Line:       lineNum,
				Complexity: complexity,
				Grade:      matches[4],
			})
		}
	}

	return complexFunctions
}
