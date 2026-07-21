package checks

import (
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
)

// expectedN8nChecks are the Phase 2 n8n check IDs.
var expectedN8nChecks = []string{
	"n8n:error-workflow",
	"n8n:silent-failures",
	"n8n:retry-config",
	"n8n:timeouts",
	"n8n:code-node-bulk",
	"n8n:pinned-data",
	"n8n:naming",
	"n8n:inline-auth",
	"n8n:in-git",
}

func TestN8nChecksInAllRegistrations(t *testing.T) {
	regs := GetAllCheckRegistrations(config.DefaultConfig())
	byID := map[string]checker.CheckRegistration{}
	for _, reg := range regs {
		byID[reg.Meta.ID] = reg
	}
	for _, id := range expectedN8nChecks {
		reg, ok := byID[id]
		if !ok {
			t.Errorf("%s missing from GetAllCheckRegistrations (a2 list checks)", id)
			continue
		}
		if reg.Meta.Detail == "" || reg.Meta.FixPrompt == "" || reg.Meta.Description == "" {
			t.Errorf("%s is missing explain metadata (Description/Detail/FixPrompt)", id)
		}
		if reg.Meta.Speed != checker.SpeedFast {
			t.Errorf("%s must be SpeedFast (static analysis only)", id)
		}
	}
}

func TestGetChecksForDetectedN8n(t *testing.T) {
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangN8n},
		Primary:   checker.LangN8n,
	}
	regs := GetChecks(config.DefaultConfig(), detected)
	got := map[string]bool{}
	for _, reg := range regs {
		got[reg.Meta.ID] = true
	}
	for _, id := range expectedN8nChecks {
		if !got[id] {
			t.Errorf("%s not returned for a detected n8n domain", id)
		}
	}
}
