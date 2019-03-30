package driver

import (
	"database/sql"
	"log"
	"shortener/configuration"

	"github.com/lib/pq"
)

// ConnectPostgreSQL returns a connection to the PostgreSQL database
func ConnectPostgreSQL() *sql.DB {
	url := configuration.GetConfiguration().PostgreSQLUrl
	pgUrl, err := pq.ParseURL(url)

	if err != nil {
		log.Panicln(err)
	}

	db, err := sql.Open("postgres", pgUrl)

	if err != nil {
		log.Panicln(err)
	}

	err = db.Ping()

	if err != nil {
		log.Panicln(err)
	}

	return db
}
