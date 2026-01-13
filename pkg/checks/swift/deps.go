package swiftcheck

import (
	"encoding/json"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsCheck checks for dependency issues in Package.resolved.
type DepsCheck struct{}

func (c *DepsCheck) ID() string   { return "swift:deps" }
func (c *DepsCheck) Name() string { return "Swift Vulnerabilities" }

// Run checks Package.resolved for outdated or potentially vulnerable dependencies.
func (c *DepsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangSwift)

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		return rb.Fail("No Package.swift found"), nil
	}

	// Check for Package.resolved
	if !safepath.Exists(path, "Package.resolved") {
		return rb.Info("No Package.resolved (run 'swift package resolve')"), nil
	}

	// Read and parse Package.resolved
	data, err := safepath.ReadFile(path, "Package.resolved")
	if err != nil {
		return rb.Warn("Cannot read Package.resolved: " + err.Error()), nil
	}

	deps, err := parsePackageResolved(data)
	if err != nil {
		return rb.Warn("Cannot parse Package.resolved: " + err.Error()), nil
	}

	// Note: Swift doesn't have a built-in vulnerability database like npm or cargo
	// We can only report dependency count and check for pinning issues
	if len(deps) == 0 {
		return rb.Pass("No dependencies"), nil
	}

	return rb.Pass(formatDepsMessage(len(deps))), nil
}

// packageResolved represents the Package.resolved file structure.
type packageResolved struct {
	Version int `json:"version"`
	// Version 1 format
	Object *struct {
		Pins []struct {
			Package string `json:"package"`
			State   struct {
				Version  string `json:"version"`
				Revision string `json:"revision"`
			} `json:"state"`
		} `json:"pins"`
	} `json:"object"`
	// Version 2 format
	Pins []struct {
		Identity string `json:"identity"`
		State    struct {
			Version  string `json:"version"`
			Revision string `json:"revision"`
		} `json:"state"`
	} `json:"pins"`
}

// parsePackageResolved parses Package.resolved and returns dependency names.
func parsePackageResolved(data []byte) ([]string, error) {
	var resolved packageResolved
	if err := json.Unmarshal(data, &resolved); err != nil {
		return nil, err
	}

	var deps []string
	if resolved.Version == 1 && resolved.Object != nil {
		for _, pin := range resolved.Object.Pins {
			deps = append(deps, pin.Package)
		}
	} else {
		// Version 2 or newer
		for _, pin := range resolved.Pins {
			deps = append(deps, pin.Identity)
		}
	}

	return deps, nil
}

// formatDepsMessage creates a message for the dependency count.
func formatDepsMessage(count int) string {
	var sb strings.Builder
	sb.WriteString("Dependencies: ")
	if count == 1 {
		sb.WriteString("1 dependency")
	} else {
		if count < 10 {
			sb.WriteByte('0' + byte(count))
		} else {
			sb.WriteString("multiple")
		}
		sb.WriteString(" dependencies")
	}
	sb.WriteString(" (no vulnerability database available)")
	return sb.String()
}
