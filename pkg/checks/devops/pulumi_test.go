package devops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type PulumiCheckTestSuite struct {
	suite.Suite
	check   *PulumiCheck
	tempDir string
}

func (s *PulumiCheckTestSuite) SetupTest() {
	s.check = &PulumiCheck{}
	tempDir, err := os.MkdirTemp("", "pulumi-test-*")
	s.Require().NoError(err)
	s.tempDir = tempDir
}

func (s *PulumiCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *PulumiCheckTestSuite) TestID() {
	s.Equal("devops:pulumi", s.check.ID())
}

func (s *PulumiCheckTestSuite) TestName() {
	s.Equal("Pulumi Configuration", s.check.Name())
}

func (s *PulumiCheckTestSuite) TestRun_NoPulumiFiles() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Pulumi files found")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiYaml() {
	pulumiYaml := filepath.Join(s.tempDir, "Pulumi.yaml")
	err := os.WriteFile(pulumiYaml, []byte(`name: my-project
runtime: nodejs
description: A minimal Pulumi program`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	// Without pulumi installed, should pass with Info
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiYml() {
	pulumiYml := filepath.Join(s.tempDir, "Pulumi.yml")
	err := os.WriteFile(pulumiYml, []byte(`name: my-project
runtime: python`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiGoYaml() {
	pulumiGoYaml := filepath.Join(s.tempDir, "Pulumi.go.yaml")
	err := os.WriteFile(pulumiGoYaml, []byte(`name: my-project
runtime: go`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiPythonYaml() {
	pulumiPythonYaml := filepath.Join(s.tempDir, "Pulumi.python.yaml")
	err := os.WriteFile(pulumiPythonYaml, []byte(`name: my-project
runtime: python`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiTsYaml() {
	pulumiTsYaml := filepath.Join(s.tempDir, "Pulumi.ts.yaml")
	err := os.WriteFile(pulumiTsYaml, []byte(`name: my-project
runtime: nodejs`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiJavaScriptYaml() {
	pulumiJsYaml := filepath.Join(s.tempDir, "Pulumi.javascript.yaml")
	err := os.WriteFile(pulumiJsYaml, []byte(`name: my-project
runtime: nodejs`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiCsYaml() {
	pulumiCsYaml := filepath.Join(s.tempDir, "Pulumi.cs.yaml")
	err := os.WriteFile(pulumiCsYaml, []byte(`name: my-project
runtime: dotnet`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiJavaYaml() {
	pulumiJavaYaml := filepath.Join(s.tempDir, "Pulumi.java.yaml")
	err := os.WriteFile(pulumiJavaYaml, []byte(`name: my-project
runtime: java`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiYamlYaml() {
	pulumiYamlYaml := filepath.Join(s.tempDir, "Pulumi.yaml.yaml")
	err := os.WriteFile(pulumiYamlYaml, []byte(`name: my-project
runtime: yaml`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiDirectory() {
	pulumiDir := filepath.Join(s.tempDir, "pulumi")
	err := os.MkdirAll(pulumiDir, 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_DotPulumiDirectory() {
	dotPulumiDir := filepath.Join(s.tempDir, ".pulumi")
	err := os.MkdirAll(dotPulumiDir, 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_PulumiInSubdirectory() {
	subDir := filepath.Join(s.tempDir, "infra")
	err := os.MkdirAll(subDir, 0755)
	s.Require().NoError(err)

	pulumiYaml := filepath.Join(subDir, "Pulumi.yaml")
	err = os.WriteFile(pulumiYaml, []byte(`name: my-project
runtime: nodejs`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_MultiplePulumiFiles() {
	// Create Pulumi.yaml in root
	pulumiYaml := filepath.Join(s.tempDir, "Pulumi.yaml")
	err := os.WriteFile(pulumiYaml, []byte(`name: my-project
runtime: nodejs`), 0644)
	s.Require().NoError(err)

	// Create .pulumi directory
	dotPulumiDir := filepath.Join(s.tempDir, ".pulumi")
	err = os.MkdirAll(dotPulumiDir, 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func (s *PulumiCheckTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func (s *PulumiCheckTestSuite) TestRun_IgnoresHiddenDirs() {
	// Create .hidden directory with Pulumi.yaml (should be ignored)
	hiddenDir := filepath.Join(s.tempDir, ".hidden")
	err := os.MkdirAll(hiddenDir, 0755)
	s.Require().NoError(err)

	pulumiYaml := filepath.Join(hiddenDir, "Pulumi.yaml")
	err = os.WriteFile(pulumiYaml, []byte(`name: hidden
runtime: nodejs`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Pulumi files found")
}

func (s *PulumiCheckTestSuite) TestRun_IgnoresVendorDirs() {
	// Create vendor directory with Pulumi.yaml (should be ignored)
	vendorDir := filepath.Join(s.tempDir, "vendor")
	err := os.MkdirAll(vendorDir, 0755)
	s.Require().NoError(err)

	pulumiYaml := filepath.Join(vendorDir, "Pulumi.yaml")
	err = os.WriteFile(pulumiYaml, []byte(`name: vendor
runtime: nodejs`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Pulumi files found")
}

func (s *PulumiCheckTestSuite) TestRun_DoesNotConfuseWithK8sManifests() {
	// Create Kubernetes manifests but no Pulumi files
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	deploymentYaml := filepath.Join(k8sDir, "deployment.yaml")
	content := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp`
	err = os.WriteFile(deploymentYaml, []byte(content), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Pulumi files found")
}

func (s *PulumiCheckTestSuite) TestRun_DoesNotConfuseWithHelmCharts() {
	// Create Helm chart but no Pulumi files
	chartsDir := filepath.Join(s.tempDir, "charts", "mychart")
	err := os.MkdirAll(chartsDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(chartsDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: mychart`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Pulumi files found")
}

func (s *PulumiCheckTestSuite) TestRun_DetectsNestedPulumiFile() {
	// Create deeply nested Pulumi.yaml
	nestedDir := filepath.Join(s.tempDir, "infra", "production", "us-east-1")
	err := os.MkdirAll(nestedDir, 0755)
	s.Require().NoError(err)

	pulumiYaml := filepath.Join(nestedDir, "Pulumi.yaml")
	err = os.WriteFile(pulumiYaml, []byte(`name: nested-project
runtime: go`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "pulumi")
}

func TestPulumiCheckTestSuite(t *testing.T) {
	suite.Run(t, new(PulumiCheckTestSuite))
}
