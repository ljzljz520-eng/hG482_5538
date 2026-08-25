package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ensureHeaders(file *excelize.File) error {
	if file == nil {
		return ErrInvalidWorkbook
	}
	for index, header := range RequiredHeaders {
		cell, err := excelize.CoordinatesToCellName(index+1, 1)
		if err != nil {
			return err
		}
		if err := file.SetCellValue(FollowUpSheet, cell, header); err != nil {
			return err
		}
	}
	return nil
}

func ReadHeaders(file *excelize.File) ([]string, error) {
	if file == nil {
		return nil, ErrInvalidWorkbook
	}
	if !containsSheet(file, FollowUpSheet) {
		return nil, ErrMissingSheet
	}
	rows, err := file.GetRows(FollowUpSheet, excelize.Options{RawCellValue: true})
	if err != nil {
		return nil, fmt.Errorf("read headers: %w", err)
	}
	if len(rows) == 0 {
		return nil, ErrMissingHeaders
	}
	got := make([]string, 0, len(rows[0]))
	for _, value := range rows[0] {
		got = append(got, strings.TrimSpace(value))
	}
	if !hasRequiredHeaders(got) {
		return got, ErrMissingHeaders
	}
	return got, nil
}

func containsSheet(file *excelize.File, wanted string) bool {
	for _, name := range file.GetSheetList() {
		if name == wanted {
			return true
		}
	}
	return false
}

func hasRequiredHeaders(got []string) bool {
	seen := make(map[string]bool, len(got))
	for _, value := range got {
		seen[value] = true
	}
	for _, required := range RequiredHeaders {
		if !seen[required] {
			return false
		}
	}
	return true
}
