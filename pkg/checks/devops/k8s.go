package devops

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

type K8sCheck struct{}

func (c *K8sCheck) ID() string   { return "devops:k8s" }
func (c *K8sCheck) Name() string { return "Kubernetes Ready" }

func (c *K8sCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var foundConfigs []string

	if c.hasHelmChart(path) {
		foundConfigs = append(foundConfigs, "Helm chart")
	}

	if c.hasKustomize(path) {
		foundConfigs = append(foundConfigs, "Kustomize")
	}

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
			break
		}
	}

	if len(foundConfigs) == 0 {
		if c.hasK8sManifestsInRoot(path) {
			foundConfigs = append(foundConfigs, "K8s manifests")
		}
	}

	if c.hasDockerCompose(path) {
		foundConfigs = append(foundConfigs, "Docker Compose")
	}

	if safepath.Exists(path, "skaffold.yaml") || safepath.Exists(path, "skaffold.yml") {
		foundConfigs = append(foundConfigs, "Skaffold")
	}

	if safepath.Exists(path, "Tiltfile") {
		foundConfigs = append(foundConfigs, "Tilt")
	}

	if len(foundConfigs) > 0 {
		return rb.Pass("Deployment config found: " + strings.Join(foundConfigs, ", ")), nil
	}

	return rb.Warn("No Kubernetes manifests or deployment config found"), nil
}

func (c *K8sCheck) hasHelmChart(path string) bool {
	if safepath.Exists(path, "Chart.yaml") || safepath.Exists(path, "Chart.yml") {
		return true
	}

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

	helmDir := filepath.Join(path, "helm")
	if info, err := os.Stat(helmDir); err == nil && info.IsDir() {
		if safepath.Exists(helmDir, "Chart.yaml") || safepath.Exists(helmDir, "Chart.yml") {
			return true
		}
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

func (c *K8sCheck) hasKustomize(path string) bool {
	kustomizeFiles := []string{
		"kustomization.yaml",
		"kustomization.yml",
		"Kustomization",
	}

	for _, file := range kustomizeFiles {
		if safepath.Exists(path, file) {
			return true
		}
	}

	kustomizeDirs := []string{"k8s", "kubernetes", "deploy", "base", "overlays"}
	for _, dir := range kustomizeDirs {
		dirPath := filepath.Join(path, dir)
		if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
			for _, file := range kustomizeFiles {
				if safepath.Exists(dirPath, file) {
					return true
				}
			}
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

func (c *K8sCheck) hasK8sManifestsInDir(path, dir string) bool {
	dirPath := filepath.Join(path, dir)
	if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
		return false
	}

	return c.containsK8sManifests(path, dirPath)
}

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

func (c *K8sCheck) containsK8sManifests(root, dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
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

		filePath := filepath.Join(dirPath, name)
		if c.isK8sManifest(root, filePath) {
			return true
		}
	}

	return false
}

func (c *K8sCheck) isK8sManifest(root, filePath string) bool {
	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return false
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

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
					continue
				}
				if hasAPIVersion {
					return true
				}
			}
		}
	}

	return false
}

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
