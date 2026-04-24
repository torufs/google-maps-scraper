package gmaps

import (
	"fmt"
	"strings"
)

// NewWriter creates a new EntryWriter based on the provided format string and path.
// Supported formats: csv, json, jsonl, excel (or xlsx), sqlite, postgres, tsv.
// The path argument is used for file-based writers; for postgres it is treated as a DSN.
func NewWriter(format, path string) (EntryWriter, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "csv":
		if path == "" {
			return nil, fmt.Errorf("csv writer requires a file path")
		}
		return NewCSVWriter(path)
	case "json":
		if path == "" {
			return nil, fmt.Errorf("json writer requires a file path")
		}
		return NewJSONWriter(path)
	case "jsonl":
		if path == "" {
			return nil, fmt.Errorf("jsonl writer requires a file path")
		}
		return NewJSONLWriter(path)
	case "excel", "xlsx":
		if path == "" {
			return nil, fmt.Errorf("excel writer requires a file path")
		}
		return NewExcelWriter(path)
	case "sqlite":
		if path == "" {
			return nil, fmt.Errorf("sqlite writer requires a file path")
		}
		return NewSQLiteWriter(path)
	case "postgres":
		if path == "" {
			return nil, fmt.Errorf("postgres writer requires a DSN")
		}
		return NewPostgresWriter(path)
	case "tsv":
		if path == "" {
			return nil, fmt.Errorf("tsv writer requires a file path")
		}
		return NewTSVWriter(path)
	default:
		return nil, fmt.Errorf("unsupported writer format: %q", format)
	}
}
