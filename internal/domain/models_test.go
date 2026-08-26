package domain

import (
	"testing"
	"time"
)

func TestClientAndRecordValidation(t *testing.T) {
	client := Client{ID: "c1", Name: "林女士", Phone: "138", PreferredChannel: ChannelWeChat}
	if err := client.Validate(); err != nil {
		t.Fatal(err)
	}
	record := FollowUpRecord{ID: "r1", ClientID: client.ID, VisitDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), NextFollowUp: time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC), Satisfaction: Satisfaction{Score: 5, Comment: "满意"}}
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
	record.NextFollowUp = record.VisitDate.AddDate(0, 0, -1)
	if record.Validate() == nil {
		t.Fatal("expected invalid date")
	}
}
