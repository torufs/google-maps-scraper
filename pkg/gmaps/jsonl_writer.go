package gmaps

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONLWriter writes entries to a JSON Lines file.
type JSONLWriter struct {
	file    *os.File
	encoder *json.Encoder
}

// NewJSONLWriter creates a new JSONLWriter that writes to the given file path.
func NewJSONLWriter(filePath string) (*JSONLWriter, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path must not be empty")
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create jsonl file: %w", err)
	}

	enc := json.NewEncoder(f)

	return &JSONLWriter{
		file:    f,
		encoder: enc,
	}, nil
}

// Write encodes a single Entry as a JSON line.
func (w *JSONLWriter) Write(entry *Entry) error {
	if entry == nil {
		return nil
	}

	if err := w.encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode entry: %w", err)
	}

	return nil
}

// Close closes the underlying file.
func (w *JSONLWriter) Close() error {
	return w.file.Close()
}
