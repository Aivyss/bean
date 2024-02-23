package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "db-file.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db, nil
}
