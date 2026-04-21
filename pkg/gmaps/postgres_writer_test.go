package gmaps

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockPostgresWriter(t *testing.T) (*PostgresWriter, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	// Expect the CREATE TABLE statement issued by createPostgresTable.
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS results`).WillReturnResult(sqlmock.NewResult(0, 0))
	if err := createPostgresTable(db); err != nil {
		t.Fatalf("createPostgresTable: %v", err)
	}
	return &PostgresWriter{db: db}, mock
}

func closeDB(t *testing.T, db *sql.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Errorf("close db: %v", err)
	}
}

func TestPostgresWriter_Write(t *testing.T) {
	w, mock := newMockPostgresWriter(t)
	defer closeDB(t, w.db)

	e := &Entry{
		Name:        "Acme Corp",
		Phone:       "+1-800-000-0000",
		Website:     "https://acme.example",
		Address:     "123 Main St",
		Rating:      4.5,
		ReviewCount: 200,
		Category:    "Services",
		Latitude:    37.7749,
		Longitude:   -122.4194,
	}

	mock.ExpectExec(`INSERT INTO results`).WithArgs(
		e.Name, e.Phone, e.Website, e.Address,
		e.Rating, e.ReviewCount, e.Category,
		e.Latitude, e.Longitude,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := w.Write(e); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestPostgresWriter_WriteNil(t *testing.T) {
	w, mock := newMockPostgresWriter(t)
	defer closeDB(t, w.db)

	if err := w.Write(nil); err != nil {
		t.Fatalf("Write(nil) should be a no-op, got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestPostgresWriter_MultipleEntries(t *testing.T) {
	w, mock := newMockPostgresWriter(t)
	defer closeDB(t, w.db)

	entries := []*Entry{
		{Name: "Place A", Rating: 3.0},
		{Name: "Place B", Rating: 4.2},
	}
	for _, e := range entries {
		mock.ExpectExec(`INSERT INTO results`).WithArgs(
			e.Name, e.Phone, e.Website, e.Address,
			e.Rating, e.ReviewCount, e.Category,
			e.Latitude, e.Longitude,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	for _, e := range entries {
		if err := w.Write(e); err != nil {
			t.Fatalf("Write(%q): %v", e.Name, err)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
