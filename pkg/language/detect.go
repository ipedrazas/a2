// Package language provides language detection for projects.
package language

import (
	"path/filepath"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DetectionResult contains detected languages and their indicators.
type DetectionResult struct {
	Languages  []checker.Language            // Detected languages
	Indicators map[checker.Language][]string // Files that indicated each language
	Primary    checker.Language              // Primary language (first detected)
	MultiLang  bool                          // True if multiple languages detected
}

// LanguageIndicators maps languages to their indicator files.
var LanguageIndicators = map[checker.Language][]string{
	checker.LangGo: {
		"go.mod",
		"go.sum",
	},
	checker.LangPython: {
		"pyproject.toml",
		"setup.py",
		"requirements.txt",
		"Pipfile",
		"poetry.lock",
		"setup.cfg",
	},
	checker.LangNode: {
		"package.json",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"bun.lockb",
	},
	checker.LangJava: {
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
		"settings.gradle",
		"settings.gradle.kts",
		"mvnw",
		"gradlew",
	},
	checker.LangRust: {
		"Cargo.toml",
		"Cargo.lock",
		"rust-toolchain.toml",
		"rust-toolchain",
	},
	checker.LangTypeScript: {
		"tsconfig.json",
		"tsconfig.base.json",
		"tsconfig.build.json",
	},
	checker.LangSwift: {
		"Package.swift",
		"Package.resolved",
	},
}

// Detect analyzes a directory and returns detected languages.
func Detect(path string) DetectionResult {
	result := DetectionResult{
		Languages:  make([]checker.Language, 0),
		Indicators: make(map[checker.Language][]string),
	}

	// Check languages in a defined order for consistent primary selection
	orderedLanguages := []checker.Language{checker.LangGo, checker.LangPython, checker.LangNode, checker.LangJava, checker.LangRust, checker.LangTypeScript, checker.LangSwift}

	for _, lang := range orderedLanguages {
		indicators := LanguageIndicators[lang]
		found := []string{}
		for _, indicator := range indicators {
			if safepath.Exists(path, indicator) {
				found = append(found, indicator)
			}
		}
		if len(found) > 0 {
			result.Languages = append(result.Languages, lang)
			result.Indicators[lang] = found
		}
	}

	// n8n is detected by file content (the export shape), not by a known
	// filename, so it gets its own pass after the indicator-file languages.
	if found := detectN8nIndicators(path); len(found) > 0 {
		result.Languages = append(result.Languages, checker.LangN8n)
		result.Indicators[checker.LangN8n] = found
	}

	// Set primary language and multi-lang flag
	result.MultiLang = len(result.Languages) > 1
	if len(result.Languages) > 0 {
		result.Primary = result.Languages[0]
	}

	return result
}

// DetectWithOverride returns detection but respects explicit language override.
func DetectWithOverride(path string, explicit []checker.Language) DetectionResult {
	if len(explicit) > 0 {
		return DetectionResult{
			Languages:  explicit,
			Indicators: nil, // No indicator files when explicit
			Primary:    explicit[0],
			MultiLang:  len(explicit) > 1,
		}
	}
	return Detect(path)
}

// HasLanguage checks if a language is in the detection result.
func (r DetectionResult) HasLanguage(lang checker.Language) bool {
	for _, l := range r.Languages {
		if l == lang {
			return true
		}
	}
	return false
}

// DetectWithSourceDirs detects languages, checking configured source directories.
// sourceDirs maps language names (e.g., "rust") to subdirectories (e.g., ["src-tauri"]).
// Each language may have multiple source directories configured.
// Languages without a configured source_dir are detected in the root path.
func DetectWithSourceDirs(path string, sourceDirs map[string][]string) DetectionResult {
	result := DetectionResult{
		Languages:  make([]checker.Language, 0),
		Indicators: make(map[checker.Language][]string),
	}

	// Check languages in a defined order for consistent primary selection
	orderedLanguages := []checker.Language{checker.LangGo, checker.LangPython, checker.LangNode, checker.LangJava, checker.LangRust, checker.LangTypeScript, checker.LangSwift}

	for _, lang := range orderedLanguages {
		indicators := LanguageIndicators[lang]
		found := []string{}

		// Determine the paths to check for this language
		dirs, ok := sourceDirs[string(lang)]
		if !ok || len(dirs) == 0 {
			// No source_dir configured, check root path
			dirs = []string{""}
		}

		for _, dir := range dirs {
			checkPath := path
			if dir != "" {
				checkPath = filepath.Join(path, dir)
			}
			for _, indicator := range indicators {
				if safepath.Exists(checkPath, indicator) {
					found = append(found, indicator)
				}
			}
		}
		if len(found) > 0 {
			result.Languages = append(result.Languages, lang)
			result.Indicators[lang] = found
		}
	}

	// n8n detection is content-based and scoped like the other languages:
	// configured n8n source_dirs, or the root when none are set.
	n8nDirs, ok := sourceDirs[string(checker.LangN8n)]
	if !ok || len(n8nDirs) == 0 {
		n8nDirs = []string{""}
	}
	var n8nFound []string
	for _, dir := range n8nDirs {
		checkPath := path
		if dir != "" {
			checkPath = filepath.Join(path, dir)
		}
		n8nFound = append(n8nFound, detectN8nIndicators(checkPath)...)
	}
	if len(n8nFound) > 0 {
		result.Languages = append(result.Languages, checker.LangN8n)
		result.Indicators[checker.LangN8n] = n8nFound
	}

	// Set primary language and multi-lang flag
	result.MultiLang = len(result.Languages) > 1
	if len(result.Languages) > 0 {
		result.Primary = result.Languages[0]
	}

	return result
}
