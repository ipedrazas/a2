package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// MigrationsCheck verifies database migrations are managed.
type MigrationsCheck struct{}

func (c *MigrationsCheck) ID() string   { return "common:migrations" }
func (c *MigrationsCheck) Name() string { return "Database Migrations" }

// Run checks for database migration configuration.
func (c *MigrationsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var found []string

	// Check for migration directories
	migrationDirs := []string{
		"migrations",
		"db/migrations",
		"db/migrate",
		"database/migrations",
		"alembic",
		"alembic/versions",
		"prisma/migrations",
		"drizzle",
		"knex/migrations",
		"sequelize/migrations",
		"typeorm/migrations",
		"flyway/sql",
		"liquibase",
		"sql/migrations",
	}
	for _, dir := range migrationDirs {
		if safepath.Exists(path, dir) {
			found = append(found, dir)
			break
		}
	}

	// Check for migration tool config files
	migrationConfigs := map[string]string{
		"alembic.ini":          "Alembic",
		"atlas.hcl":            "Atlas",
		"dbmate.yml":           "dbmate",
		"flyway.conf":          "Flyway",
		"liquibase.properties": "Liquibase",
		"knexfile.js":          "Knex",
		"knexfile.ts":          "Knex",
		"ormconfig.json":       "TypeORM",
		"ormconfig.js":         "TypeORM",
		"prisma/schema.prisma": "Prisma",
		"drizzle.config.ts":    "Drizzle",
		"drizzle.config.js":    "Drizzle",
		".sequelizerc":         "Sequelize",
		"diesel.toml":          "Diesel",
		"refinery.toml":        "Refinery",
		"goose.yaml":           "Goose",
		"dbconfig.yml":         "sql-migrate",
	}
	for file, tool := range migrationConfigs {
		if safepath.Exists(path, file) {
			if !containsString(found, tool) {
				found = append(found, tool)
			}
		}
	}

	// Check Go dependencies for migration tools
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goMigrations := map[string]string{
				"github.com/golang-migrate/migrate": "golang-migrate",
				"github.com/pressly/goose":          "Goose",
				"ariga.io/atlas":                    "Atlas",
				"github.com/rubenv/sql-migrate":     "sql-migrate",
				"entgo.io/ent":                      "Ent",
				"gorm.io/gorm":                      "GORM",
			}
			for dep, name := range goMigrations {
				if strings.Contains(string(content), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check Python dependencies
	pythonFiles := []string{"pyproject.toml", "requirements.txt", "setup.py"}
	for _, file := range pythonFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			pythonMigrations := map[string]string{
				"alembic":       "Alembic",
				"django":        "Django migrations",
				"flask-migrate": "Flask-Migrate",
				"yoyo":          "Yoyo",
				"sqlalchemy":    "SQLAlchemy",
			}
			for dep, name := range pythonMigrations {
				if strings.Contains(strings.ToLower(string(content)), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check Node.js dependencies
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			nodeMigrations := map[string]string{
				"prisma":      "Prisma",
				"drizzle-orm": "Drizzle",
				"knex":        "Knex",
				"sequelize":   "Sequelize",
				"typeorm":     "TypeORM",
				"mikro-orm":   "MikroORM",
				"@mikro-orm/": "MikroORM",
				"db-migrate":  "db-migrate",
				"umzug":       "Umzug",
			}
			for dep, name := range nodeMigrations {
				if strings.Contains(string(content), `"`+dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check Java dependencies
	javaFiles := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, file := range javaFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			javaMigrations := map[string]string{
				"flyway":    "Flyway",
				"liquibase": "Liquibase",
				"jooq":      "jOOQ",
			}
			for dep, name := range javaMigrations {
				if strings.Contains(strings.ToLower(string(content)), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check Rust dependencies
	if safepath.Exists(path, "Cargo.toml") {
		if content, err := safepath.ReadFile(path, "Cargo.toml"); err == nil {
			rustMigrations := map[string]string{
				"diesel":   "Diesel",
				"refinery": "Refinery",
				"sqlx":     "SQLx",
				"sea-orm":  "SeaORM",
			}
			for dep, name := range rustMigrations {
				if strings.Contains(string(content), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Build result
	if len(found) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Migrations: " + strings.Join(found, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No database migrations found (consider adding if using a database)"
	}

	return result, nil
}
