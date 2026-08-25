package followup

import (
	"path/filepath"
	"testing"
	"time"

	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/excel"
)

func TestPrimaryWorkflow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workflow.xlsx")
	workbook, err := excel.NewWorkbook(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := workbook.Save(); err != nil {
		t.Fatal(err)
	}
	workbook.Close()
	service := NewService(path, nil)
	if _, err := service.Start(); err != nil {
		t.Fatal(err)
	}
	record := domain.FollowUpRecord{ID: "r1", ClientID: "c1", ClientName: "林女士", ServiceTypeID: "s1", ServiceTypeName: "保洁", CaregiverID: "g1", CaregiverName: "王阿姨", VisitDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), NextFollowUp: time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC), Satisfaction: domain.Satisfaction{Score: 4, Comment: "满意"}, Status: "scheduled"}
	if err := service.AddRecord(record, "tester"); err != nil {
		t.Fatal(err)
	}
	if len(service.Records) != 1 {
		t.Fatal("record not added")
	}
}

func TestSecondaryWorkflow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workflow.xlsx")
	workbook, err := excel.NewWorkbook(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := workbook.Save(); err != nil {
		t.Fatal(err)
	}
	workbook.Close()
	service := NewService(path, nil)
	if _, err := service.Start(); err != nil {
		t.Fatal(err)
	}
	if got := service.Search("不存在"); len(got) != 0 {
		t.Fatalf("unexpected %v", got)
	}
}

func TestTertiaryWorkflow(t *testing.T) {
	if got := FormatNotice(nil); got.Blocking {
		t.Fatal("healthy notice should not block")
	}
	if got := FormatNotice(ErrStartupFormat); !got.Blocking {
		t.Fatal("error notice should block")
	}
}
