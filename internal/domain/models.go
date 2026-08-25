package domain

import (
	"errors"
	"strings"
	"time"
)

type Channel string

const (
	ChannelPhone  Channel = "phone"
	ChannelWeChat Channel = "wechat"
	ChannelSMS    Channel = "sms"
)

type Client struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Phone            string    `json:"phone"`
	Address          string    `json:"address"`
	PreferredChannel Channel   `json:"preferred_channel"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
}

type ServiceType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DefaultDays int    `json:"default_days"`
	Active      bool   `json:"active"`
}

type Caregiver struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Phone     string   `json:"phone"`
	SkillTags []string `json:"skill_tags"`
	Available bool     `json:"available"`
}

type Satisfaction struct {
	Score       int    `json:"score"`
	Comment     string `json:"comment"`
	Improvement string `json:"improvement"`
}

type FollowUpRecord struct {
	ID              string       `json:"id"`
	ClientID        string       `json:"client_id"`
	ClientName      string       `json:"client_name"`
	ServiceTypeID   string       `json:"service_type_id"`
	ServiceTypeName string       `json:"service_type_name"`
	CaregiverID     string       `json:"caregiver_id"`
	CaregiverName   string       `json:"caregiver_name"`
	VisitDate       time.Time    `json:"visit_date"`
	NextFollowUp    time.Time    `json:"next_follow_up"`
	Satisfaction    Satisfaction `json:"satisfaction"`
	Status          string       `json:"status"`
	Notes           string       `json:"notes"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type ReminderSetting struct {
	ID         string  `json:"id"`
	DaysBefore int     `json:"days_before"`
	Channel    Channel `json:"channel"`
	Enabled    bool    `json:"enabled"`
	QuietStart int     `json:"quiet_start"`
	QuietEnd   int     `json:"quiet_end"`
}

type AuditEntry struct {
	ID        int64     `json:"id"`
	RecordID  string    `json:"record_id"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (c Client) Validate() error {
	if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Name) == "" {
		return errors.New("client id and name are required")
	}
	if strings.TrimSpace(c.Phone) == "" {
		return errors.New("client phone is required")
	}
	if c.PreferredChannel != ChannelPhone && c.PreferredChannel != ChannelWeChat && c.PreferredChannel != ChannelSMS {
		return errors.New("unsupported client channel")
	}
	return nil
}

func (s ServiceType) Validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Name) == "" {
		return errors.New("service type id and name are required")
	}
	if s.DefaultDays < 1 || s.DefaultDays > 365 {
		return errors.New("service type default days must be between 1 and 365")
	}
	return nil
}

func (c Caregiver) Validate() error {
	if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Name) == "" {
		return errors.New("caregiver id and name are required")
	}
	if strings.TrimSpace(c.Phone) == "" {
		return errors.New("caregiver phone is required")
	}
	return nil
}

func (s Satisfaction) Validate() error {
	if s.Score < 1 || s.Score > 5 {
		return errors.New("satisfaction score must be 1 through 5")
	}
	if strings.TrimSpace(s.Comment) == "" {
		return errors.New("satisfaction comment is required")
	}
	return nil
}

func (r FollowUpRecord) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.ClientID) == "" {
		return errors.New("follow-up id and client id are required")
	}
	if r.VisitDate.IsZero() || r.NextFollowUp.IsZero() {
		return errors.New("visit and next follow-up dates are required")
	}
	if r.NextFollowUp.Before(r.VisitDate) {
		return errors.New("next follow-up cannot precede visit")
	}
	if err := r.Satisfaction.Validate(); err != nil {
		return err
	}
	return nil
}

func (s ReminderSetting) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("reminder setting id is required")
	}
	if s.DaysBefore < 0 || s.DaysBefore > 30 {
		return errors.New("reminder days must be between 0 and 30")
	}
	if s.QuietStart < 0 || s.QuietStart > 23 || s.QuietEnd < 0 || s.QuietEnd > 23 {
		return errors.New("quiet hours must be valid")
	}
	return nil
}
