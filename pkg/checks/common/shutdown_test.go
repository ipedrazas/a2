package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ShutdownCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *ShutdownCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "shutdown-test-*")
	s.Require().NoError(err)
}

func (s *ShutdownCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ShutdownCheckTestSuite) TestIDAndName() {
	check := &ShutdownCheck{}
	s.Equal("common:shutdown", check.ID())
	s.Equal("Graceful Shutdown", check.Name())
}

func (s *ShutdownCheckTestSuite) TestGoSignalNotify() {
	// Create go.mod
	goMod := `module myapp

go 1.21
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	// Create main.go with signal handling
	content := `package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Go signal handling")
}

func (s *ShutdownCheckTestSuite) TestGoServerShutdown() {
	goMod := `module myapp

go 1.21
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	content := `package main

import (
	"context"
	"net/http"
	"time"
)

func main() {
	srv := &http.Server{Addr: ":8080"}

	go func() {
		srv.ListenAndServe()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Go signal handling")
}

func (s *ShutdownCheckTestSuite) TestGoOklogRun() {
	goMod := `module myapp

go 1.21

require github.com/oklog/run v1.1.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Go signal handling")
}

func (s *ShutdownCheckTestSuite) TestPythonSignalHandler() {
	// Create pyproject.toml
	pyproject := `[project]
name = "myapp"
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(pyproject), 0644)
	s.Require().NoError(err)

	// Create main.py with signal handling
	content := `import signal
import sys

def signal_handler(sig, frame):
    print('Shutting down gracefully...')
    sys.exit(0)

signal.signal(signal.SIGTERM, signal_handler)
signal.signal(signal.SIGINT, signal_handler)

if __name__ == '__main__':
    # Main application logic
    pass
`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.py"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Python signal handling")
}

func (s *ShutdownCheckTestSuite) TestPythonAtexit() {
	pyproject := `[project]
name = "myapp"
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(pyproject), 0644)
	s.Require().NoError(err)

	content := `import atexit

def cleanup():
    print('Cleaning up...')

atexit.register(cleanup)

if __name__ == '__main__':
    pass
`
	err = os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Python signal handling")
}

func (s *ShutdownCheckTestSuite) TestNodeProcessOnSIGTERM() {
	// Create package.json
	packageJson := `{
  "name": "myapp",
  "version": "1.0.0"
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(packageJson), 0644)
	s.Require().NoError(err)

	// Create index.js with signal handling
	content := `const http = require('http');

const server = http.createServer((req, res) => {
  res.end('Hello World');
});

process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully');
  server.close(() => {
    console.log('Server closed');
    process.exit(0);
  });
});

server.listen(3000);
`
	err = os.WriteFile(filepath.Join(s.tempDir, "index.js"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Node.js signal handling")
}

func (s *ShutdownCheckTestSuite) TestNodeHttpTerminator() {
	packageJson := `{
  "name": "myapp",
  "dependencies": {
    "http-terminator": "^3.2.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(packageJson), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Node.js signal handling")
}

func (s *ShutdownCheckTestSuite) TestNodeTerminus() {
	packageJson := `{
  "name": "myapp",
  "dependencies": {
    "@godaddy/terminus": "^4.12.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(packageJson), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Node.js signal handling")
}

func (s *ShutdownCheckTestSuite) TestJavaShutdownHook() {
	// Create pom.xml
	pomXml := `<project>
  <groupId>com.example</groupId>
  <artifactId>myapp</artifactId>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomXml), 0644)
	s.Require().NoError(err)

	// Create Java source with shutdown hook
	srcDir := filepath.Join(s.tempDir, "src", "main", "java", "com", "example")
	err = os.MkdirAll(srcDir, 0755)
	s.Require().NoError(err)

	content := `package com.example;

public class Application {
    public static void main(String[] args) {
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("Shutting down gracefully...");
        }));
    }
}
`
	err = os.WriteFile(filepath.Join(srcDir, "Application.java"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Java shutdown hook")
}

func (s *ShutdownCheckTestSuite) TestJavaPreDestroy() {
	pomXml := `<project>
  <groupId>com.example</groupId>
  <artifactId>myapp</artifactId>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomXml), 0644)
	s.Require().NoError(err)

	srcDir := filepath.Join(s.tempDir, "src", "main", "java", "com", "example")
	err = os.MkdirAll(srcDir, 0755)
	s.Require().NoError(err)

	content := `package com.example;

import javax.annotation.PreDestroy;

public class MyService {
    @PreDestroy
    public void cleanup() {
        System.out.println("Cleaning up...");
    }
}
`
	err = os.WriteFile(filepath.Join(srcDir, "MyService.java"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Java shutdown hook")
}

func (s *ShutdownCheckTestSuite) TestK8sTerminationGracePeriod() {
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	content := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      terminationGracePeriodSeconds: 30
      containers:
        - name: myapp
          image: myapp:latest
`
	err = os.WriteFile(filepath.Join(k8sDir, "deployment.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "K8s terminationGracePeriodSeconds")
}

func (s *ShutdownCheckTestSuite) TestK8sPreStopHook() {
	k8sDir := filepath.Join(s.tempDir, "kubernetes")
	err := os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	content := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
        - name: myapp
          image: myapp:latest
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh", "-c", "sleep 5"]
`
	err = os.WriteFile(filepath.Join(k8sDir, "deployment.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "K8s preStop hook")
}

func (s *ShutdownCheckTestSuite) TestK8sInRootDirectory() {
	content := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      terminationGracePeriodSeconds: 60
      containers:
        - name: myapp
          image: myapp:latest
`
	err := os.WriteFile(filepath.Join(s.tempDir, "deployment.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "K8s terminationGracePeriodSeconds")
}

func (s *ShutdownCheckTestSuite) TestMultipleShutdownMechanisms() {
	// Go module
	goMod := `module myapp

go 1.21
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	// Go signal handling
	goContent := `package main

import (
	"os/signal"
)

func main() {
	signal.Notify(nil)
}
`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(goContent), 0644)
	s.Require().NoError(err)

	// K8s config
	k8sDir := filepath.Join(s.tempDir, "k8s")
	err = os.MkdirAll(k8sDir, 0755)
	s.Require().NoError(err)

	k8sContent := `apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      terminationGracePeriodSeconds: 30
`
	err = os.WriteFile(filepath.Join(k8sDir, "deployment.yaml"), []byte(k8sContent), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Go signal handling")
	s.Contains(result.Message, "K8s terminationGracePeriodSeconds")
}

func (s *ShutdownCheckTestSuite) TestNoShutdownHandling() {
	// Create a simple Go project without shutdown handling
	goMod := `module myapp

go 1.21
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No graceful shutdown handling found")
}

func (s *ShutdownCheckTestSuite) TestEmptyDirectory() {
	check := &ShutdownCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No graceful shutdown handling found")
}

func TestShutdownCheckTestSuite(t *testing.T) {
	suite.Run(t, new(ShutdownCheckTestSuite))
}
