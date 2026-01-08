package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// EnvCheck validates environment variable handling practices.
type EnvCheck struct{}

func (c *EnvCheck) ID() string   { return "common:env" }
func (c *EnvCheck) Name() string { return "Environment Config" }

// Run checks for proper environment variable handling.
func (c *EnvCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var findings []string
	var issues []string

	// Check for .env.example or .env.sample (documents required vars)
	envTemplates := []string{".env.example", ".env.sample", ".env.template", "example.env", ".env.local.example"}
	for _, template := range envTemplates {
		if safepath.Exists(path, template) {
			findings = append(findings, template)
		}
	}

	// Check for dotenv library usage in different languages
	dotenvLib := c.hasDotenvLibrary(path)
	if dotenvLib != "" {
		findings = append(findings, dotenvLib)
	}

	// Check if .env is in .gitignore
	envInGitignore := c.isEnvInGitignore(path)
	if envInGitignore {
		findings = append(findings, ".env in .gitignore")
	}

	// Check for potential issues
	// Issue: .env file exists in repo (should be gitignored)
	if safepath.Exists(path, ".env") && !envInGitignore {
		issues = append(issues, ".env file exists but not in .gitignore")
	}

	// Build result
	if len(findings) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Environment config: " + strings.Join(findings, ", ")
		if len(issues) > 0 {
			result.Passed = false
			result.Status = checker.Warn
			result.Message += " (warning: " + strings.Join(issues, "; ") + ")"
		}
	} else if len(issues) > 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Environment issues: " + strings.Join(issues, "; ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No environment configuration found (add .env.example to document required vars)"
	}

	return result, nil
}

// hasDotenvLibrary checks if dotenv library is used in the project.
func (c *EnvCheck) hasDotenvLibrary(path string) string {
	// Check Go
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goDotenvLibs := []string{
				"github.com/joho/godotenv",
				"github.com/caarlos0/env",
				"github.com/kelseyhightower/envconfig",
				"github.com/spf13/viper",
			}
			for _, lib := range goDotenvLibs {
				if strings.Contains(string(content), lib) {
					return "Go dotenv"
				}
			}
		}
	}

	// Check Python
	pythonConfigs := []string{"pyproject.toml", "requirements.txt", "setup.py", "Pipfile"}
	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				contentLower := strings.ToLower(string(content))
				pythonDotenvLibs := []string{
					"python-dotenv",
					"environs",
					"pydantic-settings",
					"django-environ",
					"dynaconf",
				}
				for _, lib := range pythonDotenvLibs {
					if strings.Contains(contentLower, lib) {
						return "Python dotenv"
					}
				}
			}
		}
	}

	// Check Node.js
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			nodeDotenvLibs := []string{
				"dotenv",
				"dotenv-safe",
				"dotenv-expand",
				"env-cmd",
				"cross-env",
			}
			for _, lib := range nodeDotenvLibs {
				if strings.Contains(string(content), "\""+lib+"\"") {
					return "Node.js dotenv"
				}
			}
		}
	}

	// Check Java
	javaConfigs := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, cfg := range javaConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				contentLower := strings.ToLower(string(content))
				// Check for Spring Boot's config properties or dotenv-java
				if strings.Contains(contentLower, "spring-boot") ||
					strings.Contains(contentLower, "dotenv-java") ||
					strings.Contains(contentLower, "io.github.cdimascio") {
					return "Java config"
				}
			}
		}
	}

	// Check for Spring Boot application properties
	if safepath.Exists(path, "src/main/resources/application.properties") ||
		safepath.Exists(path, "src/main/resources/application.yml") ||
		safepath.Exists(path, "src/main/resources/application.yaml") {
		return "Spring Boot config"
	}

	return ""
}

// isEnvInGitignore checks if .env is listed in .gitignore.
func (c *EnvCheck) isEnvInGitignore(path string) bool {
	if !safepath.Exists(path, ".gitignore") {
		return false
	}

	content, err := safepath.ReadFile(path, ".gitignore")
	if err != nil {
		return false
	}

	// Check each line for .env patterns
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Check for common .env ignore patterns
		if line == ".env" ||
			line == ".env*" ||
			line == ".env.*" ||
			line == "*.env" ||
			line == ".env.local" ||
			strings.HasPrefix(line, ".env") && !strings.Contains(line, "example") && !strings.Contains(line, "sample") {
			return true
		}
	}

	return false
}
