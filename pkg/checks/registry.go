package checks

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks/common"
	devopscheck "github.com/ipedrazas/a2/pkg/checks/devops"
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

// multiPathChecker runs a common check on the repo root and each configured
// source_dir so monorepos pass when language-specific code lives in subdirs
// (e.g. common:shutdown finding Go signal handling in backend/).
type multiPathChecker struct {
	checker checker.Checker
	cfg     *config.Config
}

func (m *multiPathChecker) ID() string   { return m.checker.ID() }
func (m *multiPathChecker) Name() string { return m.checker.Name() }

func (m *multiPathChecker) Run(path string) (checker.Result, error) {
	paths := m.pathsToCheck(path)
	var lastResult checker.Result
	var lastErr error
	for _, p := range paths {
		res, err := m.checker.Run(p)
		if err != nil {
			lastErr = err
			lastResult = res
			continue
		}
		if res.Passed {
			return res, nil
		}
		lastResult = res
	}
	if lastErr != nil {
		return lastResult, lastErr
	}
	return lastResult, nil
}

// pathsToCheck returns the repo root and each configured source_dir that
// exists as a directory (deduplicated).
func (m *multiPathChecker) pathsToCheck(path string) []string {
	seen := map[string]bool{path: true}
	paths := []string{path}
	for _, dir := range m.cfg.GetSourceDirs() {
		full := filepath.Join(path, dir)
		if seen[full] {
			continue
		}
		info, err := os.Stat(full)
		if err != nil || !info.IsDir() {
			continue
		}
		seen[full] = true
		paths = append(paths, full)
	}
	return paths
}

// commonChecksUsingSourceDirs are common check IDs that look at language-specific
// code and should run on root and each configured source_dir in monorepos.
var commonChecksUsingSourceDirs = map[string]bool{
	"common:migrations": true,
	"common:shutdown":   true,
	"common:dockerfile": true,
}

// wrapCommonRegistrations wraps common checks that use source_dir with multiPathChecker.
func wrapCommonRegistrations(regs []checker.CheckRegistration, cfg *config.Config) []checker.CheckRegistration {
	for i := range regs {
		if commonChecksUsingSourceDirs[regs[i].Meta.ID] {
			regs[i].Checker = &multiPathChecker{checker: regs[i].Checker, cfg: cfg}
		}
	}
	return regs
}

// GetChecks returns checks for detected languages based on configuration.
// Checks are ordered with critical checks first.
// Language-specific checks are wrapped to use the configured source_dir.
func GetChecks(cfg *config.Config, detected language.DetectionResult) []checker.CheckRegistration {
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

	// Add devops checks (language-agnostic, root path)
	registrations = append(registrations, devopscheck.Register(cfg)...)
	// Add common checks (language-agnostic, root path; some run on source_dirs too)
	registrations = append(registrations, wrapCommonRegistrations(common.Register(cfg), cfg)...)

	// Sort by order (critical checks first)
	sort.Slice(registrations, func(i, j int) bool {
		return registrations[i].Meta.Order < registrations[j].Meta.Order
	})

	// Filter disabled checks
	var enabled []checker.CheckRegistration
	for _, reg := range registrations {
		checkID := reg.Meta.ID
		if !cfg.IsCheckDisabled(checkID) {
			enabled = append(enabled, reg)
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

// defaultChecks returns checks for auto-detected languages.
// Returns empty slice if no language is detected.
func defaultChecks() []checker.CheckRegistration {
	cfg := config.DefaultConfig()
	detected := language.Detect(".")
	return GetChecks(cfg, detected)
}

// GetChecksForPath returns checks for a specific path with auto-detection.
// Returns empty checks slice if no language is detected.
func GetChecksForPath(path string, cfg *config.Config) ([]checker.CheckRegistration, language.DetectionResult) {
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

// GetSuggestions returns a map of check ID to suggestion string.
// This aggregates suggestions from all language check registrations.
func GetSuggestions(cfg *config.Config) map[string]string {
	suggestions := make(map[string]string)

	// Collect from all languages
	allRegs := GetAllCheckRegistrations(cfg)

	for _, reg := range allRegs {
		if reg.Meta.Suggestion != "" {
			suggestions[reg.Meta.ID] = reg.Meta.Suggestion
		}
	}

	return suggestions
}

// GetAllCheckRegistrations returns all check registrations from all languages.
// This is useful for listing available checks. Common checks that use source_dir
// are wrapped so single-check runs (e.g. a2 run common:shutdown) also use multi-path.
func GetAllCheckRegistrations(cfg *config.Config) []checker.CheckRegistration {
	allRegs := []checker.CheckRegistration{}
	allRegs = append(allRegs, gocheck.Register(cfg)...)
	allRegs = append(allRegs, pythoncheck.Register(cfg)...)
	allRegs = append(allRegs, nodecheck.Register(cfg)...)
	allRegs = append(allRegs, javacheck.Register(cfg)...)
	allRegs = append(allRegs, rustcheck.Register(cfg)...)
	allRegs = append(allRegs, typescriptcheck.Register(cfg)...)
	allRegs = append(allRegs, swiftcheck.Register(cfg)...)
	allRegs = append(allRegs, devopscheck.Register(cfg)...)
	allRegs = append(allRegs, wrapCommonRegistrations(common.Register(cfg), cfg)...)

	// Sort by order
	sort.Slice(allRegs, func(i, j int) bool {
		return allRegs[i].Meta.Order < allRegs[j].Meta.Order
	})

	return allRegs
}
