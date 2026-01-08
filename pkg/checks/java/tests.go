package javacheck

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TestsCheck runs Java tests.
type TestsCheck struct {
	Config *config.JavaLanguageConfig
}

func (c *TestsCheck) ID() string   { return "java:tests" }
func (c *TestsCheck) Name() string { return "Java Tests" }

// Run executes Java tests using Maven or Gradle.
func (c *TestsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangJava,
	}

	buildTool := c.detectBuildTool(path)
	if buildTool == "" {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No build tool detected"
		return result, nil
	}

	var cmd *exec.Cmd
	switch buildTool {
	case "maven":
		cmd = c.getMavenTestCommand(path)
	case "gradle":
		cmd = c.getGradleTestCommand(path)
	}

	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		result.Passed = false
		result.Status = checker.Fail
		msg := "Tests failed"
		// Try to extract test summary
		summary := extractTestSummary(output, buildTool)
		if summary != "" {
			msg += ": " + summary
		}
		result.Message = msg
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	summary := extractTestSummary(output, buildTool)
	if summary != "" {
		result.Message = summary
	} else {
		result.Message = "All tests passed"
	}
	return result, nil
}

func (c *TestsCheck) detectBuildTool(path string) string {
	if c.Config != nil && c.Config.BuildTool != "" && c.Config.BuildTool != "auto" {
		return c.Config.BuildTool
	}
	return detectBuildTool(path)
}

func (c *TestsCheck) getMavenTestCommand(path string) *exec.Cmd {
	if safepath.Exists(path, "mvnw") {
		return exec.Command("./mvnw", "test", "-q")
	}
	return exec.Command("mvn", "test", "-q")
}

func (c *TestsCheck) getGradleTestCommand(path string) *exec.Cmd {
	if safepath.Exists(path, "gradlew") {
		return exec.Command("./gradlew", "test", "--no-daemon")
	}
	return exec.Command("gradle", "test", "--no-daemon")
}

// extractTestSummary tries to extract test results from build output.
func extractTestSummary(output, buildTool string) string {
	switch buildTool {
	case "maven":
		// Maven surefire format: "Tests run: X, Failures: Y, Errors: Z, Skipped: W"
		re := regexp.MustCompile(`Tests run: (\d+), Failures: (\d+), Errors: (\d+), Skipped: (\d+)`)
		matches := re.FindAllStringSubmatch(output, -1)
		if len(matches) > 0 {
			// Sum up all test results
			var totalRun, totalFail, totalErr, totalSkip int
			for _, m := range matches {
				if len(m) >= 5 {
					totalRun += parseInt(m[1])
					totalFail += parseInt(m[2])
					totalErr += parseInt(m[3])
					totalSkip += parseInt(m[4])
				}
			}
			if totalFail > 0 || totalErr > 0 {
				return strings.TrimSpace(strings.Join([]string{
					intToStr(totalRun), "tests,",
					intToStr(totalFail), "failures,",
					intToStr(totalErr), "errors",
				}, " "))
			}
			return intToStr(totalRun) + " tests passed"
		}
	case "gradle":
		// Gradle format varies, try common patterns
		// "X tests completed, Y failed"
		re := regexp.MustCompile(`(\d+) tests? completed`)
		if m := re.FindStringSubmatch(output); len(m) > 1 {
			return m[1] + " tests completed"
		}
	}
	return ""
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
