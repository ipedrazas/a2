package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DuplicationCheck detects code duplication using jscpd.
type DuplicationCheck struct{}

func (c *DuplicationCheck) ID() string   { return "common:duplication" }
func (c *DuplicationCheck) Name() string { return "Code Duplication" }

// jscpdTotalStats represents the total statistics in jscpd JSON output.
type jscpdTotalStats struct {
	Clones          int     `json:"clones"`
	DuplicatedLines int     `json:"duplicatedLines"`
	Sources         int     `json:"sources"`
	Percentage      float64 `json:"percentage"`
}

// jscpdReport represents the jscpd JSON report file structure.
type jscpdReport struct {
	Statistics struct {
		Total jscpdTotalStats `json:"total"`
	} `json:"statistics"`
}

// Run checks for code duplication using jscpd or config file detection.
func (c *DuplicationCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Priority 1: Run jscpd if installed
	if checkutil.ToolAvailable("jscpd") {
		return c.runJscpd(path, rb)
	}

	// Priority 2: Check if duplication tooling is configured
	configured := c.findConfiguredTools(path)
	if len(configured) > 0 {
		return rb.Pass("Duplication detection configured: " + strings.Join(configured, ", ")), nil
	}

	// No tool available
	return rb.ToolNotInstalled("jscpd", "npm install -g jscpd"), nil
}

// runJscpd executes jscpd and parses the JSON report file.
func (c *DuplicationCheck) runJscpd(path string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	// jscpd requires a writable output directory for the JSON reporter
	tmpDir, err := os.MkdirTemp("", "a2-jscpd-*")
	if err != nil {
		return rb.Warn("Failed to create temp dir for jscpd output"), nil
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	args := []string{"--reporters", "json", "--output", tmpDir, "--silent", "."}

	result := checkutil.RunCommand(path, "jscpd", args...)
	output := result.CombinedOutput()

	// Read the JSON report file that jscpd writes, scoped to tmpDir
	reportData, err := safepath.ReadFile(tmpDir, "jscpd-report.json")
	if err != nil {
		if !result.Success() {
			return rb.WarnWithOutput("jscpd failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
		}
		return rb.Pass("No code duplication detected"), nil
	}

	stats, err := c.parseReport(reportData)
	if err != nil {
		return rb.WarnWithOutput("Failed to parse jscpd report: "+err.Error(), output), nil
	}

	if stats.Clones == 0 {
		return rb.PassWithOutput("No code duplication detected", output), nil
	}

	msg := fmt.Sprintf("%d %s found (%.1f%% duplicated lines)",
		stats.Clones,
		checkutil.Pluralize(stats.Clones, "clone", "clones"),
		stats.Percentage,
	)

	if stats.Percentage > 10 {
		return rb.WarnWithOutput(msg, output), nil
	}

	return rb.PassWithOutput(msg, output), nil
}

// parseReport parses the jscpd JSON report file contents.
func (c *DuplicationCheck) parseReport(data []byte) (*jscpdTotalStats, error) {
	var report jscpdReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}
	return &report.Statistics.Total, nil
}

// findConfiguredTools checks for duplication tool configuration files.
func (c *DuplicationCheck) findConfiguredTools(path string) []string {
	var found []string

	configs := []struct {
		name string
		file string
	}{
		{"jscpd", ".jscpd.json"},
		{"jscpd", ".jscpd.yaml"},
		{"jscpd", ".jscpd.yml"},
		{"PMD CPD", ".pmd"},
		{"PMD CPD", "pmd-ruleset.xml"},
		{"SonarQube", "sonar-project.properties"},
	}

	for _, cfg := range configs {
		if safepath.Exists(path, cfg.file) {
			found = append(found, cfg.name)
		}
	}

	return uniqueStrings(found)
}
