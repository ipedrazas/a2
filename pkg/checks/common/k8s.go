package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// K8sCheck verifies that Kubernetes manifests or deployment configurations exist.
type K8sCheck struct{}

func (c *K8sCheck) ID() string   { return "common:k8s" }
func (c *K8sCheck) Name() string { return "Kubernetes Ready" }

// Run checks for Kubernetes manifests, Helm charts, or Kustomize configurations.
func (c *K8sCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var foundConfigs []string

	// Check for Helm charts
	if c.hasHelmChart(path) {
		foundConfigs = append(foundConfigs, "Helm chart")
	}

	// Check for Kustomize
	if c.hasKustomize(path) {
		foundConfigs = append(foundConfigs, "Kustomize")
	}

	// Check for k8s manifest directories
	k8sDirs := []string{
		"k8s",
		"kubernetes",
		"deploy",
		"deployment",
		"deployments",
		"manifests",
		".k8s",
	}

	for _, dir := range k8sDirs {
		if c.hasK8sManifestsInDir(path, dir) {
			foundConfigs = append(foundConfigs, dir+"/ directory")
			break // Only report one directory
		}
	}

	// Check for standalone k8s manifest files in root
	if len(foundConfigs) == 0 {
		if c.hasK8sManifestsInRoot(path) {
			foundConfigs = append(foundConfigs, "K8s manifests")
		}
	}

	// Check for Docker Compose (alternative deployment)
	if c.hasDockerCompose(path) {
		foundConfigs = append(foundConfigs, "Docker Compose")
	}

	// Check for Skaffold
	if safepath.Exists(path, "skaffold.yaml") || safepath.Exists(path, "skaffold.yml") {
		foundConfigs = append(foundConfigs, "Skaffold")
	}

	// Check for Tilt
	if safepath.Exists(path, "Tiltfile") {
		foundConfigs = append(foundConfigs, "Tilt")
	}

	if len(foundConfigs) > 0 {
		return rb.Pass("Deployment config found: " + strings.Join(foundConfigs, ", ")), nil
	}

	return rb.Warn("No Kubernetes manifests or deployment config found"), nil
}

// hasHelmChart checks if a Helm chart exists.
func (c *K8sCheck) hasHelmChart(path string) bool {
	// Check root for Chart.yaml
	if safepath.Exists(path, "Chart.yaml") || safepath.Exists(path, "Chart.yml") {
		return true
	}

	// Check charts/ directory
	chartsDir := filepath.Join(path, "charts")
	if info, err := os.Stat(chartsDir); err == nil && info.IsDir() {
		entries, err := os.ReadDir(chartsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					chartFile := filepath.Join(chartsDir, entry.Name(), "Chart.yaml")
					if _, err := os.Stat(chartFile); err == nil {
						return true
					}
				}
			}
		}
	}

	// Check helm/ directory
	helmDir := filepath.Join(path, "helm")
	if info, err := os.Stat(helmDir); err == nil && info.IsDir() {
		if safepath.Exists(helmDir, "Chart.yaml") || safepath.Exists(helmDir, "Chart.yml") {
			return true
		}
		// Check subdirectories
		entries, err := os.ReadDir(helmDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					chartFile := filepath.Join(helmDir, entry.Name(), "Chart.yaml")
					if _, err := os.Stat(chartFile); err == nil {
						return true
					}
				}
			}
		}
	}

	return false
}

// hasKustomize checks if Kustomize configuration exists.
func (c *K8sCheck) hasKustomize(path string) bool {
	kustomizeFiles := []string{
		"kustomization.yaml",
		"kustomization.yml",
		"Kustomization",
	}

	// Check root
	for _, file := range kustomizeFiles {
		if safepath.Exists(path, file) {
			return true
		}
	}

	// Check common directories
	kustomizeDirs := []string{"k8s", "kubernetes", "deploy", "base", "overlays"}
	for _, dir := range kustomizeDirs {
		dirPath := filepath.Join(path, dir)
		if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
			for _, file := range kustomizeFiles {
				if safepath.Exists(dirPath, file) {
					return true
				}
			}
			// Check subdirectories (base, overlays/dev, overlays/prod, etc.)
			entries, _ := os.ReadDir(dirPath)
			for _, entry := range entries {
				if entry.IsDir() {
					subDir := filepath.Join(dirPath, entry.Name())
					for _, file := range kustomizeFiles {
						if safepath.Exists(subDir, file) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

// hasK8sManifestsInDir checks if a directory contains Kubernetes manifests.
func (c *K8sCheck) hasK8sManifestsInDir(path, dir string) bool {
	dirPath := filepath.Join(path, dir)
	if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
		return false
	}

	return c.containsK8sManifests(path, dirPath)
}

// hasK8sManifestsInRoot checks if the root directory contains Kubernetes manifests.
func (c *K8sCheck) hasK8sManifestsInRoot(path string) bool {
	k8sFiles := []string{
		"deployment.yaml",
		"deployment.yml",
		"service.yaml",
		"service.yml",
		"ingress.yaml",
		"ingress.yml",
		"configmap.yaml",
		"configmap.yml",
		"secret.yaml",
		"secret.yml",
		"pod.yaml",
		"pod.yml",
		"statefulset.yaml",
		"statefulset.yml",
		"daemonset.yaml",
		"daemonset.yml",
		"cronjob.yaml",
		"cronjob.yml",
		"job.yaml",
		"job.yml",
	}

	for _, file := range k8sFiles {
		if safepath.Exists(path, file) {
			return true
		}
	}

	return false
}

// containsK8sManifests checks if a directory contains Kubernetes manifest files.
func (c *K8sCheck) containsK8sManifests(root, dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Recursively check subdirectories
			subDir := filepath.Join(dirPath, entry.Name())
			if c.containsK8sManifests(root, subDir) {
				return true
			}
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		// Check if file contains Kubernetes manifest markers
		filePath := filepath.Join(dirPath, name)
		if c.isK8sManifest(root, filePath) {
			return true
		}
	}

	return false
}

// isK8sManifest checks if a file appears to be a Kubernetes manifest.
func (c *K8sCheck) isK8sManifest(root, filePath string) bool {
	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return false
	}
	defer func() {
		if err := file.Close(); err != nil {
			// File close errors are typically not critical in read-only scenarios
			fmt.Println("Error closing file:", err)
		}
	}()

	// K8s manifest markers
	k8sMarkers := []string{
		"apiVersion:",
		"kind: Deployment",
		"kind: Service",
		"kind: ConfigMap",
		"kind: Secret",
		"kind: Ingress",
		"kind: Pod",
		"kind: StatefulSet",
		"kind: DaemonSet",
		"kind: Job",
		"kind: CronJob",
		"kind: PersistentVolume",
		"kind: PersistentVolumeClaim",
		"kind: Namespace",
		"kind: ServiceAccount",
		"kind: Role",
		"kind: RoleBinding",
		"kind: ClusterRole",
		"kind: ClusterRoleBinding",
		"kind: HorizontalPodAutoscaler",
		"kind: NetworkPolicy",
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	hasAPIVersion := false

	for scanner.Scan() && lineCount < 50 {
		line := strings.TrimSpace(scanner.Text())
		lineCount++

		if strings.HasPrefix(line, "apiVersion:") {
			hasAPIVersion = true
		}

		for _, marker := range k8sMarkers {
			if strings.HasPrefix(line, marker) {
				if marker == "apiVersion:" {
					continue // Need to also find a kind
				}
				if hasAPIVersion {
					return true
				}
			}
		}
	}

	return false
}

// hasDockerCompose checks if Docker Compose files exist.
func (c *K8sCheck) hasDockerCompose(path string) bool {
	composeFiles := []string{
		"docker-compose.yaml",
		"docker-compose.yml",
		"compose.yaml",
		"compose.yml",
		"docker-compose.prod.yaml",
		"docker-compose.prod.yml",
		"docker-compose.production.yaml",
		"docker-compose.production.yml",
	}

	for _, file := range composeFiles {
		if safepath.Exists(path, file) {
			return true
		}
	}

	return false
}
