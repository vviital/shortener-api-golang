package repository_test

import (
	"database/sql"
	"shortener/models"
	"shortener/models/options"
	"shortener/repository"
	testutils "shortener/testUtils"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/teris-io/shortid"
)

func UserTestSuite(t *testing.T, r *testutils.Repositories) {
	if testing.Short() {
		t.Skip("Skip tests for user repository for unit tests")
	}

	user := new(models.User)
	var err error
	login, _ := shortid.Generate()

	t.Run("should create new user", func(t *testing.T) {
		user, err = r.Users.Create(models.User{
			Login: login,
		})

		require.Nil(t, err, "should create user without errors")
		assert.Equal(t, login, user.Login, "should save login correctly")
		assert.NotEmpty(t, user.ID, "should generate ID for the user")
	})

	t.Run("should throw an error when try to create user with the same login", func(t *testing.T) {
		nilUser, err := r.Users.Create(models.User{
			Login:    login,
			Password: "password",
		})

		// TODO: add custom error for such case
		assert.Error(t, err, "should return error that user could not be created")
		assert.Nil(t, nilUser, "user should be nil")
	})

	t.Run("should find user by login", func(t *testing.T) {
		foundUser, err := r.Users.FindByLogin(user.Login)

		require.Nil(t, err, "should find user by login without errors")
		assert.Equal(t, user, foundUser, "user should be found correctly")
	})

	t.Run("should find user by ID", func(t *testing.T) {
		foundUser, err := r.Users.FindByID(user.ID, options.Options{})

		require.Nil(t, err, "should find user by ID without errors")
		assert.Equal(t, user, foundUser, "user should be found correctly")
	})

	t.Run("should fetch user with links and usages count", func(t *testing.T) {
		link, err := r.Links.Create(models.Link{
			UserID: user.ID,
			URL:    "https://example.com",
		})
		require.Nil(t, err)

		for range make([]struct{}, 10) {
			_, err := r.Usages.Create(link.ID)
			assert.Nil(t, err)
		}

		foundUser, err := r.Users.FindByID(user.ID, options.Options{
			Limit: 25,
		})

		assert.Nil(t, err)

		fetchedLinks := foundUser.Links

		assert.Len(t, fetchedLinks, 1)
		assert.Equal(t, int64(10), fetchedLinks[0].UsagesCount)
	})

	t.Run("should delete user", func(t *testing.T) {
		err = r.Users.Delete(*user)

		assert.Nil(t, err, "should delete user without errors")

		foundUser, err := r.Users.FindByID(user.ID, options.Options{})

		assert.Nil(t, foundUser, "found user should be nil")
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("should throw an error when try to delete nonexisting user", func(t *testing.T) {
		err = r.Users.Delete(*user)

		assert.Equal(t, sql.ErrNoRows, err, "should throw an error when try to delete the nonexisting user")
	})
}

func TestUsersPostgres(t *testing.T) {
	var suite testutils.PostgresSuite
	suite.SetupSuite()
	defer suite.TearDownSuite()

	UserTestSuite(t, &testutils.Repositories{
		Usages: repository.NewUsageRepository(suite.GetDB()),
		Links:  repository.NewSQLLinkRepository(suite.GetDB()),
		Users:  repository.NewUserRepository(suite.GetDB()),
	})
}
