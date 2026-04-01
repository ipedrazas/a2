package common

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// NamingCheck verifies naming convention consistency across the project.
type NamingCheck struct{}

func (c *NamingCheck) ID() string   { return "common:naming" }
func (c *NamingCheck) Name() string { return "Naming Consistency" }

// namingConvention represents a detected naming style.
type namingConvention int

const (
	conventionSnakeCase  namingConvention = iota // my_file
	conventionCamelCase                          // myFile
	conventionPascalCase                         // MyFile
	conventionKebabCase                          // my-file
	conventionFlat                               // myfile (single word, no separator)
)

var (
	snakeCaseRe = regexp.MustCompile(`^[a-z][a-z0-9]*(_[a-z0-9]+)+$`)
	camelCaseRe = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)
	pascalRe    = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
	kebabCaseRe = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)+$`)
	flatRe      = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
)

// Run checks for naming convention consistency in file names.
func (c *NamingCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return rb.Fail(err.Error()), err
	}
	path = absPath

	// Priority 1: Check if naming convention enforcement is configured
	tools := c.findNamingTools(path)
	if len(tools) > 0 {
		return rb.Pass("Naming convention enforcement configured: " + strings.Join(tools, ", ")), nil
	}

	// Priority 2: Analyze file naming consistency
	conventions := c.analyzeFileNames(path)

	// Need a minimum number of classifiable files to make a meaningful assessment
	total := 0
	for _, count := range conventions {
		total += count
	}
	if total < 3 {
		return rb.Pass("Too few source files to assess naming consistency"), nil
	}

	// Find the dominant convention
	dominant, dominantCount := c.dominantConvention(conventions)
	dominantPct := float64(dominantCount) / float64(total) * 100

	if dominantPct >= 90 {
		return rb.Pass(fmt.Sprintf("Consistent file naming: %s (%.0f%% of files)", conventionName(dominant), dominantPct)), nil
	}

	if dominantPct >= 70 {
		mixed := c.mixedConventionSummary(conventions, dominant)
		return rb.Warn(fmt.Sprintf("Mostly %s file naming (%.0f%%), but also found: %s", conventionName(dominant), dominantPct, mixed)), nil
	}

	mixed := c.mixedConventionSummary(conventions, -1)
	return rb.Warn("Mixed file naming conventions detected: " + mixed), nil
}

// findNamingTools checks for naming convention enforcement tools.
func (c *NamingCheck) findNamingTools(path string) []string {
	var found []string

	// ESLint with naming-convention rule
	eslintConfigs := []string{
		".eslintrc.json", ".eslintrc.js", ".eslintrc.yml", ".eslintrc.yaml",
		".eslintrc", "eslint.config.js", "eslint.config.mjs", "eslint.config.ts",
	}
	for _, cfg := range eslintConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				if strings.Contains(string(content), "naming-convention") {
					found = append(found, "ESLint naming-convention")
					break
				}
			}
		}
	}

	// Pylint naming checks
	pylintConfigs := []string{".pylintrc", "pylintrc", "setup.cfg", "pyproject.toml"}
	for _, cfg := range pylintConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				contentStr := string(content)
				if strings.Contains(contentStr, "naming-style") || strings.Contains(contentStr, "naming-convention") {
					found = append(found, "Pylint naming")
					break
				}
			}
		}
	}

	// Rubocop naming
	rubocopConfigs := []string{".rubocop.yml", ".rubocop.yaml"}
	for _, cfg := range rubocopConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				if strings.Contains(string(content), "Naming/") {
					found = append(found, "RuboCop Naming")
					break
				}
			}
		}
	}

	// checkstyle (Java)
	if safepath.Exists(path, "checkstyle.xml") {
		if content, err := safepath.ReadFile(path, "checkstyle.xml"); err == nil {
			if strings.Contains(string(content), "NamingConvention") || strings.Contains(string(content), "MemberName") {
				found = append(found, "Checkstyle naming")
			}
		}
	}

	// Clippy (Rust) is always on for naming, so check for clippy config
	if safepath.Exists(path, "clippy.toml") || safepath.Exists(path, ".clippy.toml") {
		found = append(found, "Clippy")
	}

	return found
}

// analyzeFileNames walks the project and classifies source file names by naming convention.
func (c *NamingCheck) analyzeFileNames(path string) map[namingConvention]int {
	conventions := make(map[namingConvention]int)

	codeExtensions := map[string]bool{
		".go": true, ".py": true, ".js": true, ".ts": true,
		".jsx": true, ".tsx": true, ".java": true, ".rb": true,
		".rs": true, ".swift": true, ".kt": true, ".scala": true,
		".cs": true, ".cpp": true, ".c": true, ".h": true,
		".mjs": true, ".cjs": true, ".vue": true, ".svelte": true,
	}

	skipDirs := map[string]bool{
		"node_modules": true, "vendor": true, "__pycache__": true,
		"dist": true, "build": true, "target": true, ".git": true,
		"bin": true, "obj": true, "out": true,
	}

	_ = filepath.Walk(path, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		if !codeExtensions[ext] {
			return nil
		}

		// Get the file name without extension
		baseName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

		// Strip test suffixes/prefixes to classify the base naming style
		baseName = stripTestAffixes(baseName)

		if baseName == "" {
			return nil
		}

		conv := classifyName(baseName)
		if conv >= 0 {
			conventions[conv]++
		}

		return nil
	})

	return conventions
}

// stripTestAffixes removes common test-related suffixes and prefixes from a file name.
func stripTestAffixes(name string) string {
	// Handle Go/Python _test suffix
	name = strings.TrimSuffix(name, "_test")
	name = strings.TrimSuffix(name, "_Test")
	name = strings.TrimPrefix(name, "test_")
	name = strings.TrimPrefix(name, "Test_")

	// Handle JS/TS .test/.spec suffix (these appear as part of the name before the real extension)
	name = strings.TrimSuffix(name, ".test")
	name = strings.TrimSuffix(name, ".spec")
	name = strings.TrimSuffix(name, ".Test")
	name = strings.TrimSuffix(name, ".Spec")

	return name
}

// classifyName classifies a name into a naming convention.
// Returns -1 if the name doesn't match any known convention.
func classifyName(name string) namingConvention {
	switch {
	case snakeCaseRe.MatchString(name):
		return conventionSnakeCase
	case kebabCaseRe.MatchString(name):
		return conventionKebabCase
	case pascalRe.MatchString(name) && hasUpperAfterFirst(name):
		return conventionPascalCase
	case camelCaseRe.MatchString(name) && hasUpperAfterFirst(name):
		return conventionCamelCase
	case flatRe.MatchString(name):
		return conventionFlat
	default:
		return -1
	}
}

// hasUpperAfterFirst returns true if the name contains an uppercase letter after position 0.
// Used to distinguish multi-word PascalCase/camelCase from single-word flat names.
func hasUpperAfterFirst(name string) bool {
	for i := 1; i < len(name); i++ {
		if name[i] >= 'A' && name[i] <= 'Z' {
			return true
		}
	}
	return false
}

// dominantConvention returns the most common convention and its count.
func (c *NamingCheck) dominantConvention(conventions map[namingConvention]int) (namingConvention, int) {
	var dominant namingConvention
	maxCount := 0
	for conv, count := range conventions {
		if count > maxCount {
			maxCount = count
			dominant = conv
		}
	}
	return dominant, maxCount
}

// mixedConventionSummary returns a summary of conventions excluding the dominant one.
// If exclude is -1, all conventions are included.
func (c *NamingCheck) mixedConventionSummary(conventions map[namingConvention]int, exclude namingConvention) string {
	var parts []string
	for conv, count := range conventions {
		if conv == exclude {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s (%s)", conventionName(conv), checkutil.PluralizeCount(count, "file", "files")))
	}
	return strings.Join(parts, ", ")
}

// conventionName returns a human-readable name for a naming convention.
func conventionName(conv namingConvention) string {
	switch conv {
	case conventionSnakeCase:
		return "snake_case"
	case conventionCamelCase:
		return "camelCase"
	case conventionPascalCase:
		return "PascalCase"
	case conventionKebabCase:
		return "kebab-case"
	case conventionFlat:
		return "lowercase"
	default:
		return "unknown"
	}
}
