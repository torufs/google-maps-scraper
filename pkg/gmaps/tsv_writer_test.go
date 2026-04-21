package gmaps

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"
)

func TestTSVWriter_Write(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "*.tsv")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w, err := NewTSVWriter(tmp.Name())
	if err != nil {
		t.Fatalf("NewTSVWriter: %v", err)
	}
	defer w.Close()

	entry := &Entry{Title: "Test Place", Address: "123 Main St"}
	if err := w.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}

	w.Close()

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.Contains(content, "Test Place") {
		t.Errorf("expected 'Test Place' in output, got: %s", content)
	}

	// Verify tab-separated
	if !strings.Contains(content, "\t") {
		t.Error("expected tab separator in output")
	}
}

func TestTSVWriter_WriteNil(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "*.tsv")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w, err := NewTSVWriter(tmp.Name())
	if err != nil {
		t.Fatalf("NewTSVWriter: %v", err)
	}
	defer w.Close()

	if err := w.Write(nil); err != nil {
		t.Fatalf("Write(nil) should not error, got: %v", err)
	}
}

func TestTSVWriter_MultipleEntries(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "*.tsv")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w, err := NewTSVWriter(tmp.Name())
	if err != nil {
		t.Fatalf("NewTSVWriter: %v", err)
	}

	entries := []*Entry{
		{Title: "Place A", Address: "Addr A"},
		{Title: "Place B", Address: "Addr B"},
		{Title: "Place C", Address: "Addr C"},
	}
	for _, e := range entries {
		if err := w.Write(e); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}
	w.Close()

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(data)))
	r.Comma = '\t'
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	// header + 3 data rows
	if len(records) != 4 {
		t.Errorf("expected 4 records (header+3), got %d", len(records))
	}
}

func TestNewTSVWriter_EmptyPath(t *testing.T) {
	_, err := NewTSVWriter("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}
