package gmaps

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteWriter writes Entry records to a SQLite database.
type SQLiteWriter struct {
	db *sql.DB
}

// NewSQLiteWriter opens (or creates) a SQLite database at the given path and
// ensures the entries table exists.
// Note: WAL mode is enabled for better concurrent read performance.
func NewSQLiteWriter(path string) (*SQLiteWriter, error) {
	// Use cache=shared to allow multiple connections to share the same in-memory cache.
	db, err := sql.Open("sqlite3", path+"?cache=shared")
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	// Enable WAL mode for improved write performance and concurrent reads.
	if _, err = db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable WAL mode: %w", err)
	}

	// Increase the cache size to 8000 pages (~32 MB) for better read performance.
	if _, err = db.Exec(`PRAGMA cache_size=-8000;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("set cache size: %w", err)
	}

	const schema = `CREATE TABLE IF NOT EXISTS entries (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT,
		address     TEXT,
		phone       TEXT,
		website     TEXT,
		rating      REAL,
		review_count INTEGER,
		category    TEXT,
		latitude    REAL,
		longitude   REAL
	);`

	if _, err = db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &SQLiteWriter{db: db}, nil
}

// Write inserts a single Entry into the database.
func (w *SQLiteWriter) Write(e *Entry) error {
	if e == nil {
		return nil
	}

	const query = `INSERT INTO entries
		(name, address, phone, website, rating, review_count, category, latitude, longitude)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := w.db.Exec(query,
		e.Name, e.Address, e.Phone, e.Website,
		e.Rating, e.ReviewCount, e.Category,
		e.Latitude, e.Longitude,
	)
	if err != nil {
		return fmt.Errorf("insert entry: %w", err)
	}
	return nil
}

// Close releases the underlying database connection.
func (w *SQLiteWriter) Close() error {
	return w.db.Close()
}
