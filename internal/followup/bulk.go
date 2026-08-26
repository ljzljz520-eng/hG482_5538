package followup

import (
	"context"
	"fmt"
	"sort"
	"time"

	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/excel"
)

type ImportSummary struct {
	Accepted int
	Rejected int
	Errors   []string
}

func (s *Service) ImportRecords(records []domain.FollowUpRecord, actor string) ImportSummary {
	summary := ImportSummary{Errors: make([]string, 0)}
	for _, record := range records {
		if err := s.AddRecord(record, actor); err != nil {
			summary.Rejected++
			summary.Errors = append(summary.Errors, fmt.Sprintf("%s: %v", record.ID, err))
			continue
		}
		summary.Accepted++
	}
	return summary
}

func (s *Service) RecalculateStatuses(today time.Time) int {
	changed := 0
	for index, record := range s.Records {
		status := domain.StatusAt(record, today)
		if record.Status != status {
			s.Records[index].Status = status
			changed++
		}
	}
	return changed
}

func (s *Service) DueRecords(today time.Time, days int) []domain.FollowUpRecord {
	end := today.AddDate(0, 0, days)
	result := make([]domain.FollowUpRecord, 0)
	for _, record := range s.Records {
		if record.NextFollowUp.IsZero() {
			continue
		}
		if !record.NextFollowUp.After(end) && !record.NextFollowUp.Before(today) {
			result = append(result, record)
		}
	}
	sort.SliceStable(result, func(left, right int) bool { return result[left].NextFollowUp.Before(result[right].NextFollowUp) })
	return result
}

func (s *Service) DeleteClientRecords(clientID, actor string) (int, error) {
	ids := make([]string, 0)
	for _, record := range s.Records {
		if record.ClientID == clientID {
			ids = append(ids, record.ID)
		}
	}
	for _, id := range ids {
		if err := s.RemoveRecord(id, actor); err != nil {
			return len(ids), err
		}
	}
	return len(ids), nil
}

func (s *Service) ExportCurrent(path string, clients []domain.Client) error {
	workbook, err := excel.NewWorkbook(path)
	if err != nil {
		return err
	}
	defer workbook.Close()
	for _, record := range s.Records {
		if err := excel.AppendRecord(workbook.File, record); err != nil {
			return err
		}
	}
	if err := excel.WriteClientSheet(workbook.File, clients); err != nil {
		return err
	}
	if err := excel.WriteSummarySheet(workbook.File, s.Records); err != nil {
		return err
	}
	if err := excel.StyleWorkbook(workbook.File); err != nil {
		return err
	}
	return workbook.Save()
}

func (s *Service) SyncIndex(ctx context.Context) error {
	if s.Store == nil {
		return nil
	}
	for _, record := range s.Records {
		if err := s.Store.UpsertRecord(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func ValidateImport(records []domain.FollowUpRecord) error {
	seen := make(map[string]bool, len(records))
	for _, record := range records {
		if seen[record.ID] {
			return fmt.Errorf("duplicate imported id %s", record.ID)
		}
		if err := record.Validate(); err != nil {
			return err
		}
		seen[record.ID] = true
	}
	return nil
}
