package security

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ObfuscationCheck detects code obfuscation patterns and encoded strings.
type ObfuscationCheck struct {
	Patterns map[string][]*regexp.Regexp
}

// ID returns the unique identifier for this check.
func (c *ObfuscationCheck) ID() string {
	return "security:obfuscation"
}

// Name returns the human-readable name for this check.
func (c *ObfuscationCheck) Name() string {
	return "Code Obfuscation Detection"
}

// Run executes the obfuscation detection check.
func (c *ObfuscationCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Initialize patterns if not already done
	if c.Patterns == nil {
		c.Patterns = c.getPatterns()
	}

	// Scan all source files
	findings := c.scanDirectory(path)

	if len(findings) == 0 {
		return rb.Pass("No obfuscated code or encoded strings detected"), nil
	}

	// Format findings into message
	msg := c.formatFindings(findings)
	return rb.Fail(msg), nil // Obfuscation is suspicious -> Fail
}

// getPatterns returns patterns for detecting obfuscation.
func (c *ObfuscationCheck) getPatterns() map[string][]*regexp.Regexp {
	patterns := make(map[string][]*regexp.Regexp)

	// Patterns that apply to most languages
	patterns["common"] = []*regexp.Regexp{
		// Excessive hex escape sequences
		regexp.MustCompile(`(\\x[0-9a-fA-F]{2}\s*){5,}`), // 5+ hex escapes in sequence
		regexp.MustCompile(`(\\u[0-9a-fA-F]{4}\s*){3,}`), // 3+ unicode escapes in sequence
		// Large base64-like strings in code
		regexp.MustCompile(`["'][A-Za-z0-9+/]{40,}={0,2}["']`), // Long base64 strings
		// Character array obfuscation
		regexp.MustCompile(`\["[a-zA-Z0-9]"(?:\s*,\s*"[a-zA-Z0-9]"){10,}`), // 10+ single-char strings in array
		regexp.MustCompile(`['"][a-zA-Z0-9]['"](?:\s*,\s*['"][a-zA-Z0-9]['"]){10,}`),
		// String concatenation chains
		regexp.MustCompile(`["'][^"']{5,20}["']\s*\+\s*["'][^"']{5,20}["']\s*\+\s*["'][^"']{5,20}["']`),
	}

	// Go-specific patterns
	patterns["go"] = []*regexp.Regexp{
		// Hex-encoded strings
		regexp.MustCompile(`\\x[0-9a-fA-F]{2}`),
		// Byte slice construction
		regexp.MustCompile(`\[\]byte\s*\{[^}]{50,}\}`), // Large byte slices
		// String concatenation with runes
		regexp.MustCompile(`string\\(rune\\(0x[0-9a-fA-F]+\\)\\)`),
	}

	// Python-specific patterns
	patterns["python"] = []*regexp.Regexp{
		// Encode/decode chains
		regexp.MustCompile(`\.encode\s*\(\s*["'](?:base64|hex|rot13)`),
		regexp.MustCompile(`\.decode\s*\(\s*["'](?:base64|hex|rot13)`),
		// Bytes from hex
		regexp.MustCompile(`bytes\.fromhex\s*\(\s*["'][0-9a-fA-F]{20,}`),
		// Code object manipulation
		regexp.MustCompile(`compile\s*\(`),
		regexp.MustCompile(`types\.CodeType`),
		// String concatenation via join
		regexp.MustCompile(`["']\s*\.join\s*\(\s*\[[^\]]{50,}\]\)`),
		// List comprehension for building strings
		regexp.MustCompile(`chr\s*\(\s*0x[0-9a-fA-F]+`),
	}

	// Node/JavaScript-specific patterns
	patterns["node"] = []*regexp.Regexp{
		// atob/btoa (base64 encoding/decoding)
		regexp.MustCompile(`atob\s*\(`),
		regexp.MustCompile(`btoa\s*\(\s*[a-zA-Z_$]\w*\s*\)`),
		// Buffer operations
		regexp.MustCompile(`Buffer\.from\s*\(\s*["'][^"']{20,}["'],\s*["']`),
		regexp.MustCompile(`new Buffer\s*\(`),
		// String from char codes
		regexp.MustCompile(`String\.fromCharCode\s*\(\s*[0-9]+`),
		// Char code concatenation
		regexp.MustCompile(`\\[xu][0-9a-fA-F]{2,4}`),
	}

	// TypeScript-specific patterns (same as Node)
	patterns["typescript"] = patterns["node"]

	// Java-specific patterns
	patterns["java"] = []*regexp.Regexp{
		// Base64 encoding/decoding
		regexp.MustCompile(`Base64\.getDecoder\(\)\.decode\s*\(`),
		regexp.MustCompile(`DatatypeConverter\.parseBase64Binary\s*\(`),
		// Hex encoding
		regexp.MustCompile(`Integer\.parseInt\s*\([^)]+,\s*16\s*\)`),
		regexp.MustCompile(`Character\.toString\s*\([^)]+,\s*16\s*\)`),
		// String building with append
		regexp.MustCompile(`StringBuilder\s*\(`),
	}

	// Ruby-specific patterns
	patterns["ruby"] = []*regexp.Regexp{
		regexp.MustCompile(`\.unpack\s*\(\s*["']`),             // Binary unpacking
		regexp.MustCompile(`\.pack\s*\(\s*["']`),               // Binary packing
		regexp.MustCompile(`\[[0-9]+(?:\s*,\s*[0-9]+){10,}\]`), // Large number arrays
		regexp.MustCompile(`\*\s*\w+`),                         // String unpacking splats
	}

	// PHP-specific patterns
	patterns["php"] = []*regexp.Regexp{
		regexp.MustCompile(`base64_decode\s*\(`),
		regexp.MustCompile(`str_rot13\s*\(`),
		regexp.MustCompile(`pack\s*\(\s*["']`), // Binary packing
		regexp.MustCompile(`convert_uudecode\s*\(`),
	}

	// Rust-specific patterns
	patterns["rust"] = []*regexp.Regexp{
		regexp.MustCompile(`b\s*["'][A-Za-z0-9+/]{30,}=?["']`), // Byte strings (base64-like)
		regexp.MustCompile(`String::from_utf8_lossy\s*\(`),
		regexp.MustCompile(`as_bytes\s*\(\)\s*\.`),
	}

	// C/C++ specific patterns
	patterns["c"] = []*regexp.Regexp{
		regexp.MustCompile(`\\x[0-9a-fA-F]{2}`),                  // Hex escapes
		regexp.MustCompile(`\\[0-7]{3}`),                         // Octal escapes
		regexp.MustCompile(`char\s+\w+\[\]\s*=\s*\{[^}]{50,}\}`), // Large char arrays
	}
	patterns["cpp"] = patterns["c"]

	// Swift-specific patterns
	patterns["swift"] = []*regexp.Regexp{
		regexp.MustCompile(`\\u\{[0-9a-fA-F]+\}`),     // Unicode escapes
		regexp.MustCompile(`Data\(base64Encoded:\s*`), // Base64 encoded data
		regexp.MustCompile(`String\(data:\s*`),        // Data to string conversion
	}

	return patterns
}

// scanDirectory scans all files in directory for obfuscation patterns.
func (c *ObfuscationCheck) scanDirectory(path string) []Finding {
	var findings []Finding

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			if skipDirectories[info.Name()] {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") && info.Name() != ".github" && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this is a source file we should scan
		if !isSourceFile(filePath) {
			return nil
		}

		// Skip test files, examples, templates
		if isTestFile(filepath.Base(filePath)) || shouldSkipFile(filepath.Base(filePath)) {
			return nil
		}

		// Detect language from file extension
		language := detectLanguageFromPath(filePath)

		// Get patterns for this language (plus common patterns)
		var allPatterns []*regexp.Regexp

		// Add language-specific patterns
		if langPatterns, ok := c.Patterns[language]; ok {
			allPatterns = append(allPatterns, langPatterns...)
		}

		// Add common patterns
		if commonPatterns, ok := c.Patterns["common"]; ok {
			allPatterns = append(allPatterns, commonPatterns...)
		}

		if len(allPatterns) == 0 {
			return nil
		}

		// Scan file
		fileFindings := c.scanFile(path, filePath, language, allPatterns)
		findings = append(findings, fileFindings...)

		// Limit findings to avoid excessive output
		if len(findings) >= 50 {
			return filepath.SkipAll
		}

		return nil
	})

	_ = err
	return findings
}

// scanFile scans a single file for obfuscation patterns.
func (c *ObfuscationCheck) scanFile(root, filePath string, language string, patterns []*regexp.Regexp) []Finding {
	var findings []Finding

	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return findings
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// Get relative path for reporting
	relPath, err := filepath.Rel(root, filePath)
	if err != nil {
		relPath = filepath.Base(filePath)
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip comment lines for most patterns
		if isCommentLine(line, language) {
			// But still check for encoded content in comments
			// (malicious code sometimes hides in comments)
			if !c.isSuspiciousStringInComment(trimmed) {
				continue
			}
		}

		// Check for high-entropy strings
		if c.hasHighEntropyString(trimmed) {
			findings = append(findings, Finding{
				Type:        "obfuscation",
				File:        relPath,
				Line:        lineNum,
				Description: "high-entropy string (possible encoded data)",
				Severity:    "high",
			})
			continue
		}

		// Check each pattern
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				// Extract matched pattern for better reporting
				match := pattern.FindString(line)
				findings = append(findings, Finding{
					Type:        "obfuscation",
					File:        relPath,
					Line:        lineNum,
					Description: fmt.Sprintf("obfuscation pattern: %s", c.sanitizeMatch(match)),
					Severity:    "high",
				})
				break // One finding per line
			}
		}
	}

	return findings
}

// isSuspiciousStringInComment checks if a comment line contains suspicious encoded content.
func (c *ObfuscationCheck) isSuspiciousStringInComment(line string) bool {
	// Remove comment prefix
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
		trimmed = strings.TrimPrefix(trimmed, "//")
		trimmed = strings.TrimPrefix(trimmed, "#")
		trimmed = strings.TrimSpace(trimmed)
	}

	// Check for long base64-like strings
	if looksLikeBase64InCode(trimmed) && len(trimmed) > 40 {
		return true
	}

	// Check for many hex sequences
	if hasHexEscapeSequence(trimmed) {
		return true
	}

	return false
}

// hasHighEntropyString checks if a line contains a high-entropy string literal.
func (c *ObfuscationCheck) hasHighEntropyString(line string) bool {
	// Extract string literals from line
	strLiterals := c.extractStringLiterals(line)

	for _, s := range strLiterals {
		if len(s) >= 20 && isHighEntropy(s) {
			return true
		}
	}

	return false
}

// extractStringLiterals extracts string literals from a line.
func (c *ObfuscationCheck) extractStringLiterals(line string) []string {
	var literals []string

	// Double-quoted strings
	doubleQuoteRegex := regexp.MustCompile(`"([^"\\]|\\.)*"`)
	matches := doubleQuoteRegex.FindAllString(line, -1)
	for _, m := range matches {
		literals = append(literals, strings.Trim(m, `"`))
	}

	// Single-quoted strings
	singleQuoteRegex := regexp.MustCompile(`'([^'\\]|\\.)*'`)
	matches = singleQuoteRegex.FindAllString(line, -1)
	for _, m := range matches {
		literals = append(literals, strings.Trim(m, "'"))
	}

	// Backtick strings (Go, JavaScript)
	backtickRegex := regexp.MustCompile("`([^`]*)`")
	matches = backtickRegex.FindAllString(line, -1)
	for _, m := range matches {
		literals = append(literals, strings.Trim(m, "`"))
	}

	return literals
}

// sanitizeMatch sanitizes a matched pattern for display.
func (c *ObfuscationCheck) sanitizeMatch(match string) string {
	// Truncate very long matches
	if len(match) > 100 {
		return match[:100] + "..."
	}
	return match
}

// formatFindings formats findings into a readable message.
func (c *ObfuscationCheck) formatFindings(findings []Finding) string {
	if len(findings) == 1 {
		return fmt.Sprintf("Obfuscated code detected: %s", findings[0].String())
	}

	if len(findings) <= 3 {
		return fmt.Sprintf("Obfuscated code detected: %s", c.joinFindings(findings))
	}

	return fmt.Sprintf("%d obfuscation patterns detected (e.g., %s)", len(findings),
		c.joinFindings(findings[:3]))
}

// joinFindings joins findings for display.
func (c *ObfuscationCheck) joinFindings(findings []Finding) string {
	var strs []string
	for _, f := range findings {
		strs = append(strs, f.String())
	}
	return strings.Join(strs, ", ")
}
