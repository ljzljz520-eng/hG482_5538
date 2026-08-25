package domain

import (
	"testing"
	"time"
)

func TestRankingAndStatusCounts(t *testing.T) {
	base := FollowUpRecord{VisitDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), NextFollowUp: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), Status: "scheduled"}
	base.ID = "high"
	base.Satisfaction.Score = 1
	other := base
	other.ID = "low"
	other.Satisfaction.Score = 5
	if RankRecords([]FollowUpRecord{other, base})[0].ID != "high" {
		t.Fatal("priority ordering failed")
	}
	counts := CountByStatus([]FollowUpRecord{base, other})
	if counts["needs-attention"] != 1 || counts["scheduled"] != 1 {
		t.Fatalf("unexpected counts %#v", counts)
	}
}
