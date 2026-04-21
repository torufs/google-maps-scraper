package gmaps

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWriter_CSV(t *testing.T) {
	w, err := NewWriter(WriterConfig{Format: FormatCSV, W: &bytes.Buffer{}})
	require.NoError(t, err)
	assert.NotNil(t, w)
}

func TestNewWriter_JSON(t *testing.T) {
	w, err := NewWriter(WriterConfig{Format: FormatJSON, W: &bytes.Buffer{}})
	require.NoError(t, err)
	assert.NotNil(t, w)
}

func TestNewWriter_Excel(t *testing.T) {
	f, err := os.CreateTemp("", "factory-test-*.xlsx")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	w, err := NewWriter(WriterConfig{Format: FormatExcel, FilePath: f.Name()})
	require.NoError(t, err)
	assert.NotNil(t, w)
}

func TestNewWriter_SQLite(t *testing.T) {
	f, err := os.CreateTemp("", "factory-test-*.db")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	w, err := NewWriter(WriterConfig{Format: FormatSQLite, FilePath: f.Name()})
	require.NoError(t, err)
	assert.NotNil(t, w)
}

func TestNewWriter_ExcelMissingPath(t *testing.T) {
	_, err := NewWriter(WriterConfig{Format: FormatExcel})
	assert.Error(t, err)
}

func TestNewWriter_SQLiteMissingPath(t *testing.T) {
	_, err := NewWriter(WriterConfig{Format: FormatSQLite})
	assert.Error(t, err)
}

func TestNewWriter_Unknown(t *testing.T) {
	_, err := NewWriter(WriterConfig{Format: "parquet"})
	assert.Error(t, err)
}
