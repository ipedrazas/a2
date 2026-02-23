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

type allowRule struct {
	filePattern string
	line        int
	matchText   string
}

var (
	goSafeJoinRE        = regexp.MustCompile(`^\s*([a-zA-Z_]\w*)\s*,\s*[a-zA-Z_]\w*\s*(?:=|:=)\s*safepath\.SafeJoin\s*\(`)
	goSafeJoinSingleRE  = regexp.MustCompile(`^\s*([a-zA-Z_]\w*)\s*(?:=|:=)\s*safepath\.SafeJoin\s*\(`)
	goJoinAssignRE      = regexp.MustCompile(`^\s*([a-zA-Z_]\w*)\s*(?:=|:=)\s*filepath\.Join\s*\(\s*([a-zA-Z_]\w*)\s*,`)
	goFileOpVarRE       = regexp.MustCompile(`\b(?:ioutil|os)\.(?:Open|OpenFile|ReadFile|WriteFile|ReadDir)\s*\(\s*([a-zA-Z_]\w*)`)
	goRootVarCandidates = map[string]struct{}{
		"root":        {},
		"rootDir":     {},
		"rootPath":    {},
		"repoRoot":    {},
		"repoDir":     {},
		"projectRoot": {},
		"projectDir":  {},
		"baseDir":     {},
		"basePath":    {},
		"workspace":   {},
		"workDir":     {},
	}
)

// FileSystemCheck detects path traversal and unsafe file operations.
type FileSystemCheck struct {
	Patterns map[string][]*regexp.Regexp
	// Allowlist contains rules to suppress known-safe findings.
	Allowlist  []string
	allowRules []allowRule
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

	if c.Patterns == nil {
		c.Patterns = c.getPatterns()
	}
	if c.allowRules == nil {
		c.allowRules = c.compileAllowRules()
	}

	findings := c.scanDirectory(path)

	if len(findings) == 0 {
		return rb.Pass("No unsafe file writes or system directory access detected"), nil
	}

	msg := c.formatFindings(findings)
	return rb.Fail(msg), nil
}

// getPatterns returns language-specific patterns for detecting file system abuse.
// Focus: write operations to unexpected/system directories, path traversal.
// Excluded: ENV VAR usage (users consider this safe for config paths).
func (c *FileSystemCheck) getPatterns() map[string][]*regexp.Regexp {
	patterns := make(map[string][]*regexp.Regexp)

	// Go patterns - focus on write operations and system directory access
	patterns["go"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`ioutil\.WriteFile\s*\(\s*[a-zA-Z_]\w*\s*[\+,]`),
		regexp.MustCompile(`os\.WriteFile\s*\(\s*[a-zA-Z_]\w*\s*[\+,]`),
		regexp.MustCompile(`os\.Create\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`os\.OpenFile\s*\(\s*[a-zA-Z_]\w*\s*,\s*os\.O_WRONLY|os\.O_RDWR|os\.O_CREATE`),
		// System directory writes (hardcoded paths)
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}

	// Python patterns - focus on write operations and system directory access
	patterns["python"] = []*regexp.Regexp{
		// Write operations with variables
		regexp.MustCompile(`open\s*\(\s*[a-zA-Z_]\w*\s*,\s*["']w`),
		regexp.MustCompile(`open\s*\(\s*[a-zA-Z_]\w*\s*,\s*["']a`),
		regexp.MustCompile(`Path\s*\(\s*[a-zA-Z_]\w*\s*\)\.write_`),
		regexp.MustCompile(`Path\s*\(\s*[a-zA-Z_]\w*\s*\)[.]open\s*\(\s*["']w`),
		// Dangerous directory operations
		regexp.MustCompile(`os\.makedirs\s*\(\s*[a-zA-Z_]\w*\s*[,)]`),
		regexp.MustCompile(`shutil\.(rmtree|move)\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}

	// Node/JavaScript patterns
	patterns["node"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`fs\.writeFile\s*\(\s*[a-zA-Z_$]\w*\s*[,+]`),
		regexp.MustCompile(`fs\.writeFileSync\s*\(\s*[a-zA-Z_$]\w*\s*[,+]`),
		regexp.MustCompile(`fs\.createWriteStream\s*\(\s*[a-zA-Z_$]\w*\s*[,)]`),
		regexp.MustCompile(`fs\.(unlink|rename|mkdir|rmdir)\s*\(\s*[a-zA-Z_$]\w*\s*[,)]`),
		// fs-extra write operations
		regexp.MustCompile(`fs-extra\.(copy|move|remove|emptyDir)`),
		regexp.MustCompile(`fs\.(copy|move|remove)\s*\(`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}

	// TypeScript patterns (same as Node)
	patterns["typescript"] = patterns["node"]

	// Java patterns
	patterns["java"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`FileWriter\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`FileOutputStream\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`Files\.write\s*\(\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`Files\.copy\s*\(\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`Files\.move\s*\(\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`Files\.delete\s*\(`),
		regexp.MustCompile(`\.createNewFile\s*\(`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}

	// Ruby patterns
	patterns["ruby"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`File\.write\s*\(\s*[a-zA-Z_]\w*\s*[,+]`),
		regexp.MustCompile(`File\.open\s*\(\s*[a-zA-Z_]\w*\s*,\s*["']w`),
		regexp.MustCompile(`File\.open\s*\(\s*["'][^"']*["']\s*,\s*["']w`),
		regexp.MustCompile(`FileUtils\.(rm|rmtree|mv|cp_r)\s*\(`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}

	// PHP patterns
	patterns["php"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`file_put_contents\s*\(\s*\$?[a-zA-Z_]\w*\s*[,)]`),
		regexp.MustCompile(`fwrite\s*\(`),
		regexp.MustCompile(`unlink\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`rename\s*\(\s*\$[a-zA-Z_]\w*\s*[,)]`),
		regexp.MustCompile(`mkdir\s*\(\s*\$[a-zA-Z_]\w*\s*[,)]`),
		regexp.MustCompile(`rmdir\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}

	// C/C++ patterns
	patterns["c"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`fopen\s*\(\s*[a-zA-Z_]\w*\s*,\s*["']w`),
		regexp.MustCompile(`fopen\s*\(\s*[a-zA-Z_]\w*\s*,\s*["']a`),
		regexp.MustCompile(`remove\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`unlink\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`rename\s*\(\s*[a-zA-Z_]\w*\s*[,)]`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}
	patterns["cpp"] = patterns["c"]

	// Rust patterns
	patterns["rust"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`File::create\s*\(\s*&?[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`fs::write\s*\(\s*&?[a-zA-Z_]\w*`),
		regexp.MustCompile(`fs::remove_file\s*\(`),
		regexp.MustCompile(`fs::remove_dir\s*\(`),
		regexp.MustCompile(`fs::rename\s*\(`),
		regexp.MustCompile(`fs::copy\s*\(`),
		regexp.MustCompile(`OpenOptions::new\(\).*\.write\(true\)`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
	}

	// Swift patterns
	patterns["swift"] = []*regexp.Regexp{
		// Write operations
		regexp.MustCompile(`FileManager.*createFile\s*\(`),
		regexp.MustCompile(`FileManager.*removeItem\s*\(`),
		regexp.MustCompile(`FileManager.*moveItem\s*\(`),
		regexp.MustCompile(`FileManager.*copyItem\s*\(`),
		regexp.MustCompile(`FileHandle.*write\s*\(`),
		regexp.MustCompile(`\.write\s*\(\s*to:\s*[a-zA-Z_]\w*`),
		// System directory writes
		regexp.MustCompile(`[\"']\/(etc|var|tmp|usr|bin|home|root|sys|proc)`),
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

	relPath, err := filepath.Rel(root, filePath)
	if err != nil {
		relPath = filepath.Base(filePath)
	}
	relPath = filepath.ToSlash(relPath)

	safeVars := make(map[string]struct{})

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if language == "go" {
			c.updateGoSafeVars(line, safeVars)
		}

		if isCommentLine(line, language) {
			continue
		}

		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				match := pattern.FindString(line)

				if language == "go" && c.isGoSafeUsage(line, safeVars) {
					break
				}

				if c.isSystemPathMatch(match) {
					finding := Finding{
						Type:        "filesystem",
						File:        relPath,
						Line:        lineNum,
						Description: fmt.Sprintf("writes to system directory: %s", c.sanitizeMatch(match)),
						Severity:    "high",
					}
					if !c.isAllowedFinding(finding, line) {
						findings = append(findings, finding)
					}
					break
				}

				finding := Finding{
					Type:        "filesystem",
					File:        relPath,
					Line:        lineNum,
					Description: fmt.Sprintf("potentially unsafe file write: %s", c.sanitizeMatch(match)),
					Severity:    "medium",
				}
				if !c.isAllowedFinding(finding, line) {
					findings = append(findings, finding)
				}
				break
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

func (c *FileSystemCheck) updateGoSafeVars(line string, safeVars map[string]struct{}) {
	if m := goSafeJoinRE.FindStringSubmatch(line); len(m) == 2 {
		safeVars[m[1]] = struct{}{}
		return
	}
	if m := goSafeJoinSingleRE.FindStringSubmatch(line); len(m) == 2 {
		safeVars[m[1]] = struct{}{}
		return
	}
	if m := goJoinAssignRE.FindStringSubmatch(line); len(m) == 3 {
		target := m[1]
		base := m[2]
		if _, ok := safeVars[base]; ok {
			safeVars[target] = struct{}{}
			return
		}
		if _, ok := goRootVarCandidates[base]; ok {
			safeVars[target] = struct{}{}
			return
		}
	}
}

func (c *FileSystemCheck) isGoSafeUsage(line string, safeVars map[string]struct{}) bool {
	m := goFileOpVarRE.FindStringSubmatch(line)
	if len(m) != 2 {
		return false
	}
	_, ok := safeVars[m[1]]
	return ok
}

// isSystemPathMatch checks if the matched string indicates a system directory path.
func (c *FileSystemCheck) isSystemPathMatch(match string) bool {
	return strings.Contains(match, "/etc") ||
		strings.Contains(match, "/var") ||
		strings.Contains(match, "/tmp") ||
		strings.Contains(match, "/usr") ||
		strings.Contains(match, "/bin") ||
		strings.Contains(match, "/home") ||
		strings.Contains(match, "/root") ||
		strings.Contains(match, "/sys") ||
		strings.Contains(match, "/proc")
}

func (c *FileSystemCheck) compileAllowRules() []allowRule {
	if len(c.Allowlist) == 0 {
		return nil
	}
	var rules []allowRule
	for _, raw := range c.Allowlist {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) == 1 {
			rules = append(rules, allowRule{
				filePattern: filepath.ToSlash(strings.TrimSpace(parts[0])),
			})
			continue
		}
		filePattern := filepath.ToSlash(strings.TrimSpace(parts[0]))
		matchPart := strings.TrimSpace(parts[1])
		if filePattern == "" {
			filePattern = "*"
		}
		if matchPart == "" {
			rules = append(rules, allowRule{filePattern: filePattern})
			continue
		}
		if isAllDigits(matchPart) {
			rules = append(rules, allowRule{
				filePattern: filePattern,
				line:        atoiSafe(matchPart),
			})
			continue
		}
		rules = append(rules, allowRule{
			filePattern: filePattern,
			matchText:   matchPart,
		})
	}
	return rules
}

func (c *FileSystemCheck) isAllowedFinding(f Finding, line string) bool {
	if len(c.allowRules) == 0 {
		return false
	}
	for _, rule := range c.allowRules {
		filePattern := rule.filePattern
		if filePattern == "" {
			filePattern = "*"
		}
		if !wildcardMatch(f.File, filePattern) {
			continue
		}
		if rule.line > 0 {
			if f.Line == rule.line {
				return true
			}
			continue
		}
		if rule.matchText == "" {
			return true
		}
		if matchText(line, f, rule.matchText) {
			return true
		}
	}
	return false
}

func matchText(line string, f Finding, pattern string) bool {
	if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
		return wildcardMatch(line, pattern) || wildcardMatch(f.Description, pattern) || wildcardMatch(f.String(), pattern)
	}
	return strings.Contains(line, pattern) || strings.Contains(f.Description, pattern) || strings.Contains(f.String(), pattern)
}

func wildcardMatch(value, pattern string) bool {
	if pattern == "" {
		return value == ""
	}
	if pattern == "*" {
		return true
	}
	var b strings.Builder
	b.WriteString("^")
	for _, r := range pattern {
		switch r {
		case '*':
			b.WriteString(".*")
		case '?':
			b.WriteString(".")
		case '.', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			b.WriteString("\\")
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile(b.String())
	if err != nil {
		return false
	}
	return re.MatchString(value)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func atoiSafe(s string) int {
	n := 0
	for _, r := range s {
		n = n*10 + int(r-'0')
	}
	return n
}

// formatFindings formats findings into a readable message.
func (c *FileSystemCheck) formatFindings(findings []Finding) string {
	if len(findings) == 1 {
		return fmt.Sprintf("Unsafe file write detected:\n- %s", findings[0].String())
	}

	return fmt.Sprintf("%d unsafe file writes detected:\n%s", len(findings),
		c.joinFindingsLines(findings))
}

// joinFindingsLines joins findings with one entry per line.
func (c *FileSystemCheck) joinFindingsLines(findings []Finding) string {
	var b strings.Builder
	for i, f := range findings {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("- ")
		b.WriteString(f.String())
	}
	return b.String()
}
