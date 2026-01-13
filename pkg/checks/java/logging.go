package javacheck

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LoggingCheck detects Java structured logging practices.
type LoggingCheck struct{}

func (c *LoggingCheck) ID() string   { return "java:logging" }
func (c *LoggingCheck) Name() string { return "Java Logging" }

// Run checks for structured logging libraries and anti-patterns.
func (c *LoggingCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangJava)

	var loggers []string
	var issues []string

	// Check for logging libraries in dependencies
	if c.hasSLF4J(path) {
		loggers = append(loggers, "SLF4J")
	}
	if c.hasLogback(path) {
		loggers = append(loggers, "Logback")
	}
	if c.hasLog4j2(path) {
		loggers = append(loggers, "Log4j2")
	}

	// Check for System.out.println anti-patterns
	printlnCount := c.countSystemOutPrintln(path)
	if printlnCount > 0 {
		issues = append(issues, intToStr(printlnCount)+" System.out.println statements")
	}

	// Build result
	if len(loggers) > 0 {
		msg := "Structured logging: " + strings.Join(loggers, ", ")
		if len(issues) > 0 {
			return rb.Warn(msg + " (but found " + strings.Join(issues, ", ") + ")"), nil
		}
		return rb.Pass(msg), nil
	} else if len(issues) > 0 {
		return rb.Warn("No structured logging library found; " + strings.Join(issues, ", ")), nil
	}
	return rb.Warn("No structured logging detected (consider SLF4J with Logback or Log4j2)"), nil
}

func (c *LoggingCheck) hasSLF4J(path string) bool {
	return c.hasDependency(path, "slf4j")
}

func (c *LoggingCheck) hasLogback(path string) bool {
	// Check for logback config files
	if safepath.Exists(path, "src/main/resources/logback.xml") ||
		safepath.Exists(path, "src/main/resources/logback-spring.xml") ||
		safepath.Exists(path, "logback.xml") {
		return true
	}
	return c.hasDependency(path, "logback")
}

func (c *LoggingCheck) hasLog4j2(path string) bool {
	// Check for log4j2 config files
	if safepath.Exists(path, "src/main/resources/log4j2.xml") ||
		safepath.Exists(path, "src/main/resources/log4j2.yaml") ||
		safepath.Exists(path, "src/main/resources/log4j2.properties") ||
		safepath.Exists(path, "log4j2.xml") {
		return true
	}
	return c.hasDependency(path, "log4j")
}

func (c *LoggingCheck) hasDependency(path, depName string) bool {
	// Check Maven pom.xml
	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(strings.ToLower(string(content)), depName) {
			return true
		}
	}

	// Check Gradle build files
	if safepath.Exists(path, "build.gradle") {
		content, err := safepath.ReadFile(path, "build.gradle")
		if err == nil && strings.Contains(strings.ToLower(string(content)), depName) {
			return true
		}
	}
	if safepath.Exists(path, "build.gradle.kts") {
		content, err := safepath.ReadFile(path, "build.gradle.kts")
		if err == nil && strings.Contains(strings.ToLower(string(content)), depName) {
			return true
		}
	}

	return false
}

func (c *LoggingCheck) countSystemOutPrintln(path string) int {
	count := 0
	srcDir := filepath.Join(path, "src", "main", "java")

	// Walk through source files
	_ = filepath.Walk(srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(filePath, ".java") {
			return nil
		}

		// Count System.out.println occurrences
		count += c.countPrintlnInFile(filePath)
		return nil
	})

	return count
}

func (c *LoggingCheck) countPrintlnInFile(filePath string) int {
	file, err := os.Open(filePath) // #nosec G304 - path from filepath.Walk
	if err != nil {
		return 0
	}
	defer func() { _ = file.Close() }()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check for System.out.println, System.out.print, System.err.println
		if strings.Contains(line, "System.out.print") ||
			strings.Contains(line, "System.err.print") {
			// Skip if it's a comment
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "*") {
				count++
			}
		}
	}

	return count
}
