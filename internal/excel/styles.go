package excel

import (
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
	"homemaker-followup/internal/domain"
)

func StyleWorkbook(file *excelize.File) error {
	if file == nil {
		return ErrInvalidWorkbook
	}
	style, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Color: "FFFFFF"}, Fill: excelize.Fill{Type: "pattern", Color: []string{"123C4A"}, Pattern: 1}, Alignment: &excelize.Alignment{Horizontal: "center"}})
	if err != nil {
		return err
	}
	if err := file.SetCellStyle(FollowUpSheet, "A1", "L1", style); err != nil {
		return err
	}
	file.SetColWidth(FollowUpSheet, "A", "A", 14)
	file.SetColWidth(FollowUpSheet, "B", "E", 18)
	file.SetColWidth(FollowUpSheet, "F", "G", 14)
	file.SetColWidth(FollowUpSheet, "H", "H", 10)
	file.SetColWidth(FollowUpSheet, "I", "L", 24)
	return nil
}

func WriteClientSheet(file *excelize.File, clients []domain.Client) error {
	if file == nil {
		return ErrInvalidWorkbook
	}
	sheet := "客户目录"
	if _, err := file.NewSheet(sheet); err != nil {
		return err
	}
	headers := []string{"客户编号", "客户姓名", "电话", "地址", "首选联系", "启用"}
	for index, header := range headers {
		cell, err := excelize.CoordinatesToCellName(index+1, 1)
		if err != nil {
			return err
		}
		file.SetCellValue(sheet, cell, header)
	}
	for rowIndex, client := range domain.SortClients(clients) {
		values := []interface{}{client.ID, client.Name, client.Phone, client.Address, domain.ContactLabel(client.PreferredChannel), client.Active}
		for column, value := range values {
			cell, err := excelize.CoordinatesToCellName(column+1, rowIndex+2)
			if err != nil {
				return err
			}
			file.SetCellValue(sheet, cell, value)
		}
	}
	return nil
}

func WriteSummarySheet(file *excelize.File, records []domain.FollowUpRecord) error {
	if file == nil {
		return ErrInvalidWorkbook
	}
	sheet := "回访汇总"
	if _, err := file.NewSheet(sheet); err != nil {
		return err
	}
	counts := domain.CountByStatus(records)
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	file.SetCellValue(sheet, "A1", "状态")
	file.SetCellValue(sheet, "B1", "记录数")
	for index, key := range keys {
		file.SetCellValue(sheet, fmt.Sprintf("A%d", index+2), key)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", index+2), counts[key])
	}
	return nil
}

func ValidateSheetOrder(file *excelize.File) error {
	if file == nil || len(file.GetSheetList()) == 0 {
		return ErrInvalidWorkbook
	}
	if !containsSheet(file, FollowUpSheet) {
		return ErrMissingSheet
	}
	return nil
}
