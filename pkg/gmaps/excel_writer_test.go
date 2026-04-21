package gmaps

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestExcelWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	ew, err := NewExcelWriter(&buf)
	if err != nil {
		t.Fatalf("NewExcelWriter: %v", err)
	}

	e := &Entry{
		Title:        "Test Place",
		Category:     "Restaurant",
		Address:      "123 Main St",
		Phone:        "+1-555-0100",
		WebSite:      "https://example.com",
		ReviewCount:  42,
		ReviewRating: 4.5,
		Latitude:     37.7749,
		Longitude:    -122.4194,
	}

	if err := ew.Write(e); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := ew.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	f, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatalf("OpenReader: %v", err)
	}

	val, _ := f.GetCellValue("Results", "A1")
	if val != "Title" {
		t.Errorf("expected header 'Title', got %q", val)
	}

	val, _ = f.GetCellValue("Results", "A2")
	if val != "Test Place" {
		t.Errorf("expected 'Test Place', got %q", val)
	}
}

func TestExcelWriter_WriteNil(t *testing.T) {
	var buf bytes.Buffer
	ew, err := NewExcelWriter(&buf)
	if err != nil {
		t.Fatalf("NewExcelWriter: %v", err)
	}
	if err := ew.Write(nil); err != nil {
		t.Errorf("Write(nil) should not return error, got: %v", err)
	}
}

func TestExcelWriter_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	ew, err := NewExcelWriter(&buf)
	if err != nil {
		t.Fatalf("NewExcelWriter: %v", err)
	}

	entries := []*Entry{
		{Title: "Place A", ReviewRating: 4.0},
		{Title: "Place B", ReviewRating: 3.5},
	}
	for _, e := range entries {
		if err := ew.Write(e); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}
	if err := ew.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	f, _ := excelize.OpenReader(&buf)
	rows, _ := f.GetRows("Results")
	// 1 header + 2 data rows
	if len(rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(rows))
	}
}
