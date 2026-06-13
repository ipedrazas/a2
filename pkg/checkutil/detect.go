package checkutil

import (
	"regexp"

	"github.com/ipedrazas/a2/pkg/safepath"
)

// FirstExisting returns the first of the given relative names that exists under
// dir, or "" if none do. It standardizes the common "find the config file this
// tool actually uses" pattern (e.g. ruff.toml → .ruff.toml → pyproject.toml).
func FirstExisting(dir string, names ...string) string {
	for _, name := range names {
		if safepath.Exists(dir, name) {
			return name
		}
	}
	return ""
}

// AnyExists reports whether any of the given relative names exist under dir.
func AnyExists(dir string, names ...string) bool {
	return FirstExisting(dir, names...) != ""
}

// CountMatches returns the number of non-overlapping matches of re in s.
// It replaces the repeated `len(re.FindAllString(out, -1))` idiom used by checks
// that count findings (lint hits, unformatted files, coverage gaps) in output.
func CountMatches(re *regexp.Regexp, s string) int {
	return len(re.FindAllString(s, -1))
}

// FallbackChain returns the result of the first detector that yields a non-empty
// string. It expresses the recurring "tool → config → regex" priority fallback
// without a ladder of if/return statements; pass detectors most-authoritative
// first.
func FallbackChain(detectors ...func() string) string {
	for _, detect := range detectors {
		if v := detect(); v != "" {
			return v
		}
	}
	return ""
}
