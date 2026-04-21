package gmaps

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// CSVHeader returns the ordered column headers for CSV output.
var CSVHeader = []string{
	"title", "category", "address", "open_hours",
	"website", "phone", "plus_code",
	"review_count", "rating",
	"latitude", "longitude", "cid",
}

// CSVWriter wraps csv.Writer and writes Entry records.
type CSVWriter struct {
	w *csv.Writer
}

// NewCSVWriter creates a CSVWriter and writes the header row.
func NewCSVWriter(out io.Writer) (*CSVWriter, error) {
	w := csv.NewWriter(out)
	if err := w.Write(CSVHeader); err != nil {
		return nil, fmt.Errorf("csv: write header: %w", err)
	}
	return &CSVWriter{w: w}, nil
}

// Write serialises a single Entry as a CSV row.
func (c *CSVWriter) Write(e *Entry) error {
	row := []string{
		e.Title,
		e.Category,
		e.Address,
		e.OpenHours,
		e.Website,
		e.Phone,
		e.PlusCode,
		strconv.Itoa(e.ReviewCount),
		strconv.FormatFloat(e.Rating, 'f', 1, 64),
		strconv.FormatFloat(e.Latitude, 'f', 6, 64),
		strconv.FormatFloat(e.Longitude, 'f', 6, 64),
		e.CID,
	}
	if err := c.w.Write(row); err != nil {
		return fmt.Errorf("csv: write row: %w", err)
	}
	return nil
}

// Flush flushes any buffered data to the underlying writer.
func (c *CSVWriter) Flush() error {
	c.w.Flush()
	return c.w.Error()
}
