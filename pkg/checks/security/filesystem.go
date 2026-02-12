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

// FileSystemCheck detects path traversal and unsafe file operations.
type FileSystemCheck struct {
	Patterns map[string][]*regexp.Regexp
}

// ID returns the unique identifier for this check.
func (c *FileSystemCheck) ID() string {
	return "security:filesystem"
}

// Name returns the human-readable name for this check.
func (c *FileSystemCheck) Name() string {
	return "File System Safety"
}

// Run executes the file system safety check.
func (c *FileSystemCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Initialize patterns if not already done
	if c.Patterns == nil {
		c.Patterns = c.getPatterns()
	}

	// Scan all source files
	findings := c.scanDirectory(path)

	if len(findings) == 0 {
		return rb.Pass("No path traversal or unsafe file operations detected"), nil
	}

	// Format findings into message
	msg := c.formatFindings(findings)
	return rb.Fail(msg), nil
}

// getPatterns returns language-specific patterns for detecting file system abuse.
func (c *FileSystemCheck) getPatterns() map[string][]*regexp.Regexp {
	patterns := make(map[string][]*regexp.Regexp)

	// Go patterns
	patterns["go"] = []*regexp.Regexp{
		// File operations with variables
		regexp.MustCompile(`ioutil\.(ReadFile|WriteFile)\s*\(\s*[a-zA-Z_]\w*\s*[\+,]`),
		regexp.MustCompile(`os\.(Open|OpenFile|ReadFile|WriteFile)\s*\(\s*[a-zA-Z_]\w*\s*[\+,]`),
		regexp.MustCompile(`os\.(ReadDir|ReadFile)\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`filepath\.Join\s*\(\s*[a-zA-Z_]\w*\s*,\s*[a-zA-Z_]\w*\s*\)`), // Multiple variables
		// Path traversal patterns
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory reference
		regexp.MustCompile(`~\/`),                                    // Home directory reference
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`os\.Getenv\s*\(\s*[\"']`),                // Environment variable usage
		regexp.MustCompile(`filepath\.Abs\s*\(\s*[a-zA-Z_]\w*\s*\)`), // Absolute path from variable
	}

	// Python patterns
	patterns["python"] = []*regexp.Regexp{
		// File operations with user input
		regexp.MustCompile(`open\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
		regexp.MustCompile(`pathlib\.Path\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`Path\s*\(\s*[a-zA-Z_]\w*\s*\)[.]open\s*\(`),
		// os.path operations
		regexp.MustCompile(`os\.path\.(join|abspath)\s*\(\s*[a-zA-Z_]\w*\s*,\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`os\.(makedirs|removedirs|rename)\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
		// Path traversal patterns
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory
		regexp.MustCompile(`~\/`),                                    // Home directory
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`os\.getenv\s*\(\s*["']`),                 // Environment variables
		regexp.MustCompile(`os\.environ\[`),                          // Environment access
		// shutil operations
		regexp.MustCompile(`shutil\.(copy|move|rmtree)\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
	}

	// Node/JavaScript patterns
	patterns["node"] = []*regexp.Regexp{
		// fs operations with variables
		regexp.MustCompile(`fs\.(readFile|writeFile|readFileSync|writeFileSync|open|openSync)\s*\(\s*[a-zA-Z_$]\w*\s*[,+]`),
		regexp.MustCompile(`fs\.exists\s*\(\s*[a-zA-Z_$]\w*\s*[,+]`),
		regexp.MustCompile(`fs\.stat\s*\(\s*[a-zA-Z_$]\w*\s*`),
		// path operations
		regexp.MustCompile(`path\.(join|resolve|normalize)\s*\(\s*[a-zA-Z_$]\w*\s*[,+]`),
		regexp.MustCompile(`path\.join\s*\(\s*["']\.\.\/`),    // Explicit parent dir
		regexp.MustCompile(`path\.resolve\s*\(\s*["']\.\.\/`), // Parent dir in resolve
		// Path traversal patterns
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory
		regexp.MustCompile(`~\/`),                                    // Home directory
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`process\.env\.[a-zA-Z_$]\w*`),            // Environment access
		// child_process file operations
		regexp.MustCompile(`child_process\.(exec|spawn)\s*\(\s*[a-zA-Z_$]\w*\s*[,+]"`),
		// fs-extra operations
		regexp.MustCompile(`fs-extra\.(copy|move|remove)`),
	}

	// TypeScript patterns (same as Node plus TypeScript-specific)
	patterns["typescript"] = patterns["node"]

	// Java patterns
	patterns["java"] = []*regexp.Regexp{
		// File operations with variables
		regexp.MustCompile(`new\s+File\s*\(\s*[a-zA-Z_]\w*\s*[\+,]`),
		regexp.MustCompile(`File\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`Files\.(read|write|copy|move)\s*\(\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`FileReader|FileWriter\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`FileInputStream|FileOutputStream\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		// Path operations
		regexp.MustCompile(`Paths\.get\s*\(\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`Path\.of\s*\(\s*[a-zA-Z_]\w*`),
		// Path traversal patterns
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`System\.getenv\s*\(`),                    // Environment access
	}

	// Ruby patterns
	patterns["ruby"] = []*regexp.Regexp{
		// File operations with user input
		regexp.MustCompile(`File\.(open|read|write|delete)\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
		regexp.MustCompile(`File\.open\s*\(\s*["'].*#\{[a-zA-Z_]\w*\}`),
		regexp.MustCompile(`IO\.(read|write|foreach)\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
		// File operations
		regexp.MustCompile(`FileUtils\.(cp|mv|rm|mkdir_p)\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
		// Path traversal patterns
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`ENV\[`),                                  // Environment access
	}

	// PHP patterns
	patterns["php"] = []*regexp.Regexp{
		// File operations with variables
		regexp.MustCompile(`file_(get|put|read|write)_contents\s*\(\s*\$[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`fopen\s*\(\s*\$[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`unlink\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`is_dir|is_file\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		// Path traversal patterns
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`\$_ENV\[`),                               // Environment access
		regexp.MustCompile(`\$_SERVER\[`),                            // Server variables
	}

	// C/C++ patterns
	patterns["c"] = []*regexp.Regexp{
		regexp.MustCompile(`fopen\s*\(\s*[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`open\s*\(\s*[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`remove\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`unlink\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`getenv\s*\(`),                            // Environment access
	}
	patterns["cpp"] = patterns["c"]

	// Rust patterns
	patterns["rust"] = []*regexp.Regexp{
		regexp.MustCompile(`File::open\s*\(\s*&?[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`File::create\s*\(\s*&?[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`fs::(read|write|remove|rename)\s*\(\s*&?[a-zA-Z_]\w*`),
		regexp.MustCompile(`std::fs::.*\s*\(\s*&?[a-zA-Z_]\w*`),
		regexp.MustCompile(`PathBuf::(from|push)\s*\(\s*&?[a-zA-Z_]\w*`),
		regexp.MustCompile(`\.\.\/`),                                 // Parent directory
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root)`), // System directories
		regexp.MustCompile(`std::env::var\s*\(`),                     // Environment access
	}

	// Swift patterns
	patterns["swift"] = []*regexp.Regexp{
		regexp.MustCompile(`FileManager\.(default|s*).*\.(contents|attributes)\(atPath:\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`FileHandle\s*\(\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`\.\.\/`),                                // Parent directory
		regexp.MustCompile(`ProcessInfo\.processInfo\.environment`), // Environment access
	}

	return patterns
}

// scanDirectory scans all files in directory for file system abuse patterns.
func (c *FileSystemCheck) scanDirectory(path string) []Finding {
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

		// Get patterns for this language
		patterns, ok := c.Patterns[language]
		if !ok || len(patterns) == 0 {
			return nil
		}

		// Scan file
		fileFindings := c.scanFile(path, filePath, language, patterns)
		findings = append(findings, fileFindings...)

		// Limit findings to avoid excessive output
		if len(findings) >= 50 {
			return filepath.SkipAll
		}

		return nil
	})

	_ = err // Ignore error
	return findings
}

// scanFile scans a single file for file system abuse patterns.
func (c *FileSystemCheck) scanFile(root, filePath string, language string, patterns []*regexp.Regexp) []Finding {
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

		// Skip comment lines
		if isCommentLine(line, language) {
			continue
		}

		// Check each pattern
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				// Extract matched pattern for better reporting
				match := pattern.FindString(line)
				findings = append(findings, Finding{
					Type:        "filesystem",
					File:        relPath,
					Line:        lineNum,
					Description: fmt.Sprintf("unsafe file operation: %s", c.sanitizeMatch(match)),
					Severity:    "high",
				})
				break // One finding per line
			}
		}
	}

	return findings
}

// sanitizeMatch sanitizes a matched pattern for display.
func (c *FileSystemCheck) sanitizeMatch(match string) string {
	// Truncate very long matches
	if len(match) > 100 {
		return match[:100] + "..."
	}
	return match
}

// formatFindings formats findings into a readable message.
func (c *FileSystemCheck) formatFindings(findings []Finding) string {
	if len(findings) == 1 {
		return fmt.Sprintf("Path traversal/unsafe file operation detected: %s", findings[0].String())
	}

	if len(findings) <= 3 {
		return fmt.Sprintf("Path traversal/unsafe file operations detected: %s", c.joinFindings(findings))
	}

	return fmt.Sprintf("%d path traversal/unsafe file operations detected (e.g., %s)", len(findings),
		c.joinFindings(findings[:3]))
}

// joinFindings joins findings for display.
func (c *FileSystemCheck) joinFindings(findings []Finding) string {
	var strs []string
	for _, f := range findings {
		strs = append(strs, f.String())
	}
	return strings.Join(strs, ", ")
}
