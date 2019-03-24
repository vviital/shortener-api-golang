package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"shortener/models"
	"shortener/models/options"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// UserRepository type represent repository to work with UserRepository
type UserRepository struct {
	BaseRepository
	links LinkRepository
}

// NewUserRepository creates UserRepository repository
func NewUserRepository(db *sql.DB) UserRepository {
	var repository UserRepository

	repository.BaseRepository = BaseRepository{db}

	return repository
}

func (u *UserRepository) queryForAUserRecord(ctx context.Context, user *models.User, statement string, args ...interface{}) error {
	err := u.db.QueryRowContext(ctx, statement, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.Created,
	)

	return err
}

// FindByLogin returns user by login. This method is used only for Login action.
func (u *UserRepository) FindByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User

	statement := "select id, login, password, created from users where login = $1"

	err := u.queryForAUserRecord(ctx, &user, statement, login)

	return user, err
}

// Create saves new user object to the database
func (u *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	statement := "insert into users (login, password) values ($1, $2) returning id, created"

	log.Println("user.Login", user.Login)
	log.Println("user.Password", user.Password)
	log.Println("ctx", ctx)

	err := u.db.QueryRowContext(ctx, statement, user.Login, user.Password).Scan(&user.ID, &user.Created)

	if err != nil && strings.Contains(err.Error(), "violates unique constraint") {
		spew.Dump(err)
		return user, errors.New("User with login " + user.Login + " already exists")
	}

	user.CleanPrivateFields()

	return user, err
}

func (u *UserRepository) fetchAdditionalFieldsForUser(ctx context.Context, user *models.User, opts options.Options) error {
	var errc chan error
	var done chan bool
	var links []models.Link
	var count int64

	defer close(errc)
	defer close(done)

	var handlerSubRequestError = func(e error) {
		if e != nil {
			done <- true
		} else {
			errc <- e
		}
	}

	var goRoutineCounter = 2

	go func() {
		c, e := u.links.GetLinksCountForUser(ctx, *user)
		count = c
		handlerSubRequestError(e)
	}()

	go func() {
		l, e := u.links.GetUserLinks(ctx, *user, opts)
		links = l
		handlerSubRequestError(e)
	}()

	for i := 0; i < goRoutineCounter; i++ {
		select {
		case <-done:
			continue
		case e := <-errc:
			return e
		}
	}

	user.Links = links
	user.LinksCount = count

	return nil
}

// FindByID returns user by id. This is a preferable way to fetch user in most cases
func (u *UserRepository) FindByID(ctx context.Context, ID string, opts options.Options) (*models.User, error) {
	var user models.User

	statement := "select id, login, password, created from users where id = $1"

	err := u.queryForAUserRecord(ctx, &user, statement, ID)

	user.CleanPrivateFields()

	if err != nil {
		return nil, err
	}

	u.fetchAdditionalFieldsForUser(ctx, &user, opts)

	return &user, nil
}
