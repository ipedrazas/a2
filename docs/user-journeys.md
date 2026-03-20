# A2 - User Journeys

## Document Information

| Field | Value |
|-------|-------|
| **Product Name** | A2 (Application Analysis Tool) |
| **Version** | 1.0 |
| **Last Updated** | 2025-01-30 |
| **Status** | Active Development |

---

## Table of Contents

1. [User Personas](#user-personas)
2. [Quick Start Journey](#quick-start-journey)
3. [First-Time Setup Journey](#first-time-setup-journey)
4. [Daily Development Journey](#daily-development-journey)
5. [CI/CD Integration Journey](#cicd-integration-journey)
6. [AI-Assisted Development Journey](#ai-assisted-development-journey)
7. [Open Source Maintenance Journey](#open-source-maintenance-journey)
8. [Project Assessment Journey](#project-assessment-journey)
9. [Server Mode Journey](#server-mode-journey)
10. [Customization Journey](#customization-journey)

---

## User Personas

### 1. Sarah - Senior Software Engineer
**Background**: 8 years experience, works on a microservices API
**Goals**: Ensure code quality, catch issues early, maintain team standards
**Pain Points**: AI-generated code from junior devs lacks quality, time-consuming code reviews
**Technical Level**: Expert

### 2. Alex - AI-Assisted Developer
**Background**: 3 years experience, heavily uses AI coding assistants
**Goals**: Quickly validate AI-generated code, learn best practices
**Pain Points**: Uncertainty about AI code quality, missing edge cases
**Technical Level**: Intermediate

### 3. Jordan - DevOps Engineer
**Background**: 5 years DevOps experience, manages CI/CD pipelines
**Goals**: Integrate quality gates, automate quality checks
**Pain Points**: Fragmented tooling, inconsistent quality standards
**Technical Level**: Expert

### 4. Maria - Open Source Maintainer
**Background**: Maintains popular GitHub project, receives many PRs
**Goals**: Validate contributions efficiently, maintain project quality
**Pain Points**: AI-generated PRs varying in quality, limited review time
**Technical Level**: Expert

### 5. Carlos - Engineering Manager
**Background**: 12 years experience, leads team of 15 developers
**Goals**: Establish quality metrics, track team progress
**Pain Points**: Lack of measurable quality standards, inconsistent practices
**Technical Level**: Intermediate

---

## Quick Start Journey

### Journey: Run First Check in Under 5 Minutes

**Persona**: Any new user
**Goal**: Get immediate value from A2 with minimal setup
**Prerequisites**: Go project directory

#### Step 1: Installation (1 minute)
```bash
go install github.com/ipedrazas/a2@latest


# Verify installation
a2 --version
```

**Expected Outcome**: `a2` command available in PATH

#### Step 2: Navigate to Project (30 seconds)
```bash
cd ~/my-project
```

#### Step 3: Run First Check (2 minutes)
```bash
a2 check
```

**Expected Behavior**:
1. A2 auto-detects language (e.g., "Detected: Go")
2. Runs appropriate checks in parallel
3. Shows progress indicators
4. Displays results with color coding

**Sample Output**:
```
✓ Go: Module check passed
✓ Go: Build successful
✗ Go: Tests failed - 3 tests failing
⚠ Go: Coverage at 45% (threshold: 80%)
⚠ Go: Format issues found

Maturity Level: Development (65% score, 1 failure)
Suggestion: Run 'go test ./...' to fix failing tests
```

#### Step 4: Interpret Results (1 minute)
- **Green (✓)**: Everything good
- **Red (✗)**: Critical issue, must fix
- **Yellow (⚠)**: Recommendation, improve if time permits
- **Blue (ℹ)**: Info only (optional tool not installed)

#### Step 5: Fix and Re-run
```bash
# Fix issues, then verify
a2 check
```

**Success Criteria**:
- Installation successful
- Checks executed without errors
- User understands basic output
- Clear path to improvement provided

**Potential Friction Points**:
- Missing required tools (A2 shows Info messages)
- First run takes longer (caching effects)
- Overwhelming output for large projects

**Improvements**:
- Quick start guide in first run
- Link to documentation for each check
- `--fix` flag for automatic fixes (future)

---

## First-Time Setup Journey

### Journey: Configure A2 for a New Project

**Persona**: Sarah (Senior Software Engineer)
**Goal**: Set up A2 with optimal configuration for an API project
**Prerequisites**: New or existing project, A2 installed

#### Phase 1: Interactive Configuration

**Step 1: Initialize Configuration**
```bash
cd ~/my-api-project
a2 add -i
```

**Interactive Prompts**:
```
? What type of application is this?
  ▸ API
    CLI
    Library
    Desktop

? What maturity level describes this project?
  ▸ Production
    Mature
    Development
    Proof of Concept

? Which languages should be checked?
  ▸ Auto-detect (recommended)
    Go
    Python
    TypeScript
    Multiple (select space-separated)

? Enable external checks?
  ▸ No (skip for now)
    Yes, I'll provide custom commands
```

**Expected Outcome**: Creates `.a2.yaml` with sensible defaults

```yaml
# .a2.yaml (generated)
language:
  explicit: []

profile: api
target: production

checks:
  disabled:
    - "*:e2e"  # API profile skips E2E tests

execution:
  parallel: true
  timeout: 300

files:
  required:
    - README.md
    - LICENSE
```

#### Phase 2: Customization

**Step 2: Adjust Coverage Thresholds**
```bash
# Edit .a2.yaml
```

Sarah knows her team maintains high standards:
```yaml
language:
  go:
    coverage_threshold: 85  # Higher than default
    cyclomatic_threshold: 10  # Stricter complexity
```

**Step 3: Add Required Documentation**
```yaml
files:
  required:
    - README.md
    - LICENSE
    - CONTRIBUTING.md
    - docs/api.md
    - .env.example
```

**Step 4: Disable Irrelevant Checks**
```yaml
checks:
  disabled:
    - "*:e2e"
    - "common:k8s"  # Not using Kubernetes yet
    - "python:*"    # Only Go in this project
```

#### Phase 3: Team Adoption

**Step 5: Commit Configuration**
```bash
git add .a2.yaml
git commit -m "Add A2 quality checks configuration"
git push
```

**Step 6: Team Communication**
Sarah creates a PR description:
```
Added A2 quality checker to our workflow.

Run locally: `a2 check`

Configuration is in `.a2.yaml`. We're using:
- API profile (web service focus)
- Production target (strict standards)
- 85% coverage threshold

Checks run automatically in CI.
```

**Success Criteria**:
- Configuration file created and committed
- Team understands how to use A2
- CI pipeline configured (see CI/CD journey)
- First CI run completes successfully

**Potential Friction Points**:
- Team debates on thresholds
- Configuration too complex for new users
- Fear of blocking deployments

**Improvements**:
- Team configuration templates
- Gradual rollout (warn before fail)
- Documentation generator for custom checks

---

## Daily Development Journey

### Journey: Local Development Workflow

**Persona**: Alex (AI-Assisted Developer)
**Goal**: Ensure code quality before committing
**Frequency**: Multiple times per day

#### Scenario 1: After AI-Generated Code

**Step 1: Generate Code with AI**
Alex uses Claude/ChatGPT to generate a new API endpoint:
```
"Create a Go REST endpoint for user authentication with JWT"
```

**Step 2: Run A2 Check**
```bash
a2 check
```

**Output**:
```
✓ Go: Module check passed
✓ Go: Build successful
✗ Go: Tests failed - no tests found
⚠ Go: Format issues found
⚠ Go: Logging not structured
✓ Go: No vulnerabilities found
```

**Step 3: Fix Issues**
```bash
# Format code
go fmt ./...

# Re-run check
a2 check
```

**Output**:
```
✓ Go: Module check passed
✓ Go: Build successful
✗ Go: Tests failed - no tests found
✓ Go: Format correct
⚠ Go: Logging not structured
```

**Step 4: Add Tests**
```bash
# A2 tells Alex what's missing
a2 explain go:tests
```

**Output**:
```
Check: Go Tests
Description: Verifies test suite exists and passes
Tool: go test ./...
Requirements:
  - At least one _test.go file
  - All tests pass
  - No race conditions

Suggestion: Create auth_test.go with test cases
```

**Step 5: Re-check Before Commit**
```bash
a2 check
```

**Output**:
```
✓ Go: Module check passed
✓ Go: Build successful
✓ Go: Tests passed (5 tests)
✓ Go: Format correct
⚠ Go: Logging not structured (non-critical)

Maturity Level: Mature (92% score, 0 failures)
```

**Step 6: Commit with Confidence**
```bash
git add .
git commit -m "Add JWT authentication endpoint"
```

#### Scenario 2: Pre-Push Quality Gate

**Step 1: Work on Feature**
Alex develops a new feature for 2 hours.

**Step 2: Quick Check**
```bash
a2 check --output=toon
```

**Rapid Output** (for quick review):
```
results[8]{
  {id:"go:build",status:"pass",message:"Build successful"},
  {id:"go:tests",status:"pass",message:"12 tests passed"},
  {id:"go:coverage",status:"warn",message:"Coverage: 72% (target: 80%)"}
}
```

**Step 3: Address Warnings**
Alex sees coverage is 72% (target 80%). She adds edge case tests.

**Step 4: Final Check**
```bash
a2 check
```

**Output**:
```
✓ All checks passed!
Maturity Level: Production-Ready (100% score)
```

**Step 5: Push**
```bash
git push
```

**Success Criteria**:
- Catches quality issues before commit/push
- Fast feedback loop (< 2 minutes)
- Clear guidance on fixes
- Developer confidence increased

**Potential Friction Points**:
- Waiting for checks to complete
- False positives block work
- Unclear how to fix issues

**Improvements**:
- Incremental checking (only changed files)
- `--fix` flag for auto-fixable issues
- IDE integration with inline suggestions

---

## CI/CD Integration Journey

### Journey: Quality Gates in Production Pipeline

**Persona**: Jordan (DevOps Engineer)
**Goal**: Prevent low-quality code from reaching production
**Prerequisites**: GitHub Actions workflow, A2 configured

#### Phase 1: GitHub Actions Setup

**Step 1: Create Workflow File**
`.github/workflows/a2-check.yml`:
```yaml
name: A2 Quality Check

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  a2:
    name: Run A2 Checks
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Install A2
        run: |
          curl -sSL https://install.a2.dev | sh
          echo "$HOME/.a2/bin" >> $GITHUB_PATH

      - name: Run A2
        run: a2 check --output=json --format=json > results.json

      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: a2-results
          path: results.json

      - name: Check Exit Code
        run: |
          if [ $? -eq 2 ]; then
            echo "Critical failures detected"
            exit 1
          fi
```

**Step 2: Test Workflow**
```bash
git add .github/workflows/a2-check.yml
git commit -m "Add A2 to CI pipeline"
git push
```

**Expected Outcome**: Workflow runs on push, shows results in Actions tab

#### Phase 2: Quality Gate Configuration

**Step 3: Set Strict Enforcement**
Update workflow for production branch:
```yaml
      - name: Run A2
        run: a2 check --target=production

      - name: Fail on Warnings
        if: steps.a2.outputs.exit_code == 1
        run: |
          echo "Warnings not allowed in main branch"
          exit 1
```

**Step 4: Configure Branch Protection**
In GitHub repository settings:
- Require status check "A2 Quality Check" to pass
- Require branches to be up to date
- Enable "Require status checks to pass before merging"

#### Phase 3: Pull Request Integration

**Step 5: Add PR Comment Bot**
Create GitHub Action that comments on PRs:
```yaml
      - name: Comment Results on PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const results = JSON.parse(fs.readFileSync('results.json'));
            const score = results.maturity.score;
            const failures = results.results.filter(r => r.status === 'Fail');

            let comment = `## A2 Quality Report\n\n`;
            comment += `**Maturity Score:** ${score}%\n\n`;

            if (failures.length > 0) {
              comment += `### ❌ Failed Checks\n`;
              failures.forEach(f => {
                comment += `- **${f.name}**: ${f.message}\n`;
              });
            }

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });
```

**Sample PR Comment**:
```markdown
## A2 Quality Report

**Maturity Score:** 72%

### ❌ Failed Checks
- **Go: Tests**: 2 tests failing
  - TestCreateUserMissingEmail
  - TestCreateUserDuplicate

### ⚠️ Warnings
- **Go: Coverage**: 68% (target: 80%)
- **Go: Logging**: Unstructured logging detected

### 💡 Recommendations
Run the following to fix issues:
```bash
go test ./... -v
```

[View Full Results](https://github.com/.../actions/runs/12345)
```

#### Phase 4: Rollout Strategy

**Step 6: Gradual Rollout**
Jordan implements a phased approach:

**Week 1-2**: Report-only mode
```yaml
# Don't fail the workflow
continue-on-error: true
```

**Week 3-4**: Fail on critical issues only
```yaml
# Check for exit code 2 (failures), ignore 1 (warnings)
if: steps.a2.outputs.exit_code >= 2
```

**Week 5+**: Full enforcement
```yaml
# Fail on warnings too
if: steps.a2.outputs.exit_code >= 1
```

**Step 7: Monitor and Adjust**
Jordan tracks:
- Pipeline success rate
- Common failures
- Team feedback
- Time to resolution

**Success Criteria**:
- CI pipeline runs successfully
- Quality issues caught before merge
- Team adoption without blocking productivity
- Clear visibility into quality trends

**Potential Friction Points**:
- Too strict, blocks all PRs
- False positives frustrate developers
- Long execution times

**Improvements**:
- Caching for faster runs
- Gradual enforcement strategy
- Exception process for emergencies
- Parallel execution optimization

---

## AI-Assisted Development Journey

### Journey: Validating AI-Generated Code

**Persona**: Alex (AI-Assisted Developer)
**Use Case**: Heavy reliance on AI coding assistants
**Frequency**: 10+ times per day

#### Scenario: Building a New Feature with AI

**Step 1: AI Code Generation**
Alex prompts Claude:
```
Create a Go function that validates email addresses according to RFC 5322,
including unit tests and error handling
```

**Step 2: Immediate Validation**
Alex runs A2 before even looking at the code:
```bash
a2 check --filter=go:*
```

**Output**:
```
✓ Go: Build successful
✗ Go: Tests failed - 1 test failing
⚠ Go: Coverage at 60% for this package
⚠ Go: Cyclomatic complexity: 18 (threshold: 15)
```

**Step 3: Diagnose Issues**
```bash
a2 run go:tests --verbose
```

Shows:
```
=== FAIL: TestValidateEmail_EmptyString (0.00s)
    validator_test.go:25: Expected error, got nil
FAIL
```

**Step 4: Iterative Improvement**
Alex asks AI to fix:
```
The test is failing: TestValidateEmail_EmptyString expects an error for empty
strings but the current implementation returns nil. Fix this.
```

**Step 5: Re-validate**
```bash
a2 check
```

**Output**:
```
✓ Go: Build successful
✓ Go: Tests passed (8 tests)
⚠ Go: Coverage at 75%
⚠ Go: Cyclomatic complexity: 16
```

**Step 6: Final Polish**
Alex reduces complexity by extracting helper function.

**Step 7: Final Check**
```bash
a2 check
```

All checks pass! ✅

#### Scenario: Bulk AI Code Review

**Step 1: AI Refactors Large Module**
Alex asks AI to refactor authentication module (50+ files).

**Step 2: Comprehensive Check**
```bash
a2 check --output=json > refactor-results.json
```

**Step 3: Compare Results**
Alex compares before/after:
```bash
# Before refactoring
cat baseline.json | jq '.maturity.score'
# Output: 85%

# After refactoring
cat refactor-results.json | jq '.maturity.score'
# Output: 92%
```

**Improvement detected!** 🎉

**Step 4: Regression Testing**
Alex checks for new issues:
```bash
# Check for new failures
cat refactor-results.json | jq '.results[] | select(.status == "Fail")'
```

**Step 5: Commit**
```bash
git commit -m "Refactor auth module (A2: 92% maturity)"
```

#### Anti-Pattern: Over-Reliance on AI

**Step 1: AI Generates Code**
Alex accepts AI suggestions without review.

**Step 2: A2 Catches Issues**
```bash
a2 check
```

**Output**:
```
✗ Go: Build failed
✗ Go: Tests failed - 12 tests failing
✗ Go: Security vulnerability detected
⚠ Go: No tests for new code
⚠ Go: Missing error handling
```

**Step 3: Learning Moment**
Alex realizes:
- AI doesn't know project context
- AI makes assumptions that don't apply
- A2 catches what AI misses

**Step 4: Adjusted Workflow**
Now Alex always:
1. Generates code with AI
2. Runs `a2 check` immediately
3. Reviews and fixes issues
4. Iterates with AI if needed
5. Re-checks before commit

**Success Criteria**:
- AI code quality improved
- Fewer bugs in production
- Developer learns from A2 feedback
- Faster development cycle

**Key Insights**:
- A2 validates AI output
- Prevents "AI slop" from entering codebase
- Educational feedback loop
- Quality maintained at high velocity

---

## Open Source Maintenance Journey

### Journey: Validating External Contributions

**Persona**: Maria (Open Source Maintainer)
**Project**: Popular Go library with 1000+ stars
**Challenge**: 20+ PRs per week, varying quality

#### Scenario 1: Automated PR Quality Check

**Step 1: PR Submitted**
Contributor submits PR #234: "Add retry logic with exponential backoff"

**Step 2: CI Runs A2**
GitHub Actions workflow automatically runs:
```yaml
- name: A2 Check
  run: a2 check --profile=library --target=production
```

**Step 3: Bot Comments on PR**
```
## 🔍 A2 Quality Report for PR #234

**Maturity Score:** 75% (Target: 90%)

### ❌ Critical Issues (Must Fix)
- **Go: Tests**: 2 tests failing
- **Go: Format**: Code not formatted with gofmt

### ⚠️ Warnings (Recommended)
- **Go: Coverage**: 65% (target: 80%)
- **Go: Logging**: Unstructured logging in retry.go:45

### ℹ️ Information
- Tool `staticcheck` not installed (skipped)

### 📋 Next Steps
1. Run tests: `go test ./...`
2. Format code: `gofmt -w .`
3. Address warnings (optional)

Documentation: https://a2.dev/docs/library-profile
```

**Step 4: Maria Reviews**
Maria sees:
- Critical issues must be fixed before merge
- Clear guidance provided to contributor
- Saves Maria 30+ minutes of manual review

**Step 5: Contributor Fixes**
Contributor runs A2 locally:
```bash
a2 check
```

Fixes issues and pushes update.

**Step 6: Re-check**
CI runs again, all checks pass:
```
## 🎉 A2 Quality Report for PR #234

**Maturity Score:** 95% ✅

All checks passed! Ready for review.
```

**Step 7: Merge with Confidence**
Maria now reviews logic/architecture, knowing quality is assured.

#### Scenario 2: AI-Generated PRs

**Step 1: Suspicious PR**
PR #235: "Fix memory leak in cache" from first-time contributor.

**Step 2: A2 Analysis**
```
## 🔍 A2 Quality Report

**Maturity Score:** 40%

### ❌ Critical Issues
- **Go: Build**: Compilation failed
- **Go: Tests**: 7 tests failing
- **Go: Race**: Data race detected

### 🔒 Security Issues
- **Common: Secrets**: Potential API key detected
```

**Step 3: Immediate Rejection**
Maria sees:
- Generated code (likely AI)
- Quality too low
- Security issue (API key leaked!)

**Step 4: Template Response**
Maria comments:
```markdown
Thanks for the contribution! However, this PR has several issues:

1. Code doesn't compile
2. Tests are failing
3. Race condition detected
4. ⚠️ **Security**: Please remove the API key from line 42

Please run `a2 check` locally and address all issues before resubmitting.

We have a [contributing guide](CONTRIBUTING.md) that walks through the process.
```

#### Scenario 3: Batch PR Review

**Step 1: Weekly Review**
Maria has 15 PRs to review.

**Step 2: Prioritize by Quality**
She checks CI results:
- **Green (90%+)**: Review first (5 PRs)
- **Yellow (70-89%)**: Review second (7 PRs)
- **Red (<70%)**: Request fixes, review later (3 PRs)

**Step 3: Efficient Review**
High-quality PRs require less time:
- Green PRs: 10 minutes each (logic review only)
- Yellow PRs: 20 minutes each (logic + quality)
- Red PRs: Return to contributor

**Step 4: Time Saved**
- Before A2: 30 minutes × 15 PRs = 7.5 hours
- After A2: 10 min × 5 + 20 min × 7 + 5 min × 3 = 3 hours
- **Saved: 4.5 hours per week!**

**Success Criteria**:
- PR quality improved
- Review time reduced 60%
- Fewer bugs from contributions
- Positive contributor experience

**Potential Friction Points**:
- Contributors don't know how to use A2
- False positives discourage contributors
- Strict standards block good PRs

**Improvements**:
- Contributor guide with A2 instructions
- "Good first issue" templates
- Helpful error messages
- Mentorship program

---

## Project Assessment Journey

### Journey: Evaluating Existing Codebase

**Persona**: Carlos (Engineering Manager)
**Goal**: Assess current state of team's projects
**Use Case**: Quarterly quality assessment

#### Scenario 1: Team-Wide Quality Audit

**Step 1: Identify Projects**
Carlos's team maintains 8 microservices:
- `user-api`
- `order-service`
- `payment-gateway`
- `notification-service`
- `inventory-api`
- `shipping-service`
- `analytics-platform`
- `frontend-webapp`

**Step 2: Run Assessment**
Carlos creates a script:
```bash
#!/bin/bash
# assess-all.sh

for project in user-api order-service payment-gateway notification-service inventory-api shipping-service analytics-platform frontend-webapp; do
  echo "=== $project ==="
  cd ~/projects/$project
  a2 check --output=json --format=json > ../results/$project.json
  a2 check --output=toon --verbose > ../results/$project.txt
done
```

**Step 3: Collect Results**
```bash
./assess-all.sh
```

**Step 4: Analyze Scores**
Carlos creates summary:
```bash
for file in results/*.json; do
  name=$(basename $file .json)
  score=$(jq '.maturity.score' $file)
  echo "$name: $score%"
done
```

**Output**:
```
user-api: 92%
order-service: 78%
payment-gateway: 95%
notification-service: 65%
inventory-api: 88%
shipping-service: 72%
analytics-platform: 58%
frontend-webapp: 82%
```

**Step 5: Identify Trends**
Carlos visualizes data:
- **Top performers**: payment-gateway (95%), user-api (92%)
- **Need attention**: analytics-platform (58%), notification-service (65%)
- **Team average**: 79%

**Step 6: Deep Dive on Low Performers**
```bash
cat results/analytics-platform.json | jq '.results[] | select(.status == "Fail")'
```

**Output**:
```json
{
  "id": "go:tests",
  "name": "Go Tests",
  "status": "Fail",
  "message": "12 tests failing"
}
{
  "id": "go:coverage",
  "name": "Go Coverage",
  "status": "Fail",
  "message": "Coverage: 35% (threshold: 80%)"
}
```

**Step 7: Action Plan**
Carlos creates improvement plan:

**Project**: analytics-platform
**Current Score**: 58%
**Target Score**: 80%
**Timeline**: 6 weeks

**Week 1-2**: Fix failing tests
- Assign senior developer
- Priority: High

**Week 3-4**: Improve coverage
- Add unit tests for core logic
- Target: 60% coverage

**Week 5-6**: Reach 80% coverage
- Comprehensive test suite
- Documentation

**Step 8: Track Progress**
Carlos sets up weekly automated checks:
```yaml
# .github/workflows/weekly-quality.yml
on:
  schedule:
    - cron: '0 9 * * 1'  # Every Monday 9am

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run A2
        run: |
          a2 check --output=json > results.json
          echo "SCORE=$(jq '.maturity.score' results.json)" >> $GITHUB_OUTPUT
        id: a2

      - name: Create Issue if Low Score
        if: steps.a2.outputs.SCORE < 80
        uses: actions/github-script@v6
        with:
          script: |
            github.rest.issues.create({
              title: 'Quality Score Below Target: ${{ steps.a2.outputs.SCORE }}%',
              body: 'Automated quality check detected score below 80% target.',
              labels: ['quality', 'automated']
            });
```

#### Scenario 2: Due Diligence for Acquired Codebase

**Step 1: Access New Codebase**
Company acquires startup, Carlos needs to assess their code.

**Step 2: Quick Assessment**
```bash
cd ~/acquired-product
a2 check
```

**Output**:
```
Detected languages: TypeScript, Go, Python

Running checks...
✓ TypeScript: Build successful
✗ TypeScript: Tests failed - 45 tests failing
⚠ TypeScript: Coverage at 30%
⚠ TypeScript: Many linting issues
⚠ Go: Build successful
⚠ Go: No tests found
⚠ Python: Coverage at 45%

Overall Maturity: Development (45% score, 12 failures)
```

**Step 3: Detailed Report**
```bash
a2 check --output=json > due-diligence.json
```

**Step 4: Presentation**
Carlos creates slide deck:
- **Current State**: 45% maturity, Development level
- **Critical Issues**: 12 failures
- **Investment Needed**: 3-6 months to reach production-ready
- **Risk Assessment**: Medium-high technical debt
- **Recommendation**: Plan 2-month quality sprint before integration

**Step 5: Executive Summary**
```
Acquired Product Quality Assessment
====================================

Overall Score: 45/100 (Development Stage)

Strengths:
- Code compiles successfully
- Modern tech stack
- Some documentation

Concerns:
- 45 failing tests
- Low test coverage (30-45%)
- No security scanning
- Missing operational checks

Recommendation:
Invest 2 developer-months in quality improvements before integration.
Estimated ROI: 3x (reduced production bugs, faster onboarding)

Budget: $40,000 (2 devs × 2 months)
Timeline: 8 weeks
```

**Success Criteria**:
- Clear visibility into codebase quality
- Data-driven decisions
- Actionable improvement plans
- Trackable progress over time

**Key Metrics**:
- Maturity score per project
- Number of critical issues
- Test coverage trends
- Security vulnerabilities

---

## Server Mode Journey

### Journey: Web-Based Quality Analysis

**Persona**: Multiple stakeholders
**Goal**: Analyze repositories without local installation
**Use Case**: Quick quality checks, team dashboards

#### Scenario 1: Quick Repository Analysis

**Step 1: Access Web UI**
User navigates to: `https://a2.example.com`

**Step 2: Submit Repository**
```
GitHub URL: https://github.com/org/repo
Profile: API
Target: Production
```

**Step 3: Monitor Progress**
Real-time updates:
```
Status: Running...
Progress: ████████░░ 80%

Completed:
✓ Language detection
✓ Go build check
✓ Go test suite
✓ Go coverage analysis
⚳ Go dependency scan...
```

**Step 4: View Results**
Dashboard shows:
```
┌─────────────────────────────────────┐
│ Repository: org/repo                │
│ Maturity Score: 87%                 │
│ Level: Mature                       │
│ Duration: 2m 34s                    │
└─────────────────────────────────────┘

Summary:
✅ Passed: 42 checks
⚠️  Warnings: 5 checks
❌ Failed: 2 checks
ℹ️  Skipped: 3 checks

Critical Failures:
1. Go: Tests - 2 tests failing
2. Go: Coverage - 72% (target: 80%)

Warnings:
1. Go: Logging - Unstructured logging
2. Common: Secrets - No .gitleaksignore
3. Common: API Docs - Incomplete OpenAPI spec
4. Common: Metrics - No metrics instrumentation
5. Common: Tracing - No distributed tracing
```

**Step 5: Export Results**
```json
{
  "repository": "https://github.com/org/repo",
  "score": 87,
  "level": "Mature",
  "results": [...]
}
```

#### Scenario 2: Team Quality Dashboard

**Step 1: Admin Sets Up Monitoring**
DevOps team configures repositories:
```yaml
monitored_repositories:
  - url: https://github.com/org/user-api
    schedule: daily
  - url: https://github.com/org/order-service
    schedule: daily
  - url: https://github.com/org/payment-gateway
    schedule: hourly  # Critical service
```

**Step 2: Dashboard View**
Team sees aggregate view:
```
Team Quality Dashboard
======================

Overall Average: 82%

┌─────────────────┬───────┬──────────┬────────────┐
│ Repository      │ Score │ Trend    │ Last Check │
├─────────────────┼───────┼──────────┼────────────┤
│ user-api        │  92%  │ ↑ +3%    │ 2h ago     │
│ order-service   │  78%  │ ↓ -2%    │ 2h ago     │
│ payment-gateway │  95%  │ →        │ 1h ago     │
└─────────────────┴───────┴──────────┴────────────┘

Alerts:
⚠️ order-service: Score dropped below 80%
⚠️ order-service: New test failures detected
```

**Step 3: Drill Down**
Click on `order-service`:
```
order-service Quality Details
==============================

Current Score: 78% (↓ 2% from last check)

New Failures (since yesterday):
- Go: Tests: 2 new test failures
  - TestOrderCancellation
  - TestRefundProcessing

Regretions:
- Go: Coverage: 78% → 72% (-6%)

Recommendations:
1. Fix failing tests
2. Add tests for cancelled orders
3. Investigate coverage drop
```

**Step 4: Assign Work**
Team lead creates GitHub issue directly from dashboard:
```markdown
Alert: Quality Score Regression in order-service

**Current Score**: 78% (down from 80%)

**Issues**:
- 2 new test failures
- Coverage dropped 6%

**Action Items**:
- [ ] Fix TestOrderCancellation
- [ ] Fix TestRefundProcessing
- [ ] Restore test coverage

**Priority**: High
**Due**: This sprint
```

#### Scenario 3: Historical Trend Analysis

**Step 1: View Trends**
User selects date range: Last 90 days

**Step 2: Visualization**
```
Quality Trend - user-api
========================

100% ┤                    ╭─────
 90% ┤               ╭────╯
 80% ┤          ╭────╯
 70% ┤     ╭────╯
 60% ┤ ╭───╯
     ┼──────────────────────
     Jan  Feb  Mar  Apr  May

Events:
Feb 15: Added test suite (60% → 75%)
Mar 01: Reached 80% coverage target
Apr 10: Security scan integration
May 01: Passed all checks! 🎉
```

**Step 3: Compare Repositories**
```
Team Comparison (Last 30 Days)
===============================

                │ Start │ End   │ Change │
────────────────┼───────┼───────┼────────┼
user-api        │  89%  │  92%  │  +3%   │ ⬆️
order-service   │  76%  │  78%  │  +2%   │ ⬆️
payment-gateway │  95%  │  95%  │   0%   │ ➡️
notification    │  62%  │  65%  │  +3%   │ ⬆️

All teams improving! 🎉
```

**Success Criteria**:
- No installation required
- Fast analysis (< 5 minutes)
- Clear, actionable results
- Historical data available
- Team comparison features

**Use Cases**:
- Quick health checks
- Executive dashboards
- Team comparison
- Trend analysis
- Without local toolchain

---

## Customization Journey

### Journey: Creating Custom Quality Standards

**Persona**: Advanced users, teams with specific requirements
**Goal**: Tailor A2 to team's specific needs

#### Scenario 1: Custom Profile for Legacy Project

**Step 1: Identify Needs**
Team has legacy monolith with specific constraints:
- Can't easily add tests
- Old Go version (1.18)
- No containerization
- On-prem deployment

**Step 2: Create Custom Profile**
`.a2.yaml`:
```yaml
# Legacy monolith profile
profile: legacy-monolith

checks:
  disabled:
    # Skip container checks
    - "common:dockerfile"
    - "common:k8s"

    # Skip cloud-native checks
    - "common:health"
    - "common:metrics"
    - "common:tracing"

    # Relaxed testing requirements
    - "go:coverage"
    - "go:cyclomatic"

    # Skip modern practices
    - "common:api_docs"

language:
  go:
    # Lower expectations
    coverage_threshold: 0
    cyclomatic_threshold: 25  # Allow more complexity

target: development  # Don't expect production-ready
```

**Step 3: Gradual Improvement Plan**
Team creates roadmap:
```yaml
# Phase 1 (current): Basic checks
# - Build passes
# - Critical tests pass
# - Security scan

# Phase 2 (3 months): Add tests
# Enable coverage threshold: 40%

# Phase 3 (6 months): Improve code quality
# Enable cyclomatic check
# Coverage threshold: 60%

# Phase 4 (12 months): Modern practices
# Coverage threshold: 80%
# Add metrics
# Add health checks
```

#### Scenario 2: Industry-Specific Checks

**Step 1: Healthcare Project Requirements**
Team building healthcare app needs HIPAA compliance checks.

**Step 2: Define External Checks**
`.a2.yaml`:
```yaml
external:
  # HIPAA compliance check
  - id: hipaa-logging
    name: HIPAA Audit Logging
    command: ./scripts/check-hipaa-logging.sh
    severity: fail
    description: Verify all sensitive operations are logged

  # Data encryption check
  - id: encryption-at-rest
    name: Encryption at Rest
    command: ./scripts/check-encryption.sh
    severity: fail
    description: Verify database encryption enabled

  # PHI detection
  - id: phi-scanning
    name: PHI in Code
    command: gitleaks --config=.gitleaks-phi.config
    severity: fail
    description: Detect potential PHI in code

  # Access control
  - id: rbac-check
    name: RBAC Implementation
    command: ./scripts/check-rbac.sh
    severity: fail
    description: Verify role-based access control
```

**Step 3: Create Custom Check Scripts**
`scripts/check-hipaa-logging.sh`:
```bash
#!/bin/bash
# Check for audit logging in sensitive functions

if grep -r "audit.Log" ./internal/auth ./internal/patient; then
  exit 0  # Pass
else
  echo "Missing audit logging in sensitive operations"
  exit 2  # Fail
fi
```

**Step 4: Team-Wide Adoption**
Commit to team repo:
```bash
git add .a2.yaml scripts/
git commit -m "Add HIPAA compliance checks"
```

All team members now enforce same standards.

#### Scenario 3: Custom Severity Levels

**Step 1: Define Team Priorities**
Team decides:
- Security = Critical (fail)
- Tests = Important (warn in dev, fail in main)
- Coverage = Nice to have (warn)
- Formatting = Nice to have (warn)

**Step 2: Configure by Environment**
`.a2.yaml` (development):
```yaml
execution:
  severity_mode: relaxed  # Warnings don't fail

checks:
  # Security always fails
  critical:
    - "common:secrets"
    - "go:race"

  # Tests warn in dev
  important:
    - "go:tests"
    - "python:tests"
```

`.a2.production.yaml` (production branch):
```yaml
execution:
  severity_mode: strict  # Warnings fail

checks:
  critical:
    - "*:tests"  # All tests must pass
    - "common:secrets"
    - "go:race"
    - "*:coverage"
```

**Step 3: CI Configuration**
```yaml
# .github/workflows/a2.yml
- name: A2 Check (Development)
  if: github.ref != 'refs/heads/main'
  run: a2 check --config=.a2.yaml

- name: A2 Check (Production)
  if: github.ref == 'refs/heads/main'
  run: a2 check --config=.a2.production.yaml
```

#### Scenario 4: Custom Output Formats

**Step 1: Team Slack Integration**
Create Slack notifier:
```bash
#!/bin/bash
# a2-to-slack.sh

a2 check --output=json | jq '
  {
    repository: env.REPO_NAME,
    score: .maturity.score,
    level: .maturity.level,
    failures: [.results[] | select(.status == "Fail") | .name],
    warnings: [.results[] | select(.status == "Warn") | .name]
  }
' | curl -X POST \
  -H "Content-Type: application/json" \
  -d @- \
  "$SLACK_WEBHOOK_URL"
```

**Slack Message**:
```
🔍 A2 Quality Report

Repository: user-api
Score: 92% 🟢
Level: Production-Ready

❌ Failures: None
⚠️ Warnings: Coverage (85%)

Great job team! 🎉
```

**Step 2: HTML Report Generator**
```bash
#!/bin/bash
# a2-report.sh

a2 check --output=json | \
  jq '.' | \
  ./scripts/generate-html-report.py \
  > quality-report.html

open quality-report.html
```

Generates beautiful HTML report with charts and tables.

**Success Criteria**:
- A2 adapts to team needs
- Custom standards enforced consistently
- Industry compliance maintained
- Team-specific workflows supported

---

## Journey Summary

### Key User Journeys

| Journey | Primary Persona | Frequency | Duration | Success Metric |
|---------|----------------|-----------|----------|----------------|
| Quick Start | Any user | One-time | 5 min | First successful check |
| First-Time Setup | Senior Engineer | Project start | 15 min | Team adopts A2 |
| Daily Development | AI Developer | Daily | 2 min | Quality issues caught |
| CI/CD Integration | DevOps Engineer | Per project | 1 hour | Quality gate in place |
| AI Validation | AI Developer | 10x/day | 1 min | Less bugs in AI code |
| Open Source Review | Maintainer | Per PR | Automated | 60% review time saved |
| Project Assessment | Engineering Manager | Quarterly | 1 day | Data-driven decisions |
| Server Analysis | Any user | As needed | 5 min | Quick insights |
| Customization | Advanced Teams | Per project | 2 hours | Tailored standards |

### Common Patterns Across Journeys

1. **Installation**: Quick, single command
2. **Configuration**: Interactive with sensible defaults
3. **Execution**: Fast, parallel by default
4. **Output**: Clear, actionable, multiple formats
5. **Integration**: CI/CD, pre-commit, IDE
6. **Iteration**: Run, fix, re-run workflow
7. **Collaboration**: Team-wide standards
8. **Improvement**: Trackable progress over time

### Pain Points Solved

- ✅ Inconsistent code quality
- ✅ AI-generated code issues
- ✅ Time-consuming reviews
- ✅ Fragmented tooling
- ✅ Lack of measurable standards
- ✅ CI/CD integration complexity

### User Delight Moments

- 🎉 First check runs successfully
- ⚡ Fast feedback (seconds, not minutes)
- 💡 Clear, actionable suggestions
- 📊 Visible progress (maturity score)
- 🚀 CI catches issue before merge
- 🤖 AI code validated automatically
- 📈 Team quality improves over time

---

## Future Journey Enhancements

### Planned Improvements

**Phase 2 (Enhanced Features)**:
- Historical trend visualization
- Team comparison dashboards
- Custom check marketplace
- Automated fix suggestions

**Phase 3 (Ecosystem)**:
- IDE extensions (VS Code, JetBrains)
- Git provider integrations (GitLab, Bitbucket)
- Self-hosted enterprise version
- Custom policy engine

**Phase 4 (AI Integration)**:
- Automated PR reviews
- Learning from team patterns
- Predictive quality metrics
- Fix recommendations with AI

---

## Feedback & Iteration

### Collecting User Feedback

**Surveys**:
- Quarterly user satisfaction survey
- Feature request voting
- Persona-specific interviews

**Analytics**:
- Usage patterns (most used checks)
- Error rates and friction points
- Journey completion rates
- Time-to-success metrics

**Continuous Improvement**:
- A/B test new features
- Beta program for enhancements
- Community roadmap voting
- Responsive to user needs

### Success Indicators

**Product Metrics**:
- Adoption rate growth
- User retention
- Feature usage
- Community contributions

**Quality Metrics**:
- False positive rate < 5%
- Check execution time < 5 min
- User satisfaction > 4.5/5
- Support ticket volume

**Business Impact**:
- Production bugs reduced
- Review time saved
- Developer velocity maintained
- Team quality improved

---

## Appendix: Journey Maps

### Visual Journey Timeline

```
First Contact → Setup → Daily Use → CI/CD → Advanced
     │           │        │         │         │
     ▼           ▼        ▼         ▼         ▼
   Install    Config    Check     GitHub    Custom
   A2         .a2.yaml   & Fix    Action    Checks
     │           │        │         │         │
     └───────────┴────────┴─────────┴─────────┘
                    Continuous Quality
```

### User Lifecycle

```
New User → Active User → Power User → Champion
    │          │            │           │
 5 min     1 week        1 month     6 months
    │          │            │           │
First     Daily      Advanced    Team
Check     Use       Features     Advocate
```

### Quality Journey

```
PoC → Development → Mature → Production-Ready
 40%      60%          80%         100%
  │         │            │            │
Basic   Core      Quality     Excellence
Setup   Features   Focus      Focus
```

---

**End of User Journeys Document**
