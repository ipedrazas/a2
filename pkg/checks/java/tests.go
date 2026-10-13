package javacheck

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
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
func (c *TestsCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangJava)

	buildTool := c.detectBuildTool(path)
	if buildTool == "" {
		return rb.Fail("No build tool detected"), nil
	}

	var cmd *exec.Cmd
	switch buildTool {
	case "maven":
		cmd = c.getMavenTestCommand(ctx, path)
	case "gradle":
		cmd = c.getGradleTestCommand(ctx, path)
	}

	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		msg := "Tests failed"
		// Try to extract test summary
		summary := extractTestSummary(output, buildTool)
		if summary != "" {
			msg += ": " + summary
		}
		return rb.FailWithOutput(msg, output), nil
	}

	summary := extractTestSummary(output, buildTool)
	if summary != "" {
		return rb.PassWithOutput(summary, output), nil
	}
	return rb.PassWithOutput("All tests passed", output), nil
}

func (c *TestsCheck) detectBuildTool(path string) string {
	if c.Config != nil && c.Config.BuildTool != "" && c.Config.BuildTool != "auto" {
		return c.Config.BuildTool
	}
	return detectBuildTool(path)
}

func (c *TestsCheck) getMavenTestCommand(ctx context.Context, path string) *exec.Cmd {
	if safepath.Exists(path, "mvnw") {
		return checkutil.Command(ctx, "./mvnw", "test", "-q")
	}
	return checkutil.Command(ctx, "mvn", "test", "-q")
}

func (c *TestsCheck) getGradleTestCommand(ctx context.Context, path string) *exec.Cmd {
	if safepath.Exists(path, "gradlew") {
		return checkutil.Command(ctx, "./gradlew", "test", "--no-daemon")
	}
	return checkutil.Command(ctx, "gradle", "test", "--no-daemon")
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
