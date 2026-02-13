package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type MigrationsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *MigrationsCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "migrations-test-*")
	s.Require().NoError(err)
}

func (s *MigrationsCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *MigrationsCheckTestSuite) TestIDAndName() {
	check := &MigrationsCheck{}
	s.Equal("common:migrations", check.ID())
	s.Equal("Database Migrations", check.Name())
}

func (s *MigrationsCheckTestSuite) TestMigrationsDirectory() {
	migrationsDir := filepath.Join(s.tempDir, "migrations")
	err := os.MkdirAll(migrationsDir, 0755)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "migrations")
}

func (s *MigrationsCheckTestSuite) TestDbMigrationsDirectory() {
	migrationsDir := filepath.Join(s.tempDir, "db", "migrations")
	err := os.MkdirAll(migrationsDir, 0755)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "db/migrations")
}

func (s *MigrationsCheckTestSuite) TestAlembicConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "alembic.ini"), []byte(`
[alembic]
script_location = alembic
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Alembic")
}

func (s *MigrationsCheckTestSuite) TestFlywayConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "flyway.conf"), []byte(`
flyway.url=jdbc:postgresql://localhost/mydb
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Flyway")
}

func (s *MigrationsCheckTestSuite) TestLiquibaseConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "liquibase.properties"), []byte(`
changeLogFile=changelog.xml
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Liquibase")
}

func (s *MigrationsCheckTestSuite) TestPrismaSchema() {
	prismaDir := filepath.Join(s.tempDir, "prisma")
	err := os.MkdirAll(prismaDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(prismaDir, "schema.prisma"), []byte(`
model User {
  id Int @id
}
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Prisma")
}

func (s *MigrationsCheckTestSuite) TestDrizzleConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "drizzle.config.ts"), []byte(`
export default defineConfig({});
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Drizzle")
}

func (s *MigrationsCheckTestSuite) TestKnexConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "knexfile.js"), []byte(`
module.exports = {
  development: {}
};
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Knex")
}

func (s *MigrationsCheckTestSuite) TestSequelizeConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".sequelizerc"), []byte(`
module.exports = {};
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Sequelize")
}

func (s *MigrationsCheckTestSuite) TestDieselConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "diesel.toml"), []byte(`
[print_schema]
file = "src/schema.rs"
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Diesel")
}

func (s *MigrationsCheckTestSuite) TestGoMigrate() {
	content := `module myapp

go 1.21

require (
	github.com/golang-migrate/migrate/v4 v4.16.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "golang-migrate")
}

func (s *MigrationsCheckTestSuite) TestGoose() {
	content := `module myapp

go 1.21

require (
	github.com/pressly/goose/v3 v3.15.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Goose")
}

func (s *MigrationsCheckTestSuite) TestGORM() {
	content := `module myapp

go 1.21

require (
	gorm.io/gorm v1.25.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "GORM")
}

func (s *MigrationsCheckTestSuite) TestAtlas() {
	err := os.WriteFile(filepath.Join(s.tempDir, "atlas.hcl"), []byte(`
schema "public" {}
`), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Atlas")
}

func (s *MigrationsCheckTestSuite) TestPythonAlembic() {
	content := `alembic==1.12.0
sqlalchemy==2.0.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Alembic")
}

func (s *MigrationsCheckTestSuite) TestPythonDjango() {
	content := `Django==4.2.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Django")
}

func (s *MigrationsCheckTestSuite) TestNodePrisma() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "prisma": "^5.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Prisma")
}

func (s *MigrationsCheckTestSuite) TestNodeDrizzle() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "drizzle-orm": "^0.29.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Drizzle")
}

func (s *MigrationsCheckTestSuite) TestNodeTypeORM() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "typeorm": "^0.3.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "TypeORM")
}

func (s *MigrationsCheckTestSuite) TestJavaFlyway() {
	content := `<dependency>
    <groupId>org.flywaydb</groupId>
    <artifactId>flyway-core</artifactId>
</dependency>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Flyway")
}

func (s *MigrationsCheckTestSuite) TestRustSqlx() {
	content := `[package]
name = "myapp"

[dependencies]
sqlx = "0.7"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "SQLx")
}

func (s *MigrationsCheckTestSuite) TestRustSeaORM() {
	content := `[package]
name = "myapp"

[dependencies]
sea-orm = "0.12"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "SeaORM")
}

func (s *MigrationsCheckTestSuite) TestNoMigrationsFound() {
	content := `module myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No database migrations found")
}

func (s *MigrationsCheckTestSuite) TestEmptyDirectory() {
	check := &MigrationsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No database migrations found")
}

func TestMigrationsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationsCheckTestSuite))
}
