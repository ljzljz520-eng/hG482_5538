package reminders

import (
	"testing"
	"time"

	"homemaker-followup/internal/domain"
)

func TestReminderCalendar(t *testing.T) {
	today := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	record := domain.FollowUpRecord{ID: "r1", ClientName: "林女士", NextFollowUp: today.AddDate(0, 0, 2), Satisfaction: domain.Satisfaction{Score: 3}, Status: "scheduled"}
	items := Build([]domain.FollowUpRecord{record}, domain.ReminderSetting{DaysBefore: 3, Channel: domain.ChannelSMS, Enabled: true}, today)
	if len(items) != 1 {
		t.Fatalf("got %d reminders", len(items))
	}
	if len(SevenDayCalendar([]domain.FollowUpRecord{record}, today)) != 7 {
		t.Fatal("calendar should have seven days")
	}
}
