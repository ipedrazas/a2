package tools

import (
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

// Environment represents the detected system environment.
type Environment struct {
	OS              string   // Operating system (darwin, linux, windows)
	Arch            string   // Architecture (amd64, arm64)
	PackageManagers []string // Available package managers
}

// ToolStatus represents whether a tool is installed.
type ToolStatus struct {
	Tool      Tool
	Installed bool
	Version   string // Version if available
}

// DetectEnvironment detects the system environment.
func DetectEnvironment() Environment {
	env := Environment{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// Detect available package managers
	managers := []struct {
		name string
		cmd  string
	}{
		{"brew", "brew"},
		{"apt", "apt"},
		{"dnf", "dnf"},
		{"yum", "yum"},
		{"go", "go"},
		{"cargo", "cargo"},
		{"npm", "npm"},
		{"pip", "pip"},
		{"pip3", "pip3"},
	}

	for _, m := range managers {
		if _, err := exec.LookPath(m.cmd); err == nil {
			env.PackageManagers = append(env.PackageManagers, m.name)
		}
	}

	return env
}

// CheckInstalled checks if a tool is installed.
func CheckInstalled(tool Tool) ToolStatus {
	status := ToolStatus{
		Tool:      tool,
		Installed: false,
	}

	if len(tool.CheckCmd) == 0 {
		return status
	}

	// First, check if the binary exists in PATH
	binPath, err := exec.LookPath(tool.CheckCmd[0])
	if err != nil {
		return status
	}

	// Binary exists, mark as installed
	status.Installed = true

	// Try to run the command to get version info
	cmd := exec.Command(tool.CheckCmd[0], tool.CheckCmd[1:]...) // #nosec G204 -- command comes from hardcoded tool registry, not user input
	output, err := cmd.CombinedOutput()
	if err == nil {
		status.Version = extractVersion(string(output))
	}

	// If no version yet, try --help as fallback
	if status.Version == "" {
		cmd = exec.Command(binPath, "--help") // #nosec G204 -- binPath comes from LookPath
		output, err = cmd.CombinedOutput()
		if err == nil {
			status.Version = extractVersion(string(output))
		}
	}

	return status
}

// extractVersion tries to extract a version string from command output.
func extractVersion(output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}

	lines := strings.Split(output, "\n")
	firstLine := strings.TrimSpace(lines[0])

	// Look for version patterns in the first few lines
	for i := 0; i < len(lines) && i < 3; i++ {
		line := strings.ToLower(lines[i])
		// Check if line contains version info
		if strings.Contains(line, "version") || strings.Contains(line, " v") {
			return strings.TrimSpace(lines[i])
		}
	}

	// If first line is short enough, use it as-is (likely version output)
	if len(firstLine) < 50 {
		return firstLine
	}

	return ""
}

// CheckAllInstalled checks installation status for multiple tools.
func CheckAllInstalled(tools []Tool) []ToolStatus {
	var results []ToolStatus
	for _, t := range tools {
		results = append(results, CheckInstalled(t))
	}
	return results
}

// GetInstallCommand returns the best install command for the given environment.
func GetInstallCommand(tool Tool, env Environment) string {
	// Check package managers in order of preference
	hasManager := func(name string) bool {
		return slices.Contains(env.PackageManagers, name)
	}

	// Prefer native package managers first
	if tool.Install.Brew != "" && hasManager("brew") {
		return tool.Install.Brew
	}
	if tool.Install.Apt != "" && hasManager("apt") {
		return tool.Install.Apt
	}
	if tool.Install.Dnf != "" && hasManager("dnf") {
		return tool.Install.Dnf
	}

	// Then language-specific package managers
	if tool.Install.Go != "" && hasManager("go") {
		return tool.Install.Go
	}
	if tool.Install.Cargo != "" && hasManager("cargo") {
		return tool.Install.Cargo
	}
	if tool.Install.Npm != "" && hasManager("npm") {
		return tool.Install.Npm
	}
	if tool.Install.Pip != "" && (hasManager("pip") || hasManager("pip3")) {
		cmd := tool.Install.Pip
		// Use pip3 if pip is not available
		if !hasManager("pip") && hasManager("pip3") {
			cmd = strings.Replace(cmd, "pip ", "pip3 ", 1)
		}
		return cmd
	}

	// Fallback to manual
	if tool.Install.Manual != "" {
		return tool.Install.Manual
	}

	return ""
}

// HasPackageManager checks if a package manager is available.
func (e Environment) HasPackageManager(name string) bool {
	return slices.Contains(e.PackageManagers, name)
}

// OSName returns a human-friendly OS name.
func (e Environment) OSName() string {
	switch e.OS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return e.OS
	}
}
