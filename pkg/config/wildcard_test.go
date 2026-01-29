package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type wildcardTestCase struct {
	checkID     string
	pattern     string
	shouldMatch bool
	description string
}

func TestMatchesPattern(t *testing.T) {
	tests := []wildcardTestCase{
		// Prefix wildcard tests
		{"node:deps", "node:*", true, "prefix wildcard: node:* matches node:deps"},
		{"node:logging", "node:*", true, "prefix wildcard: node:* matches node:logging"},
		{"node:tests", "node:*", true, "prefix wildcard: node:* matches node:tests"},
		{"go:deps", "node:*", false, "prefix wildcard: node:* doesn't match go:deps"},
		{"python:pip", "python:*", true, "prefix wildcard: python:* matches python:pip"},

		// Suffix wildcard tests
		{"go:tests", "*:tests", true, "suffix wildcard: *:tests matches go:tests"},
		{"python:tests", "*:tests", true, "suffix wildcard: *:tests matches python:tests"},
		{"node:tests", "*:tests", true, "suffix wildcard: *:tests matches node:tests"},
		{"go:build", "*:tests", false, "suffix wildcard: *:tests doesn't match go:build"},
		{"common:license", "*:license", true, "suffix wildcard: *:license matches common:license"},

		// Double wildcard tests
		{"go:tests", "*:*", true, "double wildcard: *:* matches go:tests"},
		{"python:build", "*:*", true, "double wildcard: *:* matches python:build"},
		{"node:deps", "*:*", true, "double wildcard: *:* matches node:deps"},
		{"common:dockerfile", "*:*", true, "double wildcard: *:* matches common:dockerfile"},

		// Single wildcard - match everything
		{"go:tests", "*", true, "single wildcard: * matches go:tests"},
		{"python:build", "*", true, "single wildcard: * matches python:build"},
		{"anything", "*", true, "single wildcard: * matches anything"},

		// Exact match tests (no wildcards)
		{"go:tests", "go:tests", true, "exact match: go:tests matches go:tests"},
		{"node:deps", "node:deps", true, "exact match: node:deps matches node:deps"},
		{"go:tests", "go:build", false, "exact match: go:tests doesn't match go:build"},
		{"node:deps", "go:deps", false, "exact match: node:deps doesn't match go:deps"},

		// Empty and whitespace patterns
		{"go:tests", "", false, "empty pattern doesn't match"},
		{"go:tests", "   ", false, "whitespace pattern doesn't match"},
		{"go:tests", "  go:tests  ", true, "trimmed exact match"},

		// Edge cases - invalid/multiple wildcards should not match
		{"go:tests", "*:*:*", false, "multiple wildcards: *:*:* should not match"},
		{"go:tests", "*a:*", false, "embedded wildcard: *a:* should not match"},
		{"go:tests", "node:*:extra", false, "multiple wildcards: node:*:extra should not match"},

		// Check IDs without colons (edge case)
		{"notvalid", "*:*", false, "double wildcard: *:* doesn't match checkID without colon"},
		{"notvalid", "*", true, "single wildcard: * matches checkID without colon"},

		// Common check patterns
		{"go:logging", "*:logging", true, "suffix wildcard matches go:logging"},
		{"node:logging", "*:logging", true, "suffix wildcard matches node:logging"},
		{"typescript:logging", "*:logging", true, "suffix wildcard matches typescript:logging"},
		{"common:dockerfile", "common:*", true, "prefix wildcard matches common:dockerfile"},
		{"common:ci", "common:*", true, "prefix wildcard matches common:ci"},
		{"common:license", "common:*", true, "prefix wildcard matches common:license"},

		// Mixed language and check type patterns
		{"rust:tests", "*:tests", true, "suffix wildcard: *:tests matches rust:tests"},
		{"swift:tests", "*:tests", true, "suffix wildcard: *:tests matches swift:tests"},
		{"java:tests", "*:tests", true, "suffix wildcard: *:tests matches java:tests"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := matchesPattern(tt.checkID, tt.pattern)
			assert.Equal(t, tt.shouldMatch, result,
				"matchesPattern(%q, %q) should return %v", tt.checkID, tt.pattern, tt.shouldMatch)
		})
	}
}

func TestIsCheckDisabled_WithWildcards(t *testing.T) {
	tests := []struct {
		name          string
		disabled      []string
		checkID       string
		shouldDisable bool
	}{
		{
			name:          "prefix wildcard disables all node checks",
			disabled:      []string{"node:*"},
			checkID:       "node:deps",
			shouldDisable: true,
		},
		{
			name:          "prefix wildcard doesn't disable other language",
			disabled:      []string{"node:*"},
			checkID:       "go:deps",
			shouldDisable: false,
		},
		{
			name:          "suffix wildcard disables all tests",
			disabled:      []string{"*:tests"},
			checkID:       "go:tests",
			shouldDisable: true,
		},
		{
			name:          "suffix wildcard doesn't disable non-tests",
			disabled:      []string{"*:tests"},
			checkID:       "go:build",
			shouldDisable: false,
		},
		{
			name:          "multiple wildcards",
			disabled:      []string{"node:*", "*:tests"},
			checkID:       "node:tests",
			shouldDisable: true,
		},
		{
			name:          "wildcard with exact match",
			disabled:      []string{"*:logging", "go:build"},
			checkID:       "go:build",
			shouldDisable: true,
		},
		{
			name:          "double wildcard matches all",
			disabled:      []string{"*:*"},
			checkID:       "anything:here",
			shouldDisable: true,
		},
		{
			name:          "empty disabled list",
			disabled:      []string{},
			checkID:       "go:tests",
			shouldDisable: false,
		},
		{
			name:          "wildcard with alias still works",
			disabled:      []string{"*:*"},
			checkID:       "tests",
			shouldDisable: false, // "tests" doesn't have a colon, so *:* won't match
		},
		{
			name:          "alias match with wildcard present",
			disabled:      []string{"go:*", "build"},
			checkID:       "go:build",
			shouldDisable: true, // matches via go:*
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			cfg.Checks.Disabled = tt.disabled
			result := cfg.IsCheckDisabled(tt.checkID)
			assert.Equal(t, tt.shouldDisable, result,
				"IsCheckDisabled(%q) with disabled %v should return %v", tt.checkID, tt.disabled, tt.shouldDisable)
		})
	}
}

func TestIsCheckDisabled_WithAliases(t *testing.T) {
	// Test that wildcard matching works alongside existing alias functionality
	tests := []struct {
		name          string
		disabled      []string
		checkID       string
		shouldDisable bool
	}{
		{
			name:          "wildcard disables prefixed check",
			disabled:      []string{"go:*"},
			checkID:       "go:tests",
			shouldDisable: true,
		},
		{
			name:          "alias still works with wildcard in list",
			disabled:      []string{"go:*", "coverage"},
			checkID:       "go:coverage",
			shouldDisable: true, // matches via go:*
		},
		{
			name:          "alias maps to disabled check with wildcard",
			disabled:      []string{"go:*"},
			checkID:       "tests", // alias for go:tests
			shouldDisable: false,   // alias check happens separately, tests != go:tests directly
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			cfg.Checks.Disabled = tt.disabled
			result := cfg.IsCheckDisabled(tt.checkID)
			assert.Equal(t, tt.shouldDisable, result,
				"IsCheckDisabled(%q) with disabled %v should return %v", tt.checkID, tt.disabled, tt.shouldDisable)
		})
	}
}

// Benchmark for wildcard matching performance
func BenchmarkMatchesPattern(b *testing.B) {
	benchmarks := []struct {
		name    string
		pattern string
		checkID string
	}{
		{"exact_match", "go:tests", "go:tests"},
		{"prefix_wildcard", "node:*", "node:deps"},
		{"suffix_wildcard", "*:tests", "go:tests"},
		{"double_wildcard", "*:*", "go:tests"},
		{"single_wildcard", "*", "go:tests"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matchesPattern(bm.checkID, bm.pattern)
			}
		})
	}
}
