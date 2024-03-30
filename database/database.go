package database

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	session *sql.DB
	ctx     context.Context
}

func Init(path string) (Database, error) {
	session, err := sql.Open("sqlite3", path)
	if err != nil {
		return Database{}, err
	}

	database := Database{session: session, ctx: context.TODO()}
	err = database.CreateTasksTable()
	if err != nil {
		return Database{}, err
	}
	return database, nil
}

func (db Database) CreateTasksTable() error {
	stmt := `
CREATE TABLE IF NOT EXISTS TASKS (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	TITLE TEXT,
	DESCRIPTION TEXT,
	STATUS INT,
	INSERTED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	DELETED_AT TIMESTAMP
)`
	_, err := db.session.ExecContext(db.ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) InsertTask(title string, description string, status int) error {
	stmt, err := db.session.PrepareContext(
		db.ctx,
		`INSERT INTO TASKS (TITLE, DESCRIPTION, STATUS) VALUES (?, ?, ?)`,
	)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		db.ctx,
		title,
		description,
		status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) UpdateTaskStatus(id int, status int) error {
	stmt, err := db.session.PrepareContext(
		db.ctx,
		`UPDATE TASKS SET STATUS = ? WHERE id = ?`,
	)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(db.ctx, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) DeleteTask(id int) error {
	stmt, err := db.session.PrepareContext(
		db.ctx,
		`UPDATE TASKS SET DELETED_AT = CURRENT_TIMESTAMP WHERE id = ?`,
	)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(db.ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) ReadAllTasks() (*sql.Rows, error) {
	stmt := `SELECT id, title, description, status FROM TASKS WHERE DELETED_AT IS NULL`
	rows, err := db.session.QueryContext(db.ctx, stmt)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
