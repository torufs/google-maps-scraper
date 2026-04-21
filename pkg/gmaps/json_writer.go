package gmaps

import (
	"encoding/json"
	"io"
	"sync"
)

// JSONWriter writes Entry records as newline-delimited JSON (NDJSON).
type JSONWriter struct {
	mu      sync.Mutex
	encoder *json.Encoder
	count   int
}

// NewJSONWriter creates a new JSONWriter that writes to w.
func NewJSONWriter(w io.Writer) *JSONWriter {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &JSONWriter{encoder: enc}
}

// Write serialises e as a JSON object followed by a newline.
// It is safe for concurrent use.
func (j *JSONWriter) Write(e *Entry) error {
	if e == nil {
		return nil
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	if err := j.encoder.Encode(e); err != nil {
		return err
	}

	j.count++
	return nil
}

// Count returns the number of entries written so far.
func (j *JSONWriter) Count() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.count
}
