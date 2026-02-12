package common

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all common (language-agnostic) check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	registrations := []checker.CheckRegistration{
		{
			Checker: &FileExistsCheck{Files: cfg.Files.Required},
			Meta: checker.CheckMeta{
				ID:          "file_exists",
				Name:        "Required Files",
				Description: "Checks for required documentation files like README.md.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       900,
				Suggestion:  "Add missing documentation files (README.md, etc.)",
			},
		},
		{
			Checker: &DockerfileCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:dockerfile",
				Name:        "Container Ready",
				Description: "Verifies a Dockerfile exists and runs trivy security scanning on it.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       910,
				Suggestion:  "Add Dockerfile for containerization",
			},
		},
		{
			Checker: &CICheck{},
			Meta: checker.CheckMeta{
				ID:          "common:ci",
				Name:        "CI Pipeline",
				Description: "Checks for CI/CD pipeline configuration (GitHub Actions, GitLab CI, etc.).",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       920,
				Suggestion:  "Add CI pipeline configuration (.github/workflows, etc.)",
			},
		},
		{
			Checker: &HealthCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:health",
				Name:        "Health Endpoint",
				Description: "Checks for health check endpoint implementation for container orchestration.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       930,
				Suggestion:  "Add health check endpoint for production readiness",
			},
		},
		{
			Checker: &SecretsCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:secrets",
				Name:        "Secrets Detection",
				Description: "Scans code for accidentally committed secrets using gitleaks.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       940,
				Suggestion:  "Remove or secure detected secrets",
			},
		},
		{
			Checker: &EnvCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:env",
				Name:        "Environment Config",
				Description: "Checks for .env.example file to document required environment variables.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       945,
				Suggestion:  "Add .env.example for environment configuration",
			},
		},
		{
			Checker: &LicenseCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:license",
				Name:        "License Compliance",
				Description: "Verifies a LICENSE file exists for open source compliance.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       950,
				Suggestion:  "Add LICENSE file for license compliance",
			},
		},
		{
			Checker: &SASTCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:sast",
				Name:        "SAST Security Scanning",
				Description: "Runs static application security testing using semgrep.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       955,
				Suggestion:  "Fix security issues found by SAST scanning",
			},
		},
		{
			Checker: &APIDocsCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:api_docs",
				Name:        "API Documentation",
				Description: "Checks for API documentation (OpenAPI/Swagger specification).",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       960,
				Suggestion:  "Add API documentation (OpenAPI/Swagger)",
			},
		},
		{
			Checker: &ChangelogCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:changelog",
				Name:        "Changelog",
				Description: "Verifies a CHANGELOG.md file exists to track project changes.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       965,
				Suggestion:  "Add CHANGELOG.md to track changes",
			},
		},
		{
			Checker: &IntegrationCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:integration",
				Name:        "Integration Tests",
				Description: "Checks for integration test files in the project.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       980,
				Suggestion:  "Add integration tests",
			},
		},
		{
			Checker: &MetricsCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:metrics",
				Name:        "Metrics Instrumentation",
				Description: "Checks for metrics/observability instrumentation (Prometheus, OpenTelemetry).",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1010,
				Suggestion:  "Add metrics instrumentation (Prometheus, etc.)",
			},
		},
		{
			Checker: &ErrorsCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:errors",
				Name:        "Error Tracking",
				Description: "Checks for error tracking service integration (Sentry, Rollbar, etc.).",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1020,
				Suggestion:  "Add error tracking (Sentry, etc.)",
			},
		},
		{
			Checker: &PrecommitCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:precommit",
				Name:        "Pre-commit Hooks",
				Description: "Checks for pre-commit hook configuration to enforce code quality.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Optional:    true,
				Order:       1065,
				Suggestion:  "Add pre-commit hooks for code quality",
			},
		},
		{
			Checker: &K8sCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:k8s",
				Name:        "Kubernetes Ready",
				Description: "Checks for Kubernetes manifests or Helm charts for deployment.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1030,
				Suggestion:  "Add Kubernetes manifests for deployment",
			},
		},
		{
			Checker: &ShutdownCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:shutdown",
				Name:        "Graceful Shutdown",
				Description: "Checks for graceful shutdown signal handling in the codebase.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1035,
				Suggestion:  "Implement graceful shutdown handling",
			},
		},
		{
			Checker: &MigrationsCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:migrations",
				Name:        "Database Migrations",
				Description: "Checks for database migration files or migration tool configuration.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1040,
				Suggestion:  "Add database migration support",
			},
		},
		{
			Checker: &ConfigValidationCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:config_validation",
				Name:        "Config Validation",
				Description: "Checks for configuration validation implementation at startup.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1045,
				Suggestion:  "Add configuration validation",
			},
		},
		{
			Checker: &RetryCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:retry",
				Name:        "Retry Logic",
				Description: "Checks for retry/backoff logic implementation for external calls.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1050,
				Suggestion:  "Implement retry logic for external calls",
			},
		},
		{
			Checker: &ContributingCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:contributing",
				Name:        "Contributing Guidelines",
				Description: "Verifies a CONTRIBUTING.md file exists with contributor guidelines.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       970,
				Suggestion:  "Add CONTRIBUTING.md for contributor guidelines",
			},
		},
		{
			Checker: &EditorconfigCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:editorconfig",
				Name:        "Editor Config",
				Description: "Checks for .editorconfig file for consistent editor settings.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Optional:    true,
				Order:       1070,
				Suggestion:  "Add .editorconfig for consistent formatting",
			},
		},
		{
			Checker: &E2ECheck{},
			Meta: checker.CheckMeta{
				ID:          "common:e2e",
				Name:        "E2E Tests",
				Description: "Checks for end-to-end test files in the project.",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       985,
				Suggestion:  "Add end-to-end tests",
			},
		},
		{
			Checker: &TracingCheck{},
			Meta: checker.CheckMeta{
				ID:          "common:tracing",
				Name:        "Distributed Tracing",
				Description: "Checks for distributed tracing implementation (OpenTelemetry, Jaeger).",
				Languages:   []checker.Language{checker.LangCommon},
				Critical:    false,
				Order:       1015,
				Suggestion:  "Add distributed tracing support",
			},
		},
	}

	// Add external checks from config
	for _, ext := range cfg.External {
		registrations = append(registrations, checker.CheckRegistration{
			Checker: &ExternalCheck{
				CheckID:   ext.ID,
				CheckName: ext.Name,
				Command:   ext.Command,
				Args:      ext.Args,
				Severity:  ext.Severity,
				SourceDir: ext.SourceDir,
			},
			Meta: checker.CheckMeta{
				ID:        ext.ID,
				Name:      ext.Name,
				Languages: []checker.Language{checker.LangCommon},
				Critical:  ext.Severity == "fail",
				Order:     1000, // External checks run last
			},
		})
	}

	return registrations
}
