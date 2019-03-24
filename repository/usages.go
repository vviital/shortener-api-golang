package repository

import (
	"context"
	"database/sql"
	"shortener/models"
)

// UsageRepository type represents to work with usages
type UsageRepository BaseRepository

// NewUsageRepository creates users repository
func NewUsageRepository(db *sql.DB) UsageRepository {
	return UsageRepository{
		db: db,
	}
}

// Create saves new usage object to the database
func (u *UsageRepository) Create(ctx context.Context, UrlID string) (models.Usage, error) {
	usage := models.Usage{
		UrlID: UrlID,
	}
	statement := "insert into usages (link_id) values($1) returning id, created"

	err := u.db.QueryRowContext(ctx, statement, UrlID).Scan(&usage.ID, &usage.Created)

	return usage, err
}
