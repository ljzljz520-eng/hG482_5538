package excel

import (
	"path/filepath"
	"testing"
	"time"

	"homemaker-followup/internal/domain"
)

func TestRecordRowsRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.xlsx")
	workbook, err := NewWorkbook(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.FollowUpRecord{ID: "r1", ClientID: "c1", ClientName: "林女士", ServiceTypeID: "s1", ServiceTypeName: "日常保洁", CaregiverID: "g1", CaregiverName: "王阿姨", VisitDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), NextFollowUp: time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC), Satisfaction: domain.Satisfaction{Score: 4, Comment: "满意"}, Status: "scheduled"}
	if err := AppendRecord(workbook.File, record); err != nil {
		t.Fatal(err)
	}
	if err := workbook.Save(); err != nil {
		t.Fatal(err)
	}
	workbook.Close()
	opened, err := OpenWorkbook(path)
	if err != nil {
		t.Fatal(err)
	}
	defer opened.Close()
	rows, err := ReadRecords(opened.File)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].ID != "r1" {
		t.Fatalf("unexpected rows %#v", rows)
	}
}
