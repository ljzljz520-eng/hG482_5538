package dashboard

import (
	"time"

	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/reminders"
)

type Snapshot struct {
	GeneratedOn  time.Time
	TotalClients int
	TotalRecords int
	AverageScore float64
	StatusCounts map[string]int
	TopRecords   []domain.FollowUpRecord
	Reminders    []reminders.Reminder
	Calendar     []reminders.CalendarDay
}

func BuildSnapshot(clients []domain.Client, records []domain.FollowUpRecord, setting domain.ReminderSetting, today time.Time) Snapshot {
	marked := make([]domain.FollowUpRecord, 0, len(records))
	for _, record := range records {
		marked = append(marked, domain.MarkOverdue(record, today))
	}
	return Snapshot{
		GeneratedOn: today, TotalClients: countActiveClients(clients), TotalRecords: len(marked),
		AverageScore: domain.AverageScore(marked), StatusCounts: domain.CountByStatus(marked),
		TopRecords: limit(domain.RankRecords(marked), 5), Reminders: reminders.Build(marked, setting, today),
		Calendar: reminders.SevenDayCalendar(marked, today),
	}
}

func countActiveClients(clients []domain.Client) int {
	count := 0
	for _, client := range clients {
		if client.Active {
			count++
		}
	}
	return count
}

func limit(records []domain.FollowUpRecord, size int) []domain.FollowUpRecord {
	if len(records) <= size {
		return records
	}
	return records[:size]
}

func NeedsAttention(snapshot Snapshot) bool {
	return snapshot.StatusCounts["needs-attention"] > 0 || snapshot.StatusCounts["overdue"] > 0
}
