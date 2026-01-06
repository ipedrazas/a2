# Contributing to A2

Thank you for your interest in contributing to A2! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Adding New Checks](#adding-new-checks)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful and constructive in all interactions.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/a2.git
   cd a2
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/ipedrazas/a2.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- [Task](https://taskfile.dev/) (optional, but recommended for running tasks)
- Git

### Initial Setup

1. **Install dependencies**:
   ```bash
   task dependencies
   ```
   Or manually:
   ```bash
   go mod download
   go install golang.org/x/vuln/cmd/govulncheck@latest
   ```

2. **Verify setup**:
   ```bash
   task build
   task test
   ```

## Project Structure

```
a2/
├── cmd/              # CLI commands (Cobra)
├── pkg/
│   ├── checker/      # Core checker interfaces and types
│   ├── checks/       # Individual check implementations
│   ├── config/       # Configuration loading
│   ├── output/       # Output formatters (pretty, JSON)
│   └── runner/       # Check execution engine
├── main.go           # Entry point
├── Taskfile.yaml     # Task runner configuration
└── .a2.yaml          # Project configuration
```

## Development Workflow

### 1. Create a Branch

Create a feature branch from `main`:

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write clean, readable code
- Follow existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

Before submitting, ensure all tests pass:

```bash
# Run all tests
task test

# Run tests with coverage
task test:coverage

# Check coverage by package
task test:coverage:bypackage

# Run linting
task lint

# Format code
task fmt

# Run full CI checks
task ci
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add new check for X"
```

**Commit Message Format:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions/changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Testing

### Running Tests

A2 uses `testify/suite` for all unit tests. Tests are organized by package:

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/checks/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Test Coverage Requirements

- **Target**: > 80% coverage for all packages
- **Current Status**: Check with `task test:coverage:bypackage`
- **New Code**: Must include tests with > 80% coverage

### Writing Tests

All tests should use `testify/suite`:

```go
package checks

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type MyCheckTestSuite struct {
    suite.Suite
}

func (suite *MyCheckTestSuite) SetupTest() {
    // Setup code
}

func (suite *MyCheckTestSuite) TestMyCheck_Something() {
    // Test implementation
}

func TestMyCheckTestSuite(t *testing.T) {
    suite.Run(t, new(MyCheckTestSuite))
}
```

See existing test files for examples:
- `pkg/runner/runner_test.go`
- `pkg/checks/files_test.go`
- `pkg/config/config_test.go`

## Adding New Checks

### 1. Create the Check

Create a new file in `pkg/checks/`:

```go
package checks

import (
    "github.com/ipedrazas/a2/pkg/checker"
)

type MyCheck struct {
    // Configuration fields
}

func (c *MyCheck) ID() string {
    return "my_check"
}

func (c *MyCheck) Name() string {
    return "My Check"
}

func (c *MyCheck) Run(path string) (checker.Result, error) {
    // Implementation
    return checker.Result{
        Name:    c.Name(),
        ID:      c.ID(),
        Passed:  true,
        Status:  checker.Pass,
        Message: "Check passed",
    }, nil
}
```

### 2. Register the Check

Add to `pkg/checks/registry.go`:

```go
func GetChecks(cfg *config.Config) []checker.Checker {
    allChecks := []checker.Checker{
        // ... existing checks
        &MyCheck{},
    }
    // ...
}
```

### 3. Add Tests

Create `pkg/checks/mycheck_test.go` with comprehensive tests.

### 4. Update Documentation

- Update README.md with the new check
- Add to `.a2.yaml` example if configurable

## Code Style

### Go Style Guide

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (run `task fmt`)
- Follow existing code patterns
- Keep functions focused and small
- Use meaningful variable names

### Code Formatting

```bash
# Format all code
task fmt

# Check formatting
task lint
```

### Linting

```bash
# Run go vet
go vet ./...

# Check formatting
gofmt -l .
```

## Submitting Changes

### Pull Request Process

1. **Update your branch**:
   ```bash
   git checkout main
   git pull upstream main
   git checkout feature/your-feature-name
   git rebase main
   ```

2. **Ensure all checks pass**:
   ```bash
   task ci
   ```

3. **Create Pull Request**:
   - Clear title and description
   - Reference related issues
   - Include test coverage information
   - Add screenshots for UI changes

### Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added/updated and passing
- [ ] Test coverage > 80% for new code
- [ ] Documentation updated
- [ ] Commit messages are clear
- [ ] No merge conflicts
- [ ] All CI checks pass

### Review Process

- Maintainers will review your PR
- Address feedback promptly
- Keep PRs focused and reasonably sized
- Be open to suggestions and improvements

## Documentation

### Code Documentation

- Add Go doc comments for exported functions/types
- Keep comments concise and clear
- Update README.md for user-facing changes

### Example

```go
// MyCheck verifies that something is correct.
// It checks the given path and returns a Result.
type MyCheck struct {
    Threshold float64 // Minimum threshold value
}

// Run executes the check against the given path.
// Returns a Result and any error encountered during execution.
func (c *MyCheck) Run(path string) (checker.Result, error) {
    // ...
}
```

## Common Tasks

### Running A2 on the Project

```bash
# Build and run
task run

# Run with specific path
task run -- /path/to/project

# JSON output
./dist/a2 check --format json
```

### Building

```bash
# Build binary
task build

# Build for all platforms
task release
```

### Docker

```bash
# Build Docker image
task docker:build

# Run in Docker
task docker:run
```

## Getting Help

- **Issues**: Open an issue on GitHub for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check README.md and code comments

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to A2! 🎉

