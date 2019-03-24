package configuration

import (
	"os"
	"strconv"

	"github.com/subosito/gotenv"
)

type configuration struct {
	PostgreSQLUrl string
	TokenSecret   string
	AnonUserLogin string
	TokenTTL      int
}

var config configuration

// GetConfiguration from env
func GetConfiguration() configuration {
	return config
}

func init() {
	gotenv.Load()

	config.PostgreSQLUrl, _ = os.LookupEnv("SQL_DB_URL")
	config.TokenSecret, _ = os.LookupEnv("JWT_TOKEN_SECRET")
	config.AnonUserLogin, _ = os.LookupEnv("ANON_USER_LOGIN")
	tokenTTL, _ := os.LookupEnv("TOKEN_TTL")
	config.TokenTTL, _ = strconv.Atoi(tokenTTL)
}
