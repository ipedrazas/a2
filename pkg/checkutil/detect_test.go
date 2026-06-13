package checkutil

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstExisting(t *testing.T) {
	dir := t.TempDir()
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("x"), 0o600))

	// Returns the first candidate that exists, in priority order.
	assert.Equal(t, "pyproject.toml", FirstExisting(dir, "ruff.toml", ".ruff.toml", "pyproject.toml"))
	// Returns "" when none exist.
	assert.Equal(t, "", FirstExisting(dir, "missing.toml", "also-missing.toml"))
	// Earlier candidate wins.
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "ruff.toml"), []byte("x"), 0o600))
	assert.Equal(t, "ruff.toml", FirstExisting(dir, "ruff.toml", "pyproject.toml"))
}

func TestAnyExists(t *testing.T) {
	dir := t.TempDir()
	assert.False(t, AnyExists(dir, "a", "b"))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "b"), []byte("x"), 0o600))
	assert.True(t, AnyExists(dir, "a", "b"))
}

func TestCountMatches(t *testing.T) {
	re := regexp.MustCompile(`would reformat`)
	assert.Equal(t, 0, CountMatches(re, "all good"))
	assert.Equal(t, 2, CountMatches(re, "would reformat a\nwould reformat b\n"))
}

func TestFallbackChain(t *testing.T) {
	// First non-empty detector wins.
	got := FallbackChain(
		func() string { return "" },
		func() string { return "config" },
		func() string { return "regex" },
	)
	assert.Equal(t, "config", got)

	// All empty yields "".
	assert.Equal(t, "", FallbackChain(
		func() string { return "" },
		func() string { return "" },
	))
}
