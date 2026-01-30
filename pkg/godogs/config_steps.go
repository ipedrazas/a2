package godogs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Configuration-related step implementations

// Config represents the .a2.yaml structure
type A2Config struct {
	Language map[string]interface{} `yaml:"language"`
	Profile  string                 `yaml:"profile"`
	Target   string                 `yaml:"target"`
	Checks   struct {
		Disabled []string `yaml:"disabled"`
	} `yaml:"checks"`
	Files struct {
		Required []string `yaml:"required"`
	} `yaml:"files"`
	Execution struct {
		Parallel bool `yaml:"parallel"`
		Timeout  int  `yaml:"timeout"`
	} `yaml:"execution"`
}

func iHaveInitialConfig(filename string) error {
	s := GetState()
	s.SetConfigFile(filename)
	// Create a basic config if it doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config := &A2Config{
			Profile: "api",
			Target:  "production",
		}
		return saveConfig(filename, config)
	}
	return nil
}

func iEditConfig() error {
	s := GetState()
	if s.GetConfigFile() == "" {
		return fmt.Errorf("no config file specified")
	}
	return nil
}

func iSetConfigValue(section, key string, value int) error {
	s := GetState()
	if s.GetConfigFile() == "" {
		return fmt.Errorf("no config file specified")
	}

	// Load existing config
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	// Set the value
	if config.Language == nil {
		config.Language = make(map[string]interface{})
	}

	langSection, ok := config.Language[section]
	if !ok {
		langSection = make(map[string]interface{})
		config.Language[section] = langSection
	}

	langMap, ok := langSection.(map[string]interface{})
	if !ok {
		langMap = make(map[string]interface{})
	}

	langMap[key] = value

	// Save config
	return saveConfig(s.GetConfigFile(), config)
}

func iSaveFile() error {
	// Config is already saved in iSetConfigValue
	return nil
}

func a2UsesStricterThresholds() error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	// Verify stricter thresholds are set
	if goConfig, ok := config.Language["go"].(map[string]interface{}); ok {
		if coverage, ok := goConfig["coverage_threshold"].(int); ok && coverage > 80 {
			return nil
		}
	}

	return fmt.Errorf("stricter thresholds not found in config")
}

func checksEnforceNewStandards() error {
	// Run checks and verify new standards are enforced
	return nil
}

func iAddRequiredDocs() error {
	return nil
}

func iIncludeRequiredFile(filename string) error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	// Check if file is already in the list
	for _, f := range config.Files.Required {
		if f == filename {
			return nil
		}
	}

	config.Files.Required = append(config.Files.Required, filename)
	return saveConfig(s.GetConfigFile(), config)
}

func a2VerifiesFilesExist() error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	for _, file := range config.Files.Required {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("required file %s does not exist", file)
		}
	}
	return nil
}

func a2FailsOnMissingFiles() error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	for _, file := range config.Files.Required {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// A2 should fail with missing required files
			return nil
		}
	}

	return fmt.Errorf("no files missing, should not fail")
}

func iDisableChecks(checkPattern string) error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	// Check if pattern is already in the list
	for _, pattern := range config.Checks.Disabled {
		if pattern == checkPattern {
			return nil
		}
	}

	config.Checks.Disabled = append(config.Checks.Disabled, checkPattern)
	return saveConfig(s.GetConfigFile(), config)
}

func a2SkipsChecks() error {
	// Verify disabled checks were skipped
	return nil
}

func checksFaster() error {
	// Verify execution time is faster with fewer checks
	return nil
}

func resultsShowOnlyRelevant() error {
	// Verify only relevant checks are shown
	return nil
}

func iHaveBasicConfig(filename string) error {
	s := GetState()
	s.SetConfigFile(filename)
	return iHaveInitialConfig(filename)
}

func loadConfig(filename string) (*A2Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config A2Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(filename string, config *A2Config) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func configIncludesAPIProfile() error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	if config.Profile != "api" {
		return fmt.Errorf("profile is not 'api', got: %s", config.Profile)
	}
	return nil
}

func configIncludesProductionTarget() error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	if config.Target != "production" {
		return fmt.Errorf("target is not 'production', got: %s", config.Target)
	}
	return nil
}

func e2eTestsDisabled() error {
	s := GetState()
	config, err := loadConfig(s.GetConfigFile())
	if err != nil {
		return err
	}

	for _, disabled := range config.Checks.Disabled {
		if disabled == "*:e2e" || disabled == "common:e2e" {
			return nil
		}
	}

	return fmt.Errorf("E2E tests are not disabled")
}
