package reporting

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"homemaker-followup/internal/domain"
)

type CaregiverSummary struct {
	Caregiver    string
	Records      int
	AverageScore float64
	Completed    int
}

func BuildCaregiverSummaries(records []domain.FollowUpRecord) []CaregiverSummary {
	groups := make(map[string][]domain.FollowUpRecord)
	for _, record := range records {
		groups[record.CaregiverName] = append(groups[record.CaregiverName], record)
	}
	result := make([]CaregiverSummary, 0, len(groups))
	for name, items := range groups {
		completed := 0
		for _, item := range items {
			if item.Status == "completed" {
				completed++
			}
		}
		result = append(result, CaregiverSummary{Caregiver: name, Records: len(items), AverageScore: domain.AverageScore(items), Completed: completed})
	}
	sort.Slice(result, func(left, right int) bool { return result[left].AverageScore > result[right].AverageScore })
	return result
}

func WeeklyTrend(records []domain.FollowUpRecord, start time.Time) [7]int {
	var trend [7]int
	for _, record := range records {
		delta := int(record.VisitDate.Sub(start).Hours() / 24)
		if delta >= 0 && delta < len(trend) {
			trend[delta]++
		}
	}
	return trend
}

func CSV(records []domain.FollowUpRecord) string {
	lines := []string{"记录编号,客户,服务类型,阿姨,满意度,状态,下次回访"}
	for _, record := range records {
		lines = append(lines, fmt.Sprintf("%s,%s,%s,%s,%d,%s,%s", record.ID, record.ClientName, record.ServiceTypeName, record.CaregiverName, record.Satisfaction.Score, record.Status, record.NextFollowUp.Format(time.DateOnly)))
	}
	return strings.Join(lines, "\n")
}

func StableFileName(prefix string, day time.Time) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "followup"
	}
	return fmt.Sprintf("%s-%s.json", prefix, day.Format("20060102"))
}

func ValidateRecords(records []domain.FollowUpRecord) error {
	seen := make(map[string]bool)
	for _, record := range records {
		if seen[record.ID] {
			return fmt.Errorf("duplicate record %s", record.ID)
		}
		if err := record.Validate(); err != nil {
			return err
		}
		seen[record.ID] = true
	}
	return nil
}
