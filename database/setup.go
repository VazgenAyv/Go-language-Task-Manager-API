package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal(err)

	}

	createTasksTable()
	createUsersTable()
}

func createTasksTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
        completed BOOLEAN NOT NULL CHECK (completed IN (0, 1)),
		maintask INT
    );`

	statement, err := DB.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func createUsersTable() {
	query := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        role TEXT NOT NULL
    );`
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()
}
