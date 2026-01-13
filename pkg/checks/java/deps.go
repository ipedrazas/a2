package javacheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsCheck verifies Java dependency security scanning configuration.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "java:deps" }
func (c *DepsCheck) Name() string { return "Java Dependencies" }

// Run checks for dependency vulnerability scanning tools.
func (c *DepsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangJava)

	var tools []string

	// Check for OWASP Dependency-Check
	if c.hasOWASPDependencyCheck(path) {
		tools = append(tools, "OWASP Dependency-Check")
	}

	// Check for Snyk
	if c.hasSnyk(path) {
		tools = append(tools, "Snyk")
	}

	// Check for Dependabot
	if c.hasDependabot(path) {
		tools = append(tools, "Dependabot")
	}

	// Check for Renovate
	if c.hasRenovate(path) {
		tools = append(tools, "Renovate")
	}

	// Check for Maven dependency plugin
	if c.hasMavenDependencyPlugin(path) {
		tools = append(tools, "Maven Dependency Plugin")
	}

	// Check for Gradle dependency verification
	if c.hasGradleDependencyVerification(path) {
		tools = append(tools, "Gradle Dependency Verification")
	}

	if len(tools) > 0 {
		return rb.Pass("Dependency scanning configured: " + strings.Join(tools, ", ")), nil
	}
	return rb.Warn("No dependency scanning configured (consider OWASP Dependency-Check or Snyk)"), nil
}

func (c *DepsCheck) hasOWASPDependencyCheck(path string) bool {
	// Check Maven pom.xml
	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(string(content), "dependency-check") {
			return true
		}
	}

	// Check Gradle build files
	if safepath.Exists(path, "build.gradle") {
		content, err := safepath.ReadFile(path, "build.gradle")
		if err == nil && strings.Contains(string(content), "dependency-check") {
			return true
		}
	}
	if safepath.Exists(path, "build.gradle.kts") {
		content, err := safepath.ReadFile(path, "build.gradle.kts")
		if err == nil && strings.Contains(string(content), "dependency-check") {
			return true
		}
	}

	return false
}

func (c *DepsCheck) hasSnyk(path string) bool {
	// Check for Snyk config files
	if safepath.Exists(path, ".snyk") {
		return true
	}

	// Check for Snyk in build files
	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(string(content), "snyk") {
			return true
		}
	}

	return false
}

func (c *DepsCheck) hasDependabot(path string) bool {
	return safepath.Exists(path, ".github/dependabot.yml") ||
		safepath.Exists(path, ".github/dependabot.yaml")
}

func (c *DepsCheck) hasRenovate(path string) bool {
	renovateConfigs := []string{
		"renovate.json",
		"renovate.json5",
		".renovaterc",
		".renovaterc.json",
	}
	for _, config := range renovateConfigs {
		if safepath.Exists(path, config) {
			return true
		}
	}
	return false
}

func (c *DepsCheck) hasMavenDependencyPlugin(path string) bool {
	if !safepath.Exists(path, "pom.xml") {
		return false
	}
	content, err := safepath.ReadFile(path, "pom.xml")
	if err != nil {
		return false
	}
	// Check for maven-dependency-plugin with analyze goal
	return strings.Contains(string(content), "maven-dependency-plugin")
}

func (c *DepsCheck) hasGradleDependencyVerification(path string) bool {
	// Check for dependency-verification.xml in gradle directory
	return safepath.Exists(path, "gradle/verification-metadata.xml") ||
		safepath.Exists(path, "gradle/dependency-verification.xml")
}
