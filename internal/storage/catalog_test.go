package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"homemaker-followup/internal/domain"
)

func TestCatalogEntities(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "catalog.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := store.SaveClient(ctx, domain.Client{ID: "c1", Name: "林女士", Phone: "138", PreferredChannel: domain.ChannelSMS, CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveServiceType(ctx, domain.ServiceType{ID: "s1", Name: "保洁", DefaultDays: 7}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveCaregiver(ctx, domain.Caregiver{ID: "g1", Name: "王阿姨", Phone: "139"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveReminderSetting(ctx, domain.ReminderSetting{ID: "default", DaysBefore: 3, Channel: domain.ChannelPhone, QuietStart: 21, QuietEnd: 8}); err != nil {
		t.Fatal(err)
	}
}
