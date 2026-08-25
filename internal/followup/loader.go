package followup

import (
	"fmt"

	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/excel"
)

type LoadResult struct {
	Records []domain.FollowUpRecord
	Notice  StartupNotice
}

func LoadFromFile(path string) (LoadResult, error) {
	workbook, err := excel.OpenWorkbook(path)
	if err != nil {
		return LoadResult{Notice: FormatNotice(err)}, err
	}
	defer workbook.Close()
	if err := validateWorkbookHeaders(workbook); err != nil {
		return LoadResult{Notice: FormatNotice(err)}, err
	}
	records, err := excel.ReadRecords(workbook.File)
	if err != nil {
		return LoadResult{Notice: FormatNotice(err)}, err
	}
	return LoadResult{Records: records, Notice: FormatNotice(nil)}, nil
}

func validateWorkbookHeaders(workbook *excel.Workbook) error {
	if workbook == nil || !workbook.SheetExists() {
		return fmt.Errorf("%w: %s", excel.ErrMissingSheet, excel.FollowUpSheet)
	}
	_, err := excel.ReadHeaders(workbook.File)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrStartupFormat, err)
	}
	return nil
}

func PreviewFile(path string) (StartupNotice, error) {
	workbook, err := excel.OpenWorkbook(path)
	if err != nil {
		return FormatNotice(err), err
	}
	defer workbook.Close()
	if err := validateWorkbookHeaders(workbook); err != nil {
		return FormatNotice(err), err
	}
	return FormatNotice(nil), nil
}
