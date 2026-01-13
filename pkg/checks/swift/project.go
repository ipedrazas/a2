package swiftcheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ProjectCheck verifies that a Swift project exists (Package.swift).
type ProjectCheck struct{}

func (c *ProjectCheck) ID() string   { return "swift:project" }
func (c *ProjectCheck) Name() string { return "Swift Project" }

// Run checks for Package.swift and extracts project information.
func (c *ProjectCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangSwift)

	// Check for Package.swift
	if !safepath.Exists(path, "Package.swift") {
		return rb.Fail("No Package.swift found"), nil
	}

	// Read Package.swift to extract package info
	content, err := safepath.ReadFile(path, "Package.swift")
	if err != nil {
		return rb.Fail("Cannot read Package.swift: " + err.Error()), nil
	}

	// Parse basic info from Package.swift
	name := extractPackageName(string(content))

	var info []string
	if name != "" {
		info = append(info, name)
	}

	// Check for resolved dependencies
	hasResolved := safepath.Exists(path, "Package.resolved")
	if hasResolved {
		info = append(info, "(dependencies resolved)")
	}

	if len(info) > 0 {
		return rb.Pass("Package: " + strings.Join(info, " ")), nil
	}
	return rb.Pass("Package.swift found"), nil
}

// extractPackageName extracts the package name from Package.swift content.
func extractPackageName(content string) string {
	// Look for: name: "PackageName"
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "name:") {
			// Extract value between quotes
			start := strings.Index(line, "\"")
			if start != -1 {
				end := strings.Index(line[start+1:], "\"")
				if end != -1 {
					return line[start+1 : start+1+end]
				}
			}
		}
	}
	return ""
}
