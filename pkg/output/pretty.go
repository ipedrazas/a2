package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/maturity"
	"github.com/ipedrazas/a2/pkg/runner"
)

var (
	// Colors
	green  = lipgloss.Color("#22c55e")
	yellow = lipgloss.Color("#eab308")
	red    = lipgloss.Color("#ef4444")
	gray   = lipgloss.Color("#6b7280")
	white  = lipgloss.Color("#f9fafb")
	cyan   = lipgloss.Color("#06b6d4")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			MarginBottom(1)

	separatorStyle = lipgloss.NewStyle().
			Foreground(gray)

	passStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	warnStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true)

	failStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	messageStyle = lipgloss.NewStyle().
			Foreground(gray).
			PaddingLeft(4)

	statusPassStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	statusWarnStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true)

	statusFailStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	scoreStyle = lipgloss.NewStyle().
			Foreground(cyan)

	recommendStyle = lipgloss.NewStyle().
			Foreground(gray).
			Italic(true)

	maturityStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	maturityDescStyle = lipgloss.NewStyle().
				Foreground(gray)
)

// Pretty outputs the results in a formatted, colorful way.
func Pretty(result runner.SuiteResult, path string, detected language.DetectionResult) error {
	// Get project name from path
	projectName := filepath.Base(path)
	if path == "." {
		if wd, err := os.Getwd(); err == nil {
			projectName = filepath.Base(wd)
		}
	}

	// Title
	fmt.Println(titleStyle.Render(fmt.Sprintf("A2 Analysis: %s", projectName)))
	fmt.Println(separatorStyle.Render("─────────────────────────────────────"))

	// Show detected languages
	if len(detected.Languages) > 0 {
		langs := make([]string, len(detected.Languages))
		for i, l := range detected.Languages {
			langs[i] = string(l)
		}
		fmt.Println(scoreStyle.Render(fmt.Sprintf("Languages: %s", strings.Join(langs, ", "))))
	}
	fmt.Println()

	// Results
	for _, r := range result.Results {
		printResult(r)
	}

	fmt.Println()
	fmt.Println(separatorStyle.Render("─────────────────────────────────────"))

	// Status
	printStatus(result)

	// Score
	printScore(result)

	// Maturity estimation
	printMaturity(result)

	// Recommendations
	printRecommendations(result)

	if !result.Success() {
		os.Exit(1)
	}
	return nil
}

func printResult(r checker.Result) {
	var symbol, status string
	var style lipgloss.Style

	switch r.Status {
	case checker.Pass:
		symbol = "✓"
		status = "PASS"
		style = passStyle
	case checker.Warn:
		symbol = "!"
		status = "WARN"
		style = warnStyle
	case checker.Fail:
		symbol = "✗"
		status = "FAIL"
		style = failStyle
	case checker.Info:
		symbol = "ℹ"
		status = "INFO"
		style = infoStyle
	}

	// Print the check result
	fmt.Printf("%s %s %s\n",
		style.Render(symbol),
		style.Render(status),
		r.Name,
	)

	// Print message if present
	if r.Message != "" {
		fmt.Println(messageStyle.Render(r.Message))
	}
}

func printStatus(result runner.SuiteResult) {
	fmt.Println()

	if result.Failed > 0 {
		if result.Aborted {
			fmt.Println(statusFailStyle.Render("STATUS: ✗ CRITICAL FAILURE (Aborted)"))
		} else {
			fmt.Println(statusFailStyle.Render("STATUS: ✗ FAILED"))
		}
	} else if result.Warnings > 0 {
		fmt.Println(statusWarnStyle.Render("STATUS: ⚠ NEEDS ATTENTION"))
	} else {
		fmt.Println(statusPassStyle.Render("STATUS: ✓ ALL CHECKS PASSED"))
	}
}

func printScore(result runner.SuiteResult) {
	// Use ScoredChecks to exclude Info from score calculation
	scoredTotal := result.ScoredChecks()
	passed := result.Passed

	// Calculate percentage
	var pct float64
	if scoredTotal > 0 {
		pct = float64(passed) / float64(scoredTotal) * 100
	}

	fmt.Println()
	scoreMsg := fmt.Sprintf("Score: %d/%d checks passed (%.0f%%)", passed, scoredTotal, pct)
	if result.Info > 0 {
		scoreMsg += fmt.Sprintf(" + %d info", result.Info)
	}
	fmt.Println(scoreStyle.Render(scoreMsg))
}

func printMaturity(result runner.SuiteResult) {
	est := maturity.Estimate(result)

	fmt.Println()
	fmt.Println(maturityStyle.Render(fmt.Sprintf("Maturity: %s", est.Level.String())))
	fmt.Println(maturityDescStyle.Render(fmt.Sprintf("   %s", est.Level.Description())))

	if len(est.Suggestions) > 0 {
		fmt.Println()
		for _, s := range est.Suggestions {
			fmt.Println(maturityDescStyle.Render(fmt.Sprintf("   → %s", s)))
		}
	}
}

func printRecommendations(result runner.SuiteResult) {
	var recommendations []string

	for _, r := range result.Results {
		if !r.Passed {
			switch r.ID {
			// Go checks (new IDs)
			case "go:coverage":
				recommendations = append(recommendations, "→ Add more tests to improve coverage")
			case "go:format":
				recommendations = append(recommendations, "→ Run 'gofmt -w .' to format code")
			case "go:vet":
				recommendations = append(recommendations, "→ Fix issues reported by 'go vet ./...'")
			case "go:tests":
				recommendations = append(recommendations, "→ Fix failing tests before continuing")
			case "go:build":
				recommendations = append(recommendations, "→ Fix build errors before continuing")
			case "go:deps":
				recommendations = append(recommendations, "→ Update dependencies to fix vulnerabilities")
			// Python checks
			case "python:tests":
				recommendations = append(recommendations, "→ Fix failing tests before continuing")
			case "python:format":
				recommendations = append(recommendations, "→ Run formatter (black/ruff) to format code")
			case "python:lint":
				recommendations = append(recommendations, "→ Fix linting issues")
			case "python:coverage":
				recommendations = append(recommendations, "→ Add more tests to improve coverage")
			case "python:deps":
				recommendations = append(recommendations, "→ Update dependencies to fix vulnerabilities")
			// Common checks
			case "file_exists":
				recommendations = append(recommendations, "→ Add missing documentation files")
			}
		}
	}

	if len(recommendations) > 0 {
		fmt.Println()
		fmt.Println(recommendStyle.Render("Recommendations:"))
		for _, rec := range recommendations {
			fmt.Println(recommendStyle.Render(rec))
		}
	}
}
