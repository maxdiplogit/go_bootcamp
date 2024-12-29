package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to SQLite3 Database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS Users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`

	_, err := DB.Exec(createUsersTable)

	if err != nil {
		panic("Could not create Users Table in the SQLite3 Database")
	}

	createEventsTable := `
		CREATE TABLE IF NOT EXISTS Events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime DATETIME NOT NULL,
			user_id INTEGER,
			FOREIGN KEY(user_id) REFERENCES Users(id)
		)
	`

	_, err = DB.Exec(createEventsTable)

	if err != nil {
		panic("Could not create Events Table in the SQLite3 Database")
	}

	createRegistrationsTable := `
		CREATE TABLE IF NOT EXISTS Registrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER,
			user_id INTEGER,
			FOREIGN KEY(event_id) REFERENCES Events(id),
			FOREIGN KEY(user_id) REFERENCES Users(id)
		)
	`

	_, err = DB.Exec(createRegistrationsTable)

	if err != nil {
		panic("Could not create Registrations Table in SQLite3 Database")
	}
}
