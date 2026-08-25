package excel

import (
	"path/filepath"
	"testing"
)

func TestWorkbookRoundTripHeaders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "roundtrip.xlsx")
	workbook, err := NewWorkbook(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := workbook.Save(); err != nil {
		t.Fatal(err)
	}
	workbook.Close()
	opened, err := OpenWorkbook(path)
	if err != nil {
		t.Fatal(err)
	}
	defer opened.Close()
	headers, err := ReadHeaders(opened.File)
	if err != nil {
		t.Fatal(err)
	}
	if len(headers) != len(RequiredHeaders) {
		t.Fatalf("got %d headers", len(headers))
	}
}
