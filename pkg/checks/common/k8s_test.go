package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type K8sCheckTestSuite struct {
	suite.Suite
	check   *K8sCheck
	tempDir string
}

func (s *K8sCheckTestSuite) SetupTest() {
	s.check = &K8sCheck{}
	tempDir, err := os.MkdirTemp("", "k8s-test-*")
	s.Require().NoError(err)
	s.tempDir = tempDir
}

func (s *K8sCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *K8sCheckTestSuite) TestID() {
	s.Equal("common:k8s", s.check.ID())
}

func (s *K8sCheckTestSuite) TestName() {
	s.Equal("Kubernetes Ready", s.check.Name())
}

func (s *K8sCheckTestSuite) TestRun_NoK8sConfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No Kubernetes manifests")
}

func (s *K8sCheckTestSuite) TestRun_HelmChartInRoot() {
	chartYaml := filepath.Join(s.tempDir, "Chart.yaml")
	err := os.WriteFile(chartYaml, []byte("apiVersion: v2\nname: myapp"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Helm chart")
}

func (s *K8sCheckTestSuite) TestRun_HelmChartInChartsDir() {
	chartsDir := filepath.Join(s.tempDir, "charts", "myapp")
	err := os.MkdirAll(chartsDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(chartsDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte("apiVersion: v2\nname: myapp"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Helm chart")
}

func (s *K8sCheckTestSuite) TestRun_HelmChartInHelmDir() {
	helmDir := filepath.Join(s.tempDir, "helm")
	err := os.MkdirAll(helmDir, 0755)
	s.Require().NoError(err)

	chartYaml := filepath.Join(helmDir, "Chart.yaml")
	err = os.WriteFile(chartYaml, []byte("apiVersion: v2\nname: myapp"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Helm chart")
}

func (s *K8sCheckTestSuite) TestRun_KustomizeInRoot() {
	kustomizeYaml := filepath.Join(s.tempDir, "kustomization.yaml")
	err := os.WriteFile(kustomizeYaml, []byte("resources:\n  - deployment.yaml"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Kustomize")
}

func (s *K8sCheckTestSuite) TestRun_KustomizeInK8sDir() {
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	kustomizeYaml := filepath.Join(k8sDir, "kustomization.yaml")
	err = os.WriteFile(kustomizeYaml, []byte("resources:\n  - deployment.yaml"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Kustomize")
}

func (s *K8sCheckTestSuite) TestRun_KustomizeInBaseDir() {
	baseDir := filepath.Join(s.tempDir, "k8s", "base")
	err := os.MkdirAll(baseDir, 0755)
	s.Require().NoError(err)

	kustomizeYaml := filepath.Join(baseDir, "kustomization.yaml")
	err = os.WriteFile(kustomizeYaml, []byte("resources:\n  - deployment.yaml"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Kustomize")
}

func (s *K8sCheckTestSuite) TestRun_K8sManifestsDir() {
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
	s.True(result.Passed)
	s.Contains(result.Reason, "k8s/ directory")
}

func (s *K8sCheckTestSuite) TestRun_KubernetesDir() {
	kubeDir := filepath.Join(s.tempDir, "kubernetes")
	err := os.MkdirAll(kubeDir, 0755)
	s.Require().NoError(err)

	serviceYaml := filepath.Join(kubeDir, "service.yaml")
	content := `apiVersion: v1
kind: Service
metadata:
  name: myapp`
	err = os.WriteFile(serviceYaml, []byte(content), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "kubernetes/ directory")
}

func (s *K8sCheckTestSuite) TestRun_DeployDir() {
	deployDir := filepath.Join(s.tempDir, "deploy")
	err := os.MkdirAll(deployDir, 0755)
	s.Require().NoError(err)

	configmapYaml := filepath.Join(deployDir, "configmap.yaml")
	content := `apiVersion: v1
kind: ConfigMap
metadata:
  name: myconfig`
	err = os.WriteFile(configmapYaml, []byte(content), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "deploy/ directory")
}

func (s *K8sCheckTestSuite) TestRun_ManifestsInRoot() {
	deploymentYaml := filepath.Join(s.tempDir, "deployment.yaml")
	err := os.WriteFile(deploymentYaml, []byte("apiVersion: apps/v1\nkind: Deployment"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "K8s manifests")
}

func (s *K8sCheckTestSuite) TestRun_ServiceInRoot() {
	serviceYaml := filepath.Join(s.tempDir, "service.yaml")
	err := os.WriteFile(serviceYaml, []byte("apiVersion: v1\nkind: Service"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "K8s manifests")
}

func (s *K8sCheckTestSuite) TestRun_DockerCompose() {
	composeYaml := filepath.Join(s.tempDir, "docker-compose.yaml")
	err := os.WriteFile(composeYaml, []byte("version: '3'\nservices:"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Docker Compose")
}

func (s *K8sCheckTestSuite) TestRun_ComposeYml() {
	composeYml := filepath.Join(s.tempDir, "compose.yml")
	err := os.WriteFile(composeYml, []byte("version: '3'\nservices:"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Docker Compose")
}

func (s *K8sCheckTestSuite) TestRun_Skaffold() {
	skaffoldYaml := filepath.Join(s.tempDir, "skaffold.yaml")
	err := os.WriteFile(skaffoldYaml, []byte("apiVersion: skaffold/v2\nkind: Config"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Skaffold")
}

func (s *K8sCheckTestSuite) TestRun_Tilt() {
	tiltfile := filepath.Join(s.tempDir, "Tiltfile")
	err := os.WriteFile(tiltfile, []byte("k8s_yaml('k8s/deployment.yaml')"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Tilt")
}

func (s *K8sCheckTestSuite) TestRun_MultipleConfigs() {
	// Create Helm chart
	chartYaml := filepath.Join(s.tempDir, "Chart.yaml")
	err := os.WriteFile(chartYaml, []byte("apiVersion: v2\nname: myapp"), 0644)
	s.Require().NoError(err)

	// Create Docker Compose
	composeYaml := filepath.Join(s.tempDir, "docker-compose.yaml")
	err = os.WriteFile(composeYaml, []byte("version: '3'\nservices:"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Helm chart")
	s.Contains(result.Reason, "Docker Compose")
}

func (s *K8sCheckTestSuite) TestRun_NestedK8sManifests() {
	// Create nested directory structure
	subDir := filepath.Join(s.tempDir, "k8s", "apps", "myapp")
	err := os.MkdirAll(subDir, 0755)
	s.Require().NoError(err)

	deploymentYaml := filepath.Join(subDir, "deployment.yaml")
	content := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp`
	err = os.WriteFile(deploymentYaml, []byte(content), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "k8s/ directory")
}

func (s *K8sCheckTestSuite) TestRun_NonK8sYamlIgnored() {
	// Create a YAML file that is not a K8s manifest
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	configYaml := filepath.Join(k8sDir, "config.yaml")
	err = os.WriteFile(configYaml, []byte("database:\n  host: localhost"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
}

func (s *K8sCheckTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func (s *K8sCheckTestSuite) TestRun_StatefulSet() {
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	statefulsetYaml := filepath.Join(k8sDir, "statefulset.yaml")
	content := `apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mydb`
	err = os.WriteFile(statefulsetYaml, []byte(content), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
}

func (s *K8sCheckTestSuite) TestRun_DaemonSet() {
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	daemonsetYaml := filepath.Join(k8sDir, "daemonset.yaml")
	content := `apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: myagent`
	err = os.WriteFile(daemonsetYaml, []byte(content), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
}

func (s *K8sCheckTestSuite) TestRun_NetworkPolicy() {
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	networkPolicyYaml := filepath.Join(k8sDir, "networkpolicy.yaml")
	content := `apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-ingress`
	err = os.WriteFile(networkPolicyYaml, []byte(content), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
}

func TestK8sCheckTestSuite(t *testing.T) {
	suite.Run(t, new(K8sCheckTestSuite))
}
