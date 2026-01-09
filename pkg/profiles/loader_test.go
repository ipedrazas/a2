package profiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoaderTestSuite struct {
	suite.Suite
	tempDir string
}

func (suite *LoaderTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-profiles-loader-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

func (suite *LoaderTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

func (suite *LoaderTestSuite) createProfileFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

func (suite *LoaderTestSuite) TestLoadFromFile_ValidProfile() {
	content := `
name: custom
description: Custom profile for testing
disabled:
  - common:health
  - common:k8s
`
	path := suite.createProfileFile("custom.yaml", content)

	profile, err := LoadFromFile(path)

	suite.NoError(err)
	suite.NotNil(profile)
	suite.Equal("custom", profile.Name)
	suite.Equal("Custom profile for testing", profile.Description)
	suite.Equal([]string{"common:health", "common:k8s"}, profile.Disabled)
	suite.Equal(SourceUser, profile.Source)
}

func (suite *LoaderTestSuite) TestLoadFromFile_EmptyDisabled() {
	content := `
name: minimal
description: Minimal profile with no checks disabled
disabled: []
`
	path := suite.createProfileFile("minimal.yaml", content)

	profile, err := LoadFromFile(path)

	suite.NoError(err)
	suite.NotNil(profile)
	suite.Equal("minimal", profile.Name)
	suite.Empty(profile.Disabled)
}

func (suite *LoaderTestSuite) TestLoadFromFile_NoName() {
	content := `
description: Profile without a name
disabled:
  - common:sast
`
	path := suite.createProfileFile("unnamed.yaml", content)

	profile, err := LoadFromFile(path)

	suite.NoError(err)
	suite.NotNil(profile)
	suite.Equal("unnamed", profile.Name) // Should use filename
}

func (suite *LoaderTestSuite) TestLoadFromFile_InvalidYAML() {
	// Use truly invalid YAML - mixing mapping and sequence at same level
	content := `
name: broken
disabled:
  key: value
  - item
`
	path := suite.createProfileFile("broken.yaml", content)

	profile, err := LoadFromFile(path)

	suite.Error(err)
	suite.Nil(profile)
}

func (suite *LoaderTestSuite) TestLoadFromFile_FileNotFound() {
	profile, err := LoadFromFile("/nonexistent/path/profile.yaml")

	suite.Error(err)
	suite.Nil(profile)
}

func (suite *LoaderTestSuite) TestDiscoverProfiles_MultipleProfiles() {
	suite.createProfileFile("alpha.yaml", `
name: alpha
description: Alpha profile
disabled:
  - common:health
`)
	suite.createProfileFile("beta.yaml", `
name: beta
description: Beta profile
disabled:
  - common:k8s
`)

	profiles, err := DiscoverProfiles(suite.tempDir)

	suite.NoError(err)
	suite.Len(profiles, 2)
	suite.Contains(profiles, "alpha")
	suite.Contains(profiles, "beta")
}

func (suite *LoaderTestSuite) TestDiscoverProfiles_YmlExtension() {
	suite.createProfileFile("gamma.yml", `
name: gamma
description: Gamma profile
disabled: []
`)

	profiles, err := DiscoverProfiles(suite.tempDir)

	suite.NoError(err)
	suite.Len(profiles, 1)
	suite.Contains(profiles, "gamma")
}

func (suite *LoaderTestSuite) TestDiscoverProfiles_EmptyDirectory() {
	profiles, err := DiscoverProfiles(suite.tempDir)

	suite.NoError(err)
	suite.Empty(profiles)
}

func (suite *LoaderTestSuite) TestDiscoverProfiles_NonExistentDirectory() {
	profiles, err := DiscoverProfiles("/nonexistent/directory")

	suite.NoError(err) // Should not error, just return empty map
	suite.Empty(profiles)
}

func (suite *LoaderTestSuite) TestDiscoverProfiles_SkipsInvalidFiles() {
	suite.createProfileFile("valid.yaml", `
name: valid
description: Valid profile
disabled: []
`)
	// Truly invalid YAML - mixing mapping and sequence at same level
	suite.createProfileFile("invalid.yaml", `
name: invalid
disabled:
  key: value
  - broken
`)

	profiles, err := DiscoverProfiles(suite.tempDir)

	suite.NoError(err)
	suite.Len(profiles, 1)
	suite.Contains(profiles, "valid")
}

func (suite *LoaderTestSuite) TestWriteBuiltInProfiles() {
	// Skip if we can't write to the user config directory
	// This test is more of an integration test
	suite.T().Skip("Integration test - requires write access to user config directory")
}

func TestLoaderTestSuite(t *testing.T) {
	suite.Run(t, new(LoaderTestSuite))
}
