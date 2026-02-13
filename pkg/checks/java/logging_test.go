package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type LoggingCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *LoggingCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "java-logging-test-*")
	s.Require().NoError(err)
}

func (s *LoggingCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *LoggingCheckTestSuite) TestIDAndName() {
	check := &LoggingCheck{}
	s.Equal("java:logging", check.ID())
	s.Equal("Java Logging", check.Name())
}

func (s *LoggingCheckTestSuite) TestRun_NoLogging() {
	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No structured logging detected")
}

func (s *LoggingCheckTestSuite) TestRun_ResultLanguage() {
	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *LoggingCheckTestSuite) TestRun_SLF4J_Maven() {
	content := `<project>
  <dependencies>
    <dependency>
      <groupId>org.slf4j</groupId>
      <artifactId>slf4j-api</artifactId>
      <version>2.0.9</version>
    </dependency>
  </dependencies>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "SLF4J")
}

func (s *LoggingCheckTestSuite) TestRun_SLF4J_Gradle() {
	content := `plugins {
    id 'java'
}

dependencies {
    implementation 'org.slf4j:slf4j-api:2.0.9'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SLF4J")
}

func (s *LoggingCheckTestSuite) TestRun_SLF4J_GradleKts() {
	content := `plugins {
    java
}

dependencies {
    implementation("org.slf4j:slf4j-api:2.0.9")
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle.kts"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SLF4J")
}

func (s *LoggingCheckTestSuite) TestRun_Logback_ConfigFile() {
	// Create src/main/resources directory
	resourcesDir := filepath.Join(s.tempDir, "src", "main", "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	s.Require().NoError(err)

	content := `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
    <encoder>
      <pattern>%d{HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
    </encoder>
  </appender>
  <root level="INFO">
    <appender-ref ref="STDOUT" />
  </root>
</configuration>`
	err = os.WriteFile(filepath.Join(resourcesDir, "logback.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Logback")
}

func (s *LoggingCheckTestSuite) TestRun_Logback_SpringConfig() {
	resourcesDir := filepath.Join(s.tempDir, "src", "main", "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(resourcesDir, "logback-spring.xml"), []byte("<configuration/>"), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Logback")
}

func (s *LoggingCheckTestSuite) TestRun_Logback_RootConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "logback.xml"), []byte("<configuration/>"), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Logback")
}

func (s *LoggingCheckTestSuite) TestRun_Logback_Maven() {
	content := `<project>
  <dependencies>
    <dependency>
      <groupId>ch.qos.logback</groupId>
      <artifactId>logback-classic</artifactId>
      <version>1.4.11</version>
    </dependency>
  </dependencies>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Logback")
}

func (s *LoggingCheckTestSuite) TestRun_Log4j2_ConfigFile() {
	resourcesDir := filepath.Join(s.tempDir, "src", "main", "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	s.Require().NoError(err)

	content := `<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="WARN">
  <Appenders>
    <Console name="Console" target="SYSTEM_OUT">
      <PatternLayout pattern="%d{HH:mm:ss.SSS} [%t] %-5level %logger{36} - %msg%n"/>
    </Console>
  </Appenders>
  <Loggers>
    <Root level="info">
      <AppenderRef ref="Console"/>
    </Root>
  </Loggers>
</Configuration>`
	err = os.WriteFile(filepath.Join(resourcesDir, "log4j2.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Log4j2")
}

func (s *LoggingCheckTestSuite) TestRun_Log4j2_YamlConfig() {
	resourcesDir := filepath.Join(s.tempDir, "src", "main", "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(resourcesDir, "log4j2.yaml"), []byte("Configuration:"), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Log4j2")
}

func (s *LoggingCheckTestSuite) TestRun_Log4j2_PropertiesConfig() {
	resourcesDir := filepath.Join(s.tempDir, "src", "main", "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(resourcesDir, "log4j2.properties"), []byte("rootLogger.level=INFO"), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Log4j2")
}

func (s *LoggingCheckTestSuite) TestRun_Log4j2_RootConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "log4j2.xml"), []byte("<Configuration/>"), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Log4j2")
}

func (s *LoggingCheckTestSuite) TestRun_Log4j2_Maven() {
	content := `<project>
  <dependencies>
    <dependency>
      <groupId>org.apache.logging.log4j</groupId>
      <artifactId>log4j-core</artifactId>
      <version>2.22.0</version>
    </dependency>
  </dependencies>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Log4j2")
}

func (s *LoggingCheckTestSuite) TestRun_SystemOutPrintln_Warning() {
	// Create source directory structure
	srcDir := filepath.Join(s.tempDir, "src", "main", "java", "com", "example")
	err := os.MkdirAll(srcDir, 0755)
	s.Require().NoError(err)

	// Java file with System.out.println
	javaContent := `package com.example;

public class App {
    public static void main(String[] args) {
        System.out.println("Hello World");
        System.err.println("Error message");
    }
}`
	err = os.WriteFile(filepath.Join(srcDir, "App.java"), []byte(javaContent), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "System.out.println")
}

func (s *LoggingCheckTestSuite) TestRun_StructuredLogging_WithPrintln() {
	// SLF4J in pom.xml
	pomContent := `<project>
  <dependencies>
    <dependency>
      <groupId>org.slf4j</groupId>
      <artifactId>slf4j-api</artifactId>
    </dependency>
  </dependencies>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomContent), 0644)
	s.Require().NoError(err)

	// Create source with System.out.println
	srcDir := filepath.Join(s.tempDir, "src", "main", "java", "com", "example")
	err = os.MkdirAll(srcDir, 0755)
	s.Require().NoError(err)

	javaContent := `package com.example;

public class App {
    public static void main(String[] args) {
        System.out.println("Debug message");
    }
}`
	err = os.WriteFile(filepath.Join(srcDir, "App.java"), []byte(javaContent), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "SLF4J")
	s.Contains(result.Reason, "System.out.println")
}

func (s *LoggingCheckTestSuite) TestRun_MultipleLoggingFrameworks() {
	content := `<project>
  <dependencies>
    <dependency>
      <groupId>org.slf4j</groupId>
      <artifactId>slf4j-api</artifactId>
    </dependency>
    <dependency>
      <groupId>ch.qos.logback</groupId>
      <artifactId>logback-classic</artifactId>
    </dependency>
  </dependencies>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SLF4J")
	s.Contains(result.Reason, "Logback")
}

func (s *LoggingCheckTestSuite) TestRun_CommentedPrintln_NotCounted() {
	// Create source directory
	srcDir := filepath.Join(s.tempDir, "src", "main", "java", "com", "example")
	err := os.MkdirAll(srcDir, 0755)
	s.Require().NoError(err)

	// Java file with commented System.out.println
	javaContent := `package com.example;

public class App {
    public static void main(String[] args) {
        // System.out.println("This is commented");
        * System.out.println("Also in block comment");
    }
}`
	err = os.WriteFile(filepath.Join(srcDir, "App.java"), []byte(javaContent), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	// Should not find println issues because they're in comments
	s.NotContains(result.Reason, "System.out.println")
}

func TestLoggingCheckTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingCheckTestSuite))
}
