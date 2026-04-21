package gmaps

import (
	"fmt"
	"io"
	"strings"
)

// OutputFormat represents a supported output format.
type OutputFormat string

const (
	FormatCSV    OutputFormat = "csv"
	FormatJSON   OutputFormat = "json"
	FormatExcel  OutputFormat = "excel"
	FormatSQLite OutputFormat = "sqlite"
)

// WriterConfig holds configuration needed to create a writer.
type WriterConfig struct {
	Format OutputFormat
	// W is used for CSV and JSON writers.
	W io.Writer
	// FilePath is used for Excel and SQLite writers.
	FilePath string
}

// NewWriter creates the appropriate EntryWriter based on the config.
func NewWriter(cfg WriterConfig) (EntryWriter, error) {
	switch OutputFormat(strings.ToLower(string(cfg.Format))) {
	case FormatCSV:
		return NewCSVWriter(cfg.W), nil
	case FormatJSON:
		return NewJSONWriter(cfg.W), nil
	case FormatExcel:
		if cfg.FilePath == "" {
			return nil, fmt.Errorf("excel writer requires a file path")
		}
		return NewExcelWriter(cfg.FilePath)
	case FormatSQLite:
		if cfg.FilePath == "" {
			return nil, fmt.Errorf("sqlite writer requires a file path")
		}
		return NewSQLiteWriter(cfg.FilePath)
	default:
		return nil, fmt.Errorf("unsupported output format: %q", cfg.Format)
	}
}
