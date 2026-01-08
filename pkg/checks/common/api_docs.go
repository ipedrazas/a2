package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// APIDocsCheck verifies that API documentation exists.
type APIDocsCheck struct{}

func (c *APIDocsCheck) ID() string   { return "common:api_docs" }
func (c *APIDocsCheck) Name() string { return "API Documentation" }

// Run checks for API documentation files and generators.
func (c *APIDocsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var found []string

	// Check for OpenAPI/Swagger spec files
	openAPIFiles := []string{
		"openapi.yaml", "openapi.yml", "openapi.json",
		"swagger.yaml", "swagger.yml", "swagger.json",
		"api.yaml", "api.yml", "api.json",
		"docs/openapi.yaml", "docs/openapi.yml",
		"docs/swagger.yaml", "docs/swagger.yml",
		"api/openapi.yaml", "api/swagger.yaml",
	}

	for _, file := range openAPIFiles {
		if safepath.Exists(path, file) {
			found = append(found, file)
		}
	}

	// Check for API documentation directories
	docDirs := []string{
		"docs/api", "api-docs", "api/docs",
		"swagger", "swagger-ui",
	}

	for _, dir := range docDirs {
		if safepath.IsDir(path, dir) {
			found = append(found, dir+"/")
		}
	}

	// Check for documentation generators in Go
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goDocGenerators := []struct {
				pattern string
				name    string
			}{
				{"github.com/swaggo/swag", "swaggo"},
				{"github.com/go-swagger/go-swagger", "go-swagger"},
				{"github.com/grpc-ecosystem/grpc-gateway", "gRPC-Gateway"},
			}
			for _, gen := range goDocGenerators {
				if strings.Contains(string(content), gen.pattern) {
					found = append(found, gen.name)
				}
			}
		}
	}

	// Check for documentation generators in Python
	pythonConfigs := []string{"pyproject.toml", "setup.py", "requirements.txt"}
	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				pyDocGenerators := []struct {
					pattern string
					name    string
				}{
					{"drf-spectacular", "drf-spectacular"},
					{"drf-yasg", "drf-yasg"},
					{"flasgger", "flasgger"},
					{"flask-restx", "flask-restx"},
					{"fastapi", "FastAPI (auto-docs)"},
				}
				for _, gen := range pyDocGenerators {
					if strings.Contains(strings.ToLower(string(content)), gen.pattern) {
						found = append(found, gen.name)
					}
				}
			}
			break
		}
	}

	// Check for documentation generators in Node.js
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			nodeDocGenerators := []struct {
				pattern string
				name    string
			}{
				{"swagger-jsdoc", "swagger-jsdoc"},
				{"swagger-ui-express", "swagger-ui-express"},
				{"@nestjs/swagger", "@nestjs/swagger"},
				{"tsoa", "tsoa"},
				{"express-openapi", "express-openapi"},
			}
			for _, gen := range nodeDocGenerators {
				if strings.Contains(string(content), gen.pattern) {
					found = append(found, gen.name)
				}
			}
		}
	}

	// Check for GraphQL schemas (also counts as API docs)
	graphqlFiles := []string{
		"schema.graphql", "schema.gql",
		"graphql/schema.graphql", "api/schema.graphql",
	}
	for _, file := range graphqlFiles {
		if safepath.Exists(path, file) {
			found = append(found, "GraphQL schema")
			break
		}
	}

	// Check for protobuf definitions (gRPC API docs)
	if hasProtoFiles(path) {
		found = append(found, "Protocol Buffers")
	}

	// Build result
	if len(found) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "API documentation found: " + strings.Join(unique(found), ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No API documentation found (consider adding OpenAPI/Swagger specs)"
	}

	return result, nil
}

func hasProtoFiles(path string) bool {
	protoDirs := []string{"", "proto", "api/proto", "protos"}
	for _, dir := range protoDirs {
		checkPath := path
		if dir != "" {
			checkPath = path + "/" + dir
		}
		if files, err := safepath.Glob(checkPath, "*.proto"); err == nil && len(files) > 0 {
			return true
		}
	}
	return false
}

func unique(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
