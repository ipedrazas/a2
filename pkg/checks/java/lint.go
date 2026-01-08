package javacheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LintCheck detects Java static analysis tools configuration.
type LintCheck struct{}

func (c *LintCheck) ID() string   { return "java:lint" }
func (c *LintCheck) Name() string { return "Java Lint" }

// Run checks for Checkstyle, SpotBugs, and PMD configuration.
func (c *LintCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangJava,
	}

	var linters []string

	// Check for Checkstyle
	if c.hasCheckstyle(path) {
		linters = append(linters, "Checkstyle")
	}

	// Check for SpotBugs (successor to FindBugs)
	if c.hasSpotBugs(path) {
		linters = append(linters, "SpotBugs")
	}

	// Check for PMD
	if c.hasPMD(path) {
		linters = append(linters, "PMD")
	}

	// Check for Error Prone
	if c.hasErrorProne(path) {
		linters = append(linters, "Error Prone")
	}

	// Check for SonarQube/SonarLint
	if c.hasSonar(path) {
		linters = append(linters, "SonarQube")
	}

	if len(linters) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Static analysis configured: " + strings.Join(linters, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No static analysis tools configured (consider Checkstyle, SpotBugs, or PMD)"
	}

	return result, nil
}

func (c *LintCheck) hasCheckstyle(path string) bool {
	// Check for standalone config files
	checkstyleConfigs := []string{
		"checkstyle.xml",
		"config/checkstyle/checkstyle.xml",
		"checkstyle/checkstyle.xml",
	}
	for _, config := range checkstyleConfigs {
		if safepath.Exists(path, config) {
			return true
		}
	}

	// Check for plugin in build files
	return c.hasBuildPlugin(path, "checkstyle")
}

func (c *LintCheck) hasSpotBugs(path string) bool {
	// Check for standalone config files
	spotbugsConfigs := []string{
		"spotbugs.xml",
		"spotbugs-exclude.xml",
		"findbugs-exclude.xml",
		"config/spotbugs/exclude.xml",
	}
	for _, config := range spotbugsConfigs {
		if safepath.Exists(path, config) {
			return true
		}
	}

	// Check for plugin in build files
	return c.hasBuildPlugin(path, "spotbugs") || c.hasBuildPlugin(path, "findbugs")
}

func (c *LintCheck) hasPMD(path string) bool {
	// Check for standalone config files
	pmdConfigs := []string{
		"pmd.xml",
		"ruleset.xml",
		"pmd-ruleset.xml",
		"config/pmd/ruleset.xml",
	}
	for _, config := range pmdConfigs {
		if safepath.Exists(path, config) {
			return true
		}
	}

	// Check for plugin in build files
	return c.hasBuildPlugin(path, "pmd")
}

func (c *LintCheck) hasErrorProne(path string) bool {
	return c.hasBuildPlugin(path, "error-prone") || c.hasBuildPlugin(path, "errorprone")
}

func (c *LintCheck) hasSonar(path string) bool {
	// Check for SonarQube config files
	if safepath.Exists(path, "sonar-project.properties") {
		return true
	}

	// Check for plugin in build files
	return c.hasBuildPlugin(path, "sonar") || c.hasBuildPlugin(path, "sonarqube")
}

func (c *LintCheck) hasBuildPlugin(path, pluginName string) bool {
	// Check Maven pom.xml
	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(strings.ToLower(string(content)), pluginName) {
			return true
		}
	}

	// Check Gradle build.gradle
	if safepath.Exists(path, "build.gradle") {
		content, err := safepath.ReadFile(path, "build.gradle")
		if err == nil && strings.Contains(strings.ToLower(string(content)), pluginName) {
			return true
		}
	}

	// Check Gradle build.gradle.kts
	if safepath.Exists(path, "build.gradle.kts") {
		content, err := safepath.ReadFile(path, "build.gradle.kts")
		if err == nil && strings.Contains(strings.ToLower(string(content)), pluginName) {
			return true
		}
	}

	return false
}
