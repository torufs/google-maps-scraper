package gmaps

import (
	"fmt"
	"strings"
)

// NewWriter creates an EntryWriter based on the provided format and output path.
// Supported formats: csv, json, jsonl, excel, sqlite, postgres.
func NewWriter(format, outputPath, dsn string) (EntryWriter, error) {
	switch strings.ToLower(format) {
	case "csv":
		if outputPath == "" {
			return nil, fmt.Errorf("output path is required for csv format")
		}
		return NewCSVWriter(outputPath)
	case "json":
		if outputPath == "" {
			return nil, fmt.Errorf("output path is required for json format")
		}
		return NewJSONWriter(outputPath)
	case "jsonl":
		if outputPath == "" {
			return nil, fmt.Errorf("output path is required for jsonl format")
		}
		return NewJSONLWriter(outputPath)
	case "excel":
		if outputPath == "" {
			return nil, fmt.Errorf("output path is required for excel format")
		}
		return NewExcelWriter(outputPath)
	case "sqlite":
		if outputPath == "" {
			return nil, fmt.Errorf("output path is required for sqlite format")
		}
		return NewSQLiteWriter(outputPath)
	case "postgres":
		if dsn == "" {
			return nil, fmt.Errorf("dsn is required for postgres format")
		}
		return NewPostgresWriter(dsn)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
