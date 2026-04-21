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
func NewSQLiteWriter(path string) (*SQLiteWriter, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
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
