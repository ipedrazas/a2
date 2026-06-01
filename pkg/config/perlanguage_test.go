package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestChecksConfig_ParsesPerLanguageDisabled verifies that language-keyed
// disabled blocks under `checks:` are parsed alongside the global list.
func TestChecksConfig_ParsesPerLanguageDisabled(t *testing.T) {
	content := `
checks:
  disabled:
    - go:logging
    - "devops:*"
  typescript:
    disabled:
      - common:metrics
  go:
    disabled:
      - common:tracing
`
	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got := cfg.Checks.Disabled; len(got) != 2 || got[0] != "go:logging" || got[1] != "devops:*" {
		t.Fatalf("global disabled = %q", got)
	}
	if got := cfg.Checks.PerLanguage["typescript"]; len(got) != 1 || got[0] != "common:metrics" {
		t.Fatalf("typescript disabled = %q", got)
	}
	if got := cfg.Checks.PerLanguage["go"]; len(got) != 1 || got[0] != "common:tracing" {
		t.Fatalf("go disabled = %q", got)
	}
}

// TestChecksConfig_GlobalOnlyStillParses verifies backward compatibility: a
// plain `checks.disabled` list with no language blocks parses as before.
func TestChecksConfig_GlobalOnlyStillParses(t *testing.T) {
	content := `
checks:
  disabled:
    - go:logging
    - common:health
`
	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(cfg.Checks.Disabled) != 2 {
		t.Fatalf("disabled = %q", cfg.Checks.Disabled)
	}
	if len(cfg.Checks.PerLanguage) != 0 {
		t.Fatalf("expected no per-language entries, got %v", cfg.Checks.PerLanguage)
	}
}

// TestChecksConfig_UnknownLanguageKeyErrors verifies that a typo'd language key
// under `checks:` produces a helpful error.
func TestChecksConfig_UnknownLanguageKeyErrors(t *testing.T) {
	content := `
checks:
  typscript:
    disabled:
      - common:metrics
`
	var cfg Config
	err := yaml.Unmarshal([]byte(content), &cfg)
	if err == nil {
		t.Fatal("expected error for unknown language key")
	}
}

// TestLoad_PerLanguageRoundTrip verifies Load reads per-language blocks from a
// real .a2.yaml on disk.
func TestLoad_PerLanguageRoundTrip(t *testing.T) {
	dir := t.TempDir()
	content := `
checks:
  disabled:
    - go:logging
  typescript:
    disabled:
      - common:metrics
`
	if err := os.WriteFile(filepath.Join(dir, ".a2.yaml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.IsCheckDisabled("go:logging") {
		t.Error("go:logging should be globally disabled")
	}
	if cfg.IsCheckDisabled("common:metrics") {
		t.Error("common:metrics should NOT be globally disabled")
	}
	if !cfg.IsCheckDisabledForLang("common:metrics", "typescript") {
		t.Error("common:metrics should be disabled for typescript")
	}
	if cfg.IsCheckDisabledForLang("common:metrics", "go") {
		t.Error("common:metrics should NOT be disabled for go")
	}
}

// TestIsCheckDisabledForLang_Combines verifies the per-language helper combines
// global and language-specific lists and honours wildcards/aliases.
func TestIsCheckDisabledForLang_Combines(t *testing.T) {
	cfg := &Config{Checks: ChecksConfig{
		Disabled: []string{"go:logging"},
		PerLanguage: map[string][]string{
			"typescript": {"common:metrics", "*:tracing"},
		},
	}}

	// Global applies regardless of language.
	if !cfg.IsCheckDisabledForLang("go:logging", "go") {
		t.Error("global go:logging should apply to go")
	}
	// Per-language only applies to its language.
	if !cfg.IsCheckDisabledForLang("common:metrics", "typescript") {
		t.Error("common:metrics disabled for typescript")
	}
	if cfg.IsCheckDisabledForLang("common:metrics", "go") {
		t.Error("common:metrics not disabled for go")
	}
	// Wildcard inside a per-language list.
	if !cfg.IsCheckDisabledForLang("common:tracing", "typescript") {
		t.Error("*:tracing should match common:tracing for typescript")
	}
	// IsCheckDisabledForOnlyLang ignores the global list.
	if cfg.IsCheckDisabledForOnlyLang("go:logging", "typescript") {
		t.Error("go:logging is global, not a typescript-only disable")
	}
}
