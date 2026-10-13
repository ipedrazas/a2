package javacheck

import (
	"bytes"
	"context"
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
func (c *BuildCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangJava)

	buildTool := c.detectBuildTool(path)
	if buildTool == "" {
		return rb.Fail("No build tool detected (pom.xml or build.gradle)"), nil
	}

	var cmd *exec.Cmd
	switch buildTool {
	case "maven":
		cmd = c.getMavenCommand(ctx, path)
	case "gradle":
		cmd = c.getGradleCommand(ctx, path)
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

func (c *BuildCheck) getMavenCommand(ctx context.Context, path string) *exec.Cmd {
	// Prefer wrapper script
	if safepath.Exists(path, "mvnw") {
		return checkutil.Command(ctx, "./mvnw", "compile", "-q", "-DskipTests")
	}
	return checkutil.Command(ctx, "mvn", "compile", "-q", "-DskipTests")
}

func (c *BuildCheck) getGradleCommand(ctx context.Context, path string) *exec.Cmd {
	// Prefer wrapper script
	if safepath.Exists(path, "gradlew") {
		return checkutil.Command(ctx, "./gradlew", "compileJava", "-q", "--no-daemon")
	}
	return checkutil.Command(ctx, "gradle", "compileJava", "-q", "--no-daemon")
}
