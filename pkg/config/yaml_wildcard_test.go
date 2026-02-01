package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestYAML_DisabledWithWildcard verifies that wildcard patterns in disabled list
// must be quoted in YAML because * is a YAML alias character.
func TestYAML_DisabledWithWildcard(t *testing.T) {
	// Unquoted *:logging - YAML interprets * as alias reference and fails or misparses
	unquoted := `
checks:
  disabled:
    - *:logging
`
	var cfgUnquoted struct {
		Checks struct {
			Disabled []string `yaml:"disabled"`
		} `yaml:"checks"`
	}
	errUnquoted := yaml.Unmarshal([]byte(unquoted), &cfgUnquoted)
	if errUnquoted == nil {
		t.Logf("Unquoted: disabled=%q (may be wrong if alias was defined)", cfgUnquoted.Checks.Disabled)
	} else {
		t.Logf("Unquoted *:logging causes YAML error (expected): %v", errUnquoted)
	}

	// Quoted - works correctly
	quoted := `
checks:
  disabled:
    - "*:logging"
`
	var cfgQuoted struct {
		Checks struct {
			Disabled []string `yaml:"disabled"`
		} `yaml:"checks"`
	}
	errQuoted := yaml.Unmarshal([]byte(quoted), &cfgQuoted)
	if errQuoted != nil {
		t.Fatalf("Quoted wildcard should parse: %v", errQuoted)
	}
	if len(cfgQuoted.Checks.Disabled) != 1 || cfgQuoted.Checks.Disabled[0] != "*:logging" {
		t.Fatalf("Quoted disabled: got %q", cfgQuoted.Checks.Disabled)
	}
	// Verify IsCheckDisabled works with the loaded config
	cfg := &Config{Checks: ChecksConfig{Disabled: cfgQuoted.Checks.Disabled}}
	if !cfg.IsCheckDisabled("go:logging") {
		t.Error("IsCheckDisabled(go:logging) with *:logging should be true")
	}
	if !cfg.IsCheckDisabled("node:logging") {
		t.Error("IsCheckDisabled(node:logging) with *:logging should be true")
	}
}

// TestLoad_UnquotedWildcardReturnsHint verifies that Load() returns a helpful
// error when .a2.yaml contains unquoted * in checks.disabled.
func TestLoad_UnquotedWildcardReturnsHint(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".a2.yaml")
	// Unquoted *:logging causes YAML parse error
	content := `
checks:
  disabled:
    - *:logging
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load should fail for unquoted *:logging")
	}
	if !strings.Contains(err.Error(), "wildcard") || !strings.Contains(err.Error(), `"*:logging"`) {
		t.Errorf("error should suggest quoting wildcards, got: %v", err)
	}
}
