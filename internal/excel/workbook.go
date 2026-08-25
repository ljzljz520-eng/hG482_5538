package excel

import (
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type Workbook struct {
	File *excelize.File
	Path string
}

func NewWorkbook(path string) (*Workbook, error) {
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", FollowUpSheet)
	if err := ensureHeaders(file); err != nil {
		file.Close()
		return nil, err
	}
	return &Workbook{File: file, Path: path}, nil
}

func OpenWorkbook(path string) (*Workbook, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, ErrInvalidWorkbook
	}
	return &Workbook{File: file, Path: path}, nil
}

func (w *Workbook) Save() error {
	if w == nil || w.File == nil {
		return ErrInvalidWorkbook
	}
	if err := os.MkdirAll(filepath.Dir(w.Path), 0755); err != nil {
		return err
	}
	return w.File.SaveAs(w.Path)
}

func (w *Workbook) Close() error {
	if w == nil || w.File == nil {
		return nil
	}
	return w.File.Close()
}

func (w *Workbook) SheetExists() bool {
	if w == nil || w.File == nil {
		return false
	}
	for _, name := range w.File.GetSheetList() {
		if name == FollowUpSheet {
			return true
		}
	}
	return false
}
