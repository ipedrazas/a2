package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ConfigValidationCheck verifies configuration is validated at startup.
type ConfigValidationCheck struct{}

func (c *ConfigValidationCheck) ID() string   { return "common:config_validation" }
func (c *ConfigValidationCheck) Name() string { return "Config Validation" }

// Run checks for configuration validation libraries and patterns.
func (c *ConfigValidationCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var found []string

	// Check Go dependencies for config validation
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goConfigLibs := map[string]string{
				"github.com/go-playground/validator":   "validator",
				"github.com/kelseyhightower/envconfig": "envconfig",
				"github.com/spf13/viper":               "Viper",
				"github.com/knadh/koanf":               "Koanf",
				"github.com/caarlos0/env":              "env",
				"github.com/joho/godotenv":             "godotenv",
				"github.com/hashicorp/hcl":             "HCL",
			}
			for dep, name := range goConfigLibs {
				if strings.Contains(string(content), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check Node.js dependencies
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			nodeConfigLibs := map[string]string{
				"joi":             "Joi",
				"zod":             "Zod",
				"yup":             "Yup",
				"convict":         "Convict",
				"@hapi/joi":       "Joi",
				"ajv":             "AJV",
				"class-validator": "class-validator",
				"env-var":         "env-var",
				"dotenv":          "dotenv",
				"envalid":         "Envalid",
				"@nestjs/config":  "NestJS Config",
			}
			for dep, name := range nodeConfigLibs {
				if strings.Contains(string(content), `"`+dep+`"`) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check Python dependencies
	pythonFiles := []string{"pyproject.toml", "requirements.txt", "setup.py"}
	for _, file := range pythonFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			pythonConfigLibs := map[string]string{
				"pydantic":          "Pydantic",
				"pydantic-settings": "Pydantic Settings",
				"dynaconf":          "Dynaconf",
				"python-dotenv":     "python-dotenv",
				"environs":          "Environs",
				"cerberus":          "Cerberus",
				"marshmallow":       "Marshmallow",
				"attrs":             "attrs",
			}
			for dep, name := range pythonConfigLibs {
				if strings.Contains(strings.ToLower(string(content)), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check Java dependencies
	javaFiles := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, file := range javaFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			javaConfigLibs := map[string]string{
				"spring-boot-configuration-processor": "Spring Boot Config",
				"hibernate-validator":                 "Hibernate Validator",
				"jakarta.validation":                  "Jakarta Validation",
				"javax.validation":                    "Bean Validation",
				"typesafe.config":                     "Typesafe Config",
				"owner":                               "Owner",
			}
			for dep, name := range javaConfigLibs {
				if strings.Contains(string(content), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check Rust dependencies
	if safepath.Exists(path, "Cargo.toml") {
		if content, err := safepath.ReadFile(path, "Cargo.toml"); err == nil {
			rustConfigLibs := map[string]string{
				"config":    "config-rs",
				"figment":   "Figment",
				"serde":     "Serde",
				"validator": "validator",
				"clap":      "Clap",
				"dotenv":    "dotenv",
				"envy":      "envy",
			}
			for dep, name := range rustConfigLibs {
				if strings.Contains(string(content), `"`+dep+`"`) || strings.Contains(string(content), dep+" =") {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check TypeScript-specific config
	if safepath.Exists(path, "tsconfig.json") {
		// Check for TypeScript strict mode (indicates type validation)
		if content, err := safepath.ReadFile(path, "tsconfig.json"); err == nil {
			if strings.Contains(string(content), `"strict": true`) ||
				strings.Contains(string(content), `"strict":true`) {
				if !containsString(found, "TypeScript strict") {
					found = append(found, "TypeScript strict")
				}
			}
		}
	}

	// Build result
	if len(found) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Config validation: " + strings.Join(found, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No config validation found (consider adding validation at startup)"
	}

	return result, nil
}
