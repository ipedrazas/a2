package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ErrorsCheck verifies that error tracking is configured.
type ErrorsCheck struct{}

func (c *ErrorsCheck) ID() string   { return "common:errors" }
func (c *ErrorsCheck) Name() string { return "Error Tracking" }

// Run checks for error tracking SDKs and configuration.
func (c *ErrorsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var found []string

	// Check for Go error tracking libraries
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goErrorLibs := []struct {
				pattern string
				name    string
			}{
				{"getsentry/sentry-go", "Sentry"},
				{"rollbar/rollbar-go", "Rollbar"},
				{"bugsnag/bugsnag-go", "Bugsnag"},
				{"honeybadger-io/honeybadger-go", "Honeybadger"},
				{"airbrake/gobrake", "Airbrake"},
				{"raygun4go", "Raygun"},
			}
			for _, lib := range goErrorLibs {
				if strings.Contains(string(content), lib.pattern) {
					found = append(found, lib.name)
				}
			}
		}
	}

	// Check for Python error tracking libraries
	pythonConfigs := []string{"pyproject.toml", "requirements.txt", "setup.py"}
	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				contentLower := strings.ToLower(string(content))
				pythonErrorLibs := []struct {
					pattern string
					name    string
				}{
					{"sentry-sdk", "Sentry"},
					{"sentry_sdk", "Sentry"},
					{"rollbar", "Rollbar"},
					{"bugsnag", "Bugsnag"},
					{"honeybadger", "Honeybadger"},
					{"airbrake", "Airbrake"},
					{"raygun4py", "Raygun"},
				}
				for _, lib := range pythonErrorLibs {
					if strings.Contains(contentLower, lib.pattern) {
						found = append(found, lib.name)
					}
				}
			}
			break
		}
	}

	// Check for Node.js error tracking libraries
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			nodeErrorLibs := []struct {
				pattern string
				name    string
			}{
				{"@sentry/node", "Sentry"},
				{"@sentry/browser", "Sentry"},
				{"rollbar", "Rollbar"},
				{"@bugsnag/js", "Bugsnag"},
				{"bugsnag", "Bugsnag"},
				{"@honeybadger-io/js", "Honeybadger"},
				{"airbrake-js", "Airbrake"},
				{"raygun4js", "Raygun"},
			}
			for _, lib := range nodeErrorLibs {
				if strings.Contains(string(content), lib.pattern) {
					found = append(found, lib.name)
				}
			}
		}
	}

	// Check for Java error tracking libraries
	javaConfigs := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, cfg := range javaConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				contentLower := strings.ToLower(string(content))
				javaErrorLibs := []struct {
					pattern string
					name    string
				}{
					{"sentry", "Sentry"},
					{"rollbar", "Rollbar"},
					{"bugsnag", "Bugsnag"},
					{"honeybadger", "Honeybadger"},
					{"airbrake", "Airbrake"},
				}
				for _, lib := range javaErrorLibs {
					if strings.Contains(contentLower, lib.pattern) {
						found = append(found, lib.name)
					}
				}
			}
			break
		}
	}

	// Check for error tracking configuration files
	configFiles := []struct {
		file string
		name string
	}{
		{".sentryclirc", "Sentry CLI"},
		{"sentry.properties", "Sentry"},
		{".rollbar", "Rollbar"},
		{"bugsnag.json", "Bugsnag"},
	}
	for _, cfg := range configFiles {
		if safepath.Exists(path, cfg.file) {
			found = append(found, cfg.name+" config")
		}
	}

	// Check for Sentry DSN in environment examples
	envFiles := []string{".env.example", ".env.sample", "env.example"}
	for _, envFile := range envFiles {
		if safepath.Exists(path, envFile) {
			if content, err := safepath.ReadFile(path, envFile); err == nil {
				contentUpper := strings.ToUpper(string(content))
				if strings.Contains(contentUpper, "SENTRY_DSN") ||
					strings.Contains(contentUpper, "ROLLBAR_TOKEN") ||
					strings.Contains(contentUpper, "BUGSNAG_API_KEY") {
					found = append(found, "error tracking env vars")
				}
			}
		}
	}

	// Check CI config for error tracking integration
	ciFiles := []string{".github/workflows", ".gitlab-ci.yml", "Jenkinsfile"}
	for _, ciFile := range ciFiles {
		if safepath.Exists(path, ciFile) || safepath.IsDir(path, ciFile) {
			// Check for sentry-cli or similar in CI
			if safepath.IsDir(path, ciFile) {
				if files, err := safepath.Glob(path+"/"+ciFile, "*.yml"); err == nil {
					for _, f := range files {
						if content, err := safepath.ReadFileAbs(f); err == nil {
							if strings.Contains(strings.ToLower(string(content)), "sentry") {
								found = append(found, "Sentry CI integration")
								break
							}
						}
					}
				}
			}
		}
	}

	// Build result
	found = unique(found)
	if len(found) > 0 {
		return rb.Pass("Error tracking configured: " + strings.Join(found, ", ")), nil
	}
	return rb.Warn("No error tracking found (consider adding Sentry, Rollbar, or Bugsnag)"), nil
}
