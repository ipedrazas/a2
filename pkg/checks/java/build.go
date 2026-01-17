package javacheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck verifies that the Java project compiles.
type BuildCheck struct {
	Config *config.JavaLanguageConfig
}

func (c *BuildCheck) ID() string   { return "java:build" }
func (c *BuildCheck) Name() string { return "Java Build" }

// Run compiles the Java project using Maven or Gradle.
func (c *BuildCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangJava)

	buildTool := c.detectBuildTool(path)
	if buildTool == "" {
		return rb.Fail("No build tool detected (pom.xml or build.gradle)"), nil
	}

	var cmd *exec.Cmd
	switch buildTool {
	case "maven":
		cmd = c.getMavenCommand(path)
	case "gradle":
		cmd = c.getGradleCommand(path)
	}

	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	combinedOutput := stdout.String() + stderr.String()
	output := strings.TrimSpace(stdout.String())
	errOutput := strings.TrimSpace(stderr.String())

	if err != nil {
		msg := "Build failed"
		if errOutput != "" {
			// Truncate long error messages
			if len(errOutput) > 500 {
				errOutput = errOutput[:500] + "..."
			}
			msg += ": " + errOutput
		} else if output != "" {
			if len(output) > 500 {
				output = output[:500] + "..."
			}
			msg += ": " + output
		}
		return rb.FailWithOutput(msg, combinedOutput), nil
	}

	return rb.PassWithOutput("Build successful ("+buildTool+")", combinedOutput), nil
}

func (c *BuildCheck) detectBuildTool(path string) string {
	// Check config override
	if c.Config != nil && c.Config.BuildTool != "" && c.Config.BuildTool != "auto" {
		return c.Config.BuildTool
	}
	return detectBuildTool(path)
}

func (c *BuildCheck) getMavenCommand(path string) *exec.Cmd {
	// Prefer wrapper script
	if safepath.Exists(path, "mvnw") {
		return exec.Command("./mvnw", "compile", "-q", "-DskipTests")
	}
	return exec.Command("mvn", "compile", "-q", "-DskipTests")
}

func (c *BuildCheck) getGradleCommand(path string) *exec.Cmd {
	// Prefer wrapper script
	if safepath.Exists(path, "gradlew") {
		return exec.Command("./gradlew", "compileJava", "-q", "--no-daemon")
	}
	return exec.Command("gradle", "compileJava", "-q", "--no-daemon")
}
