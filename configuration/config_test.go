package configuration_test

import (
	"fmt"
	"os"
	"shortener/configuration"
	"strconv"
	"testing"
)

const sqlURL = "sql_test"
const jwtToken = "jwt_test"
const anonLogin = "anon_login"
const tokenTTL = 30000

func init() {
	fmt.Println("setting env")
	os.Setenv("SQL_DB_URL", sqlURL)
	os.Setenv("JWT_TOKEN_SECRET", jwtToken)
	os.Setenv("ANON_USER_LOGIN", anonLogin)
	os.Setenv("TOKEN_TTL", strconv.Itoa(tokenTTL))
}

func TestGetConfiguration(t *testing.T) {
	config := configuration.GetConfiguration()

	if config.AnonUserLogin != anonLogin {
		t.Errorf("Login is incorrect. Expected: " + anonLogin + " actual: " + config.AnonUserLogin)
	}

	if config.PostgreSQLUrl != sqlURL {
		t.Errorf("PostgreSQL connection url is incorrect. Expected:" + sqlURL + " actual: " + config.PostgreSQLUrl)
	}

	if config.TokenSecret != jwtToken {
		t.Errorf("JWT token secret is incorrect. Expected: " + jwtToken + " actual: " + config.TokenSecret)
	}

	if config.TokenTTL != tokenTTL {
		t.Errorf("Token TTL is not specified. Expected: " + strconv.Itoa(tokenTTL) + " actual: " + strconv.Itoa(config.TokenTTL))
	}
}
