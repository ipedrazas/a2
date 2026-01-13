package javacheck

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// FormatCheck detects Java code formatting configuration.
type FormatCheck struct{}

func (c *FormatCheck) ID() string   { return "java:format" }
func (c *FormatCheck) Name() string { return "Java Format" }

// Run checks for formatting tools configuration.
func (c *FormatCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangJava)

	var formatters []string

	// Check for google-java-format
	if c.hasGoogleJavaFormat(path) {
		formatters = append(formatters, "google-java-format")
	}

	// Check for Spotless plugin
	if c.hasSpotless(path) {
		formatters = append(formatters, "Spotless")
	}

	// Check for .editorconfig with Java settings
	if c.hasEditorConfig(path) {
		formatters = append(formatters, "EditorConfig")
	}

	// Check for IDE formatters
	if safepath.Exists(path, ".idea") {
		if c.hasIntelliJFormatter(path) {
			formatters = append(formatters, "IntelliJ")
		}
	}

	// Check for Eclipse formatter
	if safepath.Exists(path, ".settings") || safepath.Exists(path, "eclipse-formatter.xml") {
		formatters = append(formatters, "Eclipse")
	}

	if len(formatters) > 0 {
		return rb.Pass("Formatting configured: " + strings.Join(formatters, ", ")), nil
	}
	return rb.Warn("No formatter configuration found (consider Spotless or google-java-format)"), nil
}

func (c *FormatCheck) hasGoogleJavaFormat(path string) bool {
	// Check for google-java-format in Maven or Gradle config
	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(string(content), "google-java-format") {
			return true
		}
	}
	if safepath.Exists(path, "build.gradle") {
		content, err := safepath.ReadFile(path, "build.gradle")
		if err == nil && strings.Contains(string(content), "google-java-format") {
			return true
		}
	}
	if safepath.Exists(path, "build.gradle.kts") {
		content, err := safepath.ReadFile(path, "build.gradle.kts")
		if err == nil && strings.Contains(string(content), "google-java-format") {
			return true
		}
	}
	return false
}

func (c *FormatCheck) hasSpotless(path string) bool {
	// Check for Spotless plugin in Maven or Gradle
	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(string(content), "spotless") {
			return true
		}
	}
	if safepath.Exists(path, "build.gradle") {
		content, err := safepath.ReadFile(path, "build.gradle")
		if err == nil && strings.Contains(string(content), "spotless") {
			return true
		}
	}
	if safepath.Exists(path, "build.gradle.kts") {
		content, err := safepath.ReadFile(path, "build.gradle.kts")
		if err == nil && strings.Contains(string(content), "spotless") {
			return true
		}
	}
	return false
}

func (c *FormatCheck) hasEditorConfig(path string) bool {
	if !safepath.Exists(path, ".editorconfig") {
		return false
	}
	content, err := safepath.ReadFile(path, ".editorconfig")
	if err != nil {
		return false
	}
	// Check for Java-specific settings
	return strings.Contains(string(content), "*.java") ||
		strings.Contains(string(content), "[*.java]")
}

func (c *FormatCheck) hasIntelliJFormatter(path string) bool {
	ideaDir := filepath.Join(path, ".idea")
	codeStyleDir := filepath.Join(ideaDir, "codeStyles")

	// Check for code style configuration
	if info, err := os.Stat(codeStyleDir); err == nil && info.IsDir() {
		entries, err := os.ReadDir(codeStyleDir)
		if err == nil && len(entries) > 0 {
			return true
		}
	}

	// Check for codeStyleSettings.xml
	codeStyleFile := filepath.Join(ideaDir, "codeStyleSettings.xml")
	if _, err := os.Stat(codeStyleFile); err == nil {
		return true
	}

	return false
}
