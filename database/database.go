package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func Init(path string) (Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return Database{}, err
	}

	database := Database{db: db}
	return database, nil
}
