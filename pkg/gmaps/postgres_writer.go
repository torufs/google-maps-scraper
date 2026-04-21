package gmaps

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresWriter writes Entry records to a PostgreSQL database.
type PostgresWriter struct {
	db *sql.DB
}

// NewPostgresWriter opens a connection to the given PostgreSQL DSN, ensures the
// results table exists, and returns a ready-to-use PostgresWriter.
func NewPostgresWriter(dsn string) (*PostgresWriter, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres_writer: open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres_writer: ping: %w", err)
	}
	if err := createPostgresTable(db); err != nil {
		db.Close()
		return nil, err
	}
	return &PostgresWriter{db: db}, nil
}

func createPostgresTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS results (
		id          SERIAL PRIMARY KEY,
		name        TEXT,
		phone       TEXT,
		website     TEXT,
		address     TEXT,
		rating      REAL,
		review_count INTEGER,
		category    TEXT,
		latitude    REAL,
		longitude   REAL
	)`)
	if err != nil {
		return fmt.Errorf("postgres_writer: create table: %w", err)
	}
	return nil
}

// Write inserts a single Entry into the results table. Nil entries are skipped.
func (w *PostgresWriter) Write(e *Entry) error {
	if e == nil {
		return nil
	}
	_, err := w.db.Exec(
		`INSERT INTO results (name, phone, website, address, rating, review_count, category, latitude, longitude)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		e.Name, e.Phone, e.Website, e.Address,
		e.Rating, e.ReviewCount, e.Category,
		e.Latitude, e.Longitude,
	)
	if err != nil {
		return fmt.Errorf("postgres_writer: insert: %w", err)
	}
	return nil
}

// Close releases the underlying database connection.
func (w *PostgresWriter) Close() error {
	return w.db.Close()
}
