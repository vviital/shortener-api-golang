package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"shortener/models"
	"shortener/models/options"
	"strings"
)

// UserRepository type represent repository to work with UserRepository
type UserRepository struct {
	BaseRepository
	links LinksRepositoryInterface
}

// UserRepositoryInterface interface
type UserRepositoryInterface interface {
	Create(models.User) (*models.User, error)
	CreateWithContext(context.Context, models.User) (*models.User, error)
	Delete(models.User) error
	DeleteWithContext(context.Context, models.User) error
	FindByID(string, options.Options) (*models.User, error)
	FindByIDWithContext(context.Context, string, options.Options) (*models.User, error)
	FindByLogin(string) (*models.User, error)
	FindByLoginWithContext(context.Context, string) (*models.User, error)
}

// NewUserRepository creates UserRepository repository
func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	var repository UserRepository

	repository.BaseRepository = BaseRepository{db}
	repository.links = NewSQLLinkRepository(db)

	return &repository
}

func (repository *UserRepository) queryForAUserRecord(ctx context.Context, user *models.User, statement string, args ...interface{}) error {
	err := repository.db.QueryRowContext(ctx, statement, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.Created,
	)

	return err
}

// Create saves new user object to the database
func (repository *UserRepository) Create(user models.User) (*models.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.CreateWithContext(ctx, user)
}

// Create saves new user object to the database
func (repository *UserRepository) CreateWithContext(ctx context.Context, user models.User) (*models.User, error) {
	statement := "insert into users (login, password) values ($1, $2) returning id, created"

	log.Println("user.Login", user.Login)
	log.Println("user.Password", user.Password)
	log.Println("ctx", ctx)

	err := repository.db.QueryRowContext(ctx, statement, user.Login, user.Password).Scan(&user.ID, &user.Created)

	if err != nil && strings.Contains(err.Error(), "violates unique constraint") {
		return nil, errors.New("User with login " + user.Login + " already exists")
	}

	if err != nil {
		return nil, err
	}

	user.CleanPrivateFields()

	return &user, err
}

// Delete user
func (repository *UserRepository) Delete(user models.User) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.DeleteWithContext(ctx, user)
}

// DeleteWithContext user using context
func (repository *UserRepository) DeleteWithContext(ctx context.Context, user models.User) error {
	statement := "delete from users where id = $1"

	result, err := repository.db.ExecContext(ctx, statement, user.ID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	} else if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindByID returns user by id. This is a preferable way to fetch user in most cases
func (repository *UserRepository) FindByID(ID string, opts options.Options) (*models.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.FindByIDWithContext(ctx, ID, opts)
}

// FindByIDWithContext returns user by id. This is a preferable way to fetch user in most cases
func (repository *UserRepository) FindByIDWithContext(ctx context.Context, ID string, opts options.Options) (*models.User, error) {
	var user models.User

	statement := "select id, login, password, created from users where id = $1"

	err := repository.queryForAUserRecord(ctx, &user, statement, ID)

	user.CleanPrivateFields()

	if err != nil {
		return nil, err
	}

	repository.fetchAdditionalFieldsForUser(ctx, &user, opts)

	return &user, nil
}

// FindByLogin returns user by login. This method is used only for Login action.
func (repository *UserRepository) FindByLogin(login string) (*models.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.FindByLoginWithContext(ctx, login)
}

// FindByLoginWithContext returns user by login. This method is used only for Login action.
func (repository *UserRepository) FindByLoginWithContext(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	user.Links = make([]*models.Link, 0)

	statement := "select id, login, password, created from users where login = $1"

	err := repository.queryForAUserRecord(ctx, &user, statement, login)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func combineErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var builder strings.Builder

	for _, err := range errs {
		builder.WriteString(err.Error())
	}

	return errors.New(builder.String())
}

func (repository *UserRepository) fetchAdditionalFieldsForUser(ctx context.Context, user *models.User, opts options.Options) error {
	var goRoutineCounter = 2
	var links []*models.Link
	var count int64

	errc := make(chan error, goRoutineCounter)
	done := make(chan bool, goRoutineCounter)
	defer close(errc)
	defer close(done)

	var errors []error

	var handlerSubRequestError = func(e error) {
		if e == nil {
			done <- true
		} else {
			errc <- e
		}
	}

	go func() {
		c, e := repository.links.CountByUserWithContext(ctx, *user)
		count = c
		handlerSubRequestError(e)
	}()

	go func() {
		l, e := repository.links.FindAllByUserWithContext(ctx, *user, opts)
		links = l
		handlerSubRequestError(e)
	}()

	for i := 0; i < goRoutineCounter; i++ {
		select {
		case <-done:
			continue
		case e := <-errc:
			errors = append(errors, e)
		}
	}

	err := combineErrors(errors...)

	if err != nil {
		return err
	}

	user.Links = links
	user.LinksCount = count

	return nil
}
