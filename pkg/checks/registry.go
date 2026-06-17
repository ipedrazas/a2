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
	securitycheck "github.com/ipedrazas/a2/pkg/checks/security"
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
	res, err := p.checker.Run(actualPath)
	// Record which source_dir this check ran in so output can disambiguate
	// the same check running across multiple directories (monorepos).
	res.SourceDir = p.sourceDir
	return res, err
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
	for _, dirs := range m.cfg.GetSourceDirs() {
		for _, dir := range dirs {
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
// Checks already scoped to a specific source_dir by per-language disabling
// (pathResolvingChecker) are left untouched.
func wrapCommonRegistrations(regs []checker.CheckRegistration, cfg *config.Config) []checker.CheckRegistration {
	for i := range regs {
		if _, scoped := regs[i].Checker.(*pathResolvingChecker); scoped {
			continue
		}
		if commonChecksUsingSourceDirs[regs[i].Meta.ID] {
			regs[i].Checker = &multiPathChecker{checker: regs[i].Checker, cfg: cfg}
		}
	}
	return regs
}

// scopePerLanguage applies per-language disabled lists (checks.<lang>.disabled)
// to language-agnostic checks (common, devops, security). These checks normally
// run once at the repo root. When a check is disabled for only a subset of the
// detected languages, it is re-scoped to run against the source_dir(s) of the
// languages that still want it; when disabled for every detected language it is
// dropped entirely. Checks unaffected by any per-language list are returned
// unchanged, preserving the single-root-run behaviour.
func scopePerLanguage(regs []checker.CheckRegistration, cfg *config.Config, detected language.DetectionResult) []checker.CheckRegistration {
	if len(cfg.Checks.PerLanguage) == 0 || len(detected.Languages) == 0 {
		return regs
	}

	var out []checker.CheckRegistration
	for _, reg := range regs {
		id := reg.Meta.ID

		var enabledLangs []checker.Language
		anyDisabled := false
		for _, lang := range detected.Languages {
			if cfg.IsCheckDisabledForOnlyLang(id, string(lang)) {
				anyDisabled = true
			} else {
				enabledLangs = append(enabledLangs, lang)
			}
		}

		if !anyDisabled {
			// No per-language list touches this check: keep root behaviour.
			out = append(out, reg)
			continue
		}
		if len(enabledLangs) == 0 {
			// Disabled for every detected language: drop it.
			continue
		}

		// Run the check scoped to each still-enabled language's source_dir(s).
		for _, path := range scopedPathsForLangs(cfg, enabledLangs) {
			scoped := reg
			if path != "" {
				scoped.Checker = &pathResolvingChecker{checker: reg.Checker, sourceDir: path}
			}
			out = append(out, scoped)
		}
	}
	return out
}

// scopedPathsForLangs returns the deduplicated set of source_dir paths for the
// given languages. A language with no configured source_dir contributes the
// repo root (empty string).
func scopedPathsForLangs(cfg *config.Config, langs []checker.Language) []string {
	seen := make(map[string]bool)
	var paths []string
	for _, lang := range langs {
		dirs := cfg.GetSourceDirsForLang(string(lang))
		if len(dirs) == 0 {
			dirs = []string{""}
		}
		for _, d := range dirs {
			if seen[d] {
				continue
			}
			seen[d] = true
			paths = append(paths, d)
		}
	}
	return paths
}

// GetChecks returns checks for detected languages based on configuration.
// Checks are ordered with critical checks first.
// Language-specific checks are wrapped to use the configured source_dir.
func GetChecks(cfg *config.Config, detected language.DetectionResult) []checker.CheckRegistration {
	var registrations []checker.CheckRegistration

	// Get checks for each detected language
	for _, lang := range detected.Languages {
		// Per-language disabled list (checks.<lang>.disabled) applies to every
		// check run in this language's context.
		langDisabled := cfg.Checks.PerLanguage[string(lang)]
		entries := cfg.GetSourceDirEntriesForLang(string(lang))
		if len(entries) == 0 {
			// No source_dir configured, run checks from root
			regs := getChecksForLanguage(lang, cfg)
			if len(langDisabled) > 0 {
				regs = filterByDisabledList(regs, langDisabled)
			}
			registrations = append(registrations, regs...)
		} else {
			// Create a full set of checks per source_dir
			for _, entry := range entries {
				regs := getChecksForLanguage(lang, cfg)
				if len(langDisabled) > 0 {
					regs = filterByDisabledList(regs, langDisabled)
				}
				// Filter out checks disabled by the directory's profile
				if len(entry.Disabled) > 0 {
					regs = filterByDisabledList(regs, entry.Disabled)
				}
				// Override coverage threshold if set per source_dir entry
				if entry.CoverageThreshold > 0 {
					for i := range regs {
						if setter, ok := regs[i].Checker.(checker.CoverageThresholdSetter); ok {
							setter.SetCoverageThreshold(entry.CoverageThreshold)
						}
					}
				}
				for i := range regs {
					regs[i].Checker = &pathResolvingChecker{
						checker:   regs[i].Checker,
						sourceDir: entry.Path,
					}
				}
				registrations = append(registrations, regs...)
			}
		}
	}

	// Add devops checks (language-agnostic, root path)
	registrations = append(registrations, scopePerLanguage(devopscheck.Register(cfg), cfg, detected)...)
	// Add security checks (language-agnostic, root path)
	registrations = append(registrations, scopePerLanguage(securitycheck.Register(cfg), cfg, detected)...)
	// Add common checks (language-agnostic, root path; some run on source_dirs too).
	// Per-language scoping runs first so checks disabled for a subset of languages
	// run only against the still-enabled languages' source_dirs; multi-path
	// wrapping then applies to any check left running once at the repo root.
	registrations = append(registrations, wrapCommonRegistrations(scopePerLanguage(common.Register(cfg), cfg, detected), cfg)...)

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

// FilterFast returns only the fast (static/IO-light) checks, dropping any
// SpeedSlow check that spawns builds/tests/network work. It is an orthogonal
// filter applied on top of profile/target/skip selection for --quick runs.
func FilterFast(regs []checker.CheckRegistration) []checker.CheckRegistration {
	filtered := make([]checker.CheckRegistration, 0, len(regs))
	for _, reg := range regs {
		if reg.Meta.Speed == checker.SpeedFast {
			filtered = append(filtered, reg)
		}
	}
	return filtered
}

// filterByDisabledList removes checks whose ID matches any entry in the disabled list.
func filterByDisabledList(regs []checker.CheckRegistration, disabled []string) []checker.CheckRegistration {
	var filtered []checker.CheckRegistration
	for _, reg := range regs {
		skip := false
		for _, d := range disabled {
			if config.MatchDisabled(reg.Meta.ID, d) {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, reg)
		}
	}
	return filtered
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
	allRegs = append(allRegs, securitycheck.Register(cfg)...)
	allRegs = append(allRegs, wrapCommonRegistrations(common.Register(cfg), cfg)...)

	// Sort by order
	sort.Slice(allRegs, func(i, j int) bool {
		return allRegs[i].Meta.Order < allRegs[j].Meta.Order
	})

	return allRegs
}
