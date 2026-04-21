package gmaps_test

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/gosom/google-maps-scraper/pkg/gmaps"
)

func TestCSVWriter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	w, err := gmaps.NewCSVWriter(buf)
	if err != nil {
		t.Fatalf("NewCSVWriter returned error: %v", err)
	}

	entry := &gmaps.Entry{
		Title:   "Test Place",
		Address: "123 Main St",
		Phone:   "+1-555-0100",
		Website: "https://example.com",
		Rating:  4.5,
		Reviews: 100,
	}

	if err := w.Write(entry); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	if err := w.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Fatal("expected non-empty CSV output")
	}

	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV output: %v", err)
	}

	// header row + 1 data row
	if len(records) < 2 {
		t.Fatalf("expected at least 2 rows (header + data), got %d", len(records))
	}
}

func TestCSVWriter_WriteNil(t *testing.T) {
	buf := &bytes.Buffer{}
	w, err := gmaps.NewCSVWriter(buf)
	if err != nil {
		t.Fatalf("NewCSVWriter returned error: %v", err)
	}

	if err := w.Write(nil); err != nil {
		t.Fatalf("Write(nil) returned unexpected error: %v", err)
	}
}

func TestCSVWriter_MultipleEntries(t *testing.T) {
	buf := &bytes.Buffer{}
	w, err := gmaps.NewCSVWriter(buf)
	if err != nil {
		t.Fatalf("NewCSVWriter returned error: %v", err)
	}

	entries := []*gmaps.Entry{
		{Title: "Place A", Address: "1 A St", Rating: 3.0, Reviews: 10},
		{Title: "Place B", Address: "2 B Ave", Rating: 4.0, Reviews: 50},
		{Title: "Place C", Address: "3 C Blvd", Rating: 5.0, Reviews: 200},
	}

	for _, e := range entries {
		if err := w.Write(e); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(buf.String()))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV output: %v", err)
	}

	// header + 3 data rows
	expected := 1 + len(entries)
	if len(records) != expected {
		t.Fatalf("expected %d rows, got %d", expected, len(records))
	}
}
