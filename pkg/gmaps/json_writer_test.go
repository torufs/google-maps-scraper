package gmaps

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	w := NewJSONWriter(&buf)

	e := &Entry{
		Title:   "Test Place",
		Website: "https://example.com",
		Phone:   "+1-555-0100",
	}

	if err := w.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.Count() != 1 {
		t.Fatalf("expected count 1, got %d", w.Count())
	}

	var decoded Entry
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if decoded.Title != e.Title {
		t.Errorf("title mismatch: want %q, got %q", e.Title, decoded.Title)
	}
	if decoded.Website != e.Website {
		t.Errorf("website mismatch: want %q, got %q", e.Website, decoded.Website)
	}
	// also verify phone is preserved correctly
	if decoded.Phone != e.Phone {
		t.Errorf("phone mismatch: want %q, got %q", e.Phone, decoded.Phone)
	}
}

func TestJSONWriter_WriteNil(t *testing.T) {
	var buf bytes.Buffer
	w := NewJSONWriter(&buf)

	if err := w.Write(nil); err != nil {
		t.Fatalf("unexpected error on nil write: %v", err)
	}

	if w.Count() != 0 {
		t.Fatalf("expected count 0 after nil write, got %d", w.Count())
	}

	// nil write should produce no output
	if buf.Len() != 0 {
		t.Errorf("expected empty buffer after nil write, got %d bytes", buf.Len())
	}
}

func TestJSONWriter_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	w := NewJSONWriter(&buf)

	entries := []*Entry{
		{Title: "Place A"},
		{Title: "Place B"},
		{Title: "Place C"},
	}

	for _, e := range entries {
		if err := w.Write(e); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if w.Count() != len(entries) {
		t.Fatalf("expected count %d, got %d", len(entries), w.Count())
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != len(entries) {
		t.Fatalf("expected %d lines, got %d", len(entries), len(lines))
	}
}
