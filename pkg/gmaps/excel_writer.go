package gmaps

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

// ExcelWriter writes Entry records to an Excel (.xlsx) file.
type ExcelWriter struct {
	w       io.Writer
	f       *excelize.File
	sheet   string
	row     int
	headers []string
}

// NewExcelWriter creates a new ExcelWriter that writes to w.
func NewExcelWriter(w io.Writer) (*ExcelWriter, error) {
	f := excelize.NewFile()
	sheet := "Results"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"Title", "Category", "Address", "OpenHours",
		"WebSite", "Phone", "PlusCode", "ReviewCount",
		"ReviewRating", "Latitude", "Longitude",
	}

	ew := &ExcelWriter{
		w:       w,
		f:       f,
		sheet:   sheet,
		row:     1,
		headers: headers,
	}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	ew.row = 2
	return ew, nil
}

// Write appends e to the Excel sheet. Nil entries are skipped.
func (ew *ExcelWriter) Write(e *Entry) error {
	if e == nil {
		return nil
	}

	values := []interface{}{
		e.Title,
		e.Category,
		e.Address,
		e.OpenHours,
		e.WebSite,
		e.Phone,
		e.PlusCode,
		e.ReviewCount,
		e.ReviewRating,
		e.Latitude,
		e.Longitude,
	}

	for col, v := range values {
		cell, _ := excelize.CoordinatesToCellName(col+1, ew.row)
		ew.f.SetCellValue(ew.sheet, cell, v)
	}
	ew.row++
	return nil
}

// Flush writes the Excel file to the underlying writer.
func (ew *ExcelWriter) Flush() error {
	_, err := ew.f.WriteTo(ew.w)
	if err != nil {
		return fmt.Errorf("excel_writer: flush: %w", err)
	}
	return nil
}
