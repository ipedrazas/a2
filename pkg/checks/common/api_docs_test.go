package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type APIDocsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *APIDocsCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "api-docs-test-*")
	s.Require().NoError(err)
}

func (s *APIDocsCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *APIDocsCheckTestSuite) TestIDAndName() {
	check := &APIDocsCheck{}
	s.Equal("common:api_docs", check.ID())
	s.Equal("API Documentation", check.Name())
}

func (s *APIDocsCheckTestSuite) TestOpenAPIYamlExists() {
	content := `openapi: "3.0.0"
info:
  title: My API
  version: "1.0.0"
paths:
  /users:
    get:
      summary: Get users
`
	err := os.WriteFile(filepath.Join(s.tempDir, "openapi.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "openapi.yaml")
}

func (s *APIDocsCheckTestSuite) TestSwaggerJsonExists() {
	content := `{
  "swagger": "2.0",
  "info": {
    "title": "My API",
    "version": "1.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "swagger.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "swagger.json")
}

func (s *APIDocsCheckTestSuite) TestOpenAPIInDocsDir() {
	docsDir := filepath.Join(s.tempDir, "docs")
	err := os.MkdirAll(docsDir, 0755)
	s.Require().NoError(err)

	content := `openapi: "3.0.0"
info:
  title: My API
`
	err = os.WriteFile(filepath.Join(docsDir, "openapi.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "docs/openapi.yaml")
}

func (s *APIDocsCheckTestSuite) TestAPIDocsDirectory() {
	apiDocsDir := filepath.Join(s.tempDir, "api-docs")
	err := os.MkdirAll(apiDocsDir, 0755)
	s.Require().NoError(err)

	// Create a file in the directory
	err = os.WriteFile(filepath.Join(apiDocsDir, "index.html"), []byte("<html></html>"), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "api-docs/")
}

func (s *APIDocsCheckTestSuite) TestGoSwaggo() {
	content := `module myapp

go 1.21

require github.com/swaggo/swag v1.16.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "swaggo")
}

func (s *APIDocsCheckTestSuite) TestGoSwagger() {
	content := `module myapp

go 1.21

require github.com/go-swagger/go-swagger v0.30.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "go-swagger")
}

func (s *APIDocsCheckTestSuite) TestPythonFastAPI() {
	content := `[project]
name = "myapp"
dependencies = [
    "fastapi>=0.100.0",
    "uvicorn>=0.23.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "FastAPI")
}

func (s *APIDocsCheckTestSuite) TestPythonDRFSpectacular() {
	content := `[project]
name = "myapp"
dependencies = [
    "django>=4.0",
    "djangorestframework>=3.14",
    "drf-spectacular>=0.26.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "drf-spectacular")
}

func (s *APIDocsCheckTestSuite) TestNodeSwaggerJsdoc() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "express": "^4.18.0",
    "swagger-jsdoc": "^6.2.0",
    "swagger-ui-express": "^5.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "swagger-jsdoc")
}

func (s *APIDocsCheckTestSuite) TestNodeNestJSSwagger() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "@nestjs/core": "^10.0.0",
    "@nestjs/swagger": "^7.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "@nestjs/swagger")
}

func (s *APIDocsCheckTestSuite) TestGraphQLSchema() {
	content := `type Query {
  users: [User!]!
  user(id: ID!): User
}

type User {
  id: ID!
  name: String!
  email: String!
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "schema.graphql"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "GraphQL schema")
}

func (s *APIDocsCheckTestSuite) TestProtocolBuffers() {
	content := `syntax = "proto3";

package myapp;

service UserService {
  rpc GetUser (GetUserRequest) returns (User) {}
}

message GetUserRequest {
  string id = 1;
}

message User {
  string id = 1;
  string name = 2;
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "user.proto"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Protocol Buffers")
}

func (s *APIDocsCheckTestSuite) TestProtoInSubdirectory() {
	protoDir := filepath.Join(s.tempDir, "proto")
	err := os.MkdirAll(protoDir, 0755)
	s.Require().NoError(err)

	content := `syntax = "proto3";
package myapp;
`
	err = os.WriteFile(filepath.Join(protoDir, "service.proto"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Protocol Buffers")
}

func (s *APIDocsCheckTestSuite) TestMultipleAPIDocs() {
	// OpenAPI spec
	content := `openapi: "3.0.0"
info:
  title: My API
`
	err := os.WriteFile(filepath.Join(s.tempDir, "openapi.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	// GraphQL schema
	gqlContent := `type Query { users: [User] }`
	err = os.WriteFile(filepath.Join(s.tempDir, "schema.graphql"), []byte(gqlContent), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "openapi.yaml")
	s.Contains(result.Reason, "GraphQL schema")
}

func (s *APIDocsCheckTestSuite) TestNoAPIDocs() {
	// Create some other file but no API docs
	err := os.WriteFile(filepath.Join(s.tempDir, "README.md"), []byte("# My Project"), 0644)
	s.Require().NoError(err)

	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No API documentation found")
}

func (s *APIDocsCheckTestSuite) TestEmptyDirectory() {
	check := &APIDocsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No API documentation found")
}

func TestAPIDocsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(APIDocsCheckTestSuite))
}
