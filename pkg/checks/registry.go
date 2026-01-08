package checks

import (
	"sort"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks/common"
	gocheck "github.com/ipedrazas/a2/pkg/checks/go"
	javacheck "github.com/ipedrazas/a2/pkg/checks/java"
	nodecheck "github.com/ipedrazas/a2/pkg/checks/node"
	pythoncheck "github.com/ipedrazas/a2/pkg/checks/python"
	rustcheck "github.com/ipedrazas/a2/pkg/checks/rust"
	typescriptcheck "github.com/ipedrazas/a2/pkg/checks/typescript"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
)

// GetChecks returns checks for detected languages based on configuration.
// Checks are ordered with critical checks first.
func GetChecks(cfg *config.Config, detected language.DetectionResult) []checker.Checker {
	var registrations []checker.CheckRegistration

	// Get checks for each detected language
	for _, lang := range detected.Languages {
		registrations = append(registrations, getChecksForLanguage(lang, cfg)...)
	}

	// Add common checks (language-agnostic)
	registrations = append(registrations, common.Register(cfg)...)

	// Sort by order (critical checks first)
	sort.Slice(registrations, func(i, j int) bool {
		return registrations[i].Meta.Order < registrations[j].Meta.Order
	})

	// Filter disabled and extract checkers
	var enabled []checker.Checker
	for _, reg := range registrations {
		checkID := reg.Meta.ID
		if !cfg.IsCheckDisabled(checkID) {
			enabled = append(enabled, reg.Checker)
		}
	}

	return enabled
}

// getChecksForLanguage returns check registrations for a specific language.
func getChecksForLanguage(lang checker.Language, cfg *config.Config) []checker.CheckRegistration {
	switch lang {
	case checker.LangGo:
		return gocheck.Register(cfg)
	case checker.LangPython:
		return pythoncheck.Register(cfg)
	case checker.LangNode:
		return nodecheck.Register(cfg)
	case checker.LangJava:
		return javacheck.Register(cfg)
	case checker.LangRust:
		return rustcheck.Register(cfg)
	case checker.LangTypeScript:
		return typescriptcheck.Register(cfg)
	default:
		return nil
	}
}

// DefaultChecks returns checks for auto-detected languages (backward compatibility).
func DefaultChecks() []checker.Checker {
	cfg := config.DefaultConfig()
	detected := language.Detect(".")

	// Fallback to Go if nothing detected (backward compatibility)
	if len(detected.Languages) == 0 {
		detected.Languages = []checker.Language{checker.LangGo}
		detected.Primary = checker.LangGo
	}

	return GetChecks(cfg, detected)
}

// GetChecksForPath returns checks for a specific path with auto-detection.
func GetChecksForPath(path string, cfg *config.Config) ([]checker.Checker, language.DetectionResult) {
	// Detect languages or use explicit config
	var detected language.DetectionResult
	if len(cfg.Language.Explicit) > 0 {
		langs := make([]checker.Language, len(cfg.Language.Explicit))
		for i, l := range cfg.Language.Explicit {
			langs[i] = checker.Language(l)
		}
		detected = language.DetectWithOverride(path, langs)
	} else {
		detected = language.Detect(path)
	}

	// Fallback to Go if nothing detected (backward compatibility)
	if len(detected.Languages) == 0 {
		detected.Languages = []checker.Language{checker.LangGo}
		detected.Primary = checker.LangGo
	}

	return GetChecks(cfg, detected), detected
}
