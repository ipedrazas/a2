package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type LicenseCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *LicenseCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "license-test-*")
	s.Require().NoError(err)
}

func (s *LicenseCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *LicenseCheckTestSuite) TestIDAndName() {
	check := &LicenseCheck{}
	s.Equal("common:license", check.ID())
	s.Equal("License Compliance", check.Name())
}

func (s *LicenseCheckTestSuite) TestLicensrcConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".licensrc"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "licensrc config")
}

func (s *LicenseCheckTestSuite) TestLicensrcJsonConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".licensrc.json"), []byte(`{"allowedLicenses": ["MIT"]}`), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "licensrc config")
}

func (s *LicenseCheckTestSuite) TestLicenseCheckerConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "license-checker.json"), []byte(`{"production": true}`), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "license-checker config")
}

func (s *LicenseCheckTestSuite) TestFOSSAConfig() {
	content := `version: 3
project:
  name: myproject
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".fossa.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "FOSSA")
}

func (s *LicenseCheckTestSuite) TestFOSSAYamlConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "fossa.yaml"), []byte("version: 3"), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "FOSSA")
}

func (s *LicenseCheckTestSuite) TestSPDXFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "spdx.json"), []byte(`{"spdxVersion": "SPDX-2.3"}`), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SPDX")
}

func (s *LicenseCheckTestSuite) TestSPDXSBOMFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "sbom.spdx.json"), []byte(`{}`), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SPDX SBOM")
}

func (s *LicenseCheckTestSuite) TestGoLicenses() {
	content := `module myapp

go 1.21

require github.com/google/go-licenses v1.6.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "go-licenses")
}

func (s *LicenseCheckTestSuite) TestGoLichen() {
	content := `module myapp

go 1.21

require github.com/uw-labs/lichen v0.1.7
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "lichen")
}

func (s *LicenseCheckTestSuite) TestGoSyft() {
	content := `module myapp

go 1.21

require github.com/anchore/syft v0.100.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Syft")
}

func (s *LicenseCheckTestSuite) TestGoCycloneDX() {
	content := `module myapp

go 1.21

require github.com/CycloneDX/cyclonedx-gomod v1.4.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "CycloneDX")
}

func (s *LicenseCheckTestSuite) TestPythonPipLicenses() {
	content := `[project]
name = "myapp"
dependencies = [
    "pip-licenses>=4.0.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "pip-licenses")
}

func (s *LicenseCheckTestSuite) TestPythonLiccheck() {
	content := `pip-licenses>=4.0.0
liccheck>=0.9.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "liccheck")
}

func (s *LicenseCheckTestSuite) TestPythonLiccheckConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".liccheck.ini"), []byte("[Licenses]\nauthorized_licenses = MIT"), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "liccheck config")
}

func (s *LicenseCheckTestSuite) TestPythonCycloneDX() {
	content := `cyclonedx-bom>=3.0.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements-dev.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "CycloneDX")
}

func (s *LicenseCheckTestSuite) TestNodeLicenseChecker() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "license-checker": "^25.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "license-checker")
}

func (s *LicenseCheckTestSuite) TestNodeLicenseCompliance() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "license-compliance": "^2.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "license-compliance")
}

func (s *LicenseCheckTestSuite) TestNodeCycloneDX() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "@cyclonedx/bom": "^4.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "CycloneDX")
}

func (s *LicenseCheckTestSuite) TestNodeSnyk() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "snyk": "^1.1000.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Snyk")
}

func (s *LicenseCheckTestSuite) TestJavaLicenseMavenPlugin() {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <build>
        <plugins>
            <plugin>
                <groupId>com.mycila</groupId>
                <artifactId>license-maven-plugin</artifactId>
                <version>4.3</version>
            </plugin>
        </plugins>
    </build>
</project>
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	// Matches license-maven-plugin pattern
	s.Contains(result.Reason, "license-maven-plugin")
}

func (s *LicenseCheckTestSuite) TestJavaCycloneDXMaven() {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <build>
        <plugins>
            <plugin>
                <groupId>org.cyclonedx</groupId>
                <artifactId>cyclonedx-maven-plugin</artifactId>
                <version>2.7.0</version>
            </plugin>
        </plugins>
    </build>
</project>
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "CycloneDX")
}

func (s *LicenseCheckTestSuite) TestJavaGradleLicensePlugin() {
	content := `plugins {
    id 'java'
    id 'com.github.hierynomus.license' version '0.16.1'
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Gradle License Plugin")
}

func (s *LicenseCheckTestSuite) TestCILicenseScanning_GitHubActions() {
	workflowsDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowsDir, 0755)
	s.Require().NoError(err)

	content := `name: License Check
on: [push]
jobs:
  license:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: fossa-contrib/fossa-action@v2
`
	err = os.WriteFile(filepath.Join(workflowsDir, "license.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CI license scanning")
}

func (s *LicenseCheckTestSuite) TestCILicenseScanning_GoLicenses() {
	workflowsDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowsDir, 0755)
	s.Require().NoError(err)

	content := `name: CI
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: go-licenses check ./...
`
	err = os.WriteFile(filepath.Join(workflowsDir, "ci.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "CI license scanning")
}

func (s *LicenseCheckTestSuite) TestCILicenseScanning_SBOM() {
	workflowsDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowsDir, 0755)
	s.Require().NoError(err)

	content := `name: SBOM
on: [push]
jobs:
  sbom:
    runs-on: ubuntu-latest
    steps:
      - uses: anchore/sbom-action@v0
`
	err = os.WriteFile(filepath.Join(workflowsDir, "sbom.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "CI license scanning")
}

func (s *LicenseCheckTestSuite) TestMultipleFindings() {
	// Create FOSSA config
	err := os.WriteFile(filepath.Join(s.tempDir, ".fossa.yml"), []byte("version: 3"), 0644)
	s.Require().NoError(err)

	// Create Go module with go-licenses
	goMod := `module myapp

go 1.21

require github.com/google/go-licenses v1.6.0
`
	err = os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "FOSSA")
	s.Contains(result.Reason, "go-licenses")
}

func (s *LicenseCheckTestSuite) TestNoLicenseCompliance() {
	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No license compliance tooling found")
}

func (s *LicenseCheckTestSuite) TestResultLanguage() {
	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func (s *LicenseCheckTestSuite) TestGoModWithNoLicenseTools() {
	content := `module myapp

go 1.21

require github.com/gin-gonic/gin v1.9.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LicenseCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No license compliance tooling found")
}

func TestLicenseCheckTestSuite(t *testing.T) {
	suite.Run(t, new(LicenseCheckTestSuite))
}
