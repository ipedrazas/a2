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
				ID:        "file_exists",
				Name:      "Required Files",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     900,
			},
		},
		{
			Checker: &DockerfileCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:dockerfile",
				Name:      "Container Ready",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     910,
			},
		},
		{
			Checker: &CICheck{},
			Meta: checker.CheckMeta{
				ID:        "common:ci",
				Name:      "CI Pipeline",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     920,
			},
		},
		{
			Checker: &HealthCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:health",
				Name:      "Health Endpoint",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     930,
			},
		},
		{
			Checker: &SecretsCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:secrets",
				Name:      "Secrets Detection",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     940,
			},
		},
		{
			Checker: &EnvCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:env",
				Name:      "Environment Config",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     945,
			},
		},
		{
			Checker: &LicenseCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:license",
				Name:      "License Compliance",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     950,
			},
		},
		{
			Checker: &SASTCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:sast",
				Name:      "SAST Security Scanning",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     955,
			},
		},
		{
			Checker: &APIDocsCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:api_docs",
				Name:      "API Documentation",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     960,
			},
		},
		{
			Checker: &ChangelogCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:changelog",
				Name:      "Changelog",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     965,
			},
		},
		{
			Checker: &IntegrationCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:integration",
				Name:      "Integration Tests",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     980,
			},
		},
		{
			Checker: &MetricsCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:metrics",
				Name:      "Metrics Instrumentation",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1010,
			},
		},
		{
			Checker: &ErrorsCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:errors",
				Name:      "Error Tracking",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1020,
			},
		},
		{
			Checker: &PrecommitCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:precommit",
				Name:      "Pre-commit Hooks",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1065,
			},
		},
		{
			Checker: &K8sCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:k8s",
				Name:      "Kubernetes Ready",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1030,
			},
		},
		{
			Checker: &ShutdownCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:shutdown",
				Name:      "Graceful Shutdown",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1035,
			},
		},
		{
			Checker: &MigrationsCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:migrations",
				Name:      "Database Migrations",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1040,
			},
		},
		{
			Checker: &ConfigValidationCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:config_validation",
				Name:      "Config Validation",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1045,
			},
		},
		{
			Checker: &RetryCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:retry",
				Name:      "Retry Logic",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1050,
			},
		},
		{
			Checker: &ContributingCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:contributing",
				Name:      "Contributing Guidelines",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     970,
			},
		},
		{
			Checker: &EditorconfigCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:editorconfig",
				Name:      "Editor Config",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1070,
			},
		},
		{
			Checker: &E2ECheck{},
			Meta: checker.CheckMeta{
				ID:        "common:e2e",
				Name:      "E2E Tests",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     985,
			},
		},
		{
			Checker: &TracingCheck{},
			Meta: checker.CheckMeta{
				ID:        "common:tracing",
				Name:      "Distributed Tracing",
				Languages: []checker.Language{checker.LangCommon},
				Critical:  false,
				Order:     1015,
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
