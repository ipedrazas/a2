// Package security provides safety/security checks for detecting potential
// malicious code patterns.
package security

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all security check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	registrations := []checker.CheckRegistration{
		{
			Checker: &ObfuscationCheck{},
			Meta: checker.CheckMeta{
				ID:          "security:obfuscation",
				Name:        "Code Obfuscation Detection",
				Description: "Detects obfuscated code and encoded strings that may indicate malicious intent.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    true, // Obfuscation is suspicious
				Order:       40,   // Run very early
				Suggestion:  "Review obfuscated code for malicious intent and clarity",
			},
		},
		{
			Checker: &ShellInjectionCheck{},
			Meta: checker.CheckMeta{
				ID:          "security:shell_injection",
				Name:        "Shell Injection Detection",
				Description: "Detects dangerous shell/code execution patterns that could lead to command injection.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    true, // Shell injection is dangerous
				Order:       50,   // Run early (before build/test)
				Suggestion:  "Review and sanitize user input before executing",
			},
		},
		{
			Checker: &FileSystemCheck{
				Allowlist: cfg.Security.Filesystem.Allow,
			},
			Meta: checker.CheckMeta{
				ID:          "security:filesystem",
				Name:        "File System Safety",
				Description: "Detects path traversal and unsafe file operations that could access files outside the project.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    true, // File system abuse is dangerous
				Order:       55,   // Run early
				Suggestion:  "Validate file paths and restrict to project directory",
			},
		},
		{
			Checker: &NetworkCheck{},
			Meta: checker.CheckMeta{
				ID:          "security:network",
				Name:        "Network Exfiltration Detection",
				Description: "Detects suspicious network operations that could indicate data exfiltration.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false, // Network might be legit (fetch, API calls)
				Order:       60,    // Run after other security checks
				Suggestion:  "Review network endpoints and data being transmitted",
			},
		},
	}

	return registrations
}
