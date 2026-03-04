package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(path string, reset bool) *sql.DB {

	// automatically create file if missing
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(time.Hour)

	// NOTE: improve concurrency on web apps
	_, err = db.Exec(`PRAGMA journal_mode = WAL;`)
	if err != nil {
		panic(err)
	}

	schema := `
        CREATE TABLE IF NOT EXISTS lookup (
            id     INTEGER PRIMARY KEY AUTOINCREMENT,
            origin TEXT NOT NULL,
            code   TEXT NOT NULL UNIQUE
        );
        CREATE INDEX IF NOT EXISTS idx_lookup_code ON lookup(code);
    `

	_, err = db.Exec(schema)
	if err != nil {
		panic(err)
	}

	return db
}
