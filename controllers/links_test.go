package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/controllers"
	"shortener/models"
	"shortener/models/options"
	testutils "shortener/testUtils"
	"sort"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const user = "test-user"

var links = []string{"http://example.com", "https://test.com", "https://github.com"}

func TestLinkFlows(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var suite testutils.PostgresSuite
	suite.SetupSuite()
	defer suite.TearDownSuite()
	defer func() {
		statement := "delete from users where login in ($1)"
		result, _ := suite.GetDB().Exec(statement, user)
		count, _ := result.RowsAffected()
		require.Equal(t, int64(1), count)
	}()

	controller := controllers.NewLinkController(suite.GetDB())

	user := AcquireUser(suite)

	links := CreateLinks(controller, t, user)
	FetchLinks(controller, t, user, links)
	FetchLinksByIds(controller, t, links)
	FetchLinks(controller, t, user, links)
}

func AcquireUser(suite testutils.PostgresSuite) models.User {
	controller := controllers.NewUserController(suite.GetDB())

	prepare := func() (*httptest.ResponseRecorder, *http.Request) {
		json, _ := json.Marshal(controllers.UserRequest{user, "superman"})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(json))

		return w, r
	}

	w, r := prepare()
	controller.Create(w, r)
	user := new(models.User)
	json.NewDecoder(w.Body).Decode(&user)

	return *user
}

func CreateLinks(controller controllers.LinkController, t *testing.T, user models.User) []*models.Link {
	tests := []struct {
		link string
		name string
		code int
	}{
		{links[0], "should create link " + links[0], http.StatusCreated},
		{links[1], "should create link " + links[1], http.StatusCreated},
		{links[2], "should create link " + links[2], http.StatusCreated},
		{links[0], "should create link " + links[0], http.StatusCreated},
		{links[1], "should create link " + links[1], http.StatusCreated},
		{links[2], "should create link " + links[2], http.StatusCreated},
	}
	var linksResponses []*models.Link

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(models.Link{
				URL: test.link,
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonValue))
			ctx := context.WithValue(r.Context(), "user", &user)
			r = r.WithContext(ctx)
			method := controller.Create

			method(w, r)

			response := w.Result()
			assert.Equal(t, test.code, response.StatusCode)

			createdLink := new(models.Link)

			json.NewDecoder(response.Body).Decode(&createdLink)

			linksResponses = append(linksResponses, createdLink)
		})
	}

	sort.SliceStable(linksResponses, func(i, j int) bool {
		return linksResponses[i].Created.After(linksResponses[j].Created)
	})

	return linksResponses
}

func FetchLinks(controller controllers.LinkController, t *testing.T, user models.User, links []*models.Link) {
	tests := []struct {
		name   string
		links  []*models.Link
		limit  int
		offset int
		code   int
	}{
		{"should return all links", links, 25, 0, 200},
		{"should return first three links", links[:3], 3, 0, 200},
		{"should return links from second to the fourth", links[1:4], 3, 1, 200},
		{"should return all links from third", links[2:], 25, 2, 200},
		{"should return empty array", []*models.Link{}, 25, 25, 200},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := context.WithValue(r.Context(), "user", &user)
			ctx = context.WithValue(ctx, "options", options.Options{
				Limit:  test.limit,
				Offset: test.offset,
			})
			r = r.WithContext(ctx)
			method := controller.List

			method(w, r)

			values := new([]models.Link)

			response := w.Result()

			json.NewDecoder(response.Body).Decode(&values)

			require.NotNil(t, values)
			require.Len(t, *values, len(test.links))
			assert.Equal(t, http.StatusOK, response.StatusCode)
			for index, link := range test.links {
				assert.Equal(t, link.ID, (*values)[index].ID)
				assert.Equal(t, link.URL, (*values)[index].URL)
				assert.Equal(t, link.UsagesCount, (*values)[index].UsagesCount)
			}
		})
	}
}

func FetchLinksByIds(controller controllers.LinkController, t *testing.T, links []*models.Link) {
	tests := make([]struct {
		name string
		code int
		link *models.Link
	}, 0)

	for _, link := range links {
		tests = append(tests, struct {
			name string
			code int
			link *models.Link
		}{"should fetch link by ID: " + link.ID, 200, link})
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r = mux.SetURLVars(r, map[string]string{
				"id": test.link.ID,
			})
			method := controller.FetchByID

			method(w, r)

			assert.Equal(t, test.link.URL, w.Header().Get("Location"))
			test.link.UsagesCount++
		})
	}

	// Sleep for one second to update usages count
	time.Sleep(time.Second)
}
