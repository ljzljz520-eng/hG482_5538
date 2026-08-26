package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"homemaker-followup/internal/domain"
)

func (s *Store) UpsertRecord(ctx context.Context, record domain.FollowUpRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("storage is closed")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO follow_up_records
 (id, client_id, client_name, service_type_id, service_type_name, caregiver_id, caregiver_name, visit_date, next_follow_up, score, comment, improvement, status, notes, updated_at)
 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
 ON CONFLICT(id) DO UPDATE SET client_id=excluded.client_id, client_name=excluded.client_name, service_type_id=excluded.service_type_id,
 service_type_name=excluded.service_type_name, caregiver_id=excluded.caregiver_id, caregiver_name=excluded.caregiver_name,
 visit_date=excluded.visit_date, next_follow_up=excluded.next_follow_up, score=excluded.score, comment=excluded.comment,
 improvement=excluded.improvement, status=excluded.status, notes=excluded.notes, updated_at=excluded.updated_at`,
		record.ID, record.ClientID, record.ClientName, record.ServiceTypeID, record.ServiceTypeName, record.CaregiverID,
		record.CaregiverName, record.VisitDate.Format(time.RFC3339), record.NextFollowUp.Format(time.RFC3339),
		record.Satisfaction.Score, record.Satisfaction.Comment, record.Satisfaction.Improvement,
		record.Status, record.Notes, record.UpdatedAt.Format(time.RFC3339))
	return err
}

func (s *Store) GetRecord(ctx context.Context, id string) (domain.FollowUpRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return domain.FollowUpRecord{}, fmt.Errorf("storage is closed")
	}
	row := s.db.QueryRowContext(ctx, `SELECT id, client_id, client_name, service_type_id, service_type_name, caregiver_id, caregiver_name,
 visit_date, next_follow_up, score, comment, improvement, status, notes, updated_at FROM follow_up_records WHERE id=?`, id)
	return scanRecord(row)
}

func (s *Store) ListRecords(ctx context.Context) ([]domain.FollowUpRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, fmt.Errorf("storage is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, client_id, client_name, service_type_id, service_type_name, caregiver_id, caregiver_name,
 visit_date, next_follow_up, score, comment, improvement, status, notes, updated_at FROM follow_up_records ORDER BY next_follow_up, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.FollowUpRecord, 0)
	for rows.Next() {
		record, scanErr := scanRecord(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, record)
	}
	return result, rows.Err()
}

type rowScanner interface{ Scan(...any) error }

func scanRecord(row rowScanner) (domain.FollowUpRecord, error) {
	var record domain.FollowUpRecord
	var visit, next, updated string
	err := row.Scan(&record.ID, &record.ClientID, &record.ClientName, &record.ServiceTypeID, &record.ServiceTypeName,
		&record.CaregiverID, &record.CaregiverName, &visit, &next, &record.Satisfaction.Score,
		&record.Satisfaction.Comment, &record.Satisfaction.Improvement, &record.Status, &record.Notes, &updated)
	if err != nil {
		return domain.FollowUpRecord{}, err
	}
	var parseErr error
	record.VisitDate, parseErr = time.Parse(time.RFC3339, strings.TrimSpace(visit))
	if parseErr != nil {
		return domain.FollowUpRecord{}, parseErr
	}
	record.NextFollowUp, parseErr = time.Parse(time.RFC3339, strings.TrimSpace(next))
	if parseErr != nil {
		return domain.FollowUpRecord{}, parseErr
	}
	record.UpdatedAt, parseErr = time.Parse(time.RFC3339, strings.TrimSpace(updated))
	return record, parseErr
}

func (s *Store) DeleteRecord(ctx context.Context, id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("storage is closed")
	}
	_, err := s.db.ExecContext(ctx, `DELETE FROM follow_up_records WHERE id=?`, id)
	return err
}

func isNotFound(err error) bool {
	return err == sql.ErrNoRows
}
