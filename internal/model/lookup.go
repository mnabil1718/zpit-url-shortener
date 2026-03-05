package model

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/mnabil1718/zp.it/internal/cache"
)

type Lookup struct {
	ID        int       `json:"id"`
	Origin    string    `json:"origin"`
	Code      string    `json:"code"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

type ILookup interface {
	Insert(origin, code string) error
	GetOriginByCode(code string) (string, error) // faster for redirect
	GetByCode(code string) (*Lookup, error)
	IncrementClicks(code string) error
}

type SQLiteLookup struct {
	db    *sql.DB
	cache cache.ICache
}

func NewSQliteLookup(db *sql.DB, cache cache.ICache) *SQLiteLookup {
	return &SQLiteLookup{
		db:    db,
		cache: cache,
	}
}

func (l *SQLiteLookup) Insert(origin, code string) error {
	var lkp Lookup
	SQL := `insert into lookup (origin, code) values (?, ?)
			returning id, origin, code, clicks, created_at`

	if err := l.db.QueryRow(SQL, origin, code).Scan(
		&lkp.ID,
		&lkp.Origin,
		&lkp.Code,
		&lkp.Clicks,
		&lkp.CreatedAt,
	); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return ErrAlreadyExists
			}
		}

		return err
	}

	if err := l.cache.Set(context.Background(), code, lkp.Origin, 300*time.Second); err != nil {
		return err
	}

	return nil
}

// No cache hit here because cache only stores origin url
func (l *SQLiteLookup) GetByCode(code string) (*Lookup, error) {
	var lkp Lookup

	SQL := `select id, origin, code, clicks, created_at from lookup where code = ? limit 1`
	if err := l.db.QueryRow(SQL, code).Scan(
		&lkp.ID,
		&lkp.Origin,
		&lkp.Code,
		&lkp.Clicks,
		&lkp.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &lkp, nil
}

func (l *SQLiteLookup) GetOriginByCode(code string) (string, error) {
	SQL := `select origin from lookup where code = ? limit 1`
	var origin string

	v, err := l.cache.Get(context.Background(), code)
	if err == nil {
		return v, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		return "", err
	}

	if err := l.db.QueryRow(SQL, code).Scan(&origin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}

		return "", err
	}

	return origin, nil
}

func (l *SQLiteLookup) IncrementClicks(code string) error {
	SQL := `update lookup set clicks = clicks + 1 where code = ?`
	if _, err := l.db.Exec(SQL, code); err != nil {
		return err
	}

	return nil
}
