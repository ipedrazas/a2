package security

import (
	"path/filepath"
	"regexp"
	"strings"
)

// allowRule is a single parsed allowlist entry.
type allowRule struct {
	filePattern string
	line        int
	matchText   string
}

// allowlist suppresses known-safe security findings. Every scan check shares the
// same rule grammar:
//
//   - "file"            — suppress all findings in matching files
//   - "file:line"       — suppress the finding on that exact line
//   - "file:match"      — suppress findings whose line/description contains match
//
// The file part and match part both support "*" and "?" wildcards.
type allowlist struct {
	rules []allowRule
}

// newAllowlist parses raw allowlist entries into a reusable matcher. A zero
// allowlist (no rules) suppresses nothing.
func newAllowlist(raw []string) allowlist {
	return allowlist{rules: compileAllowRules(raw)}
}

// empty reports whether the allowlist has no rules. Used to drive lazy compile.
func (a allowlist) empty() bool {
	return len(a.rules) == 0
}

// allows reports whether the finding (on the given source line) matches any
// allowlist rule and should therefore be suppressed.
func (a allowlist) allows(f Finding, line string) bool {
	for _, rule := range a.rules {
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

// compileAllowRules parses raw "file", "file:line", and "file:match" entries.
func compileAllowRules(raw []string) []allowRule {
	if len(raw) == 0 {
		return nil
	}
	var rules []allowRule
	for _, entry := range raw {
		trimmed := strings.TrimSpace(entry)
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
