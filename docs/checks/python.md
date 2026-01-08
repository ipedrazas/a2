# Python Checks

This document describes all Python-specific checks available in A2.

## Overview

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `python:project` | Python Project | Yes | 100 | Verifies project config exists |
| `python:build` | Python Build | Yes | 110 | Validates package manager setup |
| `python:tests` | Python Tests | Yes | 120 | Runs tests with pytest/unittest |
| `python:format` | Python Format | No | 200 | Checks formatting with ruff/black |
| `python:lint` | Python Lint | No | 210 | Lints code with ruff/flake8/pylint |
| `python:type` | Python Type Check | No | 215 | Type checks with mypy |
| `python:coverage` | Python Coverage | No | 220 | Measures test coverage |
| `python:deps` | Python Vulnerabilities | No | 230 | Scans for vulnerabilities |
| `python:complexity` | Python Complexity | No | 240 | Analyzes cyclomatic complexity with radon |
| `python:logging` | Python Logging | No | 250 | Detects structured logging vs print() |

---

## python:project

Verifies that a Python project configuration file exists.

**Detection priority:**
1. `pyproject.toml` (preferred)
2. `setup.py` (legacy, warns)
3. `requirements.txt` (minimal, warns)

**Status:**
- **Pass**: pyproject.toml found
- **Warn**: setup.py or requirements.txt found
- **Fail**: No project configuration found

---

## python:build

Validates dependencies using the detected package manager.

**Configuration:**
```yaml
language:
  python:
    package_manager: auto  # auto, pip, poetry, pipenv
```

**Auto-detection:**
- `poetry.lock` → poetry (`poetry check`)
- `Pipfile.lock` or `Pipfile` → pipenv (`pipenv check`)
- Otherwise → pip

**Status:**
- **Pass**: Validation succeeds or tool not installed
- **Fail**: Validation fails

---

## python:tests

Runs Python tests using the detected test runner.

**Configuration:**
```yaml
language:
  python:
    test_runner: auto  # auto, pytest, unittest
```

**Auto-detection:**
- `pytest.ini`, `conftest.py`, or `[tool.pytest]` in pyproject.toml → pytest
- Otherwise → pytest (default)

**Status:**
- **Pass**: All tests pass or no tests found
- **Fail**: Tests fail

---

## python:format

Checks Python code formatting.

**Configuration:**
```yaml
language:
  python:
    formatter: auto  # auto, black, ruff
```

**Auto-detection:**
- `ruff.toml` or `.ruff.toml` → ruff
- `[tool.black]` in pyproject.toml → black
- Tries ruff first, falls back to black

**Status:**
- **Pass**: All files properly formatted
- **Warn**: Files need formatting

**Fix:** Run `ruff format .` or `black .`

---

## python:lint

Runs Python linting to check for code quality issues.

**Configuration:**
```yaml
language:
  python:
    linter: auto  # auto, ruff, flake8, pylint
```

**Auto-detection:**
- `ruff.toml`, `.ruff.toml`, or `[tool.ruff]` → ruff
- `.flake8` or `[tool.flake8]` → flake8
- `.pylintrc` or `[tool.pylint]` → pylint
- Tries ruff first, falls back to flake8

**Status:**
- **Pass**: No issues found
- **Warn**: Linting issues detected

---

## python:type

Runs mypy for static type checking. Only activates for typed Python projects.

**Typed project detection:**
- `mypy.ini` or `.mypy.ini`
- `py.typed` marker (PEP 561)
- `[mypy]` section in setup.cfg
- `[tool.mypy]` in pyproject.toml

**Status:**
- **Pass**: No type errors, not a typed project, or mypy not installed
- **Warn**: Type errors found

**Fix:** Run `mypy .`

---

## python:coverage

Measures test coverage using pytest-cov.

**Configuration:**
```yaml
language:
  python:
    coverage_threshold: 80  # Default: 80%
```

**Status:**
- **Pass**: Coverage meets or exceeds threshold
- **Warn**: Coverage below threshold or pytest-cov not installed

---

## python:deps

Scans for security vulnerabilities in Python dependencies.

**Tools (tried in order):**
1. `pip-audit`
2. `safety`

**Status:**
- **Pass**: No vulnerabilities found or scanner not installed
- **Warn**: Vulnerabilities detected

---

## python:complexity

Analyzes cyclomatic complexity of Python functions using radon.

**Requirements:** Install with `pip install radon`

**Configuration:**
```yaml
language:
  python:
    cyclomatic_threshold: 15  # Default: 15
```

**Complexity Grades:**
- A (1-5): Low complexity
- B (6-10): Moderate complexity
- C (11-20): High complexity
- D/E/F (21+): Very high complexity

**Status:**
- **Pass**: No functions exceed threshold or radon not installed
- **Warn**: Functions exceed complexity threshold

**Fix:** Break complex functions into smaller, focused functions.

---

## python:logging

Checks for proper logging practices instead of print() statements.

**Logging modules detected:**
- `logging` (standard library)
- `structlog`
- `loguru`

**Anti-patterns detected (in non-test files):**
- `print()` statements

**Status:**
- **Pass**: Uses logging module, no print() statements
- **Warn**: Uses print() for logging or no logging module detected

**Fix:** Use the `logging` module or `structlog`/`loguru`.

---

## Configuration Example

```yaml
language:
  python:
    package_manager: auto    # auto, pip, poetry, pipenv
    test_runner: auto        # auto, pytest, unittest
    formatter: auto          # auto, black, ruff
    linter: auto             # auto, ruff, flake8, pylint
    coverage_threshold: 80
    cyclomatic_threshold: 15
```
