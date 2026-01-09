package targets

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
	dir, err := os.MkdirTemp("", "a2-targets-loader-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

func (suite *LoaderTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

func (suite *LoaderTestSuite) createTargetFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

func (suite *LoaderTestSuite) TestLoadFromFile_ValidTarget() {
	content := `
name: staging
description: Staging environment target
disabled:
  - common:sast
  - go:coverage
`
	path := suite.createTargetFile("staging.yaml", content)

	target, err := LoadFromFile(path)

	suite.NoError(err)
	suite.NotNil(target)
	suite.Equal("staging", target.Name)
	suite.Equal("Staging environment target", target.Description)
	suite.Equal([]string{"common:sast", "go:coverage"}, target.Disabled)
	suite.Equal(SourceUser, target.Source)
}

func (suite *LoaderTestSuite) TestLoadFromFile_EmptyDisabled() {
	content := `
name: strict
description: Strict target with no checks disabled
disabled: []
`
	path := suite.createTargetFile("strict.yaml", content)

	target, err := LoadFromFile(path)

	suite.NoError(err)
	suite.NotNil(target)
	suite.Equal("strict", target.Name)
	suite.Empty(target.Disabled)
}

func (suite *LoaderTestSuite) TestLoadFromFile_NoName() {
	content := `
description: Target without a name
disabled:
  - common:license
`
	path := suite.createTargetFile("unnamed.yaml", content)

	target, err := LoadFromFile(path)

	suite.NoError(err)
	suite.NotNil(target)
	suite.Equal("unnamed", target.Name) // Should use filename
}

func (suite *LoaderTestSuite) TestLoadFromFile_InvalidYAML() {
	// Use truly invalid YAML - mixing mapping and sequence at same level
	content := `
name: broken
disabled:
  key: value
  - item
`
	path := suite.createTargetFile("broken.yaml", content)

	target, err := LoadFromFile(path)

	suite.Error(err)
	suite.Nil(target)
}

func (suite *LoaderTestSuite) TestLoadFromFile_FileNotFound() {
	target, err := LoadFromFile("/nonexistent/path/target.yaml")

	suite.Error(err)
	suite.Nil(target)
}

func (suite *LoaderTestSuite) TestDiscoverTargets_MultipleTargets() {
	suite.createTargetFile("dev.yaml", `
name: dev
description: Development target
disabled:
  - common:license
`)
	suite.createTargetFile("qa.yaml", `
name: qa
description: QA target
disabled:
  - common:changelog
`)

	targets, err := DiscoverTargets(suite.tempDir)

	suite.NoError(err)
	suite.Len(targets, 2)
	suite.Contains(targets, "dev")
	suite.Contains(targets, "qa")
}

func (suite *LoaderTestSuite) TestDiscoverTargets_YmlExtension() {
	suite.createTargetFile("integration.yml", `
name: integration
description: Integration target
disabled: []
`)

	targets, err := DiscoverTargets(suite.tempDir)

	suite.NoError(err)
	suite.Len(targets, 1)
	suite.Contains(targets, "integration")
}

func (suite *LoaderTestSuite) TestDiscoverTargets_EmptyDirectory() {
	targets, err := DiscoverTargets(suite.tempDir)

	suite.NoError(err)
	suite.Empty(targets)
}

func (suite *LoaderTestSuite) TestDiscoverTargets_NonExistentDirectory() {
	targets, err := DiscoverTargets("/nonexistent/directory")

	suite.NoError(err) // Should not error, just return empty map
	suite.Empty(targets)
}

func (suite *LoaderTestSuite) TestDiscoverTargets_SkipsInvalidFiles() {
	suite.createTargetFile("valid.yaml", `
name: valid
description: Valid target
disabled: []
`)
	// Truly invalid YAML - mixing mapping and sequence at same level
	suite.createTargetFile("invalid.yaml", `
name: invalid
disabled:
  key: value
  - broken
`)

	targets, err := DiscoverTargets(suite.tempDir)

	suite.NoError(err)
	suite.Len(targets, 1)
	suite.Contains(targets, "valid")
}

func TestLoaderTestSuite(t *testing.T) {
	suite.Run(t, new(LoaderTestSuite))
}
