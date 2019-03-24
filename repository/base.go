package repository

import "database/sql"

// BaseRepository type represents base repository to reuse across the app
type BaseRepository struct {
	db *sql.DB
}
