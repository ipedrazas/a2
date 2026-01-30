# BDD Check Implementation Plan

## Overview
Add a new `common:bdd` check that detects BDD/Gherkin test implementation across multiple languages (Go, Python, Node.js, Ruby, Java).

## Scope
**Detection + optional execution**:
- **Primary mode**: Detects BDD framework and feature files (fast, non-invasive)
- **Optional execution**: Users can configure test execution via `.a2.yaml` external checks

This approach provides immediate value through detection while allowing users to opt-in to test execution when needed.

## Files to Create/Modify

### 1. Create: `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/checks/common/bdd.go`
Main check implementation following the pattern from `integration.go`.

### 2. Create: `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/checks/common/bdd_test.go`
Comprehensive test suite using testify/suite.

### 3. Modify: `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/checks/common/register.go`
Add BDD check registration after `common:contributing` (order 970).

## Implementation Details

### Check Structure

```go
type BDDCheck struct{}

func (c *BDDCheck) ID() string   { return "common:bdd" }
func (c *BDDCheck) Name() string { return "BDD Tests" }

func (c *BDDCheck) Run(path string) (checker.Result, error) {
    rb := checkutil.NewResultBuilder(c, checker.LangCommon)

    var detectedFrameworks []string
    var totalFeatures int
    var hasSteps bool

    // Detect Go BDD (godog)
    if fw := c.detectGoBDD(path); fw != nil {
        detectedFrameworks = append(detectedFrameworks, fw.Name)
        totalFeatures += fw.FeatureCount
        hasSteps = hasSteps || fw.HasSteps
    }

    // Detect Python BDD (behave, pytest-bdd)
    if fw := c.detectPythonBDD(path); fw != nil {
        detectedFrameworks = append(detectedFrameworks, fw.Name)
        totalFeatures += fw.FeatureCount
        hasSteps = hasSteps || fw.HasSteps
    }

    // Detect Node.js BDD (cucumber-js)
    if fw := c.detectNodeBDD(path); fw != nil {
        detectedFrameworks = append(detectedFrameworks, fw.Name)
        totalFeatures += fw.FeatureCount
        hasSteps = hasSteps || fw.HasSteps
    }

    // Detect Ruby BDD (cucumber)
    if fw := c.detectRubyBDD(path); fw != nil {
        detectedFrameworks = append(detectedFrameworks, fw.Name)
        totalFeatures += fw.FeatureCount
        hasSteps = hasSteps || fw.HasSteps
    }

    // Detect Java BDD (cucumber-jvm)
    if fw := c.detectJavaBDD(path); fw != nil {
        detectedFrameworks = append(detectedFrameworks, fw.Name)
        totalFeatures += fw.FeatureCount
        hasSteps = hasSteps || fw.HasSteps
    }

    return c.buildResult(rb, detectedFrameworks, totalFeatures, hasSteps), nil
}
```

### Detection Patterns by Language

#### Go (godog)
- **Dependency check**: `go.mod` contains `github.com/cucumber/godog`
- **Feature files**: `**/*.feature`
- **Step definitions**: `**/*_steps.go`

#### Python (behave, pytest-bdd, radish)
- **Dependency check**: `pyproject.toml`, `requirements.txt`, `setup.py`, `Pipfile` contains `behave`, `pytest-bdd`, or `radish`
- **Feature files**: `features/**/*.feature`
- **Step definitions**: `features/steps/*.py`, `*_steps.py`

#### Node.js (cucumber-js)
- **Dependency check**: `package.json` contains `@cucumber/cucumber` or `cucumber`
- **Feature files**: `features/**/*.feature`
- **Step definitions**: `features/step-definitions/**/*.js`, `**/*.steps.js`, `**/*.steps.ts`

#### Ruby (cucumber)
- **Dependency check**: `Gemfile` contains `cucumber`
- **Feature files**: `features/**/*.feature`
- **Step definitions**: `features/step_definitions/*.rb`

#### Java (cucumber-jvm)
- **Dependency check**: `pom.xml` or `build.gradle` contains `cucumber`
- **Feature files**: `src/test/resources/features/**/*.feature`
- **Step definitions**: `**/*Steps.java`, `**/*StepDefs.java`

### Result Logic

```go
func (c *BDDCheck) buildResult(rb *checkutil.ResultBuilder,
    frameworks []string, featureCount int, hasSteps bool) checker.Result {

    if len(frameworks) == 0 {
        return rb.Info("No BDD framework detected (optional)"), nil
    }

    if featureCount == 0 {
        return rb.Info("BDD framework detected: " + strings.Join(frameworks, ", ") +
            " but no feature files found"), nil
    }

    if !hasSteps && featureCount > 0 {
        return rb.Warn(fmt.Sprintf("Found %d feature file(s) but no step definitions detected",
            featureCount)), nil
    }

    msg := fmt.Sprintf("BDD tests found: %s with %d feature file(s)",
        strings.Join(frameworks, ", "), featureCount)
    return rb.Pass(msg), nil
}
```

### Registration (in register.go)

Add after `common:contributing` (line 250):

```go
{
    Checker: &BDDCheck{},
    Meta: checker.CheckMeta{
        ID:          "common:bdd",
        Name:        "BDD Tests",
        Description: "Detects BDD/Gherkin test implementation using Cucumber, godog, behave, or cucumber-js. Validates feature files and step definitions exist.",
        Languages:   []checker.Language{checker.LangCommon},
        Critical:    false,
        Optional:    true,
        Order:       975,
        Suggestion:  "Add BDD tests using Cucumber/godog/behave for behavior verification",
    },
},
```

## Test Coverage

### Test Cases

1. **Go BDD Detection**
   - godog with feature files → Pass
   - godog without feature files → Info
   - feature files without steps → Warn
   - no Go BDD → Info

2. **Python BDD Detection**
   - behave with features/ → Pass
   - pytest-bdd dependency → Pass
   - radish framework → Pass
   - no Python BDD → Info

3. **Node.js BDD Detection**
   - cucumber-js in package.json → Pass
   - @cucumber/cucumber → Pass
   - features/ directory → Pass
   - no Node BDD → Info

4. **Ruby BDD Detection**
   - cucumber in Gemfile → Pass
   - features/support/env.rb → Pass
   - no Ruby BDD → Info

5. **Java BDD Detection**
   - cucumber-jvm in pom.xml → Pass
   - cucumber in build.gradle → Pass
   - no Java BDD → Info

6. **Multi-language Projects**
   - Go + Python BDD → Pass with both frameworks
   - Node + Java BDD → Pass with both frameworks

7. **Edge Cases**
   - Feature files without step definitions → Warn
   - Empty features/ directory → Info
   - Multiple BDD frameworks in same language → Pass

## Verification

### Manual Testing
Run against:
1. **A2 itself** (current repo) - Should detect godog with 10 feature files
2. **Plain Go project** - Should return Info
3. **Python behave project** - Should detect behave
4. **Node cucumber project** - Should detect cucumber-js

### Commands
```bash
# Build
task build

# Run all checks
./dist/a2 check

# Run just BDD check
./dist/a2 run common:bdd

# Skip BDD check
./dist/a2 check --skip common:bdd
```

## Expected Output Examples

### A2 (current repo)
```
✓ PASS BDD Tests (1ms) - common:bdd
    BDD tests found: godog with 10 feature files
```

### Project with feature files but no steps
```
⚠ WARN BDD Tests (1ms) - common:bdd
    Found 3 feature file(s) but no step definitions detected
    Suggestion: Add BDD tests using Cucumber/godog/behave for behavior verification
```

### No BDD framework
```
ℹ INFO BDD Tests (1ms) - common:bdd
    No BDD framework detected (optional)
    Suggestion: Add BDD tests using Cucumber/godog/behave for behavior verification
```

## Optional Test Execution

Users can configure BDD test execution via external checks in `.a2.yaml`:

```yaml
external:
  - id: bdd:execute:go
    name: "Run Go BDD Tests"
    command: go
    args: ["test", "./pkg/godogs/..."]
    severity: warn

  - id: bdd:execute:python
    name: "Run Python BDD Tests"
    command: behave
    args: ["features/"]
    severity: warn

  - id: bdd:execute:node
    name: "Run Node BDD Tests"
    command: npm
    args: ["test", "--", "--cucumber"]
    severity: warn
```

This allows:
- Language-specific test commands
- Integration with CI pipelines
- Custom severity levels
- Opt-in execution (only runs when configured)

## Dependencies
No new dependencies - uses existing:
- `github.com/ipedrazas/a2/pkg/checker`
- `github.com/ipedrazas/a2/pkg/checkutil`
- `github.com/ipedrazas/a2/pkg/safepath`
- Standard library: `fmt`, `strings`, `filepath`

## Implementation Sequence

1. Create `pkg/checks/common/bdd.go` with Go and Python detection
2. Add Node.js, Ruby, Java detection
3. Create `pkg/checks/common/bdd_test.go` with comprehensive tests
4. Register in `pkg/checks/common/register.go`
5. Run `task test` to verify
6. Run `./dist/a2 check` to verify integration
