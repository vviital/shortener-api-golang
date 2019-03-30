package migrator

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

// MigrateDatabase is used to migrate database to the next state
func MigrateDatabase(db *sql.DB) {
	directory := "./migrations"
	MigrateDatabaseFromDirectory(db, directory, 1)
}

// MigrateDatabaseFromDirectory is used to migrate database to the next state.
// Migrations are used from "directory"
func MigrateDatabaseFromDirectory(db *sql.DB, directory string, direction int) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+directory,
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("start migrations")

	for {
		log.Println("applying next migration")
		err = m.Steps(1 * direction)
		if err == os.ErrNotExist {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("migrations applied")
}
