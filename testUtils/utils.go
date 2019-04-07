package testutils

import (
	"database/sql"
	"os"
	"shortener/configuration"
	"shortener/driver"
	"shortener/migrator"
	"shortener/repository"
)

const correctUrl = "postgres://postgres:postgres@localhost:5432/shortener-tests?sslmode=disable"

type Repositories struct {
	Usages repository.UsageRepositoryInterface
	Links  repository.LinksRepositoryInterface
	Users  repository.UserRepositoryInterface
}

// PostgresSuite struct is used to automate logic to work with postgres database in tests
type PostgresSuite struct {
	db *sql.DB
}

// GetDB returns connection to the DB
func (s *PostgresSuite) GetDB() *sql.DB {
	return s.db
}

// SetupSuite create connections to the database and apply migrations
func (s *PostgresSuite) SetupSuite() {
	os.Setenv("SQL_DB_URL", correctUrl)
	configuration.Reload()
	s.db = driver.ConnectPostgreSQL()
	migrator.MigrateDatabaseFromDirectory(s.db, "../migrations", 1)
}

// TearDownSuite destroy connections
func (s *PostgresSuite) TearDownSuite() {
	s.db.Close()
}
