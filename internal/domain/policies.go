package domain

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var validTransitions = map[string][]string{
	"scheduled": {"overdue", "completed", "cancelled"},
	"overdue":   {"scheduled", "completed", "cancelled"},
	"completed": {"scheduled"},
	"cancelled": {"scheduled"},
}

func ValidateCatalog(clients []Client, services []ServiceType, caregivers []Caregiver) error {
	clientIDs := make(map[string]bool, len(clients))
	for _, client := range clients {
		if err := client.Validate(); err != nil {
			return err
		}
		if clientIDs[client.ID] {
			return errors.New("duplicate client id")
		}
		clientIDs[client.ID] = true
	}
	serviceIDs := make(map[string]bool, len(services))
	for _, service := range services {
		if err := service.Validate(); err != nil {
			return err
		}
		if serviceIDs[service.ID] {
			return errors.New("duplicate service type id")
		}
		serviceIDs[service.ID] = true
	}
	caregiverIDs := make(map[string]bool, len(caregivers))
	for _, caregiver := range caregivers {
		if err := caregiver.Validate(); err != nil {
			return err
		}
		if caregiverIDs[caregiver.ID] {
			return errors.New("duplicate caregiver id")
		}
		caregiverIDs[caregiver.ID] = true
	}
	return nil
}

func MergeClient(existing, incoming Client) (Client, error) {
	if incoming.ID == "" {
		incoming.ID = existing.ID
	}
	if incoming.Name == "" {
		incoming.Name = existing.Name
	}
	if incoming.Phone == "" {
		incoming.Phone = existing.Phone
	}
	if incoming.Address == "" {
		incoming.Address = existing.Address
	}
	if incoming.PreferredChannel == "" {
		incoming.PreferredChannel = existing.PreferredChannel
	}
	if incoming.CreatedAt.IsZero() {
		incoming.CreatedAt = existing.CreatedAt
	}
	if err := incoming.Validate(); err != nil {
		return Client{}, err
	}
	return incoming, nil
}

func CalculateNextFollowUp(visit time.Time, service ServiceType, requested time.Time) (time.Time, error) {
	if err := service.Validate(); err != nil {
		return time.Time{}, err
	}
	if visit.IsZero() {
		return time.Time{}, errors.New("visit date is required")
	}
	if requested.IsZero() {
		return visit.AddDate(0, 0, service.DefaultDays), nil
	}
	if requested.Before(visit) {
		return time.Time{}, errors.New("requested date is before visit")
	}
	return requested, nil
}

func Transition(record FollowUpRecord, target string) (FollowUpRecord, error) {
	target = strings.ToLower(strings.TrimSpace(target))
	if target == record.Status {
		return record, nil
	}
	allowed := false
	for _, candidate := range validTransitions[record.Status] {
		if candidate == target {
			allowed = true
			break
		}
	}
	if !allowed {
		return FollowUpRecord{}, errors.New("status transition is not allowed")
	}
	record.Status = target
	return record, nil
}

func StatusAt(record FollowUpRecord, day time.Time) string {
	if record.Status == "cancelled" || record.Status == "completed" {
		return record.Status
	}
	if record.NextFollowUp.IsZero() {
		return "unscheduled"
	}
	if record.NextFollowUp.Before(day) {
		return "overdue"
	}
	if record.Satisfaction.Score <= 2 {
		return "needs-attention"
	}
	return "scheduled"
}

func SortClients(clients []Client) []Client {
	result := append([]Client(nil), clients...)
	sort.SliceStable(result, func(left, right int) bool {
		if result[left].Active != result[right].Active {
			return result[left].Active
		}
		return strings.ToLower(result[left].Name) < strings.ToLower(result[right].Name)
	})
	return result
}

func ContactLabel(channel Channel) string {
	switch channel {
	case ChannelPhone:
		return "电话"
	case ChannelWeChat:
		return "微信"
	case ChannelSMS:
		return "短信"
	default:
		return "未设置"
	}
}

func IsHighRisk(record FollowUpRecord) bool {
	return record.Satisfaction.Score <= 2 || record.Status == "overdue" || strings.Contains(record.Notes, "投诉")
}
