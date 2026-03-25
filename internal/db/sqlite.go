package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(path string, reset bool) *sql.DB {
	// ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(fmt.Sprintf("failed to create db directory: %v", err))
	}

	// automatically create file if missing
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	pragmas := []string{
		`PRAGMA journal_mode = WAL;`,    // WAL improve concurrency on web apps
		`PRAGMA synchronous = NORMAL;`,  // safe with WAL, much faster than FULL
		`PRAGMA cache_size = -64000;`,   // 64MB page cache
		`PRAGMA temp_store = MEMORY;`,   // temp tables in RAM
		`PRAGMA mmap_size = 268435456;`, // 256MB memory-mapped I/O
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			panic(fmt.Sprintf("failed to set pragma %s: %v", p, err))
		}
	}

	schema := `
        CREATE TABLE IF NOT EXISTS lookup (
            id     INTEGER PRIMARY KEY AUTOINCREMENT,
            origin TEXT NOT NULL,
            code   TEXT NOT NULL UNIQUE,
			clicks INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_lookup_code ON lookup(code);
		CREATE INDEX IF NOT EXISTS idx_lookup_created_at ON lookup(created_at);
    `

	_, err = db.Exec(schema)
	if err != nil {
		panic(err)
	}

	return db
}
