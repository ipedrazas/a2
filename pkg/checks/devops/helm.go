package devops

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// HelmCheck verifies that Helm charts exist and validates them.
type HelmCheck struct{}

func (c *HelmCheck) ID() string   { return "devops:helm" }
func (c *HelmCheck) Name() string { return "Helm Charts" }

// Run checks for Helm charts and validates them if helm is installed.
func (c *HelmCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Check if helm charts exist
	charts := c.findHelmCharts(path)
	if len(charts) == 0 {
		return rb.Info("No Helm charts found"), nil
	}

	// Check if helm tool is installed
	if !checkutil.ToolAvailable("helm") {
		return rb.ToolNotInstalled("helm", "install from https://helm.sh"), nil
	}

	// Run helm lint on each chart
	var warnings []string
	var passed []string
	for _, chart := range charts {
		result := checkutil.RunCommand(filepath.Dir(chart), "helm", "lint", filepath.Dir(chart))
		if result.Success() {
			passed = append(passed, chart)
		} else {
			warnings = append(warnings, chart+": "+result.Output())
		}
	}

	if len(warnings) == 0 {
		if len(passed) == 1 {
			return rb.Pass("helm lint passed for chart"), nil
		}
		return rb.Pass("helm lint passed for " + checkutil.PluralizeCount(len(passed), "chart", "charts")), nil
	}

	msg := "helm lint found issues in " + checkutil.PluralizeCount(len(warnings), "chart", "charts")
	if len(passed) > 0 {
		msg += " (" + checkutil.PluralizeCount(len(passed), "chart", "charts") + " passed)"
	}
	return rb.Warn(msg + "\n" + strings.Join(warnings, "\n")), nil
}

// findHelmCharts searches for Helm chart directories.
func (c *HelmCheck) findHelmCharts(path string) []string {
	var charts []string

	// Check root for Chart.yaml
	if safepath.Exists(path, "Chart.yaml") || safepath.Exists(path, "Chart.yml") {
		charts = append(charts, path)
	}

	// Check charts/ directory
	chartsDir := filepath.Join(path, "charts")
	if info, err := os.Stat(chartsDir); err == nil && info.IsDir() {
		entries, err := os.ReadDir(chartsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					chartFile := filepath.Join(chartsDir, entry.Name(), "Chart.yaml")
					if _, err := os.Stat(chartFile); err == nil {
						charts = append(charts, filepath.Join(chartsDir, entry.Name()))
					}
					// Also check for Chart.yml
					chartFileYml := filepath.Join(chartsDir, entry.Name(), "Chart.yml")
					if _, err := os.Stat(chartFileYml); err == nil {
						charts = append(charts, filepath.Join(chartsDir, entry.Name()))
					}
				}
			}
		}
	}

	// Check helm/ directory
	helmDir := filepath.Join(path, "helm")
	if info, err := os.Stat(helmDir); err == nil && info.IsDir() {
		// Check for Chart.yaml in helm/ root
		if safepath.Exists(helmDir, "Chart.yaml") || safepath.Exists(helmDir, "Chart.yml") {
			charts = append(charts, helmDir)
		}
		// Check subdirectories
		entries, err := os.ReadDir(helmDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					chartFile := filepath.Join(helmDir, entry.Name(), "Chart.yaml")
					if _, err := os.Stat(chartFile); err == nil {
						charts = append(charts, filepath.Join(helmDir, entry.Name()))
					}
					chartFileYml := filepath.Join(helmDir, entry.Name(), "Chart.yml")
					if _, err := os.Stat(chartFileYml); err == nil {
						charts = append(charts, filepath.Join(helmDir, entry.Name()))
					}
				}
			}
		}
	}

	// Recursively check for Chart.yaml in subdirectories
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking on error
		}
		if info.IsDir() {
			// Skip hidden directories and common non-helm dirs
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == "dist" || info.Name() == "build" {
				return filepath.SkipDir
			}
			// Skip charts and helm directories as we already checked them
			if info.Name() == "charts" || info.Name() == "helm" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check for Chart.yaml files
		if info.Name() == "Chart.yaml" || info.Name() == "Chart.yml" {
			chartDir := filepath.Dir(filePath)
			// Only add if not already added
			alreadyAdded := false
			for _, c := range charts {
				if c == chartDir {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded && chartDir != path {
				charts = append(charts, chartDir)
			}
		}

		return nil
	})

	if err != nil {
		return charts
	}

	return charts
}
