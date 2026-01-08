package javacheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ProjectCheck verifies that a Java project configuration exists.
type ProjectCheck struct{}

func (c *ProjectCheck) ID() string   { return "java:project" }
func (c *ProjectCheck) Name() string { return "Java Project" }

// Run checks for Maven or Gradle project files.
func (c *ProjectCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangJava,
	}

	buildTool := detectBuildTool(path)

	switch buildTool {
	case "maven":
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Maven project (pom.xml)"
		if safepath.Exists(path, "mvnw") {
			result.Message += " with wrapper"
		}
	case "gradle":
		result.Passed = true
		result.Status = checker.Pass
		msg := "Gradle project"
		if safepath.Exists(path, "build.gradle.kts") {
			msg += " (Kotlin DSL)"
		} else {
			msg += " (Groovy DSL)"
		}
		if safepath.Exists(path, "gradlew") {
			msg += " with wrapper"
		}
		result.Message = msg
	default:
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No Java project file found (pom.xml or build.gradle)"
	}

	return result, nil
}

// detectBuildTool returns "maven", "gradle", or "" based on project files.
func detectBuildTool(path string) string {
	// Check for Gradle first (more modern)
	if safepath.Exists(path, "build.gradle") || safepath.Exists(path, "build.gradle.kts") {
		return "gradle"
	}
	// Check for Maven
	if safepath.Exists(path, "pom.xml") {
		return "maven"
	}
	// Check for wrapper scripts as fallback
	if safepath.Exists(path, "gradlew") {
		return "gradle"
	}
	if safepath.Exists(path, "mvnw") {
		return "maven"
	}
	return ""
}
