package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ConfigValidationCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *ConfigValidationCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "config-validation-test-*")
	s.Require().NoError(err)
}

func (s *ConfigValidationCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ConfigValidationCheckTestSuite) TestIDAndName() {
	check := &ConfigValidationCheck{}
	s.Equal("common:config_validation", check.ID())
	s.Equal("Config Validation", check.Name())
}

func (s *ConfigValidationCheckTestSuite) TestGoValidator() {
	content := `module myapp

go 1.21

require (
	github.com/go-playground/validator/v10 v10.16.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "validator")
}

func (s *ConfigValidationCheckTestSuite) TestGoViper() {
	content := `module myapp

go 1.21

require (
	github.com/spf13/viper v1.18.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Viper")
}

func (s *ConfigValidationCheckTestSuite) TestGoEnvconfig() {
	content := `module myapp

go 1.21

require (
	github.com/kelseyhightower/envconfig v1.4.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "envconfig")
}

func (s *ConfigValidationCheckTestSuite) TestGoKoanf() {
	content := `module myapp

go 1.21

require (
	github.com/knadh/koanf v1.5.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Koanf")
}

func (s *ConfigValidationCheckTestSuite) TestGoGodotenv() {
	content := `module myapp

go 1.21

require (
	github.com/joho/godotenv v1.5.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "godotenv")
}

func (s *ConfigValidationCheckTestSuite) TestNodeZod() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "zod": "^3.22.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Zod")
}

func (s *ConfigValidationCheckTestSuite) TestNodeJoi() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "joi": "^17.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Joi")
}

func (s *ConfigValidationCheckTestSuite) TestNodeYup() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "yup": "^1.3.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Yup")
}

func (s *ConfigValidationCheckTestSuite) TestNodeAjv() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "ajv": "^8.12.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "AJV")
}

func (s *ConfigValidationCheckTestSuite) TestNodeEnvalid() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "envalid": "^8.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Envalid")
}

func (s *ConfigValidationCheckTestSuite) TestNodeDotenv() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "dotenv": "^16.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "dotenv")
}

func (s *ConfigValidationCheckTestSuite) TestNestJSConfig() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "@nestjs/config": "^3.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "NestJS Config")
}

func (s *ConfigValidationCheckTestSuite) TestPythonPydantic() {
	content := `pydantic==2.5.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Pydantic")
}

func (s *ConfigValidationCheckTestSuite) TestPythonPydanticSettings() {
	content := `pydantic-settings==2.1.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Pydantic Settings")
}

func (s *ConfigValidationCheckTestSuite) TestPythonDynaconf() {
	content := `dynaconf==3.2.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Dynaconf")
}

func (s *ConfigValidationCheckTestSuite) TestPythonMarshmallow() {
	content := `marshmallow==3.20.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Marshmallow")
}

func (s *ConfigValidationCheckTestSuite) TestJavaHibernateValidator() {
	content := `<dependency>
    <groupId>org.hibernate.validator</groupId>
    <artifactId>hibernate-validator</artifactId>
</dependency>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Hibernate Validator")
}

func (s *ConfigValidationCheckTestSuite) TestJavaSpringBootConfig() {
	content := `implementation 'org.springframework.boot:spring-boot-configuration-processor'`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Spring Boot Config")
}

func (s *ConfigValidationCheckTestSuite) TestRustSerde() {
	content := `[package]
name = "myapp"

[dependencies]
serde = "1.0"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Serde")
}

func (s *ConfigValidationCheckTestSuite) TestRustFigment() {
	content := `[package]
name = "myapp"

[dependencies]
figment = "0.10"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Figment")
}

func (s *ConfigValidationCheckTestSuite) TestRustClap() {
	content := `[package]
name = "myapp"

[dependencies]
clap = "4.0"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Clap")
}

func (s *ConfigValidationCheckTestSuite) TestTypeScriptStrict() {
	content := `{
  "compilerOptions": {
    "strict": true,
    "target": "ES2022"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "tsconfig.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "TypeScript strict")
}

func (s *ConfigValidationCheckTestSuite) TestMultipleLibraries() {
	// Go mod with multiple config libraries
	content := `module myapp

go 1.21

require (
	github.com/spf13/viper v1.18.0
	github.com/go-playground/validator/v10 v10.16.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Viper")
	s.Contains(result.Reason, "validator")
}

func (s *ConfigValidationCheckTestSuite) TestNoConfigValidationFound() {
	content := `module myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No config validation found")
}

func (s *ConfigValidationCheckTestSuite) TestEmptyDirectory() {
	check := &ConfigValidationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No config validation found")
}

func TestConfigValidationCheckTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigValidationCheckTestSuite))
}
