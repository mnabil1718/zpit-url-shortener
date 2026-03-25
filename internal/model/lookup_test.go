package model

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/mnabil1718/zp.it/internal/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test in-memory DB setup
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS lookup (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			origin     TEXT NOT NULL,
			code       TEXT NOT NULL UNIQUE,
			clicks     INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	t.Cleanup(func() { db.Close() })
	return db
}

// ===== INSERT =====

func TestLookupInsert_Success(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)
	mc.On("Set", mock.Anything, "abc", "https://example.com", 300*time.Second).Return(nil)

	lu := NewSQliteLookup(db, mc)
	err := lu.Insert("https://example.com", "abc")
	require.NoError(t, err)

	var origin string
	db.QueryRow("SELECT origin FROM lookup WHERE code = ?", "abc").Scan(&origin)
	assert.Equal(t, "https://example.com", origin)

	mc.AssertExpectations(t)
}

func TestLookupInsert_DuplicateCode(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)
	mc.On("Set", mock.Anything, "abc", "https://example.com", 300*time.Second).Return(nil)

	lu := NewSQliteLookup(db, mc)
	require.NoError(t, lu.Insert("https://example.com", "abc"))

	err := lu.Insert("https://other.com", "abc")
	assert.ErrorIs(t, err, ErrAlreadyExists)

	mc.AssertExpectations(t) // Set called once, not twice
}

func TestLookupInsert_CacheSetFails(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)
	mc.On("Set", mock.Anything, "abc", "https://example.com", 300*time.Second).Return(errors.New("redis service is not responding"))

	lu := NewSQliteLookup(db, mc)
	err := lu.Insert("https://example.com", "abc")
	assert.Error(t, err)

	mc.AssertExpectations(t)
}

// ===== GET BY CODE =====

func TestLookupGetByCode_Found(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)
	mc.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	lu := NewSQliteLookup(db, mc)
	require.NoError(t, lu.Insert("https://example.com", "abc"))

	lkp, err := lu.GetByCode("abc")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", lkp.Origin)
	assert.Equal(t, "abc", lkp.Code)
	assert.Equal(t, 0, lkp.Clicks)
}

func TestLookupGetByCode_NotFound(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)

	lu := NewSQliteLookup(db, mc)
	_, err := lu.GetByCode("nope")
	assert.ErrorIs(t, err, ErrNotFound)
}

// ===== GET ORIGIN BY CODE =====
func TestLookupGetOriginByCode_Success(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)
	mc.On("Get", mock.Anything, "abc").Return("https://example.com", nil)

	lu := NewSQliteLookup(db, mc)
	origin, err := lu.GetOriginByCode("abc")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", origin)

	mc.AssertExpectations(t)
}

func TestLookupGetOriginByCode_NotFound(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)
	mc.On("Get", mock.Anything, "ghost").Return("", cache.ErrCacheMiss)

	lu := NewSQliteLookup(db, mc)
	_, err := lu.GetOriginByCode("ghost")
	assert.ErrorIs(t, err, ErrNotFound)
}

// ===== INCREMENT CLICKS =====

func TestLookupIncrementClicks_Success(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)
	mc.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	lu := NewSQliteLookup(db, mc)
	require.NoError(t, lu.Insert("https://example.com", "abc"))

	require.NoError(t, lu.IncrementClicks("abc"))
	require.NoError(t, lu.IncrementClicks("abc"))

	lkp, err := lu.GetByCode("abc")
	require.NoError(t, err)
	assert.Equal(t, 2, lkp.Clicks)
}

func TestIncrementClicks_UnknownCode(t *testing.T) {
	db := newTestDB(t)
	mc := new(cache.MockCache)

	lu := NewSQliteLookup(db, mc)
	// SQLite UPDATE on a missing row is not an error, it just affects 0 rows.
	err := lu.IncrementClicks("ghost")
	assert.NoError(t, err)
}
