package repository

import (
	"context"
	"database/sql"
	"shortener/models"
)

// UsageRepository type represents to work with usages
type UsageRepository BaseRepository

// UsageRepositoryInterface interface
type UsageRepositoryInterface interface {
	Create(string) (*models.Usage, error)
	CreateWithContext(context.Context, string) (*models.Usage, error)
}

// NewUsageRepository creates users repository
func NewUsageRepository(db *sql.DB) UsageRepositoryInterface {
	return &UsageRepository{
		db: db,
	}
}

// Create saves new usage object to the database
func (repository *UsageRepository) Create(UrlID string) (*models.Usage, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return repository.CreateWithContext(ctx, UrlID)
}

// CreateWithContext saves new usage object to the database
func (repository *UsageRepository) CreateWithContext(ctx context.Context, UrlID string) (*models.Usage, error) {
	usage := models.Usage{
		UrlID: UrlID,
	}
	statement := "insert into usages (link_id) values($1) returning id, created"

	err := repository.db.QueryRowContext(ctx, statement, UrlID).Scan(
		&usage.ID,
		&usage.Created,
	)

	return &usage, err
}
