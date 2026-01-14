package tools

import (
	"github.com/ipedrazas/a2/pkg/checker"
)

// Tool represents an external tool that checks may depend on.
type Tool struct {
	Name        string           // Tool name (e.g., "govulncheck")
	Description string           // What the tool does
	CheckCmd    []string         // Command to check if installed (e.g., ["govulncheck", "--version"])
	Language    checker.Language // Which language this tool is for
	CheckIDs    []string         // Which check IDs use this tool
	Required    bool             // If true, check will fail without it; if false, check will skip
	Install     InstallCommands  // Platform-specific install commands
}

// InstallCommands contains install commands for different platforms/methods.
type InstallCommands struct {
	Go     string // go install command
	Brew   string // macOS Homebrew
	Apt    string // Debian/Ubuntu
	Dnf    string // Fedora/RHEL
	Cargo  string // Rust cargo install
	Npm    string // npm install -g
	Pip    string // pip install
	Manual string // Manual instructions or download URL
}

// Registry returns all known external tools organized by language.
func Registry() []Tool {
	return []Tool{
		// Go tools
		{
			Name:        "govulncheck",
			Description: "Go vulnerability scanner",
			CheckCmd:    []string{"govulncheck", "--version"},
			Language:    checker.LangGo,
			CheckIDs:    []string{"go:deps"},
			Required:    false,
			Install: InstallCommands{
				Go: "go install golang.org/x/vuln/cmd/govulncheck@latest",
			},
		},
		{
			Name:        "gocyclo",
			Description: "Go cyclomatic complexity analyzer",
			CheckCmd:    []string{"gocyclo"},
			Language:    checker.LangGo,
			CheckIDs:    []string{"go:cyclomatic"},
			Required:    false,
			Install: InstallCommands{
				Go: "go install github.com/fzipp/gocyclo/cmd/gocyclo@latest",
			},
		},

		// Python tools
		{
			Name:        "pytest",
			Description: "Python testing framework",
			CheckCmd:    []string{"pytest", "--version"},
			Language:    checker.LangPython,
			CheckIDs:    []string{"python:tests", "python:coverage"},
			Required:    false,
			Install: InstallCommands{
				Pip: "pip install pytest",
			},
		},
		{
			Name:        "ruff",
			Description: "Fast Python linter and formatter",
			CheckCmd:    []string{"ruff", "--version"},
			Language:    checker.LangPython,
			CheckIDs:    []string{"python:lint", "python:format"},
			Required:    false,
			Install: InstallCommands{
				Pip:  "pip install ruff",
				Brew: "brew install ruff",
			},
		},
		{
			Name:        "black",
			Description: "Python code formatter",
			CheckCmd:    []string{"black", "--version"},
			Language:    checker.LangPython,
			CheckIDs:    []string{"python:format"},
			Required:    false,
			Install: InstallCommands{
				Pip: "pip install black",
			},
		},
		{
			Name:        "mypy",
			Description: "Python static type checker",
			CheckCmd:    []string{"mypy", "--version"},
			Language:    checker.LangPython,
			CheckIDs:    []string{"python:type"},
			Required:    false,
			Install: InstallCommands{
				Pip: "pip install mypy",
			},
		},
		{
			Name:        "pip-audit",
			Description: "Python dependency vulnerability scanner",
			CheckCmd:    []string{"pip-audit", "--version"},
			Language:    checker.LangPython,
			CheckIDs:    []string{"python:deps"},
			Required:    false,
			Install: InstallCommands{
				Pip: "pip install pip-audit",
			},
		},
		{
			Name:        "radon",
			Description: "Python code complexity analyzer",
			CheckCmd:    []string{"radon", "--version"},
			Language:    checker.LangPython,
			CheckIDs:    []string{"python:complexity"},
			Required:    false,
			Install: InstallCommands{
				Pip: "pip install radon",
			},
		},

		// Node.js tools
		{
			Name:        "eslint",
			Description: "JavaScript/TypeScript linter",
			CheckCmd:    []string{"eslint", "--version"},
			Language:    checker.LangNode,
			CheckIDs:    []string{"node:lint", "typescript:lint"},
			Required:    false,
			Install: InstallCommands{
				Npm: "npm install -g eslint",
			},
		},
		{
			Name:        "prettier",
			Description: "Code formatter",
			CheckCmd:    []string{"prettier", "--version"},
			Language:    checker.LangNode,
			CheckIDs:    []string{"node:format", "typescript:format"},
			Required:    false,
			Install: InstallCommands{
				Npm: "npm install -g prettier",
			},
		},
		{
			Name:        "biome",
			Description: "Fast linter and formatter for JS/TS",
			CheckCmd:    []string{"biome", "--version"},
			Language:    checker.LangNode,
			CheckIDs:    []string{"node:lint", "node:format", "typescript:lint", "typescript:format"},
			Required:    false,
			Install: InstallCommands{
				Npm:  "npm install -g @biomejs/biome",
				Brew: "brew install biome",
			},
		},

		// Rust tools
		{
			Name:        "cargo-audit",
			Description: "Rust dependency vulnerability scanner",
			CheckCmd:    []string{"cargo", "audit", "--version"},
			Language:    checker.LangRust,
			CheckIDs:    []string{"rust:deps"},
			Required:    false,
			Install: InstallCommands{
				Cargo: "cargo install cargo-audit",
			},
		},
		{
			Name:        "cargo-tarpaulin",
			Description: "Rust code coverage tool",
			CheckCmd:    []string{"cargo", "tarpaulin", "--version"},
			Language:    checker.LangRust,
			CheckIDs:    []string{"rust:coverage"},
			Required:    false,
			Install: InstallCommands{
				Cargo: "cargo install cargo-tarpaulin",
			},
		},

		// Swift tools
		{
			Name:        "swiftlint",
			Description: "Swift linter",
			CheckCmd:    []string{"swiftlint", "version"},
			Language:    checker.LangSwift,
			CheckIDs:    []string{"swift:lint"},
			Required:    false,
			Install: InstallCommands{
				Brew:   "brew install swiftlint",
				Manual: "https://github.com/realm/SwiftLint#installation",
			},
		},
		{
			Name:        "swift-format",
			Description: "Swift code formatter",
			CheckCmd:    []string{"swift-format", "--version"},
			Language:    checker.LangSwift,
			CheckIDs:    []string{"swift:format"},
			Required:    false,
			Install: InstallCommands{
				Brew:   "brew install swift-format",
				Manual: "https://github.com/apple/swift-format#getting-swift-format",
			},
		},

		// Common tools
		{
			Name:        "gitleaks",
			Description: "Secret scanner for git repos",
			CheckCmd:    []string{"gitleaks", "version"},
			Language:    checker.LangCommon,
			CheckIDs:    []string{"common:secrets"},
			Required:    false,
			Install: InstallCommands{
				Brew:   "brew install gitleaks",
				Go:     "go install github.com/gitleaks/gitleaks/v8@latest",
				Manual: "https://github.com/gitleaks/gitleaks#installing",
			},
		},
		{
			Name:        "semgrep",
			Description: "Static analysis security scanner",
			CheckCmd:    []string{"semgrep", "--version"},
			Language:    checker.LangCommon,
			CheckIDs:    []string{"common:sast"},
			Required:    false,
			Install: InstallCommands{
				Pip:    "pip install semgrep",
				Brew:   "brew install semgrep",
				Manual: "https://semgrep.dev/docs/getting-started/",
			},
		},
		{
			Name:        "trivy",
			Description: "Vulnerability scanner for containers and code",
			CheckCmd:    []string{"trivy", "--version"},
			Language:    checker.LangCommon,
			CheckIDs:    []string{"common:sast", "common:dockerfile"},
			Required:    false,
			Install: InstallCommands{
				Brew:   "brew install trivy",
				Apt:    "apt install trivy",
				Manual: "https://aquasecurity.github.io/trivy/latest/getting-started/installation/",
			},
		},
	}
}

// ByLanguage returns tools filtered by language.
func ByLanguage(lang checker.Language) []Tool {
	var result []Tool
	for _, t := range Registry() {
		if t.Language == lang {
			result = append(result, t)
		}
	}
	return result
}

// ForLanguages returns tools for the given languages plus common tools.
func ForLanguages(langs []checker.Language) []Tool {
	langSet := make(map[checker.Language]bool)
	for _, l := range langs {
		langSet[l] = true
	}
	// Always include common tools
	langSet[checker.LangCommon] = true

	var result []Tool
	for _, t := range Registry() {
		if langSet[t.Language] {
			result = append(result, t)
		}
	}
	return result
}

// ByName returns a tool by its name, or nil if not found.
func ByName(name string) *Tool {
	for _, t := range Registry() {
		if t.Name == name {
			return &t
		}
	}
	return nil
}
