package devops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type HelmCheckTestSuite struct {
	suite.Suite
	check   *HelmCheck
	tempDir string
}

func (s *HelmCheckTestSuite) SetupTest() {
	s.check = &HelmCheck{}
	tempDir, err := os.MkdirTemp("", "helm-test-*")
	s.Require().NoError(err)
	s.tempDir = tempDir
}

func (s *HelmCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *HelmCheckTestSuite) TestID() {
	s.Equal("devops:helm", s.check.ID())
}

func (s *HelmCheckTestSuite) TestName() {
	s.Equal("Helm Charts", s.check.Name())
}

func (s *HelmCheckTestSuite) TestRun_NoHelmCharts() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Helm charts found")
}

func (s *HelmCheckTestSuite) TestRun_ChartYamlInRoot() {
	chartYaml := filepath.Join(s.tempDir, "Chart.yaml")
	err := os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: myapp
version: 1.0.0`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	// Without helm: Info (tool not installed). With helm (e.g. in CI): Pass or Warn after helm lint.
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_ChartYmlInRoot() {
	chartYml := filepath.Join(s.tempDir, "Chart.yml")
	err := os.WriteFile(chartYml, []byte(`apiVersion: v2
name: myapp`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_ChartInChartsDir() {
	chartsDir := filepath.Join(s.tempDir, "charts", "mychart")
	err := os.MkdirAll(chartsDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(chartsDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: mychart`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_ChartInHelmDir() {
	helmDir := filepath.Join(s.tempDir, "helm")
	err := os.MkdirAll(helmDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(helmDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: mychart`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_ChartInHelmSubdir() {
	helmDir := filepath.Join(s.tempDir, "helm", "subchart")
	err := os.MkdirAll(helmDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(helmDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: subchart`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_ValuesYaml() {
	valuesYaml := filepath.Join(s.tempDir, "values.yaml")
	err := os.WriteFile(valuesYaml, []byte(`replicaCount: 1
image:
  repository: nginx
  tag: stable`), 0644)
	s.Require().NoError(err)

	// values.yaml alone is not enough to detect a chart
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Helm charts found")
}

func (s *HelmCheckTestSuite) TestRun_ChartWithValues() {
	chartYaml := filepath.Join(s.tempDir, "Chart.yaml")
	err := os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: myapp`), 0644)
	s.Require().NoError(err)

	valuesYaml := filepath.Join(s.tempDir, "values.yaml")
	err = os.WriteFile(valuesYaml, []byte(`replicaCount: 1`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_MultipleChartsInDifferentDirs() {
	// Create chart in root
	chart1 := filepath.Join(s.tempDir, "Chart.yaml")
	err := os.WriteFile(chart1, []byte(`apiVersion: v2
name: chart1`), 0644)
	s.Require().NoError(err)

	// Create chart in charts dir
	chartsDir := filepath.Join(s.tempDir, "charts", "chart2")
	err = os.MkdirAll(chartsDir, 0755)
	s.Require().NoError(err)

	chart2 := filepath.Join(chartsDir, "Chart.yaml")
	err = os.WriteFile(chart2, []byte(`apiVersion: v2
name: chart2`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_NestedChart() {
	// Create a nested chart directory (unusual but possible)
	nestedDir := filepath.Join(s.tempDir, "deploy", "chart")
	err := os.MkdirAll(nestedDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(nestedDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: nested`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Status == checker.Info || result.Status == checker.Pass || result.Status == checker.Warn,
		"status: %s message: %s", result.Status, result.Message)
	s.Contains(result.Message, "helm")
}

func (s *HelmCheckTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func (s *HelmCheckTestSuite) TestRun_IgnoresHiddenDirs() {
	// Create .hidden directory with a Chart.yaml (should be ignored)
	hiddenDir := filepath.Join(s.tempDir, ".hidden")
	err := os.MkdirAll(hiddenDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(hiddenDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: hidden`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Helm charts found")
}

func (s *HelmCheckTestSuite) TestRun_IgnoresVendorDirs() {
	// Create vendor directory with a Chart.yaml (should be ignored)
	vendorDir := filepath.Join(s.tempDir, "vendor")
	err := os.MkdirAll(vendorDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(vendorDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte(`apiVersion: v2
name: vendor`), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Helm charts found")
}

func (s *HelmCheckTestSuite) TestRun_DoesNotConfuseWithK8sManifests() {
	// Create Kubernetes manifests but no Helm charts
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
	s.Contains(result.Message, "No Helm charts found")
}

func TestHelmCheckTestSuite(t *testing.T) {
	suite.Run(t, new(HelmCheckTestSuite))
}
