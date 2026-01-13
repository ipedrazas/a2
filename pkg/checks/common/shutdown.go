package common

import (
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ShutdownCheck verifies that graceful shutdown handling exists.
type ShutdownCheck struct{}

func (c *ShutdownCheck) ID() string   { return "common:shutdown" }
func (c *ShutdownCheck) Name() string { return "Graceful Shutdown" }

// Run checks for graceful shutdown signal handling.
func (c *ShutdownCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var found []string

	// Check Go files for signal handling
	if c.hasGoSignalHandling(path) {
		found = append(found, "Go signal handling")
	}

	// Check Python files for signal handling
	if c.hasPythonSignalHandling(path) {
		found = append(found, "Python signal handling")
	}

	// Check Node.js files for signal handling
	if c.hasNodeSignalHandling(path) {
		found = append(found, "Node.js signal handling")
	}

	// Check Java files for shutdown hooks
	if c.hasJavaShutdownHook(path) {
		found = append(found, "Java shutdown hook")
	}

	// Check for graceful shutdown in Kubernetes configs
	if c.hasK8sGracePeriod(path) {
		found = append(found, "K8s terminationGracePeriodSeconds")
	}

	// Check for preStop hooks in K8s
	if c.hasK8sPreStopHook(path) {
		found = append(found, "K8s preStop hook")
	}

	// Build result
	found = unique(found)
	if len(found) > 0 {
		return rb.Pass("Graceful shutdown configured: " + strings.Join(found, ", ")), nil
	}
	return rb.Warn("No graceful shutdown handling found (handle SIGTERM/SIGINT for clean shutdown)"), nil
}

func (c *ShutdownCheck) hasGoSignalHandling(path string) bool {
	if !safepath.Exists(path, "go.mod") {
		return false
	}

	// Patterns indicating signal handling in Go
	patterns := []string{
		"signal.Notify",
		"os.Signal",
		"syscall.SIGTERM",
		"syscall.SIGINT",
		"context.WithCancel",
		"server.Shutdown",
		"srv.Shutdown",
		"http.Server",
	}

	// Also check for graceful shutdown packages
	gracefulPackages := []string{
		"github.com/oklog/run",
		"golang.org/x/sync/errgroup",
	}

	if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
		for _, pkg := range gracefulPackages {
			if strings.Contains(string(content), pkg) {
				return true
			}
		}
	}

	// Search in Go source files - check root and common directories
	goPatterns := []string{"*.go", "cmd/*.go", "cmd/*/*.go", "internal/*.go", "internal/*/*.go", "pkg/*.go", "pkg/*/*.go"}
	for _, pattern := range goPatterns {
		if files, err := safepath.Glob(path, pattern); err == nil {
			for _, file := range files {
				if content, err := safepath.ReadFileAbs(file); err == nil {
					for _, p := range patterns {
						if strings.Contains(string(content), p) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func (c *ShutdownCheck) hasPythonSignalHandling(path string) bool {
	pythonConfigs := []string{"pyproject.toml", "requirements.txt", "setup.py"}
	hasPython := false
	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			hasPython = true
			break
		}
	}
	if !hasPython {
		return false
	}

	// Patterns indicating signal handling in Python
	patterns := []string{
		"signal.signal",
		"signal.SIGTERM",
		"signal.SIGINT",
		"atexit.register",
		"asyncio.get_event_loop().add_signal_handler",
		"loop.add_signal_handler",
		"GracefulKiller",
		"shutdown_event",
	}

	// Search in Python source files - check root and common directories
	pyPatterns := []string{"*.py", "src/*.py", "src/*/*.py", "app/*.py", "app/*/*.py"}
	for _, pattern := range pyPatterns {
		if files, err := safepath.Glob(path, pattern); err == nil {
			for _, file := range files {
				// Skip test files and venv - check basename only
				baseName := filepath.Base(file)
				if strings.Contains(baseName, "test") || strings.Contains(file, "venv") {
					continue
				}
				if content, err := safepath.ReadFileAbs(file); err == nil {
					for _, p := range patterns {
						if strings.Contains(string(content), p) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func (c *ShutdownCheck) hasNodeSignalHandling(path string) bool {
	if !safepath.Exists(path, "package.json") {
		return false
	}

	// Patterns indicating signal handling in Node.js
	patterns := []string{
		`process.on('SIGTERM'`,
		`process.on("SIGTERM"`,
		`process.on('SIGINT'`,
		`process.on("SIGINT"`,
		`process.on('beforeExit'`,
		`process.on("beforeExit"`,
		`process.on('exit'`,
		`process.on("exit"`,
		"graceful-shutdown",
		"http-terminator",
		"stoppable",
	}

	// Also check package.json for graceful shutdown packages
	if content, err := safepath.ReadFile(path, "package.json"); err == nil {
		gracefulPackages := []string{
			"http-terminator",
			"stoppable",
			"@godaddy/terminus",
			"lightship",
		}
		for _, pkg := range gracefulPackages {
			if strings.Contains(string(content), pkg) {
				return true
			}
		}
	}

	// Search in JS/TS source files - check root and common directories
	nodePatterns := []string{"*.js", "*.ts", "src/*.js", "src/*.ts", "src/*/*.js", "src/*/*.ts", "lib/*.js", "lib/*.ts"}
	for _, pattern := range nodePatterns {
		if files, err := safepath.Glob(path, pattern); err == nil {
			for _, file := range files {
				// Skip node_modules and test files - check basename only for test
				baseName := filepath.Base(file)
				if strings.Contains(file, "node_modules") || strings.Contains(baseName, "test") {
					continue
				}
				if content, err := safepath.ReadFileAbs(file); err == nil {
					for _, p := range patterns {
						if strings.Contains(string(content), p) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func (c *ShutdownCheck) hasJavaShutdownHook(path string) bool {
	javaConfigs := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	hasJava := false
	for _, cfg := range javaConfigs {
		if safepath.Exists(path, cfg) {
			hasJava = true
			break
		}
	}
	if !hasJava {
		return false
	}

	// Patterns indicating shutdown hooks in Java
	patterns := []string{
		"Runtime.getRuntime().addShutdownHook",
		"addShutdownHook",
		"@PreDestroy",
		"DisposableBean",
		"SmartLifecycle",
		"GracefulShutdown",
	}

	// Search in Java source files - check common Maven/Gradle directory structure
	javaPatterns := []string{
		"src/main/java/*.java",
		"src/main/java/*/*.java",
		"src/main/java/*/*/*.java",
		"src/main/java/*/*/*/*.java",
		"src/main/java/*/*/*/*/*.java",
		"src/*.java",
	}
	for _, pattern := range javaPatterns {
		if files, err := safepath.Glob(path, pattern); err == nil {
			for _, file := range files {
				if content, err := safepath.ReadFileAbs(file); err == nil {
					for _, p := range patterns {
						if strings.Contains(string(content), p) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func (c *ShutdownCheck) hasK8sGracePeriod(path string) bool {
	k8sDirs := []string{"k8s", "kubernetes", "deploy", "deployment", "manifests", ""}
	for _, dir := range k8sDirs {
		searchPath := path
		if dir != "" {
			searchPath = path + "/" + dir
		}
		if files, err := safepath.Glob(searchPath, "*.yaml"); err == nil {
			for _, file := range files {
				if content, err := safepath.ReadFileAbs(file); err == nil {
					if strings.Contains(string(content), "terminationGracePeriodSeconds") {
						return true
					}
				}
			}
		}
		if files, err := safepath.Glob(searchPath, "*.yml"); err == nil {
			for _, file := range files {
				if content, err := safepath.ReadFileAbs(file); err == nil {
					if strings.Contains(string(content), "terminationGracePeriodSeconds") {
						return true
					}
				}
			}
		}
	}
	return false
}

func (c *ShutdownCheck) hasK8sPreStopHook(path string) bool {
	k8sDirs := []string{"k8s", "kubernetes", "deploy", "deployment", "manifests", ""}
	for _, dir := range k8sDirs {
		searchPath := path
		if dir != "" {
			searchPath = path + "/" + dir
		}
		if files, err := safepath.Glob(searchPath, "*.yaml"); err == nil {
			for _, file := range files {
				if content, err := safepath.ReadFileAbs(file); err == nil {
					if strings.Contains(string(content), "preStop") {
						return true
					}
				}
			}
		}
		if files, err := safepath.Glob(searchPath, "*.yml"); err == nil {
			for _, file := range files {
				if content, err := safepath.ReadFileAbs(file); err == nil {
					if strings.Contains(string(content), "preStop") {
						return true
					}
				}
			}
		}
	}
	return false
}
