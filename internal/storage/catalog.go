package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"homemaker-followup/internal/domain"
)

func (s *Store) SaveClient(ctx context.Context, client domain.Client) error {
	if err := client.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("storage is closed")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO clients(id,name,phone,address,preferred_channel,active,created_at) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,phone=excluded.phone,address=excluded.address,preferred_channel=excluded.preferred_channel,active=excluded.active`,
		client.ID, client.Name, client.Phone, client.Address, client.PreferredChannel, client.Active, client.CreatedAt.Format(time.RFC3339))
	return err
}

func (s *Store) SaveServiceType(ctx context.Context, service domain.ServiceType) error {
	if err := service.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("storage is closed")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO service_types(id,name,description,default_days,active) VALUES(?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,description=excluded.description,default_days=excluded.default_days,active=excluded.active`, service.ID, service.Name, service.Description, service.DefaultDays, service.Active)
	return err
}

func (s *Store) SaveCaregiver(ctx context.Context, caregiver domain.Caregiver) error {
	if err := caregiver.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("storage is closed")
	}
	tags := strings.Join(domain.NormalizeTags(caregiver.SkillTags), ",")
	_, err := s.db.ExecContext(ctx, `INSERT INTO caregivers(id,name,phone,skill_tags,available) VALUES(?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,phone=excluded.phone,skill_tags=excluded.skill_tags,available=excluded.available`, caregiver.ID, caregiver.Name, caregiver.Phone, tags, caregiver.Available)
	return err
}

func (s *Store) SaveReminderSetting(ctx context.Context, setting domain.ReminderSetting) error {
	if err := setting.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("storage is closed")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO reminder_settings(id,days_before,channel,enabled,quiet_start,quiet_end) VALUES(?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET days_before=excluded.days_before,channel=excluded.channel,enabled=excluded.enabled,quiet_start=excluded.quiet_start,quiet_end=excluded.quiet_end`, setting.ID, setting.DaysBefore, setting.Channel, setting.Enabled, setting.QuietStart, setting.QuietEnd)
	return err
}
