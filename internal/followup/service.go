package followup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/excel"
	"homemaker-followup/internal/storage"
)

type Service struct {
	WorkbookPath string
	Store        *storage.Store
	Records      []domain.FollowUpRecord
}

func NewService(workbookPath string, store *storage.Store) *Service {
	return &Service{WorkbookPath: workbookPath, Store: store, Records: make([]domain.FollowUpRecord, 0)}
}

func (s *Service) Start() (StartupNotice, error) {
	result, err := LoadFromFile(s.WorkbookPath)
	if err != nil {
		return result.Notice, err
	}
	s.Records = result.Records
	if s.Store != nil {
		if err := s.restoreIndex(context.Background()); err != nil {
			return FormatNotice(err), err
		}
	}
	return result.Notice, nil
}

func (s *Service) restoreIndex(ctx context.Context) error {
	for _, record := range s.Records {
		if err := s.Store.UpsertRecord(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) AddRecord(record domain.FollowUpRecord, actor string) error {
	if err := record.Validate(); err != nil {
		return err
	}
	for _, existing := range s.Records {
		if existing.ID == record.ID {
			return ErrDuplicateRecord
		}
	}
	workbook, err := excel.EnsureWorkbook(s.WorkbookPath)
	if err != nil {
		return err
	}
	defer workbook.Close()
	if err := excel.AppendRecord(workbook.File, record); err != nil {
		return err
	}
	if err := workbook.Save(); err != nil {
		return err
	}
	if s.Store != nil {
		if err := s.Store.UpsertRecord(context.Background(), record); err != nil {
			return err
		}
		if err := s.Store.AddAudit(context.Background(), domain.AuditEntry{RecordID: record.ID, Action: "create", Actor: actor, Detail: record.ClientName}); err != nil {
			return err
		}
	}
	s.Records = append(s.Records, record)
	return nil
}

func (s *Service) UpdateRecord(record domain.FollowUpRecord, actor string) error {
	if err := record.Validate(); err != nil {
		return err
	}
	workbook, err := excel.OpenWorkbook(s.WorkbookPath)
	if err != nil {
		return err
	}
	defer workbook.Close()
	if err := excel.ReplaceRecord(workbook.File, record); err != nil {
		return ErrRecordMissing
	}
	if err := workbook.Save(); err != nil {
		return err
	}
	found := false
	for index := range s.Records {
		if s.Records[index].ID == record.ID {
			s.Records[index] = record
			found = true
		}
	}
	if !found {
		return ErrRecordMissing
	}
	if s.Store != nil {
		if err := s.Store.UpsertRecord(context.Background(), record); err != nil {
			return err
		}
		return s.Store.AddAudit(context.Background(), domain.AuditEntry{RecordID: record.ID, Action: "update", Actor: actor, Detail: record.Status})
	}
	return nil
}

func (s *Service) RemoveRecord(id, actor string) error {
	workbook, err := excel.OpenWorkbook(s.WorkbookPath)
	if err != nil {
		return err
	}
	defer workbook.Close()
	if err := excel.DeleteRecord(workbook.File, id); err != nil {
		return ErrRecordMissing
	}
	if err := workbook.Save(); err != nil {
		return err
	}
	filtered := s.Records[:0]
	removed := false
	for _, record := range s.Records {
		if record.ID == id {
			removed = true
			continue
		}
		filtered = append(filtered, record)
	}
	if !removed {
		return ErrRecordMissing
	}
	s.Records = filtered
	if s.Store != nil {
		if err := s.Store.DeleteRecord(context.Background(), id); err != nil {
			return err
		}
		return s.Store.AddAudit(context.Background(), domain.AuditEntry{RecordID: id, Action: "delete", Actor: actor, Detail: "removed"})
	}
	return nil
}

func (s *Service) EnsurePaths() error {
	if s.WorkbookPath == "" {
		return fmt.Errorf("workbook path is required")
	}
	return os.MkdirAll(filepath.Dir(s.WorkbookPath), 0755)
}
