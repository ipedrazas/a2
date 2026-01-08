// Package language provides language detection for projects.
package language

import (
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
}

// Detect analyzes a directory and returns detected languages.
func Detect(path string) DetectionResult {
	result := DetectionResult{
		Languages:  make([]checker.Language, 0),
		Indicators: make(map[checker.Language][]string),
	}

	// Check languages in a defined order for consistent primary selection
	orderedLanguages := []checker.Language{checker.LangGo, checker.LangPython, checker.LangNode, checker.LangJava}

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
