package repository_test

import (
	"database/sql"
	"shortener/models"
	"shortener/models/options"
	"shortener/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnonUserLinks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip tests for links repository for unit tests")
	}

	var suite PostgresSuite
	suite.SetupSuite()
	defer suite.TearDownSuite()

	linksRepository := repository.NewSQLLinkRepository(suite.db)
	userRepository := repository.NewUserRepository(suite.db)

	user, err := userRepository.FindByLogin("anon")

	require.Nil(t, err, "anon user should be found")

	initialCount, err := linksRepository.CountByUser(*user)

	require.Nil(t, err, "initial count of links should be found without errors")

	var link1 *models.Link
	var link2 *models.Link

	t.Run("should create link for the anon user", func(t *testing.T) {
		link1, err = linksRepository.Create(models.Link{URL: "example.com", UserID: user.ID})

		require.Nil(t, err, "link should be created")

		assert.Equal(t, "example.com", link1.URL, "url should be saved correctly in the links table")
		assert.NotEmpty(t, "example.com", link1.ID, "id should be populated")
	})

	t.Run("should fetch link for the anon user", func(t *testing.T) {
		link2, err = linksRepository.FindByID(models.Link{
			ID: link1.ID,
		})

		require.Nil(t, err, "link should be fetched")

		assert.Equal(t, link1.ID, link2.ID, "ID values should be the same")
		assert.Equal(t, link1.URL, link2.URL, "URL values should be the same")
		assert.Equal(t, link1.Created, link2.Created, "Created dates should be the same")
	})

	t.Run("should increment links count", func(t *testing.T) {
		count, err := linksRepository.CountByUser(*user)

		require.Nil(t, err, "count of links should be found without errors")

		assert.Equal(t, initialCount+1, count, "link count should be incremented compared to the initial count")
	})

	t.Run("should return all user's links", func(t *testing.T) {
		links, err := linksRepository.FindAllByUser(*user, options.Options{
			Limit:  1000,
			Offset: 0,
		})

		require.Nil(t, err, "links should be fetched")

		assert.Contains(t, links, link2)
	})

	t.Run("should delete user's link", func(t *testing.T) {
		err := linksRepository.Delete(*link1)

		require.Nil(t, err, "link should be deleted")

		_, err = linksRepository.FindByID(*link1)

		assert.Equal(t, sql.ErrNoRows, err, "should not find link by ID")
	})
}
