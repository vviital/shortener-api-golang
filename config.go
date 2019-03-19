package main

import (
	"os"

	"github.com/subosito/gotenv"
)

type configuration struct {
	PostgreSQLUrl string
}

var config configuration

// GetConfiguration from env
func GetConfiguration() configuration {
	return config
}

func init() {
	gotenv.Load()

	config.PostgreSQLUrl, _ = os.LookupEnv("SQL_DB_URL")
}
