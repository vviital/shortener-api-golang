package repository_test

import (
	"database/sql"
	"os"
	"shortener/configuration"
	"shortener/driver"
	"shortener/migrator"
	"shortener/repository"
)

const correctUrl = "postgres://postgres:postgres@localhost:5432/shortener-tests?sslmode=disable"

type repositories struct {
	usages repository.UsageRepositoryInterface
	links  repository.LinksRepositoryInterface
	users  repository.UserRepositoryInterface
}
type PostgresSuite struct {
	db *sql.DB
}

func (s *PostgresSuite) SetupSuite() {
	os.Setenv("SQL_DB_URL", correctUrl)
	configuration.Reload()
	s.db = driver.ConnectPostgreSQL()
	migrator.MigrateDatabaseFromDirectory(s.db, "../migrations", 1)
}

func (s *PostgresSuite) TearDownSuite() {
	s.db.Close()
}
