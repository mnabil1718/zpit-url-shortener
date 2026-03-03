package model

import "database/sql"

type Lookup struct {
	id     int
	origin string
	code   string
}

type ILookup interface {
	Insert(origin, code string) error
	GetByCode(code string) (string, error)
}

type SQLiteLookup struct {
	db *sql.DB
}

func NewSQliteLookup(db *sql.DB) *SQLiteLookup {
	return &SQLiteLookup{db: db}
}

func (l *SQLiteLookup) Insert(origin, code string) error {
	SQL := `insert into lookup (origin, code) values (?, ?)`
	if _, err := l.db.Exec(SQL, origin, code); err != nil {
		return err
	}

	return nil
}

func (l *SQLiteLookup) GetByCode(code string) (string, error) {
	SQL := `select origin from lookup where code = ? limit 1`
	var origin string

	if err := l.db.QueryRow(SQL, code).Scan(&origin); err != nil {
		return "", err
	}

	return origin, nil
}
