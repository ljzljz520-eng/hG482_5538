package reminders

import (
	"sort"
	"time"

	"homemaker-followup/internal/domain"
)

type Reminder struct {
	RecordID   string
	ClientName string
	DueDate    time.Time
	DaysUntil  int
	Priority   int
	Channel    domain.Channel
	Reason     string
}

func Build(records []domain.FollowUpRecord, setting domain.ReminderSetting, today time.Time) []Reminder {
	reminders := make([]Reminder, 0)
	for _, record := range records {
		if !setting.Enabled || record.NextFollowUp.IsZero() {
			continue
		}
		days := daysBetween(today, record.NextFollowUp)
		if days > setting.DaysBefore && days >= 0 {
			continue
		}
		priority := domain.PriorityForRecord(record)
		reason := "scheduled follow-up"
		if days < 0 {
			reason = "overdue follow-up"
			priority++
		}
		reminders = append(reminders, Reminder{RecordID: record.ID, ClientName: record.ClientName, DueDate: record.NextFollowUp, DaysUntil: days, Priority: priority, Channel: setting.Channel, Reason: reason})
	}
	sort.SliceStable(reminders, func(left, right int) bool {
		if reminders[left].Priority != reminders[right].Priority {
			return reminders[left].Priority > reminders[right].Priority
		}
		return reminders[left].DueDate.Before(reminders[right].DueDate)
	})
	return reminders
}

func daysBetween(from, to time.Time) int {
	fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	return int(toDate.Sub(fromDate).Hours() / 24)
}

func GroupByChannel(items []Reminder) map[domain.Channel][]Reminder {
	groups := make(map[domain.Channel][]Reminder)
	for _, item := range items {
		groups[item.Channel] = append(groups[item.Channel], item)
	}
	return groups
}

func MarkReminderSent(items []Reminder, recordID string) []Reminder {
	result := make([]Reminder, 0, len(items))
	for _, item := range items {
		if item.RecordID != recordID {
			result = append(result, item)
		}
	}
	return result
}
