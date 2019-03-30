package driver_test

import (
	"os"
	"shortener/configuration"
	"shortener/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

const correctUrl = "postgres://postgres:postgres@localhost:5432/shortener-tests?sslmode=disable"
const incorrectUrl = "postgres://postgres:postgres@localhost:5432/test?sslmode=disable"

func TestConnectPostgreSQLSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
		return
	}

	os.Setenv("SQL_DB_URL", correctUrl)

	configuration.Reload()

	db := driver.ConnectPostgreSQL()

	assert.NotNil(t, db)
}

func TestConnectPostgreSQLFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "pq: database \"test\" does not exist\n", r)
			return
		}
		t.Error("driver should throw an error")
	}()

	os.Setenv("SQL_DB_URL", incorrectUrl)

	configuration.Reload()

	driver.ConnectPostgreSQL()
}
