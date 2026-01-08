# Java Checks

This document describes all Java-specific checks available in A2.

## Overview

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `java:project` | Java Project | Yes | 100 | Detects Maven (pom.xml) or Gradle (build.gradle) projects |
| `java:build` | Java Build | Yes | 110 | Compiles with Maven or Gradle (auto-detects wrapper scripts) |
| `java:tests` | Java Tests | Yes | 120 | Runs JUnit/TestNG tests |
| `java:format` | Java Format | No | 200 | Detects Spotless, google-java-format, EditorConfig, IDE formatters |
| `java:lint` | Java Lint | No | 210 | Detects Checkstyle, SpotBugs, PMD, Error Prone, SonarQube |
| `java:coverage` | Java Coverage | No | 220 | Checks JaCoCo coverage reports against threshold |
| `java:deps` | Java Dependencies | No | 230 | Detects OWASP Dependency-Check, Snyk, Dependabot, Renovate |
| `java:logging` | Java Logging | No | 250 | Detects SLF4J, Logback, Log4j2; warns on System.out.println |

---

## java:project

Detects Maven or Gradle project configuration.

**Build tool detection priority:**
1. `build.gradle` or `build.gradle.kts` → Gradle (Groovy or Kotlin DSL)
2. `pom.xml` → Maven
3. `gradlew` → Gradle wrapper
4. `mvnw` → Maven wrapper

When both Maven and Gradle are present, Gradle is preferred as the more modern build system.

**Wrapper detection:**
- `mvnw` (Maven Wrapper)
- `gradlew` (Gradle Wrapper)

**Status:**
- **Pass**: pom.xml or build.gradle found
- **Fail**: No Java project configuration found

---

## java:build

Compiles the Java project using the detected build tool.

**Configuration:**
```yaml
language:
  java:
    build_tool: auto  # auto, maven, gradle
```

**Build commands:**
- Maven: `./mvnw compile -q -DskipTests` (or `mvn` if no wrapper)
- Gradle: `./gradlew compileJava -q --no-daemon` (or `gradle` if no wrapper)

Wrapper scripts are preferred when available for reproducible builds.

**Status:**
- **Pass**: Compilation succeeds
- **Fail**: Compilation fails or no build tool found

---

## java:tests

Runs tests using the detected build tool and test runner.

**Configuration:**
```yaml
language:
  java:
    test_runner: auto  # auto, junit, testng
```

**Test commands:**
- Maven: `./mvnw test -q` (or `mvn test -q`)
- Gradle: `./gradlew test --no-daemon` (or `gradle test`)

**Status:**
- **Pass**: All tests pass
- **Warn**: No tests found
- **Fail**: Tests fail

---

## java:format

Detects code formatting configuration for Java projects.

**Formatters detected:**
- google-java-format: `google-java-format` in pom.xml/build.gradle, or `google-java-format.jar`
- Spotless: `spotless` plugin in pom.xml or build.gradle
- EditorConfig: `.editorconfig` with Java settings (`[*.java]`)
- IntelliJ IDEA: `.idea/codeStyles/` directory
- Eclipse: `.settings/org.eclipse.jdt.core.prefs`

**Status:**
- **Pass**: Formatter configuration found
- **Warn**: No formatter configuration found

**Recommendation:** Configure Spotless or google-java-format for consistent code formatting.

---

## java:lint

Detects static analysis tools configured for the project.

**Tools detected:**
- Checkstyle: `checkstyle.xml`, `checkstyle` plugin in pom.xml/build.gradle
- SpotBugs: `spotbugs.xml`, `spotbugsInclude.xml`, `spotbugs` plugin
- PMD: `pmd.xml`, `ruleset.xml`, `pmd` plugin
- Error Prone: `errorprone` in pom.xml/build.gradle
- SonarQube: `sonar-project.properties`, `sonarqube` plugin

**Status:**
- **Pass**: One or more linting tools configured
- **Warn**: No linting tools found

**Recommendation:** Configure Checkstyle, SpotBugs, or PMD for static analysis.

---

## java:coverage

Checks test coverage using JaCoCo reports.

**Configuration:**
```yaml
language:
  java:
    coverage_threshold: 80  # Default: 80%
```

**Report locations:**
- Maven: `target/site/jacoco/jacoco.xml`
- Gradle: `build/reports/jacoco/test/jacocoTestReport.xml`

Parses LINE or INSTRUCTION counter types from the XML report.

**Status:**
- **Pass**: Coverage meets or exceeds threshold
- **Warn**: Coverage below threshold or JaCoCo not configured/no report found

**Recommendation:** Add JaCoCo plugin to pom.xml or build.gradle.

---

## java:deps

Detects dependency vulnerability scanning configuration.

**Tools detected:**
- OWASP Dependency-Check: `dependency-check` plugin in pom.xml/build.gradle
- Snyk: `.snyk` file, `snyk` in CI config
- Dependabot: `.github/dependabot.yml` with maven/gradle ecosystem
- Renovate: `renovate.json`, `renovate.json5`, `.renovaterc`

**Maven plugins checked:**
- `org.owasp:dependency-check-maven`
- `snyk-maven-plugin`

**Gradle plugins checked:**
- `org.owasp.dependencycheck`
- `io.snyk.gradle.plugin`

**Status:**
- **Pass**: Dependency scanning tool configured
- **Warn**: No dependency scanning found

**Recommendation:** Configure OWASP Dependency-Check or Snyk for vulnerability scanning.

---

## java:logging

Checks for structured logging practices instead of System.out.println.

**Logging frameworks detected:**
- SLF4J: `slf4j` dependency in pom.xml/build.gradle
- Logback: `logback.xml`, `logback-spring.xml`, or `logback` dependency
- Log4j2: `log4j2.xml`, `log4j2.yaml`, `log4j2.properties`, or `log4j` dependency

**Anti-patterns detected (in src/main/java/):**
- `System.out.print` statements
- `System.err.print` statements

Comments are excluded from detection.

**Status:**
- **Pass**: Structured logging library found, no System.out.println
- **Warn**: System.out.println found, or no structured logging detected

**Recommendation:** Use SLF4J with Logback or Log4j2 for structured logging.

---

## Configuration Example

```yaml
language:
  java:
    build_tool: auto         # auto, maven, gradle
    test_runner: auto        # auto, junit, testng
    coverage_threshold: 80
```
