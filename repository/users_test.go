package repository_test

import (
	"database/sql"
	"shortener/models"
	"shortener/models/options"
	"shortener/repository"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/teris-io/shortid"
)

func TestUsersOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip tests for user repository for unit tests")
	}

	var suite PostgresSuite
	suite.SetupSuite()
	defer suite.TearDownSuite()
	user := new(models.User)
	var err error
	usages := repository.NewUsageRepository(suite.db)
	links := repository.NewSQLLinkRepository(suite.db)
	repository := repository.NewUserRepository(suite.db)

	login, _ := shortid.Generate()

	t.Run("should create new user", func(t *testing.T) {
		user, err = repository.Create(models.User{
			Login: login,
		})

		require.Nil(t, err, "should create user without errors")
		assert.Equal(t, login, user.Login, "should save login correctly")
		assert.NotEmpty(t, user.ID, "should generate ID for the user")
	})

	t.Run("should throw an error when try to create user with the same login", func(t *testing.T) {
		nilUser, err := repository.Create(models.User{
			Login:    login,
			Password: "password",
		})

		// TODO: add custom error for such case
		assert.Error(t, err, "should return error that user could not be created")
		assert.Nil(t, nilUser, "user should be nil")
	})

	t.Run("should find user by login", func(t *testing.T) {
		foundUser, err := repository.FindByLogin(user.Login)

		require.Nil(t, err, "should find user by login without errors")
		assert.Equal(t, user, foundUser, "user should be found correctly")
	})

	t.Run("should find user by ID", func(t *testing.T) {
		foundUser, err := repository.FindByID(user.ID, options.Options{})

		require.Nil(t, err, "should find user by ID without errors")
		assert.Equal(t, user, foundUser, "user should be found correctly")
	})

	t.Run("should fetch user with links and usages count", func(t *testing.T) {
		link, err := links.Create(models.Link{
			UserID: user.ID,
			URL:    "https://example.com",
		})
		require.Nil(t, err)

		for range make([]struct{}, 10) {
			_, err := usages.Create(link.ID)
			assert.Nil(t, err)
		}

		foundUser, err := repository.FindByID(user.ID, options.Options{
			Limit: 25,
		})

		assert.Nil(t, err)

		fetchedLinks := foundUser.Links

		assert.Len(t, fetchedLinks, 1)
		assert.Equal(t, int64(10), fetchedLinks[0].UsagesCount)
	})

	t.Run("should delete user", func(t *testing.T) {
		err = repository.Delete(*user)

		assert.Nil(t, err, "should delete user without errors")

		foundUser, err := repository.FindByID(user.ID, options.Options{})

		assert.Nil(t, foundUser, "found user should be nil")
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("should throw an error when try to delete nonexisting user", func(t *testing.T) {
		err = repository.Delete(*user)

		assert.Equal(t, sql.ErrNoRows, err, "should throw an error when try to delete the nonexisting user")
	})
}
