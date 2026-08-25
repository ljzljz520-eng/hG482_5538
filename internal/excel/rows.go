package excel

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"homemaker-followup/internal/domain"
)

func AppendRecord(file *excelize.File, record domain.FollowUpRecord) error {
	if file == nil {
		return ErrInvalidWorkbook
	}
	if err := record.Validate(); err != nil {
		return err
	}
	rows, err := file.GetRows(FollowUpSheet)
	if err != nil {
		return err
	}
	row := len(rows) + 1
	values := recordValues(record)
	for index, value := range values {
		cell, cellErr := excelize.CoordinatesToCellName(index+1, row)
		if cellErr != nil {
			return cellErr
		}
		if setErr := file.SetCellValue(FollowUpSheet, cell, value); setErr != nil {
			return setErr
		}
	}
	return nil
}

func ReplaceRecord(file *excelize.File, record domain.FollowUpRecord) error {
	if file == nil {
		return ErrInvalidWorkbook
	}
	rows, err := file.GetRows(FollowUpSheet)
	if err != nil {
		return err
	}
	for rowIndex, row := range rows[1:] {
		if len(row) > 0 && strings.TrimSpace(row[0]) == record.ID {
			for index, value := range recordValues(record) {
				cell, cellErr := excelize.CoordinatesToCellName(index+1, rowIndex+2)
				if cellErr != nil {
					return cellErr
				}
				if setErr := file.SetCellValue(FollowUpSheet, cell, value); setErr != nil {
					return setErr
				}
			}
			return nil
		}
	}
	return fmt.Errorf("record %s not found", record.ID)
}

func ReadRecords(file *excelize.File) ([]domain.FollowUpRecord, error) {
	if file == nil {
		return nil, ErrInvalidWorkbook
	}
	rows, err := file.GetRows(FollowUpSheet, excelize.Options{RawCellValue: true})
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return []domain.FollowUpRecord{}, nil
	}
	records := make([]domain.FollowUpRecord, 0, len(rows)-1)
	for _, row := range rows[1:] {
		for len(row) < len(RequiredHeaders) {
			row = append(row, "")
		}
		record, parseErr := parseRecord(row)
		if parseErr != nil {
			return nil, parseErr
		}
		records = append(records, record)
	}
	return records, nil
}

func recordValues(record domain.FollowUpRecord) []interface{} {
	return []interface{}{
		record.ID, record.ClientID, record.ClientName, record.ServiceTypeName, record.CaregiverName,
		record.VisitDate.Format(time.DateOnly), record.NextFollowUp.Format(time.DateOnly),
		record.Satisfaction.Score, record.Satisfaction.Comment, record.Satisfaction.Improvement,
		record.Status, record.Notes,
	}
}

func parseRecord(row []string) (domain.FollowUpRecord, error) {
	if len(row) < len(RequiredHeaders) {
		return domain.FollowUpRecord{}, fmt.Errorf("回访记录列数不足: %d", len(row))
	}
	visit, err := time.Parse(time.DateOnly, strings.TrimSpace(row[5]))
	if err != nil {
		return domain.FollowUpRecord{}, fmt.Errorf("服务日期无效: %w", err)
	}
	next, err := time.Parse(time.DateOnly, strings.TrimSpace(row[6]))
	if err != nil {
		return domain.FollowUpRecord{}, fmt.Errorf("下次回访日期无效: %w", err)
	}
	score, err := strconv.Atoi(strings.TrimSpace(row[7]))
	if err != nil {
		return domain.FollowUpRecord{}, fmt.Errorf("满意度无效: %w", err)
	}
	return domain.FollowUpRecord{
		ID: row[0], ClientID: row[1], ClientName: row[2], ServiceTypeName: row[3], CaregiverName: row[4],
		VisitDate: visit, NextFollowUp: next,
		Satisfaction: domain.Satisfaction{Score: score, Comment: row[8], Improvement: row[9]},
		Status:       row[10], Notes: row[11], UpdatedAt: visit.UTC(),
	}, nil
}
