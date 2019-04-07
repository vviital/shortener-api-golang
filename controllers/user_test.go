package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/controllers"
	testutils "shortener/testUtils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userTestResult struct {
	code     int
	hasToken bool
}

var users = []string{"test-user", "test-user-2"}

func TestUserFlows(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var suite testutils.PostgresSuite
	suite.SetupSuite()
	defer suite.TearDownSuite()
	defer func() {
		statement := "delete from users where login in ($1, $2)"
		result, _ := suite.GetDB().Exec(statement, users[0], users[1])
		count, _ := result.RowsAffected()
		require.Equal(t, int64(2), count)
	}()

	controller := controllers.NewUserController(suite.GetDB())

	userRequests := CreateUser(controller, t)

	for _, request := range userRequests {
		AuthorizeUser(controller, t, request)
	}
}

func CreateUser(controller controllers.UserController, t *testing.T) []controllers.UserRequest {
	tests := []struct {
		login          string
		password       string
		name           string
		expectedResult userTestResult
	}{
		{users[0], "superman", "should create a user", userTestResult{201, false}},
		{users[0], "batman", "should thrown an error because user exists", userTestResult{400, false}},
		{users[1], "spiderman", "should create one more user", userTestResult{201, false}},
	}
	var userRequests []controllers.UserRequest

	for _, test := range tests {
		jsonValue, _ := json.Marshal(controllers.UserRequest{test.login, test.password})

		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonValue))
			method := controller.Create

			method(w, r)
			response := w.Result()

			assert.Equal(t, test.expectedResult.code, response.StatusCode)

			if test.expectedResult.code == http.StatusCreated {
				userRequests = append(userRequests, controllers.UserRequest{test.login, test.password})
			}
		})
	}

	return userRequests
}

func AuthorizeUser(controller controllers.UserController, t *testing.T, user controllers.UserRequest) {
	tests := []struct {
		request        controllers.UserRequest
		name           string
		expectedResult userTestResult
	}{
		{user, "should authorize user", userTestResult{http.StatusOK, true}},
		{controllers.UserRequest{user.Login, user.Password + user.Password}, "should throw an error is user has wrong credentials", userTestResult{http.StatusUnauthorized, false}},
	}

	for _, test := range tests {
		jsonValue, _ := json.Marshal(test.request)

		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonValue))
			method := controller.Authorize

			method(w, r)

			response := w.Result()
			assert.Equal(t, test.expectedResult.code, response.StatusCode)

			respv := make(map[string]interface{})
			json.NewDecoder(response.Body).Decode(&respv)

			if test.expectedResult.hasToken {
				token, ok := respv["token"]
				require.True(t, ok)
				tokenString, ok := token.(string)
				require.True(t, ok)
				require.NotEmpty(t, tokenString)
			}
		})
	}
}
