package configuration

import (
	"os"
	"strconv"
	"sync"
)

type configuration struct {
	PostgreSQLUrl string
	TokenSecret   string
	AnonUserLogin string
	TokenTTL      int
	initialized   bool
}

var config configuration
var once sync.Once
var mutex sync.Mutex

func loadEnv() {
	config.PostgreSQLUrl, _ = os.LookupEnv("SQL_DB_URL")
	config.TokenSecret, _ = os.LookupEnv("JWT_TOKEN_SECRET")
	config.AnonUserLogin, _ = os.LookupEnv("ANON_USER_LOGIN")
	tokenTTL, _ := os.LookupEnv("TOKEN_TTL")
	config.TokenTTL, _ = strconv.Atoi(tokenTTL)
}

// GetConfiguration from env
func GetConfiguration() *configuration {
	once.Do(loadEnv)

	return &config
}

// Reload function replace configuration. Please be advise and use it only for test purpose
func Reload() {
	mutex.Lock()
	defer mutex.Unlock()
	once = sync.Once{}
	loadEnv()
}
