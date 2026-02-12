// Package security provides safety/security checks for detecting potential
// malicious code patterns including shell injection, file system abuse,
// network exfiltration, and code obfuscation.
package security

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/safepath"
)

// File extensions by language
var languageExtensions = map[string][]string{
	"go":         {".go"},
	"python":     {".py", ".pyx", ".pyi"},
	"node":       {".js", ".jsx", ".mjs"},
	"typescript": {".ts", ".tsx"},
	"java":       {".java"},
	"rust":       {".rs"},
	"swift":      {".swift"},
	"ruby":       {".rb"},
	"php":        {".php"},
	"c":          {".c", ".h"},
	"cpp":        {".cpp", ".cc", ".cxx", ".hpp", ".hxx"},
}

// Comment patterns by language
var commentPatterns = map[string][]string{
	"go":         {"//"},
	"python":     {"#"},
	"node":       {"//", "/*"},
	"typescript": {"//", "/*"},
	"java":       {"//", "/*"},
	"rust":       {"//", "/*"},
	"swift":      {"//", "/*"},
	"ruby":       {"#"},
	"php":        {"//", "#", "/*"},
	"c":          {"//", "/*"},
	"cpp":        {"//", "/*"},
}

// Directories to skip during scanning
var skipDirectories = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	".git":         true,
	"__pycache__":  true,
	".venv":        true,
	"venv":         true,
	"dist":         true,
	"build":        true,
	".idea":        true,
	".vscode":      true,
	".next":        true,
	".nuxt":        true,
	"target":       true,
	"bin":          true,
	"obj":          true,
	"out":          true,
	".terraform":   true,
	"coverage":     true,
	".cache":       true,
	".gradle":      true,
	".mypy_cache":  true,
	"pytest_cache": true,
	".egg-info":    true,
	".tox":         true,
}

// calculateEntropy calculates the Shannon entropy of a string.
// Higher entropy (typically >7.5) suggests encoded or random data.
func calculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	// Count character frequencies
	freq := make(map[rune]float64)
	for _, ch := range s {
		freq[ch]++
	}

	// Calculate entropy
	var entropy float64
	length := float64(len(s))
	for _, count := range freq {
		p := count / length
		entropy -= p * math.Log2(p)
	}

	return entropy
}

// isHighEntropy checks if a string has suspiciously high entropy.
// Threshold of 7.5 is commonly used to detect encoded/encrypted content.
func isHighEntropy(s string) bool {
	if len(s) < 16 {
		return false // Too short to meaningfully calculate entropy
	}
	return calculateEntropy(s) > 7.5
}

// isBase64 checks if a string appears to be base64 encoded.
func isBase64(s string) bool {
	// Base64 strings should only contain valid characters
	matched, _ := regexp.MatchString(`^[A-Za-z0-9+/=]+$`, s)
	if !matched {
		return false
	}

	// Should have length that's a multiple of 4 (with padding)
	if len(s)%4 != 0 && len(s)%4 != 1 { // Allow for missing padding
		return false
	}

	// Try to decode
	_, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		return true
	}

	// Try with padding
	_, err = base64.StdEncoding.DecodeString(s + strings.Repeat("=", (4-len(s)%4)%4))
	return err == nil
}

// isHex checks if a string appears to be hex encoded.
func isHex(s string) bool {
	// Hex strings should only contain hex characters
	matched, _ := regexp.MatchString(`^[0-9a-fA-F]+$`, s)
	if !matched {
		return false
	}

	// Hex encoded content is typically even length
	return len(s)%2 == 0 && len(s) >= 8
}

// looksLikeBase64InCode checks if a string literal looks like base64 in code.
// This is more lenient than strict base64 validation as code might have
// string concatenation or other quirks.
func looksLikeBase64InCode(s string) bool {
	// Remove quotes and whitespace
	s = strings.Trim(s, `"'`)
	s = strings.TrimSpace(s)

	// Check for base64-like characteristics
	if len(s) < 16 {
		return false
	}

	// Count base64 characters
	base64Chars := 0
	for _, ch := range s {
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') || ch == '+' || ch == '/' || ch == '=' {
			base64Chars++
		}
	}

	ratio := float64(base64Chars) / float64(len(s))
	return ratio > 0.9
}

// hasHexEscapeSequence checks for excessive hex escape sequences.
func hasHexEscapeSequence(s string) bool {
	// Count \xNN or \uNNNN patterns
	hexEscape := regexp.MustCompile(`\\[xu][0-9a-fA-F]{2,4}`)
	matches := hexEscape.FindAllString(s, -1)
	if len(matches) > 3 {
		return true
	}

	// Check for many hex-like patterns in a row
	consecutiveHex := regexp.MustCompile(`(?:\\[xu][0-9a-fA-F]{2,4}\s*){3,}`)
	return consecutiveHex.MatchString(s)
}

// detectLanguageFromPath attempts to detect the primary language from a file path.
func detectLanguageFromPath(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	for lang, exts := range languageExtensions {
		for _, e := range exts {
			if ext == e {
				return lang
			}
		}
	}

	return ""
}

// isCommentLine checks if a line is a comment for the given language.
func isCommentLine(line, language string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}

	patterns, ok := commentPatterns[language]
	if !ok {
		// Default to common patterns
		patterns = []string{"//", "#", "/*"}
	}

	for _, pattern := range patterns {
		if strings.HasPrefix(trimmed, pattern) {
			return true
		}
	}

	return false
}

// isTestFile checks if a file is a test file based on its name.
func isTestFile(fileName string) bool {
	lower := strings.ToLower(fileName)
	return strings.Contains(lower, "_test.") ||
		strings.Contains(lower, ".test.") ||
		strings.Contains(lower, ".spec.") ||
		strings.HasPrefix(lower, "test_")
}

// shouldSkipFile checks if a file should be skipped during scanning.
func shouldSkipFile(fileName string) bool {
	lower := strings.ToLower(fileName)

	// Skip example/template files
	if strings.Contains(lower, "example") ||
		strings.Contains(lower, "sample") ||
		strings.Contains(lower, "template") ||
		strings.Contains(lower, "mock") ||
		strings.Contains(lower, "fixture") {
		return true
	}

	// Skip documentation
	if strings.HasSuffix(lower, ".md") ||
		strings.HasSuffix(lower, ".txt") ||
		strings.HasSuffix(lower, ".rst") {
		return true
	}

	return false
}

// isSourceFile checks if a file is a source file based on extension.
func isSourceFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	for _, exts := range languageExtensions {
		for _, e := range exts {
			if ext == e {
				return true
			}
		}
	}

	// Also include config files
	configExts := []string{".yaml", ".yml", ".json", ".toml", ".xml", ".sh", ".bash"}
	for _, e := range configExts {
		if ext == e {
			return true
		}
	}

	return false
}

// Finding represents a detected security issue.
type Finding struct {
	Type        string
	File        string
	Line        int
	Description string
	Severity    string // "critical", "high", "medium", "low"
}

// String returns a formatted representation of the finding.
func (f Finding) String() string {
	return fmt.Sprintf("%s in %s:%d", f.Description, f.File, f.Line)
}

// FileScanner handles scanning of files for security patterns.
type FileScanner struct {
	RootPath string
}

// NewFileScanner creates a new file scanner.
func NewFileScanner(rootPath string) *FileScanner {
	return &FileScanner{RootPath: rootPath}
}

// ScanFile scans a single file and returns lines matching the given patterns.
func (fs *FileScanner) ScanFile(filePath string, patternGroups map[string][]*regexp.Regexp, language string) []Finding {
	var findings []Finding

	file, err := filepath.Abs(filePath)
	if err != nil {
		return findings
	}

	// Open file relative to root
	relPath, err := filepath.Rel(fs.RootPath, file)
	if err != nil {
		return findings
	}

	// Use safepath to open the file
	content, err := readFile(fs.RootPath, relPath)
	if err != nil {
		return findings
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comment lines for most patterns
		if language != "" && isCommentLine(line, language) {
			continue
		}

		// Check each pattern group
		for groupName, patterns := range patternGroups {
			for _, pattern := range patterns {
				if pattern.MatchString(line) {
					findings = append(findings, Finding{
						Type:        groupName,
						File:        relPath,
						Line:        lineNum,
						Description: groupName,
						Severity:    "high",
					})
					break // One match per line per group
				}
			}
		}
	}

	return findings
}

// readFile reads a file relative to root path safely.
func readFile(root, relPath string) (string, error) {
	content, err := safepath.ReadFile(root, relPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// FormatFindings formats a slice of findings into a readable message.
func FormatFindings(findings []Finding, maxItems int) string {
	if len(findings) == 0 {
		return "No issues found"
	}

	var items []string
	for i, f := range findings {
		if i >= maxItems {
			break
		}
		items = append(items, f.String())
	}

	msg := strings.Join(items, ", ")
	if len(findings) > maxItems {
		msg += fmt.Sprintf(" (%d more)", len(findings)-maxItems)
	}

	return msg
}
