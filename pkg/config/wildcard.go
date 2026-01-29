package config

import (
	"strings"
)

// matchesPattern checks if a checkID matches a wildcard pattern.
// Patterns:
//   - "*:suffix"  matches anything ending with ":suffix" (e.g., "*:tests" matches "go:tests", "node:tests")
//   - "prefix:*"  matches anything starting with "prefix:" (e.g., "node:*" matches "node:deps", "node:logging")
//   - "*:*"       matches anything containing a colon
//   - "*"         matches everything
//   - exact match if no wildcards
func matchesPattern(checkID string, pattern string) bool {
	// Trim whitespace from pattern
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return false
	}

	// Exact match (no wildcards)
	if !strings.Contains(pattern, "*") {
		return checkID == pattern
	}

	// Single wildcard - match everything
	if pattern == "*" {
		return true
	}

	// Pattern with both parts wildcard: *:*
	if pattern == "*:*" {
		return strings.Contains(checkID, ":")
	}

	// Suffix wildcard: *:something
	if strings.HasPrefix(pattern, "*:") {
		suffix := strings.TrimPrefix(pattern, "*:")
		return strings.HasSuffix(checkID, ":"+suffix)
	}

	// Prefix wildcard: something:*
	if strings.HasSuffix(pattern, ":*") {
		prefix := strings.TrimSuffix(pattern, ":*")
		return strings.HasPrefix(checkID, prefix+":")
	}

	// Multiple wildcards or invalid patterns - treat as no match
	// This prevents ambiguous patterns like *:*:* from matching unintended things
	return false
}
