package reporting

import (
	"path/filepath"
	"testing"

	"homemaker-followup/internal/domain"
)

func TestReportExport(t *testing.T) {
	record := domain.FollowUpRecord{ID: "r1", ServiceTypeName: "保洁", Satisfaction: domain.Satisfaction{Score: 5}}
	path := filepath.Join(t.TempDir(), "report.json")
	if err := ExportJSON(path, []domain.FollowUpRecord{record}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateExport(path); err != nil {
		t.Fatal(err)
	}
	if got := SatisfactionDistribution([]domain.FollowUpRecord{record})[5]; got != 1 {
		t.Fatalf("got %d", got)
	}
}
