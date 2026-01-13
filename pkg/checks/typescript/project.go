package typescriptcheck

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ProjectCheck detects TypeScript project configuration.
type ProjectCheck struct{}

func (c *ProjectCheck) ID() string   { return "typescript:project" }
func (c *ProjectCheck) Name() string { return "TypeScript Project" }

// TSConfig represents the structure of tsconfig.json.
type TSConfig struct {
	CompilerOptions struct {
		Target    string `json:"target"`
		Module    string `json:"module"`
		Strict    bool   `json:"strict"`
		OutDir    string `json:"outDir"`
		RootDir   string `json:"rootDir"`
		ESModules bool   `json:"esModuleInterop"`
	} `json:"compilerOptions"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
	Extends string   `json:"extends"`
}

// PackageJSON represents package.json structure (partial).
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

// Run checks for TypeScript project configuration.
func (c *ProjectCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangTypeScript)

	// Check for tsconfig.json
	tsconfigs := []string{"tsconfig.json", "tsconfig.base.json", "tsconfig.build.json"}
	var foundConfig string
	for _, cfg := range tsconfigs {
		if safepath.Exists(path, cfg) {
			foundConfig = cfg
			break
		}
	}

	if foundConfig == "" {
		return rb.Fail("No tsconfig.json found"), nil
	}

	// Parse tsconfig.json
	content, err := safepath.ReadFile(path, foundConfig)
	if err != nil {
		return rb.Fail(fmt.Sprintf("Cannot read %s: %v", foundConfig, err)), nil
	}

	var tsconfig TSConfig
	if err := json.Unmarshal(content, &tsconfig); err != nil {
		// tsconfig might use comments or extends, try basic detection
		return rb.Pass(fmt.Sprintf("Found %s", foundConfig)), nil
	}

	// Build info message
	var info []string
	info = append(info, foundConfig)

	if tsconfig.CompilerOptions.Target != "" {
		info = append(info, fmt.Sprintf("target: %s", tsconfig.CompilerOptions.Target))
	}
	if tsconfig.CompilerOptions.Module != "" {
		info = append(info, fmt.Sprintf("module: %s", tsconfig.CompilerOptions.Module))
	}
	if tsconfig.CompilerOptions.Strict {
		info = append(info, "strict mode")
	}

	// Check for TypeScript version in package.json
	if safepath.Exists(path, "package.json") {
		if pkgContent, err := safepath.ReadFile(path, "package.json"); err == nil {
			var pkg PackageJSON
			if err := json.Unmarshal(pkgContent, &pkg); err == nil {
				if tsVersion, ok := pkg.DevDependencies["typescript"]; ok {
					info = append(info, fmt.Sprintf("TypeScript %s", tsVersion))
				} else if tsVersion, ok := pkg.Dependencies["typescript"]; ok {
					info = append(info, fmt.Sprintf("TypeScript %s", tsVersion))
				}
			}
		}
	}

	return rb.Pass(strings.Join(info, ", ")), nil
}

// ParsePackageJSON reads and parses package.json from the given path.
func ParsePackageJSON(path string) (*PackageJSON, error) {
	content, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(content, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}
