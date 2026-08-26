package dashboard

import (
	"testing"
	"time"

	"homemaker-followup/internal/domain"
)

func TestSnapshotAndFilter(t *testing.T) {
	today := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	record := domain.FollowUpRecord{ID: "r1", ClientName: "林女士", CaregiverName: "王阿姨", ServiceTypeName: "保洁", NextFollowUp: today, Satisfaction: domain.Satisfaction{Score: 2}, Status: "scheduled"}
	snapshot := BuildSnapshot([]domain.Client{{ID: "c1", Active: true}}, []domain.FollowUpRecord{record}, domain.ReminderSetting{Enabled: true, DaysBefore: 3}, today)
	if !NeedsAttention(snapshot) {
		t.Fatal("attention should be visible")
	}
	if len(ApplyFilter([]domain.FollowUpRecord{record}, Filter{MinimumScore: 3})) != 0 {
		t.Fatal("filter should exclude low score")
	}
}
