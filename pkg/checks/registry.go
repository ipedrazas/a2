package checks

import (
	"path/filepath"
	"sort"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks/common"
	gocheck "github.com/ipedrazas/a2/pkg/checks/go"
	javacheck "github.com/ipedrazas/a2/pkg/checks/java"
	nodecheck "github.com/ipedrazas/a2/pkg/checks/node"
	pythoncheck "github.com/ipedrazas/a2/pkg/checks/python"
	rustcheck "github.com/ipedrazas/a2/pkg/checks/rust"
	swiftcheck "github.com/ipedrazas/a2/pkg/checks/swift"
	typescriptcheck "github.com/ipedrazas/a2/pkg/checks/typescript"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
)

// pathResolvingChecker wraps a checker to resolve source directories per language.
type pathResolvingChecker struct {
	checker   checker.Checker
	sourceDir string // Subdirectory to use (empty means use root)
}

func (p *pathResolvingChecker) ID() string {
	return p.checker.ID()
}

func (p *pathResolvingChecker) Name() string {
	return p.checker.Name()
}

func (p *pathResolvingChecker) Run(path string) (checker.Result, error) {
	// Resolve the actual path to use
	actualPath := path
	if p.sourceDir != "" {
		actualPath = filepath.Join(path, p.sourceDir)
	}
	return p.checker.Run(actualPath)
}

// GetChecks returns checks for detected languages based on configuration.
// Checks are ordered with critical checks first.
// Language-specific checks are wrapped to use the configured source_dir.
func GetChecks(cfg *config.Config, detected language.DetectionResult) []checker.Checker {
	var registrations []checker.CheckRegistration

	// Get checks for each detected language
	for _, lang := range detected.Languages {
		regs := getChecksForLanguage(lang, cfg)
		// Get source directory for this language
		sourceDir := cfg.GetSourceDir(string(lang))
		// If source_dir is configured, wrap checks to use it
		if sourceDir != "" {
			for i := range regs {
				regs[i].Checker = &pathResolvingChecker{
					checker:   regs[i].Checker,
					sourceDir: sourceDir,
				}
			}
		}
		registrations = append(registrations, regs...)
	}

	// Add common checks (language-agnostic, always use root path)
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
	case checker.LangSwift:
		return swiftcheck.Register(cfg)
	default:
		return nil
	}
}

// DefaultChecks returns checks for auto-detected languages.
// Returns empty slice if no language is detected.
func DefaultChecks() []checker.Checker {
	cfg := config.DefaultConfig()
	detected := language.Detect(".")
	return GetChecks(cfg, detected)
}

// GetChecksForPath returns checks for a specific path with auto-detection.
// Returns empty checks slice if no language is detected.
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
		// Auto-detect languages, checking configured source directories
		detected = language.DetectWithSourceDirs(path, cfg.GetSourceDirs())
	}

	return GetChecks(cfg, detected), detected
}
