package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"homemaker-followup/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.sqlite")
	record := domain.FollowUpRecord{ID: "r1", ClientID: "c1", ClientName: "林女士", ServiceTypeID: "s1", ServiceTypeName: "日常保洁", CaregiverID: "g1", CaregiverName: "王阿姨", VisitDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), NextFollowUp: time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC), Satisfaction: domain.Satisfaction{Score: 5, Comment: "满意"}, Status: "scheduled", UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.UpsertRecord(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	got, err := reopened.GetRecord(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ClientName != record.ClientName || got.Satisfaction.Score != 5 {
		t.Fatalf("unexpected record %#v", got)
	}
}
