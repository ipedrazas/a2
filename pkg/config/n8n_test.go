package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestN8nConfigParsing(t *testing.T) {
	dir := t.TempDir()
	yaml := `
n8n:
  code_node_line_limit: 120
  naming_regex: "^[a-z][a-z0-9 -]+$"
  error_workflow_fail_ratio: 0.75
  duplication_min_sequence: 4
  untrusted_triggers_extra:
    - "acme.customTrigger"
  write_tools_extra:
    - "acme.crmTool:update"
checks:
  n8n:
    disabled:
      - n8n:naming
`
	if err := os.WriteFile(filepath.Join(dir, ".a2.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.N8n.CodeNodeLineLimit != 120 {
		t.Errorf("CodeNodeLineLimit = %d, want 120", cfg.N8n.CodeNodeLineLimit)
	}
	if cfg.N8n.NamingRegex != "^[a-z][a-z0-9 -]+$" {
		t.Errorf("NamingRegex = %q", cfg.N8n.NamingRegex)
	}
	if cfg.N8n.ErrorWorkflowFailRatio != 0.75 {
		t.Errorf("ErrorWorkflowFailRatio = %v, want 0.75", cfg.N8n.ErrorWorkflowFailRatio)
	}
	if cfg.N8n.DuplicationMinSequence != 4 {
		t.Errorf("DuplicationMinSequence = %d, want 4", cfg.N8n.DuplicationMinSequence)
	}
	if len(cfg.N8n.UntrustedTriggersExtra) != 1 || cfg.N8n.UntrustedTriggersExtra[0] != "acme.customTrigger" {
		t.Errorf("UntrustedTriggersExtra = %v", cfg.N8n.UntrustedTriggersExtra)
	}
	if len(cfg.N8n.WriteToolsExtra) != 1 || cfg.N8n.WriteToolsExtra[0] != "acme.crmTool:update" {
		t.Errorf("WriteToolsExtra = %v", cfg.N8n.WriteToolsExtra)
	}
	if got := cfg.Checks.PerLanguage["n8n"]; len(got) != 1 || got[0] != "n8n:naming" {
		t.Errorf("per-language disabled for n8n = %v", got)
	}
}
