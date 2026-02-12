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

// ShellInjectionCheck detects dangerous shell/code execution patterns.
type ShellInjectionCheck struct {
	Patterns map[string][]*regexp.Regexp
}

// ID returns the unique identifier for this check.
func (c *ShellInjectionCheck) ID() string {
	return "security:shell_injection"
}

// Name returns the human-readable name for this check.
func (c *ShellInjectionCheck) Name() string {
	return "Shell Injection Detection"
}

// Run executes the shell injection detection check.
func (c *ShellInjectionCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Initialize patterns if not already done
	if c.Patterns == nil {
		c.Patterns = c.getPatterns()
	}

	// Scan all source files
	findings := c.scanDirectory(path)

	if len(findings) == 0 {
		return rb.Pass("No dangerous shell/code execution patterns detected"), nil
	}

	// Format findings into message
	msg := c.formatFindings(findings)
	return rb.Fail(msg), nil
}

// getPatterns returns language-specific patterns for detecting shell injection.
func (c *ShellInjectionCheck) getPatterns() map[string][]*regexp.Regexp {
	patterns := make(map[string][]*regexp.Regexp)

	// Go patterns
	patterns["go"] = []*regexp.Regexp{
		// exec.Command with string concatenation/formatting
		regexp.MustCompile("exec\\.Command\\s*\\(\\s*[\"']sh[\"'].*,\\s*[\"']-c"),
		regexp.MustCompile("exec\\.Command\\s*\\(\\s*[\"']bash[\"'].*,\\s*[\"']-c"),
		regexp.MustCompile(`exec\.Command\s*\(\s*\w+\s*\+`), // Variable concatenation
		regexp.MustCompile(`exec\.Command\s*\(\s*fmt\.Sprintf`),
		regexp.MustCompile(`exec\.Command\s*\(\s*strings\.Join`),
		// Direct command execution
		regexp.MustCompile(`exec\.Command\s*\(\s*[a-zA-Z_]\w*\s*,`), // Variable as command
		regexp.MustCompile(`\bexec\b\.\w*`),                         // Any exec package usage
		regexp.MustCompile(`os\.Command\s*\(`),                      // Deprecated os.Command
		regexp.MustCompile(`syscall\.Exec\s*\(`),                    // syscall.Exec
	}

	// Python patterns
	patterns["python"] = []*regexp.Regexp{
		// eval/exec with variables
		regexp.MustCompile(`eval\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`exec\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`execfile\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`compile\s*\(\s*[a-zA-Z_]\w*`),
		// subprocess with shell=True
		regexp.MustCompile(`subprocess\.(call|run|Popen|check_output)\([^)]*shell\s*=\s*True`),
		// os.system, os.popen
		regexp.MustCompile(`os\.(system|popen|spawn[lpe])\s*\(\s*[a-zA-Z_]\w*\s*`),
		// pickle/marshal with loads
		regexp.MustCompile(`(pickle|marshal|cPickle)\.loads\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`__import__\s*\(\s*[a-zA-Z_]\w*\s*\)`), // Dynamic imports
		// Commands.getoutput
		regexp.MustCompile(`commands\.getoutput\s*\(`),
	}

	// Node/JavaScript patterns
	patterns["node"] = []*regexp.Regexp{
		// eval with variables
		regexp.MustCompile(`eval\s*\(\s*[a-zA-Z_$]\w*\s*\)`),
		regexp.MustCompile(`eval\s*\(\s*["']` + `["'].*\+`), // String concatenation in eval
		regexp.MustCompile(`Function\s*\(\s*[a-zA-Z_$]\w*`),
		regexp.MustCompile(`new Function\s*\(`),
		// child_process with variables
		regexp.MustCompile(`child_process\.(exec|execSync|spawn)\s*\(\s*[a-zA-Z_$]\w*\s*,`),
		regexp.MustCompile(`require\s*\(\s*["']child_process["']`), // Check if imported
		// vm module
		regexp.MustCompile(`vm\.(runInThisContext|runInNewContext|runInContext|compileFunction|Script)\s*\(\s*[a-zA-Z_$]\w*`),
		// setTimeout/setInterval with string
		regexp.MustCompile(`(setTimeout|setInterval)\s*\(\s*["'][^"']+["']\s*\+`),
	}

	// TypeScript patterns (same as Node plus TypeScript-specific)
	patterns["typescript"] = patterns["node"]

	// Java patterns
	patterns["java"] = []*regexp.Regexp{
		// Runtime.exec with concatenation
		regexp.MustCompile(`Runtime\.getRuntime\(\)\.exec\s*\(\s*[a-zA-Z_]\w*\s*\+`),
		regexp.MustCompile(`Runtime\.getRuntime\(\)\.exec\s*\(\s*String\[`),
		// ProcessBuilder with variables
		regexp.MustCompile(`new\s+ProcessBuilder\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`ProcessBuilder\s*\(\s*[a-zA-Z_]\w*\.split\s*\(`),
		// ScriptEngine
		regexp.MustCompile(`ScriptEngine\.eval\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`ScriptEngineManager\.getEngineBy`),
	}

	// Rust patterns
	patterns["rust"] = []*regexp.Regexp{
		// Command::new with variable
		regexp.MustCompile(`Command::new\s*\(\s*&?[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`std::process::Command::new\s*\(\s*&?[a-zA-Z_]\w*\s*\)`),
		// .arg with user input
		regexp.MustCompile(`\.arg\s*\(\s*&?[a-zA-Z_]\w*\s*\)`),
		// libc::system
		regexp.MustCompile(`libc::system\s*\(`),
	}

	// Ruby patterns
	patterns["ruby"] = []*regexp.Regexp{
		// eval, class_eval, instance_eval with variables
		regexp.MustCompile(`(eval|class_eval|instance_eval|module_eval)\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		// system, exec, backtick with variables
		regexp.MustCompile(`system\s*\(\s*["'].*#\{[a-zA-Z_]\w*\}`),
		regexp.MustCompile(`exec\s*\(\s*["'].*#\{[a-zA-Z_]\w*\}`),
		regexp.MustCompile(`\s+` + "`" + `[a-zA-Z_]\w*\s*`), // Backtick execution
		// Open3 with variables
		regexp.MustCompile(`Open3\.(popen3|capture2|capture3)\s*\(\s*["'].*#\{`),
		// %x, %X execution
		regexp.MustCompile(`%x\s*\(\s*["'].*#\{`),
	}

	// PHP patterns
	patterns["php"] = []*regexp.Regexp{
		// eval with variables
		regexp.MustCompile(`eval\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		// system, exec, shell_exec, passthru with variables
		regexp.MustCompile(`(system|exec|shell_exec|passthru)\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		// backtick operator (needs line scanning)
		regexp.MustCompile(`proc_open\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		// preg_replace with /e modifier (eval)
		regexp.MustCompile(`preg_replace\s*\([^)]*/e\s*\)`),
	}

	// C/C++ patterns
	patterns["c"] = []*regexp.Regexp{
		regexp.MustCompile(`system\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`popen\s*\(\s*[a-zA-Z_]\w*\s*`),
		regexp.MustCompile(`exec[lv][pe]?\s*\(\s*[a-zA-Z_]\w*\s*`),
	}
	patterns["cpp"] = patterns["c"]

	// Swift patterns
	patterns["swift"] = []*regexp.Regexp{
		regexp.MustCompile(`Process\s*\(\s*executable:.*arguments:\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`NSPipe\s*\(\s*\)`),
		regexp.MustCompile(`NSTask\s*\(\s*\)`),
	}

	return patterns
}

// scanDirectory scans all files in the directory for shell injection patterns.
func (c *ShellInjectionCheck) scanDirectory(path string) []Finding {
	var findings []Finding

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			// Check if we should skip this directory
			if skipDirectories[info.Name()] {
				return filepath.SkipDir
			}
			// Skip hidden directories except .github
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

		// Scan the file
		fileFindings := c.scanFile(path, filePath, language, patterns)
		findings = append(findings, fileFindings...)

		// Limit findings to avoid excessive output
		if len(findings) >= 50 {
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		// Log error but continue
		_ = err
	}

	return findings
}

// scanFile scans a single file for shell injection patterns.
func (c *ShellInjectionCheck) scanFile(root, filePath string, language string, patterns []*regexp.Regexp) []Finding {
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
				// Extract the function name for better reporting
				match := pattern.FindString(line)
				findings = append(findings, Finding{
					Type:        "shell_injection",
					File:        relPath,
					Line:        lineNum,
					Description: fmt.Sprintf("%s detected (%s)", language, match),
					Severity:    "critical",
				})
				break // One finding per line
			}
		}
	}

	return findings
}

// formatFindings formats findings into a readable message.
func (c *ShellInjectionCheck) formatFindings(findings []Finding) string {
	if len(findings) == 1 {
		return fmt.Sprintf("Dangerous pattern detected: %s", findings[0].String())
	}

	if len(findings) <= 3 {
		return fmt.Sprintf("Dangerous patterns detected: %s", c.joinFindings(findings))
	}

	return fmt.Sprintf("%d dangerous patterns detected (e.g., %s)", len(findings),
		c.joinFindings(findings[:3]))
}

// joinFindings joins findings for display.
func (c *ShellInjectionCheck) joinFindings(findings []Finding) string {
	var strs []string
	for _, f := range findings {
		strs = append(strs, f.String())
	}
	return strings.Join(strs, ", ")
}
