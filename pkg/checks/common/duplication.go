package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DuplicationCheck detects code duplication using jscpd.
type DuplicationCheck struct{}

func (c *DuplicationCheck) ID() string   { return "common:duplication" }
func (c *DuplicationCheck) Name() string { return "Code Duplication" }

// jscpdStatistics represents the statistics section of jscpd JSON output.
type jscpdStatistics struct {
	Clones     int     `json:"clones"`
	Duplicates int     `json:"duplicatedLines"`
	Sources    int     `json:"sources"`
	Percentage float64 `json:"percentage"`
}

// jscpdOutput represents the top-level jscpd JSON output.
type jscpdOutput struct {
	Statistics jscpdStatistics `json:"statistics"`
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

// runJscpd executes jscpd and parses the JSON output.
func (c *DuplicationCheck) runJscpd(path string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	args := []string{"--reporters", "json", "--output", "/dev/null", "--silent", "."}

	result := checkutil.RunCommand(path, "jscpd", args...)
	output := result.CombinedOutput()

	// Try to parse JSON from stdout
	stats, err := c.parseOutput(result.Stdout)
	if err != nil {
		// jscpd may exit non-zero when duplicates are found but still produce output
		stats, err = c.parseOutput(output)
		if err != nil {
			if !result.Success() {
				return rb.WarnWithOutput("jscpd failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
			}
			return rb.Pass("No code duplication detected"), nil
		}
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

// parseOutput attempts to parse jscpd JSON output.
func (c *DuplicationCheck) parseOutput(output string) (*jscpdStatistics, error) {
	if output == "" {
		return nil, fmt.Errorf("empty output")
	}

	var result jscpdOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, err
	}

	return &result.Statistics, nil
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
