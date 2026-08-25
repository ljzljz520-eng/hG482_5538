package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"homemaker-followup/internal/domain"
)

func (s *Store) ListClients(ctx context.Context) ([]domain.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, fmt.Errorf("storage is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,name,phone,address,preferred_channel,active,created_at FROM clients ORDER BY active DESC,name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Client, 0)
	for rows.Next() {
		var item domain.Client
		var active int
		var created string
		if err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Address, &item.PreferredChannel, &active, &created); err != nil {
			return nil, err
		}
		item.Active = active != 0
		item.CreatedAt, err = time.Parse(time.RFC3339, created)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) ListServiceTypes(ctx context.Context) ([]domain.ServiceType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, fmt.Errorf("storage is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,name,description,default_days,active FROM service_types ORDER BY active DESC,name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.ServiceType, 0)
	for rows.Next() {
		var item domain.ServiceType
		var active int
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.DefaultDays, &active); err != nil {
			return nil, err
		}
		item.Active = active != 0
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) ListCaregivers(ctx context.Context, onlyAvailable bool) ([]domain.Caregiver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, fmt.Errorf("storage is closed")
	}
	query := `SELECT id,name,phone,skill_tags,available FROM caregivers ORDER BY name`
	if onlyAvailable {
		query = `SELECT id,name,phone,skill_tags,available FROM caregivers WHERE available=1 ORDER BY name`
	}
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Caregiver, 0)
	for rows.Next() {
		var item domain.Caregiver
		var tags string
		var available int
		if err := rows.Scan(&item.ID, &item.Name, &item.Phone, &tags, &available); err != nil {
			return nil, err
		}
		item.SkillTags = splitTags(tags)
		item.Available = available != 0
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) GetReminderSetting(ctx context.Context, id string) (domain.ReminderSetting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return domain.ReminderSetting{}, fmt.Errorf("storage is closed")
	}
	var item domain.ReminderSetting
	var enabled int
	err := s.db.QueryRowContext(ctx, `SELECT id,days_before,channel,enabled,quiet_start,quiet_end FROM reminder_settings WHERE id=?`, id).Scan(&item.ID, &item.DaysBefore, &item.Channel, &enabled, &item.QuietStart, &item.QuietEnd)
	item.Enabled = enabled != 0
	return item, err
}

func (s *Store) SearchRecords(ctx context.Context, query string) ([]domain.FollowUpRecord, error) {
	query = "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, fmt.Errorf("storage is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,client_id,client_name,service_type_id,service_type_name,caregiver_id,caregiver_name,visit_date,next_follow_up,score,comment,improvement,status,notes,updated_at FROM follow_up_records WHERE lower(id) LIKE ? OR lower(client_name) LIKE ? OR lower(caregiver_name) LIKE ? ORDER BY next_follow_up`, query, query, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.FollowUpRecord, 0)
	for rows.Next() {
		item, scanErr := scanRecord(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) CountRecords(ctx context.Context, status string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return 0, fmt.Errorf("storage is closed")
	}
	var count int
	var err error
	if status == "" {
		err = s.db.QueryRowContext(ctx, `SELECT count(*) FROM follow_up_records`).Scan(&count)
	} else {
		err = s.db.QueryRowContext(ctx, `SELECT count(*) FROM follow_up_records WHERE status=?`, status).Scan(&count)
	}
	return count, err
}

func splitTags(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			result = append(result, strings.TrimSpace(part))
		}
	}
	return result
}
