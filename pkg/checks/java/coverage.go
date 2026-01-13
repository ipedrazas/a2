package javacheck

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// CoverageCheck verifies Java test coverage using JaCoCo reports.
type CoverageCheck struct {
	Threshold float64
}

func (c *CoverageCheck) ID() string   { return "java:coverage" }
func (c *CoverageCheck) Name() string { return "Java Coverage" }

// Run checks for JaCoCo coverage reports and validates against threshold.
func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangJava)

	threshold := c.Threshold
	if threshold <= 0 {
		threshold = 80.0
	}

	// Check if JaCoCo is configured
	if !c.hasJaCoCo(path) {
		return rb.Warn("JaCoCo not configured (add jacoco plugin to enable coverage)"), nil
	}

	// Try to find and parse JaCoCo reports
	coverage, found := c.findCoverage(path)
	if !found {
		return rb.Warn("JaCoCo configured but no coverage report found (run tests first)"), nil
	}

	if coverage >= threshold {
		return rb.Pass(fmt.Sprintf("Coverage: %.1f%% (threshold: %.1f%%)", coverage, threshold)), nil
	}
	return rb.Warn(fmt.Sprintf("Coverage %.1f%% below threshold %.1f%%", coverage, threshold)), nil
}

func (c *CoverageCheck) hasJaCoCo(path string) bool {
	// Check Maven pom.xml
	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(string(content), "jacoco") {
			return true
		}
	}

	// Check Gradle build files
	if safepath.Exists(path, "build.gradle") {
		content, err := safepath.ReadFile(path, "build.gradle")
		if err == nil && strings.Contains(string(content), "jacoco") {
			return true
		}
	}
	if safepath.Exists(path, "build.gradle.kts") {
		content, err := safepath.ReadFile(path, "build.gradle.kts")
		if err == nil && strings.Contains(string(content), "jacoco") {
			return true
		}
	}

	return false
}

func (c *CoverageCheck) findCoverage(path string) (float64, bool) {
	// Common JaCoCo report locations
	reportPaths := []string{
		// Maven
		"target/site/jacoco/jacoco.xml",
		"target/jacoco.xml",
		// Gradle
		"build/reports/jacoco/test/jacocoTestReport.xml",
		"build/jacoco/test.xml",
		// Multi-module (aggregate)
		"target/site/jacoco-aggregate/jacoco.xml",
		"build/reports/jacoco/jacocoAggregatedReport.xml",
	}

	for _, reportPath := range reportPaths {
		fullPath := filepath.Join(path, reportPath)
		coverage, ok := c.parseCoverageReport(fullPath)
		if ok {
			return coverage, true
		}
	}

	return 0, false
}

// JaCoCoReport represents the JaCoCo XML report structure.
type JaCoCoReport struct {
	XMLName  xml.Name        `xml:"report"`
	Counters []JaCoCoCounter `xml:"counter"`
}

// JaCoCoCounter represents a counter element in JaCoCo report.
type JaCoCoCounter struct {
	Type    string `xml:"type,attr"`
	Missed  int    `xml:"missed,attr"`
	Covered int    `xml:"covered,attr"`
}

func (c *CoverageCheck) parseCoverageReport(reportPath string) (float64, bool) {
	data, err := os.ReadFile(reportPath) // #nosec G304 - path constructed from known locations
	if err != nil {
		return 0, false
	}

	var report JaCoCoReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return 0, false
	}

	// Look for LINE or INSTRUCTION coverage (prefer LINE)
	for _, counter := range report.Counters {
		if counter.Type == "LINE" || counter.Type == "INSTRUCTION" {
			total := counter.Missed + counter.Covered
			if total > 0 {
				return float64(counter.Covered) / float64(total) * 100, true
			}
		}
	}

	return 0, false
}
