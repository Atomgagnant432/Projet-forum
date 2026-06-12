package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open(
		"sqlite3",
		"./back-end/database/forum.db",
	)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}