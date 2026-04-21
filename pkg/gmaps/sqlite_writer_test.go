package gmaps

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tmpDB(t *testing.T) (string, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "gmaps-test-*.db")
	require.NoError(t, err)
	f.Close()
	return f.Name(), func() { os.Remove(f.Name()) }
}

func TestSQLiteWriter_Write(t *testing.T) {
	path, cleanup := tmpDB(t)
	defer cleanup()

	w, err := NewSQLiteWriter(path)
	require.NoError(t, err)
	defer w.Close()

	e := &Entry{
		Name:        "Test Place",
		Address:     "123 Main St",
		Phone:       "+1-555-0100",
		Website:     "https://example.com",
		Rating:      4.5,
		ReviewCount: 42,
		Category:    "Restaurant",
		Latitude:    37.7749,
		Longitude:   -122.4194,
	}

	err = w.Write(e)
	require.NoError(t, err)

	var name string
	row := w.db.QueryRow("SELECT name FROM entries WHERE id = 1")
	require.NoError(t, row.Scan(&name))
	assert.Equal(t, "Test Place", name)
}

func TestSQLiteWriter_WriteNil(t *testing.T) {
	path, cleanup := tmpDB(t)
	defer cleanup()

	w, err := NewSQLiteWriter(path)
	require.NoError(t, err)
	defer w.Close()

	assert.NoError(t, w.Write(nil))
}

func TestSQLiteWriter_MultipleEntries(t *testing.T) {
	path, cleanup := tmpDB(t)
	defer cleanup()

	w, err := NewSQLiteWriter(path)
	require.NoError(t, err)
	defer w.Close()

	entries := []*Entry{
		{Name: "Place A", Rating: 3.0},
		{Name: "Place B", Rating: 4.0},
		{Name: "Place C", Rating: 5.0},
	}

	for _, e := range entries {
		require.NoError(t, w.Write(e))
	}

	var count int
	require.NoError(t, w.db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&count))
	assert.Equal(t, 3, count)
}
