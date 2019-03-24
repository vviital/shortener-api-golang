package repository

import (
	"context"
	"database/sql"
	"shortener/models"
	"shortener/models/options"

	"github.com/davecgh/go-spew/spew"
)

// LinkRepository type represents to work with usages
type LinkRepository BaseRepository

// Create saves user's link to the database
func (l *LinkRepository) Create(ctx context.Context, link models.Link) (models.Link, error) {
	statement := "insert into links (url, user_id) values($1, $2) returning id, created"

	err := l.db.QueryRowContext(ctx, statement, link.URL, link.UserID).Scan(&link.ID, &link.Created)

	return link, err
}

// NewLinkRepository creates LinkRepository repository
func NewLinkRepository(db *sql.DB) LinkRepository {
	return LinkRepository{db}
}

// GetUserLinks returns user's links
func (l *LinkRepository) GetUserLinks(ctx context.Context, user models.User, opts options.Options) ([]models.Link, error) {
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

	spew.Dump("user", user)
	spew.Dump("options", opts)
	rows, err := l.db.QueryContext(ctx, statement, user.ID, opts.Limit, opts.Offset)

	spew.Dump("err 111", err)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []models.Link

	for rows.Next() {
		var link models.Link

		if err = rows.Scan(&link.ID, &link.URL, &link.Created, &link.UsagesCount); err == nil {
			spew.Dump("link", link)
			links = append(links, link)
		} else {
			spew.Dump("wtf", err)
			return nil, err
		}
	}

	return links, nil
}

// GetLinksCountForUser return total count of user's links
func (l *LinkRepository) GetLinksCountForUser(ctx context.Context, user models.User) (int64, error) {
	var count int64
	statement := "select count(*) from links where user_id = $1"
	err := l.db.QueryRowContext(ctx, statement, user.ID).Scan(&count)
	return count, err
}

// FindByID returns link by link id
func (l *LinkRepository) FindByID(ctx context.Context, link models.Link) (models.Link, error) {
	statement := "select id, url, created from links where id = $1"

	err := l.db.QueryRowContext(ctx, statement, link.ID).Scan(
		&link.ID,
		&link.URL,
		&link.Created,
	)

	return link, err
}
