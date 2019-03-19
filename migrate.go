package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	pg "github.com/lib/pq"
)

var maxSteps = 100

// MigrateDatabase is used to migrate database to the
func MigrateDatabase() {
	config = GetConfiguration()

	url, err := pg.ParseURL(config.PostgreSQLUrl)

	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal(err)
	}

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

	for {
		log.Println("applying next migration")
		err = m.Steps(1)
		if err == os.ErrNotExist {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
}
