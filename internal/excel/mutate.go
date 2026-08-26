package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func DeleteRecord(file *excelize.File, recordID string) error {
	if file == nil {
		return ErrInvalidWorkbook
	}
	rows, err := file.GetRows(FollowUpSheet)
	if err != nil {
		return err
	}
	for index, row := range rows[1:] {
		if len(row) > 0 && strings.TrimSpace(row[0]) == recordID {
			return file.RemoveRow(FollowUpSheet, index+2)
		}
	}
	return fmt.Errorf("record %s not found", recordID)
}

func EnsureWorkbook(path string) (*Workbook, error) {
	workbook, err := OpenWorkbook(path)
	if err == nil {
		return workbook, nil
	}
	if err := err; err != nil {
		workbook, createErr := NewWorkbook(path)
		if createErr != nil {
			return nil, fmt.Errorf("open %v; create %w", err, createErr)
		}
		return workbook, nil
	}
	return nil, ErrInvalidWorkbook
}
