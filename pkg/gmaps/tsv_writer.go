package gmaps

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// TSVWriter writes entries to a tab-separated values file.
type TSVWriter struct {
	writer *csv.Writer
	file   *os.File
	header bool
}

// NewTSVWriter creates a new TSVWriter that writes to the given path.
// If path is empty, it returns an error.
func NewTSVWriter(path string) (*TSVWriter, error) {
	if path == "" {
		return nil, fmt.Errorf("tsv writer: path is required")
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("tsv writer: create file: %w", err)
	}

	w := csv.NewWriter(f)
	w.Comma = '\t'

	return &TSVWriter{
		writer: w,
		file:   f,
	}, nil
}

// Write writes a single Entry to the TSV file.
// If entry is nil, it is a no-op.
func (t *TSVWriter) Write(entry *Entry) error {
	if entry == nil {
		return nil
	}

	if !t.header {
		if err := t.writer.Write(csvHeader()); err != nil {
			return fmt.Errorf("tsv writer: write header: %w", err)
		}
		t.header = true
	}

	row := entryToCSVRow(entry)
	if err := t.writer.Write(row); err != nil {
		return fmt.Errorf("tsv writer: write row: %w", err)
	}

	t.writer.Flush()
	return t.writer.Error()
}

// Close flushes and closes the underlying file.
func (t *TSVWriter) Close() error {
	t.writer.Flush()
	if err := t.writer.Error(); err != nil {
		_ = t.file.Close()
		return err
	}
	return t.file.Close()
}

// compile-time check
var _ io.Closer = (*TSVWriter)(nil)
