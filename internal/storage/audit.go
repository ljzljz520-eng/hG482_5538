package storage

import (
	"context"
	"fmt"
	"time"

	"homemaker-followup/internal/domain"
)

func (s *Store) AddAudit(ctx context.Context, entry domain.AuditEntry) error {
	if entry.RecordID == "" || entry.Action == "" || entry.Actor == "" {
		return fmt.Errorf("audit record, action and actor are required")
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Unix(0, 0).UTC()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("storage is closed")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_entries(record_id,action,actor,detail,created_at) VALUES(?,?,?,?,?)`, entry.RecordID, entry.Action, entry.Actor, entry.Detail, entry.CreatedAt.Format(time.RFC3339))
	return err
}

func (s *Store) ListAudit(ctx context.Context, recordID string) ([]domain.AuditEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, fmt.Errorf("storage is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,record_id,action,actor,detail,created_at FROM audit_entries WHERE record_id=? ORDER BY id`, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.AuditEntry, 0)
	for rows.Next() {
		var item domain.AuditEntry
		var created string
		if err := rows.Scan(&item.ID, &item.RecordID, &item.Action, &item.Actor, &item.Detail, &created); err != nil {
			return nil, err
		}
		item.CreatedAt, err = time.Parse(time.RFC3339, created)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
