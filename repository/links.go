package repository

import (
	"context"
	"database/sql"
	"shortener/models"
	"shortener/models/options"
)

// LinksRepositoryInterface interface
type LinksRepositoryInterface interface {
	CountByUser(models.User) (int64, error)
	CountByUserWithContext(context.Context, models.User) (int64, error)
	Create(models.Link) (*models.Link, error)
	CreateWithContext(context.Context, models.Link) (*models.Link, error)
	Delete(models.Link) error
	DeleteWithContext(context.Context, models.Link) error
	FindAllByUser(models.User, options.Options) ([]*models.Link, error)
	FindAllByUserWithContext(context.Context, models.User, options.Options) ([]*models.Link, error)
	FindByID(models.Link) (*models.Link, error)
	FindByIDWithContext(context.Context, models.Link) (*models.Link, error)
}

// LinkRepository type represents to work with usages
type LinkRepository BaseRepository

// NewSQLLinkRepository creates LinkRepository repository
func NewSQLLinkRepository(db *sql.DB) LinksRepositoryInterface {
	return &LinkRepository{db: db}
}

// CountByUser return total count of user's links
func (repository *LinkRepository) CountByUser(user models.User) (int64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.CountByUserWithContext(ctx, user)
}

// CountByUserWithContext return total count of user's links
func (repository *LinkRepository) CountByUserWithContext(ctx context.Context, user models.User) (int64, error) {
	var count int64
	statement := "select count(*) from links where user_id = $1"
	err := repository.db.QueryRowContext(ctx, statement, user.ID).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create saves user's link to the database
func (repository *LinkRepository) Create(link models.Link) (*models.Link, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.CreateWithContext(ctx, link)
}

// CreateWithContext saves user's link to the database
func (repository *LinkRepository) CreateWithContext(ctx context.Context, link models.Link) (*models.Link, error) {
	statement := "insert into links (url, user_id) values($1, $2) returning id, created"
	err := repository.db.QueryRowContext(ctx, statement, link.URL, link.UserID).Scan(&link.ID, &link.Created)

	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (repository *LinkRepository) Delete(link models.Link) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.DeleteWithContext(ctx, link)
}

func (repository *LinkRepository) DeleteWithContext(ctx context.Context, link models.Link) error {
	statement := "delete from links where id = $1"

	_, err := repository.db.ExecContext(ctx, statement, link.ID)

	return err
}

// FindAllByUser returns user's links
func (repository *LinkRepository) FindAllByUser(user models.User, opts options.Options) ([]*models.Link, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.FindAllByUserWithContext(ctx, user, opts)
}

// FindAllByUserWithContext returns user's links
func (repository *LinkRepository) FindAllByUserWithContext(ctx context.Context, user models.User, opts options.Options) ([]*models.Link, error) {
	statement := `
		select l.id, l.url, l.created, count(u.id) as usagesCount
		from links l
		left join usages u
		on l.id = u.link_id
		where user_id = $1
		group by l.id
		order by l.created desc
		limit $2
		offset $3
		`

	rows, err := repository.db.QueryContext(ctx, statement, user.ID, opts.Limit, opts.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	links := make([]*models.Link, 0)

	for rows.Next() {
		var link models.Link

		if err = rows.Scan(&link.ID, &link.URL, &link.Created, &link.UsagesCount); err == nil {
			links = append(links, &link)
		} else {
			return nil, err
		}
	}

	return links, nil
}

// FindByID returns link by link id
func (repository *LinkRepository) FindByID(link models.Link) (*models.Link, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.FindByIDWithContext(ctx, link)
}

// FindByIDWithContext returns link by link id
func (repository *LinkRepository) FindByIDWithContext(ctx context.Context, link models.Link) (*models.Link, error) {
	statement := "select id, url, created from links where id = $1"
	err := repository.db.QueryRowContext(ctx, statement, link.ID).Scan(
		&link.ID,
		&link.URL,
		&link.Created,
	)

	return &link, err
}
