package gmaps

import (
	"fmt"
	"strings"
)

// NewWriter creates a new EntryWriter based on the provided format and path.
// Supported formats: csv, json, jsonl, excel (xlsx), sqlite, postgres, tsv.
func NewWriter(format, path, dsn string) (EntryWriter, error) {
	switch strings.ToLower(format) {
	case "csv":
		return NewCSVWriter(path)
	case "json":
		return NewJSONWriter(path)
	case "jsonl":
		return NewJSONLWriter(path)
	case "excel", "xlsx":
		if path == "" {
			return nil, fmt.Errorf("writer factory: excel format requires a file path")
		}
		return NewExcelWriter(path)
	case "sqlite":
		return NewSQLiteWriter(path)
	case "postgres":
		return NewPostgresWriter(dsn)
	case "tsv":
		return NewTSVWriter(path)
	default:
		return nil, fmt.Errorf("writer factory: unsupported format %q", format)
	}
}
