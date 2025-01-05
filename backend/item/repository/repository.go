package repository

import "github.com/ido50/sqlz"

type Repository struct {
	db *sqlz.DB
}

func New(db *sqlz.DB) *Repository {

	return &Repository{db: db}
}
