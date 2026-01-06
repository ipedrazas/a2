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
