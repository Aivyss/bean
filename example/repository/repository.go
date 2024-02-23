package repository

import (
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (UserRepository, error) {
	if db == nil {
		return UserRepository{}, errors.New("no db")
	}
	return UserRepository{
		db: db,
	}, nil
}
