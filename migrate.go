package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

// MigrateDatabase is used to migrate database to the
func MigrateDatabase(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("start migrations")

	for {
		log.Println("applying next migration")
		err = m.Steps(1)
		if err == os.ErrNotExist {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("migrations applied")
}
